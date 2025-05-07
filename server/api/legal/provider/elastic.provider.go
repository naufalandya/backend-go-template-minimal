package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"modular_monolith/server/config/elastic"
)

func InsertTextToElastic(text string) (string, error) {
	// Prepare the document
	doc := map[string]interface{}{
		"content": text,
	}

	// Serialize the document into JSON format
	docJSON, err := json.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal document: %v", err)
	}

	// Index the document into Elasticsearch (no need to specify document type)
	res, err := elastic.Client.Index(
		"documents",                              // Index name
		bytes.NewReader(docJSON),                 // Document body
		elastic.Client.Index.WithRefresh("true"), // Optionally refresh the index
	)
	if err != nil {
		return "", fmt.Errorf("failed to insert document: %v", err)
	}
	defer res.Body.Close()

	// Check if the response is an error
	if res.IsError() {
		return "", fmt.Errorf("error indexing document: %s", res.String())
	}

	// Parse the response to get the document ID
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse Elasticsearch response: %v", err)
	}

	// Extract the document ID from the response
	docID, ok := result["_id"].(string)
	if !ok {
		return "", fmt.Errorf("document ID not found in response")
	}

	return docID, nil
}
