package http

import (
	"insinyur-radius/domain"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

// ResponseError ...
type ResponseError struct {
	Message string `json:"error"`
}

// RadusergroupHandler ...
type RadusergroupHandler struct {
	ucase domain.RadusergroupUsecase
}

// NewRadusergroupHandler ...
func NewRadusergroupHandler(e *echo.Echo, uc domain.RadusergroupUsecase) {
	handler := &RadusergroupHandler{ucase: uc}

	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString(`administrator.key`)),
	})

	group := e.Group("/api/rad/usergroup", isLoggedIn)
	group.GET("", handler.Fetch)
	group.GET("/load/:username", handler.LoadProfile)
	group.POST("", handler.Save)
	group.PUT("", handler.SaveProfile)
	group.DELETE("/:username", handler.Delete)

}

// Fetch ...
func (hn *RadusergroupHandler) Fetch(e echo.Context) error {

	res, err := hn.ucase.Fetch(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}
// LoadProfile ...
func (hn *RadusergroupHandler) LoadProfile(e echo.Context) error {
	username := e.Param("username")
	if username == "" {
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: "required field"})
	}

	res, err := hn.ucase.LoadProfile(e.Request().Context(), username)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{ Message: err.Error() })
	}

	return e.JSON(http.StatusOK, res)
}

// SaveProfile ....
func (hn *RadusergroupHandler) SaveProfile(e echo.Context) error {
	var profile domain.Profile
	err := e.Bind(&profile)
	if err != nil {
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: err.Error()})
	}

	nameArr := strings.Split(profile.ProfileName, " ")
	name := strings.Join(nameArr, "")
	profile.ProfileName = name

	result, err := hn.ucase.SaveProfile(e.Request().Context(), &profile)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusCreated, result)
}

// Save ...
func (hn *RadusergroupHandler) Save(e echo.Context) error {
	var radusergroup domain.Radusergroup
	err := e.Bind(&radusergroup)
	if err != nil {
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: err.Error()})
	}

	err = hn.ucase.Save(e.Request().Context(), &radusergroup)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusCreated, radusergroup)
}

// Delete ...
func (hn *RadusergroupHandler) Delete(e echo.Context) error {
	username := e.Param("username")
	if username == "" {
		return e.JSON(http.StatusFailedDependency, ResponseError{Message: "required field"})
	}

	resRug, err := hn.ucase.Get(e.Request().Context(), username)
	if err != nil {
		return e.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	err = hn.ucase.Delete(e.Request().Context(), username)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return e.JSON(http.StatusAccepted, resRug)
}
