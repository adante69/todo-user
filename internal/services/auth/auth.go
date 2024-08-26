package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
	"todo-sso/internal/domain/models"
	"todo-sso/internal/lib/jwt"
	"todo-sso/internal/storage"
	"todo-sso/internal/storage/postgres"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (user models.User, err error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (app models.App, err error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

// New returns new instane of the Auth service
func New(log *slog.Logger, provider AppProvider, saver *postgres.Storage, userProvider UserProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    saver,
		userProvider: userProvider,
		appProvider:  provider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string string, err error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("attempting to login")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found")

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid password")

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully logged in")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Warn("failed to generate token")

		return "", fmt.Errorf("%s: %w", op, slog.String("token not generated", err.Error()))
	}
	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (uid int64, err error) {
	const op = "auth.registerNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("pass hashed")
	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists")
		}
		return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
	}
	log.Info("registered user")
	return id, nil
}
