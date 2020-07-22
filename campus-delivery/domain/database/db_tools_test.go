package database

import (
	_ "campus-delivery/domain"
	"campus-delivery/domain/model"
	"fmt"
	"log"
	"testing"
)

var (
	id = int64(2577015)
	db = NewDBClient(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"94.250.250.51", 5432, "ruslan", "rulsan-206-astapov", "delivery"))
	user = model.User{
		Id:         id,
		NickName:   "Diroman",
		FirstName:  "Роман",
		SecondName: "Дубатов",
		Latitude:   34.123235,
		Longitude:  -45.23442,
		Rating:     0,
	}
	courier = model.Courier{
		User: user,
		Shop: "Перек",
	}
)

func TestDBClient_AddUser(t *testing.T) {
	db.Connect()
	if err := db.AddUser(user); err != nil {
		t.Errorf("AddUser() error = %v", err)
	}

	db.CloseConnection()
}

func TestDBClient_DeleteUser(t *testing.T) {
	db.Connect()
	courier, err := db.GetUserWithOrder(user)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(courier)

	db.CloseConnection()
}

func TestDBClient_DeleteUserByTimer(t *testing.T) {
	db.Connect()
	if err := db.AddUserWithOrder(&courier); err != nil {
		log.Print(err)
	}

	db.CloseConnection()
}
