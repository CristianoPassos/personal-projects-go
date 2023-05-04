package cassandra

import (
	"testing"
)

func TestLoadUserData(t *testing.T) {
	userData := SelectUserData(3834807)

	if userData.BankAccountStatus != "PENDING" {
		t.Errorf("Invalid Bank Account Status")
	}

	if userData.VerificationState != "PENDING" {
		t.Errorf("Invalid Verification State")
	}
}
