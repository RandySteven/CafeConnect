package repositories

import (
	"database/sql"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
)

type Repositories struct {
	UserRepository                repository_interfaces.UserRepository
	PointRepository               repository_interfaces.PointRepository
	AddressRepository             repository_interfaces.AddressRepository
	AddressUserRepository         repository_interfaces.AddressUserRepository
	ReferralRepository            repository_interfaces.ReferralRepository
	CafeRepository                repository_interfaces.CafeRepository
	CafeFranchiseRepository       repository_interfaces.CafeFranchiseRepository
	ProductRepository             repository_interfaces.ProductRepository
	ProductCategoryRepository     repository_interfaces.ProductCategoryRepository
	CafeProductRepository         repository_interfaces.CafeProductRepository
	ReviewRepository              repository_interfaces.ReviewRepository
	CartRepository                repository_interfaces.CartRepository
	RoleRepository                repository_interfaces.RoleRepository
	TransactionHeaderRepository   repository_interfaces.TransactionHeaderRepository
	TransactionDetailRepository   repository_interfaces.TransactionDetailRepository
	MidtransTransactionRepository repository_interfaces.MidtransTransactionRepository
	VerifyTokenRepository         repository_interfaces.VerifyTokenRepository
	Transaction                   repository_interfaces.Transaction
}

func NewRepositories(db *sql.DB) *Repositories {
	transaction, dbx := newTransaction(db)
	return &Repositories{
		UserRepository:                newUserRepository(dbx),
		PointRepository:               newPointRepository(dbx),
		AddressRepository:             newAddressRepository(dbx),
		AddressUserRepository:         newAddressUserRepository(dbx),
		ReferralRepository:            newReferralRepository(dbx),
		CafeFranchiseRepository:       newCafeFranchiseRepository(dbx),
		CafeRepository:                newCafeRepository(dbx),
		ProductCategoryRepository:     newProductCategoryRepository(dbx),
		ProductRepository:             newProductRepository(dbx),
		CafeProductRepository:         newCafeProductRepository(dbx),
		ReviewRepository:              newReviewRepository(dbx),
		CartRepository:                newCartRepository(dbx),
		RoleRepository:                newRoleRepository(dbx),
		TransactionHeaderRepository:   newTransactionHeaderRepository(dbx),
		TransactionDetailRepository:   newTransactionDetailRepository(dbx),
		MidtransTransactionRepository: newMidtransTransactionRepository(dbx),
		VerifyTokenRepository:         newVerifyTokenRepository(dbx),
		Transaction:                   transaction,
	}
}
