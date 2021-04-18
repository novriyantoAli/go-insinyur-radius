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
type customerHandler struct {
	ucase domain.CustomerUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.CustomerUsecase) {
	handler := &customerHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/customer", isLoggedIn)
	group.GET("", handler.Fetch)
	group.POST("", handler.Save)
	group.PUT("/:id", handler.Refill)
	group.DELETE("/:id", handler.Delete)
}

// Fetch ...
func (hn *customerHandler) Fetch(e echo.Context) error {
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

// Update ..
func (hn *customerHandler) Refill(e echo.Context) error {
	id := e.Param("id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	err = hn.ucase.Refill(e.Request().Context(), idInt64)

	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, 1)
}


func (hn *customerHandler) Save(e echo.Context) error {
	var customer domain.Customer
	err := e.Bind(&customer)
	if err != nil {
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: err.Error()})
	}

	err = hn.ucase.Insert(e.Request().Context(), customer)
	if err != nil {
		if err == domain.ErrNotFound {
			return e.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
		} else if err == domain.ErrConflict {
			return e.JSON(http.StatusConflict, ResponseError{Message: err.Error()})
		} else {
			return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		}
	}

	return e.JSON(http.StatusCreated, 1)
}

// Delete ...
func (hn *customerHandler) Delete(e echo.Context) error {

	id := e.Param("id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	err = hn.ucase.Delete(e.Request().Context(), idInt64)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusAccepted, 1)
}
