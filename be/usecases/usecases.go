package usecases

import (
	"github.com/RandySteven/CafeConnect/be/caches"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	aws_client "github.com/RandySteven/CafeConnect/be/pkg/aws"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	temporal_client "github.com/RandySteven/CafeConnect/be/pkg/temporal"
	"github.com/RandySteven/CafeConnect/be/repositories"
	"github.com/RandySteven/CafeConnect/be/topics"
	auto_transfer_usecases "github.com/RandySteven/CafeConnect/be/usecases/auto"
	midtrans_usecases "github.com/RandySteven/CafeConnect/be/usecases/midtrans"
	transactions_usecases "github.com/RandySteven/CafeConnect/be/usecases/transactions"
)

type Usecases struct {
	AddressUsecase       usecase_interfaces.AddressUsecase
	OnboardingUsecase    usecase_interfaces.OnboardingUsecase
	CafeUsecase          usecase_interfaces.CafeUsecase
	ProductUsecase       usecase_interfaces.ProductUsecase
	ReviewUsecase        usecase_interfaces.ReviewUsecase
	CartUsecase          usecase_interfaces.CartUsecase
	RoleUsecase          usecase_interfaces.RoleUsecase
	TransactionUsecase   usecase_interfaces.TransactionUsecase
	TransactionWorkflow  transactions_usecases.TransactionWorkflow
	MidtransWorkflow     midtrans_usecases.MidtransWorkflow
	AutoTransferWorkflow auto_transfer_usecases.AutoTransferWorkflow
}

func NewUsecases(repo *repositories.Repositories,
	cache *caches.Caches,
	aws aws_client.AWS,
	topics *topics.Topics,
	midtrans midtrans_client.Midtrans,
	temporal temporal_client.Temporal,
) *Usecases {
	return &Usecases{
		AddressUsecase:       newAddressUsecase(repo.AddressRepository, repo.AddressUserRepository, repo.UserRepository, cache.AddressCache),
		OnboardingUsecase:    newOnboardingUsecase(repo.UserRepository, repo.PointRepository, repo.AddressRepository, repo.AddressUserRepository, repo.ReferralRepository, repo.VerifyTokenRepository, repo.Transaction, cache.OnboardCache, topics.OnboardingTopic, aws),
		RoleUsecase:          newRoleUsecase(repo.RoleRepository),
		CafeUsecase:          newCafeUsecase(repo.CafeRepository, repo.CafeFranchiseRepository, repo.AddressRepository, repo.Transaction, aws, cache.CafeCache),
		ProductUsecase:       newProductUsecase(repo.CafeRepository, repo.CafeFranchiseRepository, repo.CafeProductRepository, repo.ProductRepository, repo.ProductCategoryRepository, aws, repo.Transaction, cache.ProductCache),
		ReviewUsecase:        newReviewUsecase(repo.ReviewRepository, repo.CafeRepository, repo.UserRepository),
		CartUsecase:          newCartUsecase(repo.CafeRepository, repo.CafeProductRepository, repo.CartRepository, repo.ProductRepository, repo.UserRepository, repo.CafeFranchiseRepository, cache.ProductCache, repo.Transaction),
		TransactionUsecase:   newTransactionUsecase(repo.TransactionHeaderRepository, repo.TransactionDetailRepository, repo.AddressRepository, repo.CartRepository, repo.UserRepository, repo.CafeRepository, repo.CafeFranchiseRepository, repo.ProductRepository, repo.MidtransTransactionRepository, repo.CafeProductRepository, repo.Transaction, cache.TransactionCache, cache.ProductCache, cache.CheckoutCache, topics.TransactionTopic, midtrans),
		TransactionWorkflow:  transactions_usecases.NewTransactionWorkflow(repo.TransactionHeaderRepository, repo.TransactionDetailRepository, repo.AddressRepository, repo.CartRepository, repo.UserRepository, repo.CafeRepository, repo.CafeFranchiseRepository, repo.ProductRepository, repo.CafeProductRepository, repo.Transaction, topics.TransactionTopic, repo.MidtransTransactionRepository, midtrans, cache.TransactionCache, cache.ProductCache, cache.CheckoutCache, temporal_client.NewWorkflowExecutionData(temporal), temporal),
		MidtransWorkflow:     midtrans_usecases.NewMidtransWorkflow(temporal, temporal_client.NewWorkflowExecutionData(temporal), repo.TransactionHeaderRepository, repo.MidtransTransactionRepository, repo.TransactionDetailRepository, repo.CafeProductRepository, repo.ProductRepository, midtrans),
		AutoTransferWorkflow: auto_transfer_usecases.NewAutoTransferWorkflow(temporal, temporal_client.NewWorkflowExecutionData(temporal), repo.TransactionHeaderRepository, repo.TransactionDetailRepository, repo.AddressRepository, repo.CartRepository, repo.UserRepository, repo.CafeRepository, repo.CafeFranchiseRepository, repo.ProductRepository, repo.CafeProductRepository, repo.Transaction, topics.TransactionTopic, repo.MidtransTransactionRepository, midtrans, cache.TransactionCache, cache.ProductCache, cache.CheckoutCache),
	}
}
