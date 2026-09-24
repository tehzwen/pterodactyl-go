package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	pterodactyl "github.com/tehzwen/pterodactyl-go/pkg/application"
	"github.com/tehzwen/pterodactyl-go/pkg/application/types"
)

func main() {
	PTERO_HOST_ADDR := ""
	API_KEY := ""

	ctx := context.Background()
	ptero, err := pterodactyl.NewApplicationApi(PTERO_HOST_ADDR, pterodactyl.WithApiKey(API_KEY))
	if err != nil {
		log.Fatal().Err(err).Msg("could not create new application client")
	}

	mcEgg, err := ptero.Nests.GetEggDetails(ctx, 1, 4, []string{"variables"})
	if err != nil {
		log.Fatal().Err(err).Msg("could not lookup mc egg")
	}

	testServer, err := ptero.Servers.CreateServer(ctx, types.CreateServerRequest{
		Name: "zach test server",
		User: 1,
		Egg:  mcEgg.Attributes.ID,
		Limits: types.Limits{
			Memory: 2048,
			Disk:   2000,
			IO:     500,
			CPU:    100,
		},
		DockerImage: mcEgg.Attributes.DockerImage,
		Startup:     mcEgg.Attributes.Startup,
		// Allocation: &types.CreateServerAllocation{
		// 	Default: 1,
		// },
		Environment: map[string]string{
			"SERVER_JARFILE":  "server.jar",
			"VANILLA_VERSION": "latest",
		},
		Deploy: &types.CreateServerDeploy{
			Locations:   []int{1},
			DedicatedIp: false,
			PortRange:   []int{},
		},
		StartOnCompletion: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("could not create server")
	}

	servers, err := ptero.Servers.ListServers(ctx, &types.PteroListRequest{})
	if err != nil {
		log.Fatal().Err(err).Msg("could not list servers")
	}

	fmt.Printf("%+v\n", servers[1])

	// update this server as a test
	if err := ptero.Servers.UpdateServerBuild(ctx, types.UpdateServerBuildRequest{
		ServerId:      testServer.Attributes.Id,
		Allocation:    testServer.Attributes.Allocation,
		Memory:        100,
		CPU:           testServer.Attributes.Limits.CPU,
		Disk:          testServer.Attributes.Limits.DiskMB,
		FeatureLimits: testServer.Attributes.FeatureLimits,
		IO:            testServer.Attributes.Limits.Io,
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to update server")
	}

	if err := ptero.Servers.UpdateServerDetails(ctx, types.UpdateServerDetailsRequest{
		ServerId: testServer.Attributes.Id,
		UserId:   testServer.Attributes.User,
		Name:     "my minecraft test - updated",
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to update server details")
	}

	if err := ptero.Servers.UpdateServerStartup(ctx, types.UpdateServerStartupRequest{
		ServerId:    testServer.Attributes.Id,
		Startup:     testServer.Attributes.Container.StartupCommand,
		Environment: testServer.Attributes.Container.Environment,
		Egg:         testServer.Attributes.Egg,
		SkipScripts: testServer.Attributes.Container.SkipScripts,
		Image:       testServer.Attributes.Container.Image,
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to update server startup")
	}

	if err := ptero.Servers.SuspendServer(ctx, testServer.Attributes.Id); err != nil {
		log.Fatal().Err(err).Msg("failed to suspend server")
	}

	if err := ptero.Servers.UnsuspendServer(ctx, testServer.Attributes.Id); err != nil {
		log.Fatal().Err(err).Msg("failed to suspend server")
	}

	// if err := ptero.Servers.ReinstallServer(ctx, testServer.Attributes.Id); err != nil {
	// 	log.Fatal().Err(err).Msg("failed to reinstall server")
	// }

	if err := ptero.Servers.DeleteServer(ctx, servers[1].Attributes.Id); err != nil {
		log.Fatal().Err(err).Msg("failed to delete server")
	}
}
