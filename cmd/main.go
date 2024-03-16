package main

import (
	"crypto/rand"
	"crypto/rsa"
	usersRepo "github.com/SawitProRecruitment/UserService/repository/users"
	"github.com/SawitProRecruitment/UserService/usecase/users"
	"os"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	userRepository, err := usersRepo.NewUserRepository(repository.NewRepositoryOptions{Dsn: os.Getenv("DATABASE_URL")})
	if err != nil {
		panic(err)
	}

	// for demo purposes, we'll just generate a key with the same lifetime as the server,
	// meaning all JWT will be invalidated on restart
	rsaPrivateKey, _ := rsa.GenerateKey(rand.Reader, 1024)

	userUsecase := users.NewUserUsecases(users.NewUserUsecasesOptions{
		UserRepo:  userRepository,
		JwtSecret: rsaPrivateKey,
		JwtTtl:    10 * time.Minute,
	})

	opts := handler.NewServerOptions{
		UserUsecase: userUsecase,
	}
	return handler.NewServer(opts)
}
