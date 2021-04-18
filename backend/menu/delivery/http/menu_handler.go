package http

import (
	"insinyur-radius/domain"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

// ResponseError ...
type ResponseError struct {
	Message string `json:"error"`
}

// MenuHandler ...
type MenuHandler struct {
	ucase domain.MenuUsecase
}

// NewMenuHandler ...
func NewMenuHandler(e *echo.Echo, uc domain.MenuUsecase) {
	handler := &MenuHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/menu", isLoggedIn)
	group.GET("", handler.Fetch)
	group.POST("", handler.Store)
	group.PUT("/:id", handler.Update)
	group.DELETE("/:id", handler.Delete)
}

// Fetch ...
func (hn *MenuHandler) Fetch(e echo.Context) error {

	res, err := hn.ucase.Fetch(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

// Update ..
func (hn *MenuHandler) Update(e echo.Context) error {
	var menu domain.Menu
	err := e.Bind(&menu)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	id := e.Param("id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	err = hn.ucase.Update(e.Request().Context(), idInt64, &menu)

	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, menu)
}

// Store ...
func (hn *MenuHandler) Store(e echo.Context) error {

	var menu domain.Menu
	err := e.Bind(&menu)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	err = hn.ucase.Insert(e.Request().Context(), &menu)

	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, menu)

}

// Delete ...
func (hn *MenuHandler) Delete(e echo.Context) error {

	id := e.Param("id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	res, err := hn.ucase.Get(e.Request().Context(), idInt64)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	if res == (domain.Menu{}) {
		return e.JSON(http.StatusNotFound, ResponseError{Message: "not found item"})
	}

	err = hn.ucase.Delete(e.Request().Context(), idInt64)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusAccepted, res)
}
