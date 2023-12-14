package server

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/jonasleonhard/go-htmx-time/src/templates/pages"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Static("/static", "static")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", s.IndexHandler)
	e.GET("/user/create", s.UserCreateHandler)
	e.GET("/health", s.HealthHandler)

	return e
}

func asHtml(c echo.Context, component templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) IndexHandler(c echo.Context) error {
	user, err := s.db.GetUser("test", c.Request().Context())

	if err != nil {
		log.Println(err)
	}

	component := pages.IndexPage(user)

	return asHtml(c, component)
}

func (s *Server) UserCreateHandler(c echo.Context) error {
	user, err := s.db.CreateUser("test", "test@test.com", nil, c.Request().Context())

	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Server) HealthHandler(c echo.Context) error {
	status := map[string]interface{}{
		"db":       s.db.Health(),
		"tailwind": "compiled",
	}
	return c.JSON(http.StatusOK, status)
}
