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

// ReportHandler ...
type reportHandler struct {
	ucase domain.ReportUsecase
}

// NewHandler ...
func NewHandler(e *echo.Echo, uc domain.ReportUsecase) {
	handler := &reportHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/report", isLoggedIn)
	group.GET("/finance/day", handler.FinanceCurrentDay)
	group.GET("/finance/month", handler.FinanceCurrentMonth)
	group.GET("/finance/year", handler.FinanceCurrentYear)
	group.GET("/expiration/users/day", handler.ReportExpirationCurrent)
}

// FinanceCurrentDay ...
func (hn *reportHandler) FinanceCurrentDay(e echo.Context) error {

	res, err := hn.ucase.ReportFinanceCurrent(e.Request().Context(), 1)
	if err != nil {
		return e.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

// FinanceCurrentMonth ...
func (hn *reportHandler) FinanceCurrentMonth(e echo.Context) error {

	res, err := hn.ucase.ReportFinanceCurrent(e.Request().Context(), 2)
	if err != nil {
		return e.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

// FinanceCurrentYear ...
func (hn *reportHandler) FinanceCurrentYear(e echo.Context) error {

	res, err := hn.ucase.ReportFinanceCurrent(e.Request().Context(), 3)
	if err != nil {
		return e.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

func (hn *reportHandler) ReportExpirationCurrent(e echo.Context) error {

	res, err := hn.ucase.ReportExpirationCurrent(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}
