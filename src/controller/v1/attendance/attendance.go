package attendance

import (
	"log"
	"net/http"
	"strconv"

	"go-rest-api/src/constant"
	entity "go-rest-api/src/http"
	"go-rest-api/src/pkg/jwt"
	"go-rest-api/src/pkg/pagination"
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
// @Param Page query string true "page"
// @Param Limit query string true "limit"
// @Param Filter query string true "string enums" Enums(day, week, month, year)
// @Success 200 {object} http.GetAttendance
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/attendance/history [get]
func (ctrl *Controller) Get(ctx *gin.Context) {
	accountID, err := jwt.ExtractID(ctx.GetHeader("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	limitString := ctx.Query("Limit")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		log.Print(err)
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"limit": constant.ErrInvalidFormat.Error()})
		return
	}
	limit = int(limit)
	pageString := ctx.Query("Page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		log.Print(err)
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"page": constant.ErrInvalidFormat.Error()})
		return
	}
	page = int(page)
	if limit < 0 || limit == 0 {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"limit": constant.ErrInvalidFormat.Error()})
		return
	}
	if page < 0 || page == 0 {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"page": constant.ErrInvalidFormat.Error()})
		return
	}
	pgn := pagination.Pagination{
		Limit:  limit,
		Page:   page,
	}
	pgn.Paginate()

	filter := ctx.Query("Filter")
	log.Print(filter)
	if filter != constant.FilterByDay && filter != constant.FilterByWeek && filter != constant.FilterByMonth && filter != constant.FilterByYear {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"filter": constant.ErrInvalidFormat.Error()})
		return
	}

	response, err := ctrl.svc.FindAttendanceHistory(accountID, pgn, filter)
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
// @Param Page query string true "page"
// @Param Limit query string true "limit"
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

	limitString := ctx.Query("Limit")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		log.Print(err)
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"limit": constant.ErrInvalidFormat.Error()})
		return
	}
	limit = int(limit)
	pageString := ctx.Query("Page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		log.Print(err)
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"page": constant.ErrInvalidFormat.Error()})
		return
	}
	page = int(page)
	if limit < 0 || limit == 0 {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"limit": constant.ErrInvalidFormat.Error()})
		return
	}
	if page < 0 || page == 0 {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"page": constant.ErrInvalidFormat.Error()})
		return
	}
	pgn := pagination.Pagination{
		Limit:  limit,
		Page:   page,
	}
	pgn.Paginate()

	response, err := ctrl.svc.FindByLocation(accountID, pgn)
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
// @Summary Add Attendance
// @Description Add Attendance
// @Tags Attendance
// @Param Authorization header string true "Bearer Token"
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
	} else if errors.Is(err, constant.ErrLocationNotExist) {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"location": constant.ErrLocationNotExist.Error()})
	} else if errors.Is(err, constant.ErrInvalidStatusAttendance) {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"status": constant.ErrInvalidStatusAttendance.Error()})
	} else if err != nil {
		log.Println("attendance name:", err.Error())
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
	} else {
		rest.ResponseMessage(ctx, http.StatusCreated)
	}
}