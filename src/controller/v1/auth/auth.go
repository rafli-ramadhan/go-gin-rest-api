package auth

import (
	"fmt"
	"log"
	"net/http"

	"go-rest-api/src/constant"
	entity "go-rest-api/src/http"
	"go-rest-api/src/pkg/bcrypt"
	"go-rest-api/src/service/v1/account"
	"go-rest-api/src/pkg/jwt"
	"github.com/forkyid/go-utils/v1/aes"
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

// @Summary User Login
// @Description User Login
// @Tags Auth
// @Produce application/json
// @Param Payload body http.Auth true "Payload"
// @Success 200 {object} http.Token
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/auth [post]
func (ctrl *Controller) Login(ctx *gin.Context) {
	request := entity.Auth{}
	err := rest.BindJSON(ctx, &request)
	if err != nil {
		log.Println("bind json:", err, "request:", request)
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"body": constant.ErrInvalidFormat.Error()})
		return
	}

	if err := validation.Validator.Struct(request); err != nil {
		log.Println("validate struct:", err, "request:", request)
		rest.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	exist, _ := ctrl.svc.CheckAccountByUsername(request.Username)
	if !exist {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"accounts": constant.ErrAccountNotRegistered.Error()})
		return
	}

	account, err:= ctrl.svc.TakeAccountByUsername(request.Username)
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		return
	}

	err = bcrypt.ComparePassword(account.Password, request.Password)
	if err != nil {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"accounts": constant.ErrInvalidPassword.Error()})
		return
	}

	token, err := jwt.GenerateJWT(aes.Encrypt(int(account.ID)))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		return
	}

	rest.ResponseData(ctx, http.StatusOK, entity.Token{
		Token: fmt.Sprintf("Bearer %v", token),
	})
}

// @Summary Update User Password
// @Description Update User Password
// @Tags Auth
// @Produce application/json
// @Param Payload body http.ForgotPassword true "Payload"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/auth/forgot [patch]
func (ctrl *Controller) ForgotPassword(ctx *gin.Context) {
	request := entity.ForgotPassword{}
	err := rest.BindJSON(ctx, &request)
	if err != nil {
		log.Println("bind json:", err, "request:", request)
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"body": constant.ErrInvalidFormat.Error()})
		return
	}

	if err := validation.Validator.Struct(request); err != nil {
		log.Println("validate struct:", err, "request:", request)
		rest.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	err = ctrl.svc.UpdatePassword(request)
	if err != nil {
		if errors.Is(err, constant.ErrPasswordCannotBeEmpty) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrPasswordCannotBeEmpty.Error()})
			return
		}
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("update account password:", err)
		return
	}

	rest.ResponseMessage(ctx, http.StatusOK)
}