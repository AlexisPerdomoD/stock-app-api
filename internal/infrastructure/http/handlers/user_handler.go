package handlers

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/application/usecases"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/http/mappers"
	"github.com/alexisPerdomoD/stock-app-api/pkg"
	"github.com/alexisPerdomoD/stock-app-api/pkg/auth"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	getStocks     *usecases.GetStocks
	register      *usecases.RegisterUser
	login         *usecases.Login
	registerStock *usecases.RegisterUserStock
	removeStock   *usecases.RemoveUserStock
}

func (sc *UserHandler) GetStocksHandler(c *gin.Context) {
	userID := c.GetUint64("user_id")
	if userID <= 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := c.Request.Context()
	filters := mappers.MapGetStocksFilter(c)

	stocks, err := sc.getStocks.Execute(ctx, filters, &userID)
	if err != nil {
		res := mappers.MapHttpErr(err)
		c.JSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, stocks)
}

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
		"user":    usr,
	})
}

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
		"user":    user,
	})
}

func (uc *UserHandler) RegisterStockHandler(c *gin.Context) {
	stockID, ok := c.Params.Get("stockID")
	if !ok {
		res := mappers.MapHttpErr(pkg.BadRequest("stockID is required"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	parseStockID, err := strconv.ParseUint(stockID, 10, 0)
	if err != nil {
		res := mappers.MapHttpErr(pkg.BadRequest("stockID is invalid"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	userID := c.GetUint64("user_id")

	ctx := c.Request.Context()
	if err := uc.registerStock.Execute(ctx, userID, parseStockID); err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ok": true, "message": "user stock registered"})
}

func (uc *UserHandler) RemoveStockHandler(c *gin.Context) {
	stockID, ok := c.Params.Get("stockID")
	if !ok {
		res := mappers.MapHttpErr(pkg.BadRequest("stockID is required"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	parseStockID, err := strconv.ParseUint(stockID, 10, 0)
	if err != nil {
		res := mappers.MapHttpErr(pkg.BadRequest("stockID is invalid"))
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	userID := c.GetUint64("user_id")

	ctx := c.Request.Context()
	if err := uc.removeStock.Execute(ctx, userID, parseStockID); err != nil {
		res := mappers.MapHttpErr(err)
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "message": "user stock removed"})
}

func NewUserHandler(
	getStocksUC *usecases.GetStocks,
	registerUC *usecases.RegisterUser,
	loginUC *usecases.Login,
	registerStockUC *usecases.RegisterUserStock,
	removeStockUC *usecases.RemoveUserStock,
) *UserHandler {
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

	return &UserHandler{getStocksUC, registerUC, loginUC, registerStockUC, removeStockUC}
}
