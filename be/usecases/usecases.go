package usecases

import (
	"github.com/RandySteven/CafeConnect/be/caches"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	aws_client "github.com/RandySteven/CafeConnect/be/pkg/aws"
	kafka_client "github.com/RandySteven/CafeConnect/be/pkg/kafka"
	midtrans_client "github.com/RandySteven/CafeConnect/be/pkg/midtrans"
	storage_client "github.com/RandySteven/CafeConnect/be/pkg/storage"
	"github.com/RandySteven/CafeConnect/be/repositories"
)

type Usecases struct {
	OnboardingUsecase  usecase_interfaces.OnboardingUsecase
	CafeUsecase        usecase_interfaces.CafeUsecase
	ProductUsecase     usecase_interfaces.ProductUsecase
	ReviewUsecase      usecase_interfaces.ReviewUsecase
	CartUsecase        usecase_interfaces.CartUsecase
	TransactionUsecase usecase_interfaces.TransactionUsecase
}

func NewUsecases(repo *repositories.Repositories,
	cache *caches.Caches,
	googleStorage storage_client.GoogleStorage,
	aws aws_client.AWS,
	pub kafka_client.Publisher,
	sub kafka_client.Consumer,
	midtrans midtrans_client.Midtrans) *Usecases {
	return &Usecases{
		OnboardingUsecase:  newOnboardingUsecase(repo.UserRepository, repo.PointRepository, repo.AddressRepository, repo.AddressUserRepository, repo.ReferralRepository, repo.Transaction, cache.OnboardCache, pub, aws),
		CafeUsecase:        newCafeUsecase(repo.CafeRepository, repo.CafeFranchiseRepository, repo.AddressRepository, repo.Transaction, googleStorage, aws, cache.CafeCache),
		ProductUsecase:     newProductUsecase(repo.CafeRepository, repo.CafeFranchiseRepository, repo.CafeProductRepository, repo.ProductRepository, repo.ProductCategoryRepository, aws, repo.Transaction, cache.ProductCache),
		ReviewUsecase:      newReviewUsecase(repo.ReviewRepository, repo.CafeRepository, repo.UserRepository),
		CartUsecase:        newCartUsecase(repo.CafeRepository, repo.CafeProductRepository, repo.CartRepository, repo.ProductRepository, repo.UserRepository, repo.Transaction),
		TransactionUsecase: newTransactionUsecase(repo.TransactionHeaderRepository, repo.TransactionDetailRepository, repo.CartRepository, repo.UserRepository, repo.CafeRepository, repo.ProductRepository, repo.CafeProductRepository, repo.Transaction, cache.TransactionCache, midtrans),
	}
}
