package cassandra_addressbook

import (
	"log"
	"testing"
)

func TestLoadUserData(t *testing.T) {
	userData := SelectUsersEventsFromFile()

	log.Print("Size: %n:", len(userData))
}
