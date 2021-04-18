package http

import (
	"insinyur-radius/domain"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

// ResponseError ...
type ResponseError struct {
	Message string `json:"message"`
}

type nasHandler struct {
	ucase domain.NasUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.NasUsecase) {
	handler := &nasHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/rad/nas", isLoggedIn)
	group.GET("", handler.Get)
	group.POST("", handler.Store)
}

// Get ...
func (hn *nasHandler) Get(e echo.Context) error {

	res, err := hn.ucase.Get(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

// Store ...
func (hn *nasHandler) Store(e echo.Context) error {
	nasname := e.FormValue("nasname")
	secret := e.FormValue("secret")

	if &nasname == nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: "form ip address required"})
	}

	if nasname == "" {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: "form ip address required"})
	}

	if &secret == nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: "form secret required"})
	}

	if secret == "" {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: "form secret required"})
	}

	res, err := hn.ucase.Upsert(e.Request().Context(), nasname, secret)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusCreated, res)
}
