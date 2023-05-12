package opp

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const baseURL = "https://api.onlinebetaalplatform.nl"
const allMerchantsEndpoint = baseURL + "/v1/merchants"

func GetMerchant(merchantId string, queryParams string) Merchant {
	apiKey, err := os.ReadFile(os.Getenv("OPP_API_PATH_TO_KEY"))

	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("GET", allMerchantsEndpoint+"/"+merchantId+"?"+queryParams, nil)
	req.Header.Set("Authorization", string(apiKey))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	var result Merchant
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func GetMerchantFromPages(pageNumber int) ApiListResponse[Merchant] {
	req, _ := http.NewRequest("GET", allMerchantsEndpoint, nil)
	req.Header.Set("Authorization", "")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	var result ApiListResponse[Merchant]
	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}

	return result
}

type ApiListResponse[T any] struct {
	HasMore     bool `json:"has_more"`
	CurrentPage int  `json:"current_page"`
	Data        []T  `json:"data"`
}

type Merchant struct {
	Uid        string `json:"uid"`
	Object     string `json:"object"`
	Created    int64  `json:"created"`
	Updated    int64  `json:"updated"`
	Status     string `json:"status"`
	Compliance struct {
		Level  int    `json:"level"`
		Status string `json:"status"`
	} `json:"compliance"`
	Metadata Metadata `json:"metadata"`
	Type     string   `json:"type"`
}

type Metadata []struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
