package transactions_usecases

import (
	"context"
	"fmt"
	"log"

	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	temporal_client "github.com/RandySteven/CafeConnect/be/pkg/temporal"
	"github.com/google/uuid"
)

const (
	checkUserActivity             = "CheckUser"
	checkCafeActivity             = "CheckCafe"
	checkFranchiseActivity        = "CheckFranchise"
	saveTransactionHeaderActivity = "SaveTransactionHeader"
	publishTransactionActivity    = "PublishTransaction"
	checkStatusActivity           = "CheckStatus"
)

type (
	TransactionWorkflow interface {
		CheckoutTransactionV3(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError)
		PaymentConfirmation(ctx context.Context, request *requests.PaymentConfirmationRequest) (result []*responses.PaymentConfirmationResponse, customErr *apperror.CustomError)
	}

	transactionWorkflow struct {
		workflow                      temporal_client.Workflow
		transactionHeaderRepository   repository_interfaces.TransactionHeaderRepository
		transactionDetailRepository   repository_interfaces.TransactionDetailRepository
		addressRepository             repository_interfaces.AddressRepository
		cartRepository                repository_interfaces.CartRepository
		userRepository                repository_interfaces.UserRepository
		cafeRepository                repository_interfaces.CafeRepository
		cafeFranchiseRepository       repository_interfaces.CafeFranchiseRepository
		productRepository             repository_interfaces.ProductRepository
		cafeProductRepository         repository_interfaces.CafeProductRepository
		transaction                   repository_interfaces.Transaction
		transactionTopic              topics_interfaces.TransactionTopic
		midtransTransactionRepository repository_interfaces.MidtransTransactionRepository
		midtrans                      midtrans_client.Midtrans
		transactionCache              cache_interfaces.TransactionCache
		productCache                  cache_interfaces.ProductCache
		checkoutCache                 cache_interfaces.CheckoutCache
	}
)

func (t *transactionWorkflow) registerWorkflowAndActivities() {
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkUserActivity,
		Fn:   t.checkUser,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkCafeActivity,
		Fn:   t.checkCafe,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkFranchiseActivity,
		Fn:   t.checkFranchise,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: saveTransactionHeaderActivity,
		Fn:   t.saveTransactionHeader,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: publishTransactionActivity,
		Fn:   t.publishTransaction,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkStatusActivity,
		Fn:   t.checkStatus,
	})
	t.workflow.RegisterWorkflow(temporal_client.WorkflowDefinition{
		Name: "CreateTransaction",
		Fn:   t.transactionWorkflow,
	})
}

func (t *transactionWorkflow) CheckoutTransactionV3(ctx context.Context, request *requests.CreateTransactionRequest) (result *responses.TransactionReceiptResponse, customErr *apperror.CustomError) {
	userID, ok := ctx.Value(enums.UserID).(uint64)
	if !ok {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get user id from context`, fmt.Errorf("user id not found in context"))
	}

	correlationID := uuid.NewString()

	// 1. Start Transaction workflow and wait for it to complete
	txWorkflowID := fmt.Sprintf("CreateTransaction-%s", correlationID)
	txRun, err := t.workflow.StartWorkflow(ctx, temporal_client.StartWorkflowOptions{
		WorkflowID: txWorkflowID,
	}, t.transactionWorkflow, userID, request)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to start transaction workflow`, err)
	}
	log.Println("Transaction workflow started:", txRun.GetID())

	var txResult transactionResult
	err = t.workflow.GetWorkflowResult(context.Background(), txRun.GetID(), txRun.GetRunID(), &txResult)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get transaction workflow result`, err)
	}

	// 2. Start Midtrans workflow (it waits for signal)
	midtransWorkflowID := fmt.Sprintf("MidtransTransaction-%s", correlationID)
	midtransRun, err := t.workflow.StartWorkflow(ctx, temporal_client.StartWorkflowOptions{
		WorkflowID: midtransWorkflowID,
	}, "MidtransTransaction")
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to start midtrans workflow`, err)
	}
	log.Println("Midtrans workflow started:", midtransWorkflowID)

	// 3. Signal the Midtrans workflow with the transaction data
	err = t.workflow.SignalWorkflow(ctx, midtransRun.GetID(), midtransRun.GetRunID(), "MidtransTransaction", &messages.TransactionMidtransMessage{
		UserID:            txResult.UserID,
		FName:             txResult.FName,
		LName:             txResult.LName,
		Email:             txResult.Email,
		Phone:             txResult.Phone,
		TransactionCode:   txResult.Receipt.TransactionCode,
		CafeID:            txResult.CafeID,
		CafeFranchiseName: txResult.CafeFranchiseName,
		CheckoutList:      request.Checkouts,
	})
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to signal midtrans workflow`, err)
	}

	// 4. Wait for Midtrans workflow to complete and get the result
	var midtransResponse *midtrans_client.MidtransResponse
	err = t.workflow.GetWorkflowResult(context.Background(), midtransRun.GetID(), midtransRun.GetRunID(), &midtransResponse)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get midtrans workflow result`, err)
	}

	// 5. Attach Midtrans response to receipt
	result = txResult.Receipt
	result.MidtransResponse = midtransResponse

	return result, nil
}

