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
	Message string `json:"message"`
}

type usersHandler struct {
	ucase domain.UsersUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.UsersUsecase) {
	handler := &usersHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/users", isLoggedIn)
	group.GET("", handler.Fetch)
	group.POST("", handler.Save)
}

func (hn *usersHandler) Fetch(e echo.Context) error {
	// get query param
	idString := e.QueryParam("id")
	limitString := e.QueryParam("limit")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		id = 0
	}

	limit, err := strconv.ParseInt(limitString, 10, 64)
	if err != nil {
		limit = 10
	}
	res, err := hn.ucase.Fetch(e.Request().Context(), id, limit)
	if err != nil {
		return e.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

func (hn *usersHandler) Save(e echo.Context) error {
	var users domain.Users
	err := e.Bind(&users)
	if err != nil {
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: err.Error()})
	}

	res, err := hn.ucase.Save(e.Request().Context(), &users)
	if err != nil {
		if err == domain.ErrNotFound {
			return e.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
		} else if err == domain.ErrConflict {
			return e.JSON(http.StatusConflict, ResponseError{Message: err.Error()})
		} else {
			return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		}
	}

	return e.JSON(http.StatusCreated, res)
}
