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

	err = bcrypt.CheckPasswordHarsh(account.Password, request.Password)
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

// @Summary Get User Data
// @Description Get User Data
// @Tags Auth
// @Produce application/json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} http.GetUser
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/auth/self [get]
func (ctrl *Controller) AuthSelf(ctx *gin.Context) {
	accountID, err := jwt.ExtractID(ctx.GetHeader("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	response, err := ctrl.svc.TakeAccountByID(accountID)
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("get account by id:", err)
		return
	}

	rest.ResponseData(ctx, http.StatusOK, response)
}