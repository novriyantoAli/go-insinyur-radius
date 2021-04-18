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

// ResponseMessage ...
type ResponseMessage struct {
	Message string `json:"message"`
}

type ResponseBalance struct {
	Message int64 `json:"balance"`
}

// ResponseMessage ...

// TransactionHandler ...
type TransactionHandler struct {
	ucase domain.TransactionUsecase
}

// NewTransactionHandler ...
func NewTransactionHandler(e *echo.Echo, uc domain.TransactionUsecase) {
	handler := &TransactionHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/transaction", isLoggedIn)
	group.GET("", handler.Fetch)
	group.GET("/balance/:id", handler.Balance)
	group.POST("/refill", handler.Refill)
	group.POST("/report", handler.Report)
}

// Fetch ...
func (hn *TransactionHandler) Fetch(e echo.Context) error {

	res, err := hn.ucase.Fetch(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

// Report ...
func (hn *TransactionHandler) Report(e echo.Context) error {

	dateStart := e.FormValue("date_start")
	dateEnd := e.FormValue("date_end")

	res, err := hn.ucase.Report(e.Request().Context(), dateStart, dateEnd)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

// Balance ...
func (hn *TransactionHandler) Balance(e echo.Context) error {
	id := e.Param("id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	balance, err := hn.ucase.Balance(e.Request().Context(), idInt64)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, ResponseBalance{Message: balance})

}

// Refill ...
func (hn *TransactionHandler) Refill(e echo.Context) error {

	id := e.FormValue("id_reseller")
	balance := e.FormValue("balance")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	balanceInt64, err := strconv.ParseInt(balance, 10, 64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	_, err = hn.ucase.Refill(e.Request().Context(), idInt64, balanceInt64)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, ResponseMessage{Message: "success"})
}
