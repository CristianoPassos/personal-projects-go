package main

import (
	"encoding/json"
	"extractor/cassandra"
	"os"
	"testing"
)

func TestApiResponse(t *testing.T) {

	data, _ := os.ReadFile("merchants.json")

	var result cassandra.ApiListResponse[cassandra.Merchant]
	err := json.Unmarshal(data, &result)

	if err != nil {
		t.Errorf("%s", err)
	}

	if len(result.Data) != 10 {
		t.Errorf("Missing merchants")
	}

	if result.HasMore == false {
		t.Errorf("Should have more merchants")
	}

	value := result.Data[2].Metadata[0].Value

	if value != "68768055" {
		t.Errorf("user_id should be present")
	}

}
