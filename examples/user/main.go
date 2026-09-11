package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	pterodactyl "github.com/tehzwen/pterodactyl-go/pkg/application"
	"github.com/tehzwen/pterodactyl-go/pkg/application/types"
)

var (
	API_KEY         = ""
	PTERO_HOST_ADDR = "" // ie: http://localhost
)

func main() {
	ctx := context.Background()
	ptero, err := pterodactyl.NewClient(PTERO_HOST_ADDR, pterodactyl.WithApiKey(API_KEY))
	if err != nil {
		log.Fatal().Err(err)
	}

	users, err := ptero.Users.ListUsers(ctx, types.ListUsersFilters{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", users)

	newUser, err := ptero.Users.CreateUser(ctx, types.CreateUserRequest{
		Email:     "fake@fake.com",
		Username:  "fake",
		FirstName: "Fake",
		LastName:  "Fiction",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(newUser)

	if err := ptero.Users.UpdateUser(ctx, types.UpdateUserRequest{
		UserId:    newUser.Attributes.ID,
		Email:     "fake@fake.com",
		Username:  "fake",
		FirstName: "Science",
		LastName:  "Fiction",
	}); err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)

	if err := ptero.Users.DeleteUser(ctx, newUser.Attributes.ID); err != nil {
		panic(err)
	}
}
