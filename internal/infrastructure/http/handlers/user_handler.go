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

// LoginUserHandler godoc
// @Summary Login de usuario
// @Description Autentica un usuario y retorna un JWT de sesión
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserLoginDTO true "Credenciales de usuario"
// @Success 200 {object} map[string]interface{}
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

	c.JSON(http.StatusOK, gin.H{
		"Ok":      true,
		"message": "user logged in properly",
		"session": session,
		"user":    user,
	})
}

// RegisterUserHandler godoc
// @Summary Registro de usuario
// @Description Registra un usuario nuevo y retorna un JWT
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.RegisterUserDTO true "Datos de registro"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} mappers.HttpErrResponse
// @Failure 409 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/users [post]
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

// GetStocksHandler godoc
// @Summary Obtener stocks del usuario
// @Description Retorna la lista de stocks asociados al usuario autenticado
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param authorization header string true "Esquema JWT. Usar: \"Bearer {token}\""
// @in header
// @name authorization
// @Success 200 {array} models.StockView
// @Failure 401 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/users/stocks [get]
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

// RegisterStockHandler godoc
// @Summary Registrar stock al usuario
// @Description Asocia un stock existente al usuario autenticado
// @Tags users
// @Security BearerAuth
// @Param stockID path int true "ID del stock"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} mappers.HttpErrResponse
// @Failure 401 {object} mappers.HttpErrResponse
// @Failure 404 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/users/stocks/{stockID} [post]
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

	c.JSON(http.StatusCreated, gin.H{
		"Ok":      true,
		"message": "user registered properly",
		"session": session,
		"user":    usr,
	})
}

// RemoveStockHandler godoc
// @Summary Eliminar stock del usuario
// @Description Remueve un stock asociado al usuario autenticado
// @Tags users
// @Security BearerAuth
// @Param stockID path int true "ID del stock"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} mappers.HttpErrResponse
// @Failure 401 {object} mappers.HttpErrResponse
// @Failure 404 {object} mappers.HttpErrResponse
// @Failure 500 {object} mappers.HttpErrResponse
// @Router /api/v1/users/stocks/{stockID} [delete]
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
