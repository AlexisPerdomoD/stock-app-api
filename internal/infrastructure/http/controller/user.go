package controller

import (
	"log"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	registerUC *usecase.RegisterUserUseCase
	loginUC    *usecase.LoginUseCase
}

func (uc *UserController) RegisterUserHandler(c *gin.Context) {

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
