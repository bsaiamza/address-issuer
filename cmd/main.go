package main

import (
	"address_issuer/api"
	"address_issuer/pkg/client"
	"address_issuer/pkg/config"
	"address_issuer/pkg/log"
	"address_issuer/pkg/server"
	"address_issuer/pkg/utils"
)

func main() {
	runServer()
}

func runServer() {
	config := config.LoadConfig()
	acapyClient := client.NewClient(config.GetAcapyURL())
	cache := utils.NewBigCache()

	srv := server.NewServer().
		WithAddress(config.GetServerAddress()).
		WithRouter(api.NewRouter(config, acapyClient, cache)).
		WithErrorLogger(log.ServerError)

	go func() {
		log.ServerInfo.Println("-----------------------------------------")
		log.ServerInfo.Println("|		Address Issuer		|")
		log.ServerInfo.Println("-----------------------------------------")
		log.ServerInfo.Println("")
		log.ServerInfo.Printf("Server started on: %s", config.GetServerAddress())
		if err := srv.Start(); err != nil {
			log.ServerError.Fatal(err)
		}
	}()

	utils.GracefulServerExit(func() {
		if err := srv.Stop(); err != nil {
			log.ServerError.Printf("Failed to stop server: %s", err.Error())
		}
	})
}
