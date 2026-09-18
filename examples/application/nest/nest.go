package main

import (
	"context"
	"fmt"

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
	ptero, err := pterodactyl.NewApplicationApi(PTERO_HOST_ADDR, pterodactyl.WithApiKey(API_KEY))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect")
	}
	nests, err := ptero.Nests.ListNests(ctx, &types.PteroListRequest{Include: []string{"eggs"}})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list nests")
	}
	fmt.Printf("%+v\n\n", nests)

	nest, err := ptero.Nests.GetNestDetails(ctx, nests[0].Attributes.ID, []string{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get nest details")
	}
	fmt.Printf("%+v\n\n", nest)

	nestEggs, err := ptero.Nests.ListNestEggs(ctx, nest.Attributes.ID, &types.PteroListRequest{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list nest eggs")
	}
	fmt.Printf("%+v\n\n", nestEggs)

	eggDetails, err := ptero.Nests.GetEggDetails(ctx, nest.Attributes.ID, nestEggs[0].Attributes.ID, []string{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get egg details")
	}
	fmt.Printf("%+v\n\n", eggDetails)
}
