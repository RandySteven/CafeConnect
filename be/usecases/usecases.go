package usecases

import (
	"github.com/RandySteven/CafeConnect/be/caches"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	aws_client "github.com/RandySteven/CafeConnect/be/pkg/aws"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	"github.com/RandySteven/CafeConnect/be/repositories"
	"github.com/RandySteven/CafeConnect/be/topics"
)

type Usecases struct {
	AddressUsecase     usecase_interfaces.AddressUsecase
	OnboardingUsecase  usecase_interfaces.OnboardingUsecase
	CafeUsecase        usecase_interfaces.CafeUsecase
	ProductUsecase     usecase_interfaces.ProductUsecase
	ReviewUsecase      usecase_interfaces.ReviewUsecase
	CartUsecase        usecase_interfaces.CartUsecase
	RoleUsecase        usecase_interfaces.RoleUsecase
	TransactionUsecase usecase_interfaces.TransactionUsecase
}

func NewUsecases(repo *repositories.Repositories,
	cache *caches.Caches,
	aws aws_client.AWS,
	topics *topics.Topics,
	midtrans midtrans_client.Midtrans) *Usecases {
	return &Usecases{
		AddressUsecase:     newAddressUsecase(repo.AddressRepository, repo.AddressUserRepository, repo.UserRepository, cache.AddressCache),
		OnboardingUsecase:  newOnboardingUsecase(repo.UserRepository, repo.PointRepository, repo.AddressRepository, repo.AddressUserRepository, repo.ReferralRepository, repo.VerifyTokenRepository, repo.Transaction, cache.OnboardCache, topics.OnboardingTopic, aws),
		RoleUsecase:        newRoleUsecase(repo.RoleRepository),
		CafeUsecase:        newCafeUsecase(repo.CafeRepository, repo.CafeFranchiseRepository, repo.AddressRepository, repo.Transaction, aws, cache.CafeCache),
		ProductUsecase:     newProductUsecase(repo.CafeRepository, repo.CafeFranchiseRepository, repo.CafeProductRepository, repo.ProductRepository, repo.ProductCategoryRepository, aws, repo.Transaction, cache.ProductCache),
		ReviewUsecase:      newReviewUsecase(repo.ReviewRepository, repo.CafeRepository, repo.UserRepository),
		CartUsecase:        newCartUsecase(repo.CafeRepository, repo.CafeProductRepository, repo.CartRepository, repo.ProductRepository, repo.UserRepository, repo.CafeFranchiseRepository, cache.ProductCache, repo.Transaction),
		TransactionUsecase: newTransactionUsecase(repo.TransactionHeaderRepository, repo.TransactionDetailRepository, repo.AddressRepository, repo.CartRepository, repo.UserRepository, repo.CafeRepository, repo.CafeFranchiseRepository, repo.ProductRepository, repo.MidtransTransactionRepository, repo.CafeProductRepository, repo.Transaction, cache.TransactionCache, cache.ProductCache, cache.CheckoutCache, topics.TransactionTopic, midtrans),
	}
}
