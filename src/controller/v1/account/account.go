package account

import (
	"log"
	"strings"
	"net/http"

	"go-rest-api/src/constant"
	"go-rest-api/src/pkg/jwt"
	entity "go-rest-api/src/http"
	"go-rest-api/src/service/v1/account"
	// "github.com/forkyid/go-utils/v1/aes"
	"github.com/forkyid/go-utils/v1/rest"
	"github.com/forkyid/go-utils/v1/validation"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Controller struct {
	svc account.Servicer
}

func NewController(
	servicer account.Servicer,
) *Controller {
	return &Controller{
		svc: servicer,
	}
}

// Register godoc
// @Summary Register Account
// @Description Register Account
// @Tags Accounts
// @Param Payload body http.RegisterUser true "Payload"
// @Success 201 {object} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 409 {string} string "Resource Conflict"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/accounts/register [post]
func (ctrl *Controller) Register(ctx *gin.Context) {
	req := entity.RegisterUser{}
	if err := rest.BindJSON(ctx, &req); err != nil {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"body": constant.ErrInvalidFormat.Error()})
		return
	}

	req.Username = strings.ToLower(req.Username)
	err := ctrl.svc.Create(req)
	if errors.Is(err, constant.ErrAccountExist) {
		rest.ResponseMessage(ctx, http.StatusConflict, errors.Cause(err).Error())
	} else if err != nil {
		log.Println("register:", err.Error())
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
	} else {
		rest.ResponseMessage(ctx, http.StatusCreated)
	}
}

// Update godoc
// @Summary Update Account
// @Description Update Account
// @Tags Accounts
// @Param Payload body http.UpdateUser true "Payload"
// @Param account_id query string false "account_id"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/accounts [patch]
func (ctrl *Controller) Update(ctx *gin.Context) {
	request := entity.UpdateUser{}
	// int di isi dengan string maka akan return invalid format
	err := rest.BindJSON(ctx, &request)
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

	accountID, err := jwt.ExtractID(ctx.GetHeader("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	err = ctrl.svc.Update(accountID, request)
	if err != nil {
		if errors.Is(err, constant.ErrAccountNotRegistered) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrAccountNotRegistered.Error()})
			return
		} else if errors.Is(err, constant.ErrUsernameAlreadyExist) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrUsernameAlreadyExist.Error()})
			return
		} else if errors.Is(err, constant.ErrEmailAlreadyExist) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrEmailAlreadyExist.Error()})
			return
		} else if errors.Is(err, constant.ErrKTPNumberAlreadyExist) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrKTPNumberAlreadyExist.Error()})
			return
		} else if errors.Is(err, constant.ErrPhoneNumberAlreadyExist) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrPhoneNumberAlreadyExist.Error()})
			return
		}
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("update account: ", err.Error())
		return
	}

	rest.ResponseMessage(ctx, http.StatusOK)
}

// Delete godoc
// @Summary Delete Account
// @Description Delete Account By User Itself
// @Tags Accounts
// @Param Authorization header string true "Bearer Token"
// @Param account_id query string false "account_id"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Bad Request"
// @Failure 409 {string} string "Resource Conflict"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/accounts [delete]
func (ctrl *Controller) Delete(ctx *gin.Context) {
	accountID, err := jwt.ExtractID(ctx.GetHeader("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	err = ctrl.svc.Delete(accountID)
	if err != nil {
		if errors.Is(err, constant.ErrAccountNotRegistered) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrAccountNotRegistered.Error()})
			return
		}
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("delete account: ", err.Error())
		return
	}
		
	rest.ResponseMessage(ctx, http.StatusOK)
}
