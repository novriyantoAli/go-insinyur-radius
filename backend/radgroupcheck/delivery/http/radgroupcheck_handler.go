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
	Message string `json:"error"`
}

// RadgroupcheckHandler ...
type RadgroupcheckHandler struct {
	ucase domain.RadgroupcheckUsecase
}

// NewRadgroupcheckHandler ...
func NewRadgroupcheckHandler(e *echo.Echo, uc domain.RadgroupcheckUsecase) {
	handler := &RadgroupcheckHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/rad/groupcheck", isLoggedIn)
	group.POST("/find", handler.Find)
}

// Find ...
func (hn *RadgroupcheckHandler) Find(e echo.Context) error {

	var bp domain.Radgroupcheck
	err := e.Bind(&bp)
	if err != nil {
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: err.Error()})
	}

	res, err := hn.ucase.Find(e.Request().Context(), bp)
	if err != nil {
		return e.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}
