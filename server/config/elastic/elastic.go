package elastic

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

// Global Elasticsearch client instance
var Client *elasticsearch.Client

// InitElasticSearch initializes the Elasticsearch client
func InitElasticSearch() error {
	// Configure Elasticsearch client
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // Elasticsearch server address
		},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
		return err
	}

	// Check the connection by pinging Elasticsearch
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()

	// If successful, assign to global Client
	Client = es

	// Print basic info about the connection
	fmt.Println("Connected to Elasticsearch")
	return nil
}
