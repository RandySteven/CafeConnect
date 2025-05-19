package usecases

import (
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/apperror"
	"github.com/RandySteven/CafeConnect/be/entities/models"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"github.com/google/uuid"
	"time"
)

type cartUsecase struct {
	cafeRepository        repository_interfaces.CafeRepository
	cafeProductRepository repository_interfaces.CafeProductRepository
	cartRepository        repository_interfaces.CartRepository
	productRepository     repository_interfaces.ProductRepository
	userRepository        repository_interfaces.UserRepository
	transaction           repository_interfaces.Transaction
}

func (c *cartUsecase) AddToCart(ctx context.Context, request *requests.AddToCartRequest) (result *responses.AddCartResponse, customErr *apperror.CustomError) {
	userId := ctx.Value(enums.UserID).(uint64)
	cart := &models.Cart{
		UserID: userId,
	}
	var err error
	action := `CREATED`

	if customErr = c.transaction.RunInTx(ctx, func(ctx context.Context) (customErr *apperror.CustomError) {
		//pre cond check duplicate
		//if it duplicate or cafe_product already exists on cart then the qty will update
		cart, err = c.cartRepository.FindByUserIDAndCafeProductID(ctx, userId, request.CafeProductID)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cart`, err)
		}
		if cart != nil {
			cart.Qty += request.Qty
			cart.UpdatedAt = time.Now()
			cart, err = c.cartRepository.Update(ctx, cart)
			if err != nil {
				return apperror.NewCustomError(apperror.ErrInternalServer, `failed to update cart`, err)
			}
			action = `UPDATED`
			return nil
		}

		//1. check product stock ready
		cafeProduct, err := c.cafeProductRepository.FindByID(ctx, request.CafeProductID)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe product`, err)
		}
		if cafeProduct.Stock <= request.Qty {
			return apperror.NewCustomError(apperror.ErrBadRequest, `the stock is less than qty`, fmt.Errorf(`insufficient stock`))
		}

		//2. insert to cart
		cart.CafeProductID = request.CafeProductID
		cart.Qty = request.Qty
		cart, err = c.cartRepository.Save(ctx, cart)
		if err != nil {
			return apperror.NewCustomError(apperror.ErrInternalServer, `failed to create cart`, err)
		}

		return nil
	}); customErr != nil {
		return nil, customErr
	}

	result = &responses.AddCartResponse{
		ID:        uuid.NewString(),
		Action:    action,
		CreatedAt: time.Now(),
	}
	return result, nil
}

func (c *cartUsecase) GetUserCart(ctx context.Context) (result *responses.ListCartResponse, customErr *apperror.CustomError) {
	userId := ctx.Value(enums.UserID).(uint64)

	carts, err := c.cartRepository.FindByUserID(ctx, userId)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get carts`, err)
	}

	items := make([]*responses.CartItem, len(carts))
	for index, cart := range carts {
		cafeProduct, err := c.cafeProductRepository.FindByID(ctx, cart.ID)
		if err != nil {
			return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get cafe product`, err)
		}

		product, err := c.productRepository.FindByID(ctx, cafeProduct.ProductID)
		if err != nil {
			return nil, apperror.NewCustomError(apperror.ErrInternalServer, `failed to get product`, err)
		}

		items[index] = &responses.CartItem{
			ProductID:    cafeProduct.ID,
			ProductName:  product.Name,
			ProductPrice: cafeProduct.Price,
			Qty:          cart.Qty,
			CreatedAt:    cart.CreatedAt,
			UpdatedAt:    cart.UpdatedAt,
			DeletedAt:    cart.DeletedAt,
		}
	}

	result = &responses.ListCartResponse{
		UserID: userId,
		Items:  items,
	}
	return result, nil
}

var _ usecase_interfaces.CartUsecase = &cartUsecase{}

func newCartUsecase(
	cafeRepository repository_interfaces.CafeRepository,
	cafeProductRepository repository_interfaces.CafeProductRepository,
	cartRepository repository_interfaces.CartRepository,
	productRepository repository_interfaces.ProductRepository,
	userRepository repository_interfaces.UserRepository,
	transaction repository_interfaces.Transaction,
) *cartUsecase {
	return &cartUsecase{
		cartRepository:        cartRepository,
		cafeRepository:        cafeRepository,
		cafeProductRepository: cafeProductRepository,
		productRepository:     productRepository,
		userRepository:        userRepository,
		transaction:           transaction,
	}
}
