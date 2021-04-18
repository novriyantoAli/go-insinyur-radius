package http

import (
	"insinyur-radius/domain"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

// ResponseError ...
type ResponseError struct {
	Message string `json:"error"`
}

type paymentHandler struct {
	ucase domain.PaymentUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.PaymentUsecase) {
	handler := &paymentHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/payment", isLoggedIn)
	group.GET("", handler.Fetch)
	group.POST("", handler.Save)
	group.POST("/find", handler.Find)
	group.PUT("", handler.Update)
	group.DELETE("/:id", handler.Delete)
}

// Fetch ...
func (hn *paymentHandler) Fetch(e echo.Context) error {
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
func (hn *paymentHandler) Update(e echo.Context) error {
	var payment domain.Payment
	err := e.Bind(&payment)
	if err != nil {
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: err.Error()})
	}

	err = hn.ucase.Update(e.Request().Context(), &payment)

	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, 1)
}

func (hn *paymentHandler) Find(e echo.Context) error {
	var payment domain.Payment
	err := e.Bind(&payment)
	if err != nil {
		logrus.Error(err)
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: err.Error()})
	}

	res, err := hn.ucase.Find(e.Request().Context(), &payment)
	if err != nil {
		logrus.Error(err)
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusCreated, res)
}

func (hn *paymentHandler) Save(e echo.Context) error {
	var payment domain.Payment
	err := e.Bind(&payment)
	if err != nil {
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: err.Error()})
	}

	err = hn.ucase.Save(e.Request().Context(), &payment)
	if err != nil {
		if err == domain.ErrNotFound {
			return e.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
		} else if err == domain.ErrConflict {
			return e.JSON(http.StatusConflict, ResponseError{Message: err.Error()})
		} else {
			return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		}
	}

	return e.JSON(http.StatusCreated, payment)
}

// Delete ...
func (hn *paymentHandler) Delete(e echo.Context) error {

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
