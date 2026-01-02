package handlers

import (
	"log"
	"net/http"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/models"
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/mappers"
	"github.com/alexisPerdomoD/stock-app-api/pkg/auth"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	register *usecases.RegisterUser
	login    *usecases.Login
}

// LoginUserHandler godoc
// @Summary Login de usuario
// @Description Autentica un usuario y retorna un JWT de sesión
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserLoginDTO true "Credenciales de usuario"
// @Success 200 {object} models.UserLoginView
// @Failure 400 {object} mappers.HttpErrResponse
// @Failure 401 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/login [post]
func (uc *UserHandler) LoginUserHandler(c *gin.Context) {
	userLogin, err := mappers.MapUserLoginDTO(c)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}
	ctx := c.Request.Context()
	user, err := uc.login.Execute(ctx, userLogin)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	session, err := auth.GenerateSessionToken(user.ID)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, models.UserLoginView{
		Message: "login success",
		Session: session,
		User:    user,
	})
}

// RegisterUserHandler godoc
// @Summary Registro de usuario
// @Description Registra un usuario nuevo y retorna un JWT
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.RegisterUserDTO true "Datos de registro"
// @Success 201 {object} models.UserLoginView
// @Failure 400 {object} mappers.HttpErrResponse
// @Failure 409 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/users [post]
func (uc *UserHandler) RegisterUserHandler(c *gin.Context) {
	args, err := mappers.MapRegisterUserDTO(c)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	usr, err := uc.register.Execute(c.Request.Context(), args)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	session, err := auth.GenerateSessionToken(usr.ID)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusCreated, models.UserLoginView{
		Message: "user registered properly",
		Session: session,
		User:    usr,
	})
}

func NewUserHandler(
	registerUC *usecases.RegisterUser,
	loginUC *usecases.Login,
) *UserHandler {

	if registerUC == nil {
		log.Fatalln("[UserController]: registerUC provided as nil")
	}

	if loginUC == nil {
		log.Fatalln("[UserController]: loginUC provided as nil")
	}

	return &UserHandler{registerUC, loginUC}
}
