package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/jonasleonhard/go-htmx-time/src/database/ent"
	"github.com/jonasleonhard/go-htmx-time/src/database/ent/user"
	_ "github.com/mattn/go-sqlite3"
)

type Service interface {
	Health() map[string]string
	Client() *ent.Client
	CreateUser(name string, email string, age *int, ctx context.Context) (*ent.User, error)
	GetUser(name string, ctx context.Context) (*ent.User, error)
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
func (s *service) CreateUser(name string, email string, age *int, ctx context.Context) (*ent.User, error) {
	query := s.db.User.Create().SetName(name).SetEmail(email)

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
