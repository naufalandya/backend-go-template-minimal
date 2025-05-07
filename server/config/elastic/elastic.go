package elastic

import (
	"fmt"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

var Client *elasticsearch.Client

func InitElasticSearch() error {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://elasticsearch:9200"
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esURL},
	})

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
		return err
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()

	Client = es

	fmt.Println("Connected to Elasticsearch")
	return nil
}
