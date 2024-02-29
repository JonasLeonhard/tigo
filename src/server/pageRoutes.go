package server

import (
  "github.com/labstack/echo/v4"
  "github.com/jonasleonhard/go-htmx-time/src/templates/pages"
)

func (s *Server) IndexHandler(c echo.Context) error {
	user, _ := s.db.GetUserFromRequestToken(c)
	component := pages.IndexPage(user)

	return s.asHtml(c, component)
}

func (s *Server) LoginHandler(c echo.Context) error {
	user, _ := s.db.GetUserFromRequestToken(c)
	component := pages.LoginPage(user, "")

	return s.asHtml(c, component)
}

func (s *Server) RegisterHandler(c echo.Context) error {
	user, _ := s.db.GetUserFromRequestToken(c)
	component := pages.RegisterPage(user, "", "", "")

	return s.asHtml(c, component)
}
