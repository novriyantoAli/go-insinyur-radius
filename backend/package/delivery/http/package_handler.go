package http

import (
	"insinyur-radius/domain"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ResponseError ...
type ResponseError struct {
	Message string `json:"message"`
}

// ResponseMessage ...
type ResponseMessage struct {
	Message string `json:"message"`
}

// PackageHandler ...
type PackageHandler struct {
	ucase domain.PackageUsecase
}

// NewPackageHandler ...
func NewPackageHandler(e *echo.Echo, uc domain.PackageUsecase) {
	handler := &PackageHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/package", isLoggedIn)
	group.GET("", handler.Fetch)
	group.POST("", handler.Store)
	group.PUT("/:id", handler.Update)
	group.DELETE("/:id", handler.Delete)
}

// Fetch ...
func (hn *PackageHandler) Fetch(e echo.Context) error {

	res, err := hn.ucase.Fetch(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

// Store ...
func (hn *PackageHandler) Store(e echo.Context) error {
	var packages domain.Package
	err := e.Bind(&packages)
	if err != nil {
		logrus.Error(err)
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	err = hn.ucase.Insert(e.Request().Context(), &packages)
	if err != nil {
		logrus.Error(err)
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusCreated, packages)
}

// Update ...
func (hn *PackageHandler) Update(e echo.Context) error {
	var packages domain.Package
	err := e.Bind(&packages)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, &packages)
	}

	id := e.Param("id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return e.JSON(http.StatusNotFound, ResponseError{err.Error()})
	}

	err = hn.ucase.Update(e.Request().Context(), idInt64, &packages)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{err.Error()})
	}

	return e.JSON(http.StatusAccepted, packages)
}

// Delete ...
func (hn *PackageHandler) Delete(e echo.Context) error {
	id := e.Param("id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return e.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	err = hn.ucase.Delete(e.Request().Context(), idInt64)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusAccepted, ResponseMessage{Message: "package deleted"})
}
