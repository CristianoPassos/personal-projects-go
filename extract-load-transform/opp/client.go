package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const BaseURL = "https://api.onlinebetaalplatform.nl"
const AllMerchantsUrl = BaseURL + "/v1/merchants?expand[]=metadata"

func getMerchantFromPages(pageNumber int) ApiListResponse[Merchant] {
	req, _ := http.NewRequest("GET", AllMerchantsUrl, nil)
	req.Header.Set("Authorization", "")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	var result ApiListResponse[Merchant]
	json.NewDecoder(resp.Body).Decode(&result)

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
	Created    int    `json:"created"`
	Updated    int    `json:"updated"`
	Status     string `json:"status"`
	Compliance struct {
		Level  int    `json:"level"`
		Status string `json:"status"`
	} `json:"compliance"`
	Type     string `json:"type"`
	Metadata []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"metadata"`
}