func (t *transactionWorkflow) PaymentConfirmation(ctx context.Context, request *requests.PaymentConfirmationRequest) (result []*responses.PaymentConfirmationResponse, customErr *apperror.CustomError) {
	transactionHeader, err := t.transactionHeaderRepository.FindByTransactionCode(ctx, request.TransactionCode)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get transaction header`, err)
	}

	transactionDetails, err := t.transactionDetailRepository.FindByTransactionId(ctx, transactionHeader.ID)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get transaction details`, err)
	}

	paymentConfirmationResponses := make([]*responses.PaymentConfirmationResponse, len(transactionDetails))
	for index, detail := range transactionDetails {
		cafeProduct, err := t.cafeProductRepository.FindByID(ctx, detail.CafeProductID)
		if err != nil {
			return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe product`, err)
		}
		product, err := t.productRepository.FindByID(ctx, cafeProduct.ProductID)
		if err != nil {
			return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get product`, err)
		}
		paymentConfirmationResponses[index] = &responses.PaymentConfirmationResponse{
			CafeProductID:   cafeProduct.ID,
			ProductName:     product.Name,
			ProductPerPrice: cafeProduct.Price,
			ProductPrice:    cafeProduct.Price * detail.Qty,
			ProductImage:    product.PhotoURL,
			CurrentStock:    cafeProduct.Stock,
			PrevStock:       cafeProduct.Stock + detail.Qty,
			Qty:             detail.Qty,
		}
	}
	return paymentConfirmationResponses, nil
}

func NewTransactionWorkflow(
	transactionHeaderRepository repository_interfaces.TransactionHeaderRepository,
	transactionDetailRepository repository_interfaces.TransactionDetailRepository,
	addressRepository repository_interfaces.AddressRepository,
	cartRepository repository_interfaces.CartRepository,
	userRepository repository_interfaces.UserRepository,
	cafeRepository repository_interfaces.CafeRepository,
	cafeFranchiseRepository repository_interfaces.CafeFranchiseRepository,
	productRepository repository_interfaces.ProductRepository,
	cafeProductRepository repository_interfaces.CafeProductRepository,
	transaction repository_interfaces.Transaction,
	transactionTopic topics_interfaces.TransactionTopic,
	midtransTransactionRepository repository_interfaces.MidtransTransactionRepository,
	midtrans midtrans_client.Midtrans,
	transactionCache cache_interfaces.TransactionCache,
	productCache cache_interfaces.ProductCache,
	checkoutCache cache_interfaces.CheckoutCache,
	workflow temporal_client.Workflow,
) TransactionWorkflow {
	tw := &transactionWorkflow{
		transactionHeaderRepository:   transactionHeaderRepository,
		transactionDetailRepository:   transactionDetailRepository,
		addressRepository:             addressRepository,
		cartRepository:                cartRepository,
		userRepository:                userRepository,
		cafeRepository:                cafeRepository,
		cafeFranchiseRepository:       cafeFranchiseRepository,
		productRepository:             productRepository,
		cafeProductRepository:         cafeProductRepository,
		transaction:                   transaction,
		transactionTopic:              transactionTopic,
		midtransTransactionRepository: midtransTransactionRepository,
		midtrans:                      midtrans,
		transactionCache:              transactionCache,
		productCache:                  productCache,
		checkoutCache:                 checkoutCache,
		workflow:                      workflow,
	}
	tw.registerWorkflowAndActivities()
	return tw
}
