package controller

import (
	"log"
	"net/http"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecase"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/dto"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	registerUC *usecase.RegisterUserUseCase
	loginUC    *usecase.LoginUseCase
}

func NewUserController(
	registerUC *usecase.RegisterUserUseCase,
	loginUC *usecase.LoginUseCase,
) *UserController {
	if registerUC == nil {
		log.Fatalln("bad impl: RegisterUserUseCase was nil for NewUserController")

	}
	if loginUC == nil {
		log.Fatalln("bad impl: LoginUseCase was nil for NewUserController")
	}

	return &UserController{registerUC, loginUC}
}

func (uc *UserController) RegisterUserHandler(c *gin.Context) {
	user, err := dto.MapNewUserForm(c)
	if err != nil {
		issues := dto.GetValidationErrors(err)
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"name": "bad_request", "message": "invalid format", "issues": issues},
		)
		return
	}
	ctx := c.Request.Context()
	session, err := uc.registerUC.Execute(ctx, user)
	if err != nil {
		res := pkg.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"Ok":      true,
		"message": "user registered properly",
		"session": session,
	})
}

func (uc *UserController) LoginUserHandler(c *gin.Context) {
	credentials := &dto.UserDto{}
	if err := c.ShouldBindBodyWithJSON(credentials); err != nil {
		issues := dto.GetValidationErrors(err)
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"name": "bad_request", "message": "invalid credentials", "issues": issues},
		)
		return
	}
	ctx := c.Request.Context()
	session, err := uc.loginUC.Execute(ctx, credentials.Email, credentials.Password)
	if err != nil {
		res := pkg.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Ok":      true,
		"message": "user logged in properly",
		"session": session,
	})
}

func (uc *UserController) SetRoutes(r *gin.Engine) {
	group := r.Group("/users")

	group.POST("", uc.RegisterUserHandler)
	group.POST("/login", uc.LoginUserHandler)
}
