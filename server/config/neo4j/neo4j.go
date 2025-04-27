package neo4j

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var Neo4jDriver neo4j.Driver

// InitNeo4j initializes the Neo4j connection.
func InitNeo4j() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read Neo4j credentials from environment variables
	neo4jURI := os.Getenv("NEO4J_URI")
	neo4jUsername := os.Getenv("NEO4J_USERNAME")
	neo4jPassword := os.Getenv("NEO4J_PASSWORD")

	// Default values if not set in .env
	if neo4jURI == "" {
		neo4jURI = "bolt://localhost:7687" // default Neo4j connection URL
	}
	if neo4jUsername == "" {
		neo4jUsername = "neo4j" // default username
	}
	if neo4jPassword == "" {
		neo4jPassword = "password" // default password
	}

	// Create a Neo4j driver instance
	var errDriver error
	Neo4jDriver, errDriver = neo4j.NewDriver(neo4jURI, neo4j.BasicAuth(neo4jUsername, neo4jPassword, ""))
	if errDriver != nil {
		log.Fatalf("Failed to create Neo4j driver: %v", errDriver)
	}

	// Check if the connection is successful
	err = Neo4jDriver.VerifyConnectivity()
	if err != nil {
		log.Fatalf("Failed to verify Neo4j connectivity: %v", err)
	}

	fmt.Println("Successfully connected to Neo4j~ ₍˶ˆ꒳ˆ˶₎✨")
}

// CloseNeo4j closes the Neo4j connection.
func CloseNeo4j() {
	if Neo4jDriver != nil {
		err := Neo4jDriver.Close()
		if err != nil {
			log.Fatalf("Failed to close Neo4j connection: %v", err)
		}
	}
}
