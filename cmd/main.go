package main

import (
	"encoding/json"
	"fmt"
)

type SuperComplex struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Metadata    Metadata        `json:"metadata"`
	Relations   []Relation      `json:"relations"`
	NestedArray [][]NestedItem  `json:"nested_array"`
	DeepMap     map[string]Deep `json:"deep_map"`
}

type Metadata struct {
	Version     string                 `json:"version"`
	Properties  map[string]interface{} `json:"properties"`
	Permissions []Permission           `json:"permissions"`
}

type Permission struct {
	Role  string `json:"role"`
	Level int    `json:"level"`
}

type Relation struct {
	Type    string `json:"type"`
	RefID   int    `json:"ref_id"`
	Details struct {
		ConnectedAt string `json:"connected_at"`
		Active      bool   `json:"active"`
	} `json:"details"`
}

type NestedItem struct {
	Index int    `json:"index"`
	Value string `json:"value"`
}

type Deep struct {
	Keys map[string]struct {
		Data string `json:"data"`
		Code int    `json:"code"`
	} `json:"keys"`
}

func main() {
	data := []SuperComplex{
		{
			ID:   1,
			Name: "ComplexObject1",
			Metadata: Metadata{
				Version: "1.0",
				Properties: map[string]interface{}{
					"created_by": "User1",
					"size":       1024,
					"features":   []string{"fast", "secure"},
				},
				Permissions: []Permission{
					{"admin", 10},
					{"user", 5},
				},
			},
			Relations: []Relation{
				{
					Type:  "friend",
					RefID: 101,
					Details: struct {
						ConnectedAt string `json:"connected_at"`
						Active      bool   `json:"active"`
					}{
						ConnectedAt: "2024-02-03",
						Active:      true,
					},
				},
			},
			NestedArray: [][]NestedItem{
				{
					{Index: 0, Value: "A"},
					{Index: 1, Value: "B"},
				},
				{
					{Index: 2, Value: "C"},
					{Index: 3, Value: "D"},
				},
			},
			DeepMap: map[string]Deep{
				"config": {
					Keys: map[string]struct {
						Data string `json:"data"`
						Code int    `json:"code"`
					}{
						"setting1": {"Enabled", 1},
						"setting2": {"Disabled", 0},
					},
				},
			},
		},
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonData))
}
