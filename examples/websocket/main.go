package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	pterodactyl "github.com/tehzwen/pterodactyl-go/pkg/client"
	"github.com/tehzwen/pterodactyl-go/pkg/client/types"
	"github.com/tehzwen/pterodactyl-go/pkg/websocket"
)

var (
	API_KEY         = ""
	PTERO_HOST_ADDR = "" // ie: http://localhost
)

func main() {
	ctx := context.Background()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	signal.Notify(sigs, syscall.SIGINT)

	// in order to get websocket connectivity, we need token from the client api
	ptero, err := pterodactyl.NewClientApi(PTERO_HOST_ADDR, pterodactyl.WithApiKey(API_KEY))
	if err != nil {
		panic(err)
	}

	servers, err := ptero.Servers.ListServers(ctx, types.ListServersParams{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", servers[0])

	socketDetails, err := ptero.Servers.GetConsoleAccess(ctx, servers[0].Attributes.Identifier)
	if err != nil {
		panic(err)
	}

	ws, err := websocket.NewWebsocketConnection(socketDetails.Data.Token, socketDetails.Data.Socket)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := ws.AuthAndListen(func(received websocket.SocketMessage) error {
			fmt.Printf("Received: %+v\n", received)
			return nil
		}); err != nil {
			panic(err)
		}
	}()

	startMessage := websocket.SocketMessage{
		Event: websocket.SocketEventTypeSetState,
		Args:  []string{"start"},
	}

	if err := ws.SendCommand(startMessage); err != nil {
		panic(err)
	}
	fmt.Printf("successfully sent message %+v\n", startMessage)

	signal := <-sigs
	fmt.Printf("caught signal %s, gracefully shutting down", signal)
	ws.Shutdown()
}
