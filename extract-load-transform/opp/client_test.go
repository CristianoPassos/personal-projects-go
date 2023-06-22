package opp

import (
	"encoding/json"
	"os"
	"testing"
)

func TestGetMerchant(t *testing.T) {
	err := os.Setenv("OPP_API_PATH_TO_KEY", "PATH_TO_OPP_KEY")
	if err != nil {
		panic(err)
	}

	apiResponse := GetMerchant("mer_id", "expand[]=metadata")

	if apiResponse.Uid == "" {
		t.Errorf("error")
	}
}

func TestApiResponse2(t *testing.T) {
	apiResponse := GetMerchantFromPages(1)

	if apiResponse.HasMore != true {
		t.Errorf("error")
	}
}

func TestApiResponse(t *testing.T) {

	data, _ := os.ReadFile("merchants.json")

	var result ApiListResponse[Merchant]
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
