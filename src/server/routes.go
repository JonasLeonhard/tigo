package server

import (
  "os"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Static("/static", "static")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("COOKIE_SECRET")))))

  // pageRoutes:
	e.GET("/", s.IndexHandler)
	e.GET("/login", s.LoginHandler)
	e.GET("/register", s.RegisterHandler)
	e.GET("/logout", s.UserLogoutHandler)
	e.GET("/user/dashboard", s.UserDashboardHandler)

	// userRoutes
	e.POST("/user/register", s.UserRegisterHandler)
	e.POST("/user/login", s.UserLoginHandler)

	// ...components for htmx
	// e.GET("header/usermenu", s.HeaderUserMenuHandler)

  // miscRoutes:
	e.GET("/health", s.HealthHandler)

	return e
}

func (s *Server) asHtml(c echo.Context, component templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return component.Render(c.Request().Context(), c.Response().Writer)
}
