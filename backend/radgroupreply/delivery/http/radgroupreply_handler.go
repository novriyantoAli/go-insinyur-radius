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

// RadgroupreplyHandler ...
type RadgroupreplyHandler struct {
	ucase domain.RadgroupreplyUsecase
}

// NewRadgroupreplyHandler ...
func NewRadgroupreplyHandler(e *echo.Echo, uc domain.RadgroupreplyUsecase) {
	handler := &RadgroupreplyHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/rad/groupreply", isLoggedIn)
	group.POST("/find", handler.Find)
}

// Find ...
func (hn *RadgroupreplyHandler) Find(e echo.Context) error {

	var bp domain.Radgroupreply
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
