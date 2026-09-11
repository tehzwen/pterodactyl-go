package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	pterodactyl "github.com/tehzwen/pterodactyl-go/pkg/client"
	"github.com/tehzwen/pterodactyl-go/pkg/client/types"
)

var (
	API_KEY         = ""
	PTERO_HOST_ADDR = "" // ie: http://localhost
)

func main() {
	ctx := context.Background()
	ptero, err := pterodactyl.NewClientApi(PTERO_HOST_ADDR, pterodactyl.WithApiKey(API_KEY))
	if err != nil {
		log.Fatal().Err(err)
	}

	servers, err := ptero.Servers.ListServers(ctx, types.ListServersParams{})
	if err != nil {
		log.Fatal().Err(err)
	}

	fmt.Printf("%+v\n", servers[0])
}
