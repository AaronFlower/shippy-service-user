package main

import (
	"fmt"
	"log"

	pb "github.com/aaronflower/shippy-service-user/proto/auth"
	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/mdns"
)

func main() {
	// create a database connection and handles
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	// Automatically migrates the user struct into database columns/types etc.
	// This will check for changes and migrate them each itme this service is restrarted.
	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}
	tokenService := &TokenService{repo}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(
		micro.Name("shippy.auth"),
		micro.Version("latest"),
	)

	srv.Init()

	// Get instance of the broker using our defaults.
	publisher := micro.NewPublisher("user.created", srv.Client())

	pb.RegisterAuthHandler(srv.Server(), &service{repo, tokenService, publisher})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
