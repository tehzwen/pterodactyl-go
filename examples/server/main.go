package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	pterodactyl "github.com/tehzwen/pterodactyl-go/pkg/application"
	pterodactyl_server "github.com/tehzwen/pterodactyl-go/pkg/application/server"
	pterodactyl_application_types "github.com/tehzwen/pterodactyl-go/pkg/application/types"
)

func main() {
	PTERO_HOST_ADDR := "http://localhost"
	API_KEY := os.Getenv("PTERO_API_KEY")

	ctx := context.Background()
	ptero, err := pterodactyl.NewClient(PTERO_HOST_ADDR, pterodactyl.WithApiKey(API_KEY))
	if err != nil {
		log.Fatal().Err(err).Msg("could not create new application client")
	}

	testServer, err := ptero.Servers.CreateServer(ctx, pterodactyl_server.CreateServerRequest{
		Name: "zach test server",
		User: 1,
		Egg:  1,
		Limits: pterodactyl_application_types.Limits{
			Memory: 100,
			Disk:   100,
			IO:     10,
		},
		DockerImage: "testimage",
		Startup:     "echo 'hi'",
		Allocation: struct {
			Default int "json:\"default\""
			Backups int "json:\"backups\""
		}{
			Default: 2,
		},
		Environment: map[string]string{
			"SPONGE_VERSION": "1.12.2-7.3.0",
			"SERVER_JARFILE": "server.jar",
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("could not create server")
	}

	servers, err := ptero.Servers.ListServers(ctx, &pterodactyl_application_types.PteroListRequest{})
	if err != nil {
		log.Fatal().Err(err).Msg("could not list servers")
	}

	fmt.Printf("%+v\n", servers[1])

	// update this server as a test
	if err := ptero.Servers.UpdateServerBuild(ctx, pterodactyl_server.UpdateServerBuildRequest{
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

	if err := ptero.Servers.UpdateServerDetails(ctx, pterodactyl_server.UpdateServerDetailsRequest{
		ServerId: testServer.Attributes.Id,
		UserId:   testServer.Attributes.User,
		Name:     "my minecraft test - updated",
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to update server details")
	}

	if err := ptero.Servers.UpdateServerStartup(ctx, pterodactyl_server.UpdateServerStartupRequest{
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
