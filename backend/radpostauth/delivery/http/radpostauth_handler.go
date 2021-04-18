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

type radpostauthHandler struct {
	ucase domain.RadpostauthUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.RadpostauthUsecase) {
	handler := &radpostauthHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/rad/postauth", isLoggedIn)
	group.GET("", handler.Get)
}

// Get ...
func (hn *radpostauthHandler) Get(e echo.Context) error {

	username := e.QueryParam("username")

	id := e.QueryParam("id")

	limit := e.QueryParam("limit")

	if &username == nil || username == "" {
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: "username not defined"})
	}

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		idInt64 = 0
	}

	limitInt64, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		limitInt64 = 10
	}

	res, err := hn.ucase.Get(e.Request().Context(), username, idInt64, limitInt64)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}
