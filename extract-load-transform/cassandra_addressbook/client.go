package cassandra_addressbook

import (
	"bufio"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"os"
)

var cassandraSession *gocql.Session

func SelectUsersEvents() []UserEvent {
	if cassandraSession == nil {
		openSession()
	}

	var userId int
	var eventType string
	var userEvents []UserEvent

	iter := cassandraSession.Query("SELECT user_id, event_type FROM address_book.event_by_user_id").Iter()

	for iter.Scan(&userId, &eventType) {
		userEvents = append(userEvents, UserEvent{userId, eventType})
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	return userEvents
}

func openSession() {
	var err error

	cluster := gocql.NewCluster("localhost:9042")

	cluster.Keyspace = "keyspace"
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "", Password: ""}

	cassandraSession, err = cluster.CreateSession()

	if err != nil {
		log.Fatal(err)
	}

}

func SelectUsersEventsFromFile() []UserEvent {
	var userEvents []UserEvent

	file, err := os.Open("FILE_PATH")

	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return userEvents
}

type UserEvent struct {
	UserId    int
	EventType string
}
