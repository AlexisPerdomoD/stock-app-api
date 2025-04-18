package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecase"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/dto"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/middleware"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	getStocksUC     *usecase.GetStocksUseCase
	registerUC      *usecase.RegisterUserUseCase
	loginUC         *usecase.LoginUseCase
	registerStockUC *usecase.RegisterUserStockUseCase
	removeStockUC   *usecase.RemoveUserStockUseCase
}

func NewUserController(
	getStocksUC *usecase.GetStocksUseCase,
	registerUC *usecase.RegisterUserUseCase,
	loginUC *usecase.LoginUseCase,
	registerStockUC *usecase.RegisterUserStockUseCase,
	removeStockUC *usecase.RemoveUserStockUseCase,
) *UserController {
	if getStocksUC == nil {
		log.Fatalln("[UserController]: getStocksUC provided as nil")
	}

	if registerUC == nil {
		log.Fatalln("[UserController]: registerUC provided as nil")
	}

	if loginUC == nil {
		log.Fatalln("[UserController]: loginUC provided as nil")
	}

	if registerStockUC == nil {
		log.Fatalln("[UserController]: registerStockUC provided as nil")
	}

	if removeStockUC == nil {
		log.Fatalln("[UserController]: removeStockUC provided as nil")
	}

	return &UserController{getStocksUC, registerUC, loginUC, registerStockUC, removeStockUC}
}

func (sc *UserController) GetStocksHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID <= 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := c.Request.Context()
	filters := dto.MapGetStocksFilter(c)

	stocks, err := sc.getStocksUC.Execute(ctx, *filters, &userID)
	if err != nil {
		res := pkg.MapHttpErr(err)
		c.JSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, stocks)
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

func (uc *UserController) RegisterStockHandler(c *gin.Context) {
	stockID, ok := c.Params.Get("stockID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"name":    "Bad Request",
			"message": "stockID not provided",
		})
		return
	}

	parseStockID, err := strconv.Atoi(stockID)
	if err != nil || parseStockID <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"name":    "Bad Request",
			"message": "stockID invalid",
		})
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"name":    "Unauthorized",
			"message": "userID not provided",
		})
		return
	}

	ctx := c.Request.Context()
	if err := uc.registerStockUC.Execute(ctx, userID, uint(parseStockID)); err != nil {
		res := pkg.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": true, "message": "user stock registered"})
}

func (uc *UserController) RemoveStockHandler(c *gin.Context) {
	stockID, ok := c.Params.Get("stockID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"name":    "Bad Request",
			"message": "stockID not provided",
		})
		return
	}

	parseStockID, err := strconv.Atoi(stockID)
	if err != nil || parseStockID <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"name":    "Bad Request",
			"message": "stockID invalid",
		})
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"name":    "Unauthorized",
			"message": "userID not provided",
		})
		return
	}

	ctx := c.Request.Context()
	if err := uc.removeStockUC.Execute(ctx, userID, uint(parseStockID)); err != nil {
		res := pkg.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "message": "user stock removed"})
}

func (uc *UserController) SetRoutes(r *gin.Engine) {
	group := r.Group("/users")

	group.POST("", uc.RegisterUserHandler)
	group.POST("/login", uc.LoginUserHandler)
	group.GET("/stocks", middleware.UserSessionMiddleware, uc.GetStocksHandler)
	group.POST("/stocks/:stockID", middleware.UserSessionMiddleware, uc.RegisterStockHandler)
	group.DELETE("/stocks/:stockID", middleware.UserSessionMiddleware, uc.RemoveStockHandler)
}
