package client

import (
	"log"
	"modular_monolith/protobuf/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Clients *ServiceClients

type ServiceClients struct {
	HelloWorldClient   api.HelloWorldClient
	LegalServiceClient api.FileServiceClient
}

func Connect() (*ServiceClients, error) {
	conn, err := grpc.NewClient("localhost:3551", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		if err != nil {
			log.Printf("Failed to connect to HelloWorld service: %v", err)
			return nil, err
		}
	}

	connLegals, err := grpc.NewClient("103.196.155.16:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		if err != nil {
			log.Printf("Failed to connect to Legal service: %v", err)
			return nil, err
		}
	}

	return &ServiceClients{
		HelloWorldClient:   api.NewHelloWorldClient(conn),
		LegalServiceClient: api.NewFileServiceClient(connLegals),
	}, nil
}

// // ServiceClients holds all your service clients
// type ServiceClients struct {
// 	HelloWorldClient pbapi.HelloWorldClient
// 	UserServiceClient pbapi.UserServiceClient  // Another service
// 	ProductServiceClient pbapi.ProductServiceClient // Example for Product Service
// }

// // Connect initializes and returns service clients
// func Connect() (*ServiceClients, error) {
// 	// First connection to HelloWorld service
// 	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Printf("Failed to connect to HelloWorld service: %v", err)
// 		return nil, err
// 	}

// 	// You can add more connections for different services (e.g., UserService, ProductService)
// 	connUser, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Printf("Failed to connect to UserService: %v", err)
// 		return nil, err
// 	}

// 	connProduct, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Printf("Failed to connect to ProductService: %v", err)
// 		return nil, err
// 	}

// 	// Returning a struct containing all clients
// 	return &ServiceClients{
// 		HelloWorldClient: pbapi.NewHelloWorldClient(conn),
// 		UserServiceClient: pbapi.NewUserServiceClient(connUser),
// 		ProductServiceClient: pbapi.NewProductServiceClient(connProduct),
// 	}, nil
// }
