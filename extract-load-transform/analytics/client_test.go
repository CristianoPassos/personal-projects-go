package analytics

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"os"
	"testing"
	"time"
)

func TestLoadMissingKycUsers(t *testing.T) {
	os.Setenv("BIG_QUERY_PATH_TO_FILE_KEY", "C:\\dev\\git\\gcp_key.json")

	merchantIds := SelectMissingKycUsers()

	if len(merchantIds) == 0 {
		t.Errorf("Users not loaded")
	}
}

func TestInsertMerchantKyc(t *testing.T) {
	os.Setenv("BIG_QUERY_PATH_TO_FILE_KEY", "C:\\dev\\git\\gcp_key.json")

	now := civil.DateTimeOf(time.Now())
	merchant := MerchantKyc{
		UserId:                  1,
		ComplianceLevel:         "null",
		ComplianceStatus:        "null",
		BankAccountStatus:       bigquery.NullString{},
		IdVerificationStatus:    bigquery.NullString{},
		MerchantId:              "null",
		MerchantStatus:          "null",
		MerchantCreatedDatetime: now,
		MerchantUpdatedDatetime: now,
		UpdateDatetime:          now,
	}

	err := InsertMerchantKyc(merchant)

	if err != nil {
		t.Errorf(err.Error())
	}
}
