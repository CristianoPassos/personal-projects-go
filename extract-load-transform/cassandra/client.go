package cassandra

import (
	"github.com/gocql/gocql"
	"time"
)

var CassandraSession *gocql.Session

func SelectUserData(userId int) UserData {
	if CassandraSession == nil {
		openSession()
	}

	userData := UserData{}

	CassandraSession.Query("SELECT verification_state,  bank_account_status FROM user_data WHERE user_id = ?", userId).Consistency(gocql.Quorum).Scan(&userData.VerificationState, &userData.BankAccountStatus)

	return userData
}

func openSession() {
	cluster := gocql.NewCluster("localhost:9042") //replace PublicIP with the IP addresses used by your cluster.
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "", Password: ""}
	cluster.Keyspace = "payment_and_shipping"

	CassandraSession, _ = cluster.CreateSession()
}

type UserData struct {
	VerificationState string
	BankAccountStatus string
}
