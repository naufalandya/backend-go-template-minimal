package scylla

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitScylla() {
	cluster := gocql.NewCluster(os.Getenv("SCYLLA_HOSTS")) // example: "127.0.0.1"
	cluster.Keyspace = os.Getenv("SCYLLA_KEYSPACE")        // example: "your_keyspace"
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second

	// Optional: if your ScyllaDB requires authentication
	username := os.Getenv("SCYLLA_USERNAME")
	password := os.Getenv("SCYLLA_PASSWORD")
	if username != "" && password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: username,
			Password: password,
		}
	}

	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Unable to connect to ScyllaDB: %v", err)
	}

	fmt.Println("Connected to ScyllaDB~! (｡♥‿♥｡) ✨")
}

func CloseScylla() {
	if Session != nil {
		Session.Close()
		fmt.Println("Disconnected from ScyllaDB nya~ 🐾")
	}
}
