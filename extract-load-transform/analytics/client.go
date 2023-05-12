package analytics

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"os"
)

var ctx = context.Background()
var bigQueryClient *bigquery.Client

func SelectMissingKycUsers() []MissingKycUser {
	if bigQueryClient == nil {
		openSession()
	}

	missingKycUsers := make([]MissingKycUser, 0)
	query := bigQueryClient.Query("SELECT users.user_id, users.merchant_id, users.verification_state, users.bank_account_status " +
		"FROM `analytics.users` AS users " +
		"LEFT JOIN `analytics.kyc` AS kyc ON users.user_id = kyc.user_id " +
		"WHERE kyc.user_id IS NULL")

	rows, err := query.Read(context.Background())

	if err != nil {
		panic(err)
	}

	for {
		var row MissingKycUser
		err := rows.Next(&row)
		if err == iterator.Done {
			break
		}

		if err != nil {
			panic(err)
		}

		missingKycUsers = append(missingKycUsers, row)
	}

	return missingKycUsers
}

func InsertMerchantKyc(kyc MerchantKyc) error {
	if bigQueryClient == nil {
		openSession()
	}

	inserter := bigQueryClient.Dataset("analytics").Table("kyc").Inserter()

	if err := inserter.Put(ctx, kyc); err != nil {
		return err
	}

	return nil
}

func openSession() {
	pathToFileKey := os.Getenv("BIG_QUERY_PATH_TO_FILE_KEY")

	client, err := bigquery.NewClient(ctx, "projectId", option.WithCredentialsFile(pathToFileKey))

	if err != nil {
		panic(err)
	}

	bigQueryClient = client
}

type MissingKycUser struct {
	UserId            int                 `bigquery:"user_id"`
	MerchantId        string              `bigquery:"merchant_id"`
	VerificationState bigquery.NullString `bigquery:"verification_state"`
	BankAccountStatus bigquery.NullString `bigquery:"bank_account_status"`
}

type MerchantKyc struct {
	UserId                  int                 `bigquery:"user_id"`
	ComplianceLevel         string              `bigquery:"compliance_level"`
	ComplianceStatus        string              `bigquery:"compliance_status"`
	BankAccountStatus       bigquery.NullString `bigquery:"bank_account_status"`
	IdVerificationStatus    bigquery.NullString `bigquery:"id_verification_status"`
	MerchantId              string              `bigquery:"merchant_id"`
	MerchantStatus          string              `bigquery:"merchant_status"`
	MerchantCreatedDatetime civil.DateTime      `bigquery:"merchant_created_datetime"`
	MerchantUpdatedDatetime civil.DateTime      `bigquery:"merchant_updated_datetime"`
	UpdateDatetime          civil.DateTime      `bigquery:"update_datetime"`
}
