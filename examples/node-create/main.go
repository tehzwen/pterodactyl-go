package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	pterodactyl "github.com/tehzwen/pterodactyl-go/pkg"
)

var (
	API_KEY         = ""
	PTERO_HOST_ADDR = "" // ie: http://localhost
)

func main() {
	ctx := context.Background()
	ptero, err := pterodactyl.NewApplicationClient(PTERO_HOST_ADDR, pterodactyl.WithApiKey(API_KEY))
	if err != nil {
		log.Fatal().Err(err)
	}

	newNode, err := ptero.CreateNode(ctx, pterodactyl.CreateNodeRequest{
		Name:       "test-node",
		LocationId: 1,
		FQDN:       "test.localdomain",
		Memory:     1000,
		Disk:       1000,
		Scheme:     "http",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create node")
	}
	fmt.Println(newNode)

	if err := ptero.CreateNodeAllocation(ctx, pterodactyl.CreateNodeAllocationRequest{
		NodeId: newNode.Attributes.ID,
		Ip:     "192.168.0.170",
		Ports:  []string{"20200"},
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to create node allocation")
	}

	// wait then delete the node
	time.Sleep(time.Second * 10)
	if err := ptero.DeleteNode(ctx, newNode.Attributes.ID); err != nil {
		log.Fatal().Err(err).Msg("failed to delete node")
	}
}
