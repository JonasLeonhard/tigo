package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) HealthHandler(c echo.Context) error {
	user, _ := s.db.GetUserFromRequestToken(c)

	usernameOrAnonymous := "anonymous"
	if user != nil {
		usernameOrAnonymous = user.Name
	}

	status := map[string]interface{}{
		"db":         s.db.Health(),
		"tailwind":   "compiled",
		"loggedInAs": usernameOrAnonymous,
	}
	return c.JSON(http.StatusOK, status)
}
