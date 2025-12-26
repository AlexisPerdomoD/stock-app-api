package handlers

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/mappers"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/middleware"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/models"
	"github.com/alexisPerdomoD/stock-app-api/pkg/auth"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type UserController struct {
	getStocksUC     *usecases.GetStocks
	registerUC      *usecases.RegisterUser
	loginUC         *usecases.Login
	registerStockUC *usecases.RegisterUserStock
	removeStockUC   *usecases.RemoveUserStock
}

func (sc *UserController) GetStocksHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID <= 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := c.Request.Context()
	filters := mappers.MapGetStocksFilter(c)

	stocks, err := sc.getStocksUC.Execute(ctx, filters, &userID)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.JSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, stocks)
}

func (uc *UserController) RegisterUserHandler(c *gin.Context) {
	args, err := mappers.MapRegisterUserDTO(c)
	if err != nil {
		issues := mappers.MapValidationErrors(err)
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"name": "bad_request", "message": "invalid format", "issues": issues},
		)
		return
	}
	ctx := c.Request.Context()
	usr := args.ToDomain()
	if err := uc.registerUC.Execute(ctx, usr); err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	session, err := auth.GenerateSessionToken(usr)
	if err != nil {
		res := mappers.MapHttpErr(err)
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
	credentials := &models.UserLoginDTO{}
	if err := c.ShouldBindBodyWithJSON(credentials); err != nil {
		issues := mappers.MapValidationErrors(err)
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"name": "bad_request", "message": "invalid credentials", "issues": issues},
		)
		return
	}
	ctx := c.Request.Context()
	user, err := uc.loginUC.Execute(ctx, credentials.Email, []byte(credentials.Password))
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	session, err := auth.GenerateSessionToken(user)
	if err != nil {
		res := mappers.MapHttpErr(err)
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
	println(userID)
	if userID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"name":    "Unauthorized",
			"message": "userID not provided",
		})
		return
	}

	ctx := c.Request.Context()
	if err := uc.registerStockUC.Execute(ctx, userID, uint(parseStockID)); err != nil {
		res := mappers.MapHttpErr(err)
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
		res := mappers.MapHttpErr(err)
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

func NewUserController(
	getStocksUC *usecases.GetStocks,
	registerUC *usecases.RegisterUser,
	loginUC *usecases.Login,
	registerStockUC *usecases.RegisterUserStock,
	removeStockUC *usecases.RemoveUserStock,
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
