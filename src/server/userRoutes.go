package server

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/gorilla/sessions"
	"github.com/jonasleonhard/go-htmx-time/src/templates/pages"
	"github.com/labstack/echo-contrib/session"
)

// /user/<route> handlers

func (s *Server) UserRegisterHandler(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm-password")

	var usernameError, emailError, passwordError string

	if username == "" {
		usernameError = "Username is required"
	}

	if email == "" {
		emailError = "Email is required"
	}

	if password != confirmPassword {
		passwordError = "Passwords do not match"
	}

	user, err := s.db.CreateUser(username, email, password, nil, c.Request().Context())

	if err != nil {
		log.Println(err)
		emailError = "Email is already taken"
	} else if user != nil && usernameError == "" && emailError == "" && passwordError == "" {
		return c.Redirect(http.StatusSeeOther, "/login?register=success")
	}

	component := pages.RegisterPage(user, usernameError, emailError, passwordError)
	return s.asHtml(c, component)
}

func (s *Server) UserLogoutHandler(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusSeeOther, "/?logout=success")
}


// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (s *Server) UserLoginHandler(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if !s.db.CheckForUserLogin(c.Request().Context(), email, password) {
		component := pages.LoginPage(nil, "Username or password is incorrect")
		return s.asHtml(c, component)
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	// set the token in the cookie of the request
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["token"] = t
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusSeeOther, "/?login=success")
}


func (s *Server) UserDashboardHandler(c echo.Context) error {
	user, _ := s.db.GetUserFromRequestToken(c)
	return c.String(http.StatusOK, "Welcome "+user.Name+"!")
}
