package attendance

import (
	"log"
	"net/http"

	"go-rest-api/src/constant"
	entity "go-rest-api/src/http"
	"go-rest-api/src/pkg/jwt"
	"go-rest-api/src/service/v1/attendance"

	"github.com/forkyid/go-utils/v1/rest"
	"github.com/forkyid/go-utils/v1/validation"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Controller struct {
	svc attendance.Servicer
}

func NewController(
	servicer attendance.Servicer,
) *Controller {
	return &Controller{
		svc: servicer,
	}
}

// @Summary Get Attendance  History
// @Description Get Attendance  History
// @Tags Attendance 
// @Produce application/json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} http.GetAttendance
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/attendance [get]
func (ctrl *Controller) Get(ctx *gin.Context) {
	accountID, err := jwt.ExtractID(ctx.GetHeader("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	response, err := ctrl.svc.Find(accountID)
	if err != nil {
		if errors.Is(err, constant.ErrAccountNotRegistered) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"account_id": constant.ErrAccountNotRegistered.Error()})
			return
		}
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("get attendance by account id:", err)
		return
	}

	rest.ResponseData(ctx, http.StatusOK, response)
}

// @Summary Get Attendance Status By Locations
// @Description Get Attendance Status By Locations
// @Tags Attendance 
// @Produce application/json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} http.GetAttendanceByLocation
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/attendance/locations [get]
func (ctrl *Controller) GetByLocation(ctx *gin.Context) {
	accountID, err := jwt.ExtractID(ctx.GetHeader("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	response, err := ctrl.svc.FindByLocation(accountID)
	if err != nil {
		if errors.Is(err, constant.ErrAccountNotRegistered) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"account_id": constant.ErrAccountNotRegistered.Error()})
			return
		}
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("get attendance by account id and by location:", err)
		return
	}

	rest.ResponseData(ctx, http.StatusOK, response)
}

// Create godoc
// @Summary Create Attendance
// @Description Create Attendance
// @Tags Attendance 
// @Param Payload body http.AddAttendance true "Payload"
// @Success 201 {object} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 409 {string} string "Resource Conflict"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/attendance [post]
func (ctrl *Controller) Add(ctx *gin.Context) {
	accountID, err := jwt.ExtractID(ctx.GetHeader("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	req := entity.AddAttendance{}
	if err := rest.BindJSON(ctx, &req); err != nil {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"body": constant.ErrInvalidFormat.Error()})
		return
	}

	// required tapi tidak diisi akan return bad request
	if err := validation.Validator.Struct(req); err != nil {
		log.Println("validate struct:", err, "request:", req)
		rest.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	err = ctrl.svc.Add(accountID, req)
	if errors.Is(err, constant.ErrAccountNotRegistered) {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"account_id": constant.ErrAccountNotRegistered.Error()})
		return
	} else if err != nil {
		log.Println("attendance name:", err.Error())
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
	} else {
		rest.ResponseMessage(ctx, http.StatusCreated)
	}
}