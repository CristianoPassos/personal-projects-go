package main

import (
	"cloud.google.com/go/civil"
	"extractor/analytics"
	"extractor/opp"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const maxConcurrentJobs = 20

func main() {
	waitChan := make(chan struct{}, maxConcurrentJobs)
	var wg sync.WaitGroup
	executionTime := civil.DateTimeOf(time.Now())

	err := os.Setenv("OPP_API_PATH_TO_KEY", "PATH_TO_OPP_KEY")
	if err != nil {
		panic(err)
	}

	err = os.Setenv("BIG_QUERY_PATH_TO_FILE_KEY", "PATH_TO_GCP_KEY")
	if err != nil {
		panic(err)
	}

	missingKycUsers := analytics.SelectMissingKycUsers()

	for _, value := range missingKycUsers {
		wg.Add(1)
		waitChan <- struct{}{}

		go func(user analytics.MissingKycUser) {
			defer wg.Done()

			merchantOpp := opp.GetMerchant(user.MerchantId, "expand[]=metadata")

			userId, err := mapUserId(merchantOpp)

			if err != nil {
				log.Printf(err.Error()+"\nUser: %#v", user)

				return
			}

			merchantKyc := analytics.MerchantKyc{
				UserId:                  userId,
				ComplianceLevel:         mapKycLevel(merchantOpp.Compliance.Level),
				ComplianceStatus:        mapVerificationState(merchantOpp.Compliance.Status),
				BankAccountStatus:       user.BankAccountStatus,
				IdVerificationStatus:    user.VerificationState,
				MerchantId:              user.MerchantId,
				MerchantStatus:          mapMerchantStatus(merchantOpp.Status),
				MerchantCreatedDatetime: mapEpoch(merchantOpp.Created),
				MerchantUpdatedDatetime: mapEpoch(merchantOpp.Updated),
				UpdateDatetime:          executionTime,
			}

			err = analytics.InsertMerchantKyc(merchantKyc)
			if err != nil {
				panic(err)
			}

			<-waitChan
		}(value)
	}

	wg.Wait()
}

func mapKycLevel(complianceLevel int) string {
	switch complianceLevel {
	case 200:
		return "LOW_KYC"
	case 400:
		return "HIGH_KYC"
	case 500:
		return "HIGH_KYC"
	default:
		return "NO_KYC"
	}
}

func mapUserId(merchant opp.Merchant) (int, error) {
	for _, value := range merchant.Metadata {
		if value.Key == "user_id" {
			userId, err := strconv.Atoi(value.Value)

			if err != nil {
				panic(err)
			}

			return userId, nil
		}
	}

	return -1, fmt.Errorf("userId not found for merchant: %#v", merchant)
}

func mapVerificationState(verification string) string {
	switch verification {
	case "pending":
		return "PENDING"
	case "verified":
		return "VERIFIED"
	default:
		return "UNVERIFIED"
	}
}

func mapMerchantStatus(status string) string {
	return strings.ToUpper(status)
}

func mapEpoch(epoch int64) civil.DateTime {
	return civil.DateTimeOf(time.Unix(epoch, 0))
}
