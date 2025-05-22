package apis

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/enums"
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ProductApi struct {
	usecase usecase_interfaces.ProductUsecase
}

func (p *ProductApi) AddProduct(w http.ResponseWriter, r *http.Request) {
	var (
		rID                  = uuid.NewString()
		ctx                  = context.WithValue(r.Context(), enums.RequestID, rID)
		productCategoryID, _ = strconv.Atoi(r.FormValue("product_category_id"))
		franchiseID, _       = strconv.Atoi(r.FormValue("franchise_id"))
		price, _             = strconv.Atoi(r.FormValue("price"))
		request              = &requests.AddProductRequest{
			Name:              r.FormValue("name"),
			ProductCategoryID: uint64(productCategoryID),
			FranchiseID:       uint64(franchiseID),
			Price:             uint64(price),
		}
		dataKey = `result`
	)

	imageFile, fileHeader, err := r.FormFile("photo")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer imageFile.Close()

	request.Photo = imageFile
	if err = utils.BindRequest(r, request); err != nil {
		utils.ResponseHandler(w, http.StatusBadRequest, `failed to proceed request`, nil, nil, err)
		return
	}
	ctx2 := context.WithValue(ctx, enums.FileHeader, fileHeader)
	ctx = ctx2

	result, customErr := p.usecase.AddProduct(ctx, request)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusCreated, `success register product`, &dataKey, result, nil)
}

func (p *ProductApi) GetListOfProducts(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		request = &requests.GetProductListByCafeIDRequest{}
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		dataKey = `menus`
	)
	if err := utils.BindRequest(r, request); err != nil {
		utils.ResponseHandler(w, http.StatusBadRequest, `failed to proceed request`, nil, nil, err)
		return
	}
	result, customErr := p.usecase.GetProductByCafe(ctx, request)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusOK, `success get menu list`, &dataKey, result, nil)

}

func (p *ProductApi) GetProductDetail(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		vars    = mux.Vars(r)
		dataKey = `product`
	)
	id, _ := strconv.Atoi(vars[`id`])
	result, customErr := p.usecase.GetProductDetail(ctx, uint64(id))
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusOK, `success get product`, &dataKey, result, nil)
}

var _ api_interfaces.ProductApi = &ProductApi{}

func newProductApi(usecase usecase_interfaces.ProductUsecase) *ProductApi {
	return &ProductApi{
		usecase: usecase,
	}
}
