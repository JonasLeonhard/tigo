package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/jonasleonhard/go-htmx-time/src/database/ent"
	"github.com/jonasleonhard/go-htmx-time/src/database/ent/user"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Health() map[string]string
	Client() *ent.Client
	CreateUser(name string, email string, password string, age *int, ctx context.Context) (*ent.User, error)
	GetUser(name string, ctx context.Context) (*ent.User, error)
	GetUserFromRequestToken(c echo.Context) (*ent.User, error)
	CheckForUserLogin(c context.Context, username string, password string) bool
}

type service struct {
	db *ent.Client
}

var (
	dburl = os.Getenv("DB_URL")
)

func New() Service {
	db, err := ent.Open("sqlite3", fmt.Sprintf("file:%s?_fk=1", dburl))
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	// Run the auto migration tool.
	if err := db.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	s := &service{db: db}
	return s
}

func (s *service) Close() {
	if err := s.db.Close(); err != nil {
		log.Fatalf("failed closing database connection: %v", err)
	}
}

func (s *service) Client() *ent.Client {
	return s.db
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Run a simple query to check the database connection
	_, err := s.db.User.Query().Exist(ctx)

	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"status": "sqlite running",
	}
}

// --- Schema specific methods:
func (s *service) CreateUser(name string, email string, password string, age *int, ctx context.Context) (*ent.User, error) {
	passwordBytes := []byte(password)
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	if err != nil {
		return nil, fmt.Errorf("failed hashing password: %w", err)
	}

	hashedPassword := string(hashedPasswordBytes)

	query := s.db.User.Create().SetName(name).SetEmail(email).SetPasswordHash(hashedPassword)

	if age != nil {
		query.SetAge(*age)
	}

	user, err := query.Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	return user, nil
}

func (s *service) GetUser(name string, ctx context.Context) (*ent.User, error) {
	user, err := s.db.User.Query().Where(user.Name(name)).Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}

	return user, nil
}

// TODO: this is shared in routes.go
// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func (s *service) GetUserFromRequestToken(c echo.Context) (*ent.User, error) {
	sess, err := session.Get("session", c)
	tokenVal := sess.Values["token"]
	if err != nil || tokenVal == nil {
		return nil, fmt.Errorf("failed getting token from session: %w", err)
	}

	tokenStr := sess.Values["token"].(string)
	token, err := jwt.ParseWithClaims(tokenStr, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // TODO: use real secret here
	})

	if !token.Valid {
		return nil, fmt.Errorf("failed parsing token - invalid: %w", err)
	}

	claims := token.Claims.(*jwtCustomClaims)

	user, err := s.GetUser(claims.Name, c.Request().Context())

	return user, nil
}

func (s *service) CheckForUserLogin(ctx context.Context, username string, password string) bool {
	user, err := s.db.User.Query().Where(user.Name(username)).Only(ctx)

	if err != nil {
		return false
	}

	if user == nil {
		return false
	}

	pwErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if pwErr != nil {
		return false
	}

	return true
}
