package main

import (
	"context"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/repository/users"
)

func main() {
	ctx := context.Background()

	repo, err := users.NewUserRepository(repository.NewRepositoryOptions{Dsn: "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable"})
	if err != nil {
		panic(err)
	}

	//res, err := repo.CreateUser(ctx, repository.CreateUserInput{
	//	PhoneNo:      "+6281214173377",
	//	FullName:     "Fairuz Astra Karta",
	//	PasswordHash: []byte("password-hash-sample-2"),
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(res.ID)

	//res, err := repo.GetUser(ctx, repository.GetUserInput{PhoneNo: "+6281214173388"})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(res.ID)
	//fmt.Println(res.PhoneNo)
	//fmt.Println(res.FullName)
	//fmt.Println(res.PasswordHash)
	//fmt.Println(res.SuccessfulLoginCount)

	_, err = repo.UpdateUser(ctx, repository.UpdateUserInput{
		ID:                   15,
		PhoneNo:              "+6281214173377",
		FullName:             "",
		SuccessfulLoginCount: 0,
	})
	if err != nil {
		panic(err)
	}
}
