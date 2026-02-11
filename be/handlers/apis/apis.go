package apis

import (
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	"github.com/RandySteven/CafeConnect/be/usecases"
)

type APIs struct {
	DevApi         api_interfaces.DevApi
	AddressApi     api_interfaces.AddressApi
	OnboardingApi  api_interfaces.OnboardingApi
	CafeApi        api_interfaces.CafeApi
	ProductApi     api_interfaces.ProductApi
	ReviewApi      api_interfaces.ReviewApi
	CartApi        api_interfaces.CartApi
	RoleApi        api_interfaces.RoleApi
	TransactionApi api_interfaces.TransactionApi
}

func NewAPIs(usecases *usecases.Usecases) *APIs {
	return &APIs{
		DevApi:         newDevApi(),
		AddressApi:     newAddressApi(usecases.AddressUsecase),
		OnboardingApi:  newOnboardingApi(usecases.OnboardingUsecase),
		CafeApi:        newCafeApi(usecases.CafeUsecase),
		ProductApi:     newProductApi(usecases.ProductUsecase),
		ReviewApi:      newReviewApi(usecases.ReviewUsecase),
		RoleApi:        newRoleApi(usecases.RoleUsecase),
		CartApi:        newCartApi(usecases.CartUsecase),
		TransactionApi: newTransactionApi(usecases.TransactionUsecase, usecases.TransactionWorkflow, usecases.AutoTransferWorkflow),
	}
}
