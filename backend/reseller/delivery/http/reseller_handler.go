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

// ResellerHandler ...
type ResellerHandler struct {
	ucase domain.ResellerUsecase
	ucasem domain.MessageUsecase
}

// NewResellerUsecase ...
func NewResellerUsecase(e *echo.Echo, uc domain.ResellerUsecase, ucm domain.MessageUsecase) {
	handler := &ResellerHandler{ucase: uc, ucasem: ucm}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/reseller", isLoggedIn)
	group.GET("", handler.Fetch)
	group.POST("/activated", handler.Activated)
	group.DELETE("/:id", handler.Delete)
}

// Fetch ...
func (hn *ResellerHandler) Fetch(e echo.Context) error {

	res, err := hn.ucase.Fetch(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

// Delete ...
func (hn *ResellerHandler) Delete(e echo.Context) error {
	id := e.Param("id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	res, err := hn.ucase.Get(e.Request().Context(), idInt64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	if res == (domain.Reseller{}) {
		return e.JSON(http.StatusNotFound, ResponseError{Message: "not found"})
	}

	err = hn.ucase.Delete(e.Request().Context(), idInt64)
	if err != nil {
		return e.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusAccepted, res)
}

// Activated ...
func (hn *ResellerHandler) Activated(e echo.Context) error {

	id := e.FormValue("id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	res, err := hn.ucase.Get(e.Request().Context(), idInt64)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	if res == (domain.Reseller{}) {
		return e.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	def := "yes"
	res.Active = &def

	err = hn.ucase.Update(e.Request().Context(), idInt64, &res)

	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	var defMessageID int64 = 0
	received := "no"
	mes := "akun anda telah di aktivasi, ketik /menu untuk mulai menggunakan..."
	message := domain.Message{}
	message.ChatID = res.ChatID
	message.MessageID = &defMessageID
	message.Received = &received
	message.Message = &mes

	hn.ucasem.Insert(e.Request().Context(), message)

	return e.JSON(http.StatusAccepted, res)
}
