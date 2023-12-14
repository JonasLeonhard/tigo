package server

import (
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/jonasleonhard/go-htmx-time/src/templates/pages"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Static("/static", "static")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret")))) // TODO: change to secret from env file

	e.GET("/", s.IndexHandler)
	e.GET("/login", s.LoginHandler)
	e.GET("/logout", s.UserLogoutHandler)

	e.GET("/user/register", s.UserRegisterHandler)
	e.POST("/user/login", s.UserLoginHandler)
	e.GET("/user/dashboard", s.UserDashboardHandler)

	e.GET("/health", s.HealthHandler)

	return e
}

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

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func asHtml(c echo.Context, component templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) IndexHandler(c echo.Context) error {
	user, _ := s.db.GetUserFromRequestToken(c)
	component := pages.IndexPage(user)

	return asHtml(c, component)
}

func (s *Server) LoginHandler(c echo.Context) error {
	component := pages.LoginPage(nil)

	return asHtml(c, component)
}

func (s *Server) UserLoginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throw unauthorized error if username & password dont match in db: TODO: implement db auth check here.
	if username == "test" || password == "test" {
		return c.Redirect(http.StatusSeeOther, "/login?login=unauthorized")
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		username, // TODO: use username here
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

func (s *Server) UserRegisterHandler(c echo.Context) error {
	user, err := s.db.CreateUser("test2", "test2@test.com", nil, c.Request().Context())

	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Server) UserDashboardHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
