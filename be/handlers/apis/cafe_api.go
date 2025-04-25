package apis

import (
	"bytes"
	"context"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/requests"
	"github.com/RandySteven/CafeConnect/be/enums"
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strconv"
)

type CafeApi struct {
	usecase usecase_interfaces.CafeUsecase
}

func (c *CafeApi) GetListCafeFranchise(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		dataKey = `franchises`
	)
	result, customErr := c.usecase.GetListCafeFranchises(ctx)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusOK, `success get all franchises`, &dataKey, result, nil)
}

func (c *CafeApi) RegisterCafeAndFranchise(w http.ResponseWriter, r *http.Request) {
	latitude, _ := strconv.ParseFloat(r.FormValue(`latitude`), 64)
	longitude, _ := strconv.ParseFloat(r.FormValue(`longitude`), 64)
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		request = &requests.RegisterCafeAndFranchiseRequest{
			Name:      r.FormValue(`name`),
			Address:   r.FormValue(`address`),
			Latitude:  latitude,
			Longitude: longitude,
			OpenHour:  r.FormValue(`open_hour`),
			CloseHour: r.FormValue(`close_hour`),
		}
		dataKey = `result`
	)
	logoFile, fileHeader, err := r.FormFile(`logo_file`)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer logoFile.Close()

	request.LogoFile = logoFile

	var photoFiles []io.Reader
	files := r.MultipartForm.File["photo_urls[]"]
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			http.Error(w, "failed to open uploaded photo", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		var buf bytes.Buffer
		_, err = io.Copy(&buf, file)
		if err != nil {
			http.Error(w, "failed to read uploaded photo", http.StatusInternalServerError)
			return
		}
		photoFiles = append(photoFiles, bytes.NewReader(buf.Bytes()))
	}
	request.PhotoFiles = photoFiles

	if err := utils.BindRequest(r, request); err != nil {
		utils.ResponseHandler(w, http.StatusBadRequest, `bad request`, nil, nil, err)
		return
	}

	ctx2 := context.WithValue(ctx, enums.FileHeader, fileHeader)
	ctx = ctx2

	result, customErr := c.usecase.RegisterCafeAndFranchise(ctx, request)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusOK, `success register new cafe`, &dataKey, result, nil)
}

func (c *CafeApi) GetCafeDetail(w http.ResponseWriter, r *http.Request) {
}

func (c *CafeApi) GetCafeProducts(w http.ResponseWriter, r *http.Request) {
}

func (c *CafeApi) GetListOfCafeBasedOnRadius(w http.ResponseWriter, r *http.Request) {
}

var _ api_interfaces.CafeApi = &CafeApi{}

func newCafeApi(usecase usecase_interfaces.CafeUsecase) *CafeApi {
	return &CafeApi{
		usecase: usecase,
	}
}
