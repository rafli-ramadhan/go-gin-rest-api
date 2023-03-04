package location

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"strconv"

	"go-rest-api/src/constant"
	entity "go-rest-api/src/http"
	"go-rest-api/src/pkg/jwt"
	"go-rest-api/src/service/v1/location"

	"github.com/forkyid/go-utils/v1/rest"
	"github.com/forkyid/go-utils/v1/validation"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Controller struct {
	svc location.Servicer
}

func NewController(
	servicer location.Servicer,
) *Controller {
	return &Controller{
		svc: servicer,
	}
}

// @Summary Get Locations Data
// @Description Get Locations Data
// @Tags Locations
// @Produce application/json
// @Param Authorization header string true "Bearer Token"
// @Param location_ids query string false "location_ids separated by comma, example: 1,2,3,4,5"
// @Success 200 {object} http.GetLocation
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/locations [get]
func (ctrl *Controller) Get(ctx *gin.Context) {
	_, err := jwt.ExtractID(ctx.GetHeader("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	locationIDsStr := ctx.Query("location_ids")
	if locationIDsStr == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"location_id": constant.ErrInvalidID.Error()})
		return
	}

	locationIDsStrArr := strings.Split(locationIDsStr, ",")
	var locationIDs []int
	for i := range locationIDsStrArr {
	    locationID, err := strconv.Atoi(locationIDsStrArr[i])
	    if err != nil {
		    rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			    "location_id": fmt.Sprintf(constant.ErrInvalidID.Error() + " : " + locationIDsStrArr[i])})
		    return
	    }
		locationIDs = append(locationIDs, locationID)
	}

	response, err := ctrl.svc.Find(locationIDs)
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("get location by id:", err)
		return
	}

	rest.ResponseData(ctx, http.StatusOK, response)
}

// Create godoc
// @Summary Create Location
// @Description Create Location
// @Tags Locations
// @Param Authorization header string true "Bearer Token"
// @Param Payload body http.CreateLocation true "Payload"
// @Success 201 {object} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 409 {string} string "Resource Conflict"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/locations [post]
func (ctrl *Controller) Create(ctx *gin.Context) {
	_, err := jwt.ExtractID(ctx.GetHeader("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	req := entity.CreateLocation{}
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

	err = ctrl.svc.Create(req)
	if errors.Is(err, constant.ErrInvalidLocationName) {
		rest.ResponseError(ctx, http.StatusConflict, map[string]string{
			"location": constant.ErrInvalidLocationName.Error()})
	} else if errors.Is(err, constant.ErrInvalidAddress) {
		rest.ResponseError(ctx, http.StatusConflict, map[string]string{
			"location": constant.ErrInvalidAddress.Error()})
	} else if errors.Is(err, constant.ErrLocationAlreadyExist) {
		rest.ResponseError(ctx, http.StatusConflict, map[string]string{
			"location": constant.ErrLocationAlreadyExist.Error()})
	} else if err != nil {
		log.Println("location:", err.Error())
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
	} else {
		rest.ResponseMessage(ctx, http.StatusCreated)
	}
}

// Update godoc
// @Summary Update Location
// @Description Update Location
// @Tags Locations
// @Param Authorization header string true "Bearer Token"
// @Param location_id query string false "location_id"
// @Param Payload body http.UpdateLocation true "Payload"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/locations [patch]
func (ctrl *Controller) Update(ctx *gin.Context) {
	_, err := jwt.ExtractID(ctx.GetHeader("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	request := entity.UpdateLocation{}
	// int di isi dengan string maka akan return invalid format
	err = rest.BindJSON(ctx, &request)
	if err != nil {
		log.Println("bind json:", err, "request:", request)
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"body": constant.ErrInvalidFormat.Error()})
		return
	}

	// required tapi tidak diisi akan return bad request
	if err := validation.Validator.Struct(request); err != nil {
		log.Println("validate struct:", err, "request:", request)
		rest.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	locationIDStr := ctx.Query("location_id")
	if locationIDStr == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"location_id": constant.ErrInvalidID.Error()})
		return
	}

	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"location_id": constant.ErrInvalidID.Error()})
		return
	}

	err = ctrl.svc.Update(locationID, request)
	if err != nil {
		if errors.Is(err, constant.ErrLocationNotExist) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"locations": constant.ErrLocationNotExist.Error()})
			return
		} else if errors.Is(err, constant.ErrLocationNameAlreadyExist) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"locations": constant.ErrLocationNameAlreadyExist.Error()})
			return
		}
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("update location: ", err.Error())
		return
	}

	rest.ResponseMessage(ctx, http.StatusOK)
}

// Delete godoc
// @Summary Delete Location
// @Description Delete Location By Location Itself
// @Tags Locations
// @Param Authorization header string true "Bearer Token"
// @Param location_id query string false "location_id"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Bad Request"
// @Failure 409 {string} string "Resource Conflict"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/locations [delete]
func (ctrl *Controller) Delete(ctx *gin.Context) {
	_, err := jwt.ExtractID(ctx.GetHeader("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}
	
	locationIDStr := ctx.Query("location_id")
	if locationIDStr == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"location_id": constant.ErrInvalidID.Error()})
		return
	}

	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"location_id": constant.ErrInvalidID.Error()})
		return
	}

	err = ctrl.svc.Delete(locationID)
	if err != nil {
		if errors.Is(err, constant.ErrLocationNotExist) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"locations": constant.ErrLocationNotExist.Error()})
			return
		}
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("delete location: ", err.Error())
		return
	}
		
	rest.ResponseMessage(ctx, http.StatusOK)
}
