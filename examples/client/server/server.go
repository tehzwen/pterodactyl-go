package main

import (
	"context"
	"fmt"

	pterodactyl "github.com/tehzwen/pterodactyl-go/pkg/client"
	"github.com/tehzwen/pterodactyl-go/pkg/client/types"
)

var (
	API_KEY         = "ptlc_bQlBpIt3IfsDk1iliEa4R5PqbnedQNHRzOaXxgI8OaM"
	PTERO_HOST_ADDR = "http://sailor:81" // ie: http://localhost
)

func main() {
	ctx := context.Background()
	ptero, err := pterodactyl.NewClientApi(PTERO_HOST_ADDR, pterodactyl.WithApiKey(API_KEY))
	if err != nil {
		panic(err)
	}

	servers, err := ptero.Servers.ListServers(ctx, types.ListServersParams{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n\n", servers[0])

	server, err := ptero.Servers.GetServerDetails(ctx, servers[0].Attributes.Identifier)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n\n", server)

	resources, err := ptero.Servers.GetServerResources(ctx, server.Attributes.Identifier)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n\n", resources)

	consoleDetails, err := ptero.Servers.GetConsoleAccess(ctx, server.Attributes.Identifier)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n\n", consoleDetails)

	logs, err := ptero.Servers.GetServerActivity(ctx, server.Attributes.Identifier)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n\n", logs)

	// if err := ptero.Servers.ReinstallServer(ctx, server.Attributes.Identifier); err != nil {
	// 	panic(err)
	// }

	// if err := ptero.Servers.ManagePower(ctx, server.ManagePowerRequest{
	// 	ServerIdentifier: servers[0].Attributes.Identifier,
	// 	Signal:           server.PowerActionStop,
	// }); err != nil {
	// 	panic(err)
	// }
}
