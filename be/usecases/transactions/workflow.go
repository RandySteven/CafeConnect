package transactions_usecases

import (
	"context"
	"fmt"
	"log"

	"github.com/RandySteven/CafeConnect/be/apperror"
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
	checkUserActivity             = "TransactionCheckUser"
	checkCafeActivity             = "TransactionCheckCafe"
	checkFranchiseActivity        = "TransactionCheckFranchise"
	saveTransactionHeaderActivity = "TransactionSaveTransactionHeader"
	publishTransactionActivity    = "TransactionPublishTransaction"
	checkStatusActivity           = "TransactionCheckStatus"
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
		Fn:   t.transactionCheckUser,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkCafeActivity,
		Fn:   t.transactionCheckCafe,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkFranchiseActivity,
		Fn:   t.transactionCheckFranchise,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: saveTransactionHeaderActivity,
		Fn:   t.transactionSaveTransactionHeader,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: publishTransactionActivity,
		Fn:   t.publishTransaction,
	})
	t.workflow.RegisterActivity(temporal_client.ActivityDefinition{
		Name: checkStatusActivity,
		Fn:   t.transactionCheckUser,
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

	// 2. Wait for the workflow to complete (includes child MidtransTransaction)
	var txResult transactionResult
	err = t.workflow.GetWorkflowResult(context.Background(), txRun.GetID(), txRun.GetRunID(), &txResult)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get transaction workflow result`, err)
	}

	return txResult.Receipt, nil
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
