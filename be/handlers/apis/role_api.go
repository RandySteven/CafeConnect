package apis

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/enums"
	api_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/handlers/apis"
	usecase_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/usecases"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/google/uuid"
	"net/http"
)

type RoleApi struct {
	usecase usecase_interfaces.RoleUsecase
}

func (r2 *RoleApi) GetRoleList(w http.ResponseWriter, r *http.Request) {
	var (
		rID     = uuid.NewString()
		ctx     = context.WithValue(r.Context(), enums.RequestID, rID)
		dataKey = `roles`
	)

	result, customErr := r2.usecase.GetRoleList(ctx)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.LogMessage, nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusOK, `success get roles`, &dataKey, result, nil)
}

var _ api_interfaces.RoleApi = &RoleApi{}

func newRoleApi(
	usecase usecase_interfaces.RoleUsecase,
) *RoleApi {
	return &RoleApi{
		usecase: usecase,
	}
}
