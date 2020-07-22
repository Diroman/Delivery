package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"campus-delivery/domain"
	"campus-delivery/domain/model"
)

// FormatTimestamp, ParseTimestamp

const (
	User_Table        = "Users"
	Subscribers_Table = "Subscribers"
	Orders_Table      = "UserWithOrder"
)

type DBClient struct {
	url         string
	connection  *pgxpool.Pool
	deleteTimer map[int64]struct {
		table   string
		endTime time.Time
	}
}

func NewDBClient(url string) domain.DataBase {
	return &DBClient{
		url: url,
		deleteTimer: map[int64]struct {
			table   string
			endTime time.Time
		}{},
	}
}

func (db *DBClient) Connect() error {
	conn, err := pgxpool.Connect(context.Background(), db.url)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return err
	}

	db.connection = conn
	return nil
}

func (db *DBClient) CloseConnection() {
	db.connection.Close()
}

func (db *DBClient) DeleteUserByTimer() {
	couriers := db.GetAllCourier()
	for _, courier := range couriers {
		db.deleteTimer[courier.User.Id] = struct {
			table   string
			endTime time.Time
		}{
			table:   "",
			endTime: courier.TimeTo.Add(4 * time.Hour),
		}
	}
	for {
		time.Sleep(time.Minute)
		deleteKey := []int64{}
		//log.Printf("Now check time: %v", time.Now().Add(7*time.Hour))

		for key, value := range db.deleteTimer {
			if time.Now().Add(7 * time.Hour).After(value.endTime) {
				log.Printf("Delete user <%v>", key)
				deleteString := "DELETE FROM UserWithOrder WHERE user_id = $1"
				ctx := context.Background()

				_, err := db.connection.Exec(ctx, deleteString, key)
				if err != nil {
					log.Printf("Error to delete <%v> from <%v>", key, value.table)
				}
				deleteKey = append(deleteKey, key)
			}

			for _, key := range deleteKey {
				delete(db.deleteTimer, key)
			}
		}
	}
}

func (db *DBClient) AddUser(user model.User) error {
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()

	selectString := "SELECT * FROM Users WHERE id = $1"
	insert := "INSERT INTO Users VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := conn.Exec(ctx, insert,
		user.Id, user.NickName, user.FirstName, user.SecondName, user.Latitude, user.Longitude, user.Notification)
	if err != nil {
		log.Printf("Add user error: %v", err)
		if conn.QueryRow(ctx, selectString, user.Id) == nil {
			return err
		}
	}

	return nil
}

func (db *DBClient) DeleteUser(id int64) error {
	deleteString := "DELETE FROM Users WHERE id = $1"
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()

	_, err := conn.Exec(ctx, deleteString, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBClient) AddUserWithOrder(courier *model.Courier) error {
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()

	selectString := "SELECT * FROM UserWithOrder WHERE user_id = $1"
	insert := "INSERT INTO UserWithOrder VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := conn.Exec(ctx, insert, courier.User.Id, courier.Shop,
		courier.Description, courier.TimeFrom, courier.TimeTo, courier.Link, courier.ChatId)
	if err != nil {
		log.Printf("Add courier error: %v", err)
		if db.connection.QueryRow(ctx, selectString, courier.User.Id, courier.Shop,
			courier.Description, courier.TimeFrom, courier.TimeTo, courier.Link, courier.ChatId) == nil {
			return err
		}
	}

	db.deleteTimer[courier.User.Id] = struct {
		table   string
		endTime time.Time
	}{table: Orders_Table, endTime: courier.TimeTo.Add(4 * time.Hour)}
	log.Printf("Add delete timer: %v", courier.TimeTo.Add(4*time.Hour))
	return nil
}

func (db *DBClient) DeleteUserWithOrder(id int64) error {
	deleteString := "DELETE FROM UserWithOrder WHERE user_id = $1"
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()

	_, err := conn.Exec(ctx, deleteString, id)
	if err != nil {
		return err
	}

	delete(db.deleteTimer, id)
	return nil
}

func (db *DBClient) GetUserWithOrder(user model.User) ([]model.Courier, error) {
	//var timeFrom, timeTo int64
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()

	selectString := "SELECT user_id, first_name, shop, description, time_from, time_to, link, chat_id " +
		"FROM UserWithOrder LEFT JOIN Users s ON UserWithOrder.user_id = s.id " +
		"WHERE point($1, $2) <@> point(s.latitude, s.longitude) < 0.124374238 AND user_id <> $3"

	rows, err := conn.Query(ctx, selectString, user.Latitude, user.Longitude, user.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var couriers []model.Courier
	for rows.Next() {
		courier := model.Courier{User: model.User{}}
		if err := rows.Scan(&courier.User.Id, &courier.User.FirstName, &courier.Shop, &courier.Description,
			&courier.TimeFrom, &courier.TimeTo, &courier.Link, &courier.ChatId); err != nil {
			return nil, err
		}

		//courier.TimeFrom = time.Unix(timeFrom, 0)
		//courier.TimeTo = time.Unix(timeTo, 0)
		couriers = append(couriers, courier)
	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		return nil, err
	}
	return couriers, nil
}

func (db *DBClient) ChangeNotificationStatus(id int64, status bool) error {
	updateString := "UPDATE Users SET notification = $1 WHERE id = $2"
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()

	_, err := conn.Exec(ctx, updateString, status, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *DBClient) ChangeLocation(id int64, lat float64, long float64) error {
	updateString := "UPDATE Users SET latitude = $1, longitude = $2 WHERE id = $3"
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()

	_, err := conn.Exec(ctx, updateString, lat, long, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *DBClient) GetNotificationUser(user model.User) []model.User {
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()

	selectString := "SELECT id FROM Users WHERE point($1, $2) <@> point(latitude, longitude) < 0.124374238 AND id <> $3 " +
		"AND notification = True"

	rows, err := conn.Query(ctx, selectString, user.Latitude, user.Longitude, user.Id)
	if err != nil {
		log.Printf("Get notification error: %v", err)
		return []model.User{}
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		user := model.User{}
		if err := rows.Scan(&user.Id); err != nil {
			log.Printf("Get notification error: %v", err)
			return []model.User{}
		}

		//courier.TimeFrom = time.Unix(timeFrom, 0)
		//courier.TimeTo = time.Unix(timeTo, 0)
		users = append(users, user)
	}
	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		log.Printf("Get notification error: %v", err)
		return []model.User{}
	}
	return users
}

func (db *DBClient) GetUserInfo(id int64) (model.User, error) {
	var user model.User
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()
	selectString := "SELECT * FROM Users WHERE id = $1"

	row := conn.QueryRow(ctx, selectString, id)

	if err := row.Scan(&user.Id, &user.NickName, &user.FirstName, &user.SecondName,
		&user.Latitude, &user.Longitude, &user.Notification); err != nil {
		return user, err
	}
	return user, nil
}

func (db *DBClient) CheckCourier(id int64) (bool, error) {
	var count int
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()
	selectString := "SELECT count(*) FROM UserWithOrder WHERE user_id = $1"

	row := conn.QueryRow(ctx, selectString, id)

	if err := row.Scan(&count); err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (db *DBClient) AddRating(rating model.Rating) error {
	var ratCount int
	var rat float32
	ctx := context.Background()
	selectString := "SELECT rat_count, rating FROM Users WHERE id = $1"
	updateString := "UPDATE User SET rating = $1 WHERE id = $2"

	row := db.connection.QueryRow(ctx, selectString, rating.Id)
	if err := row.Scan(&ratCount, &rat); err != nil {
		return err
	}

	mean := (float32(ratCount)*rat + rating.Rating) / float32(ratCount)
	ratCount++

	_, err := db.connection.Exec(ctx, updateString, mean, rating.Id)
	if err != nil {
		return err
	}

	return nil
}

func (db *DBClient) GetAllCourier() []model.Courier {
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()

	selectString := "SELECT user_id, time_to FROM UserWithOrder"

	rows, err := conn.Query(ctx, selectString)
	if err != nil {
		return []model.Courier{}
	}
	defer rows.Close()

	var couriers []model.Courier
	for rows.Next() {
		courier := model.Courier{User: model.User{}}
		if err := rows.Scan(&courier.User.Id, &courier.TimeTo); err != nil {
			return []model.Courier{}
		}
		couriers = append(couriers, courier)
	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		return []model.Courier{}
	}
	return couriers
}

func (db *DBClient) GetAllUsers() []model.User {
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()

	selectString := "SELECT id FROM Users;"

	rows, err := conn.Query(ctx, selectString)
	if err != nil {
		log.Printf("Error to get users: %v", err)
		return []model.User{}
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		user := model.User{}
		if err := rows.Scan(user); err != nil {
			log.Printf("Error to parse user id: %v", err)
			return []model.User{}
		}

		//courier.TimeFrom = time.Unix(timeFrom, 0)
		//courier.TimeTo = time.Unix(timeTo, 0)
		users = append(users, user)
	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		return []model.User{}
	}
	return users
}

func (db *DBClient) GetAllNotificationUser() []model.User {
	ctx := context.Background()
	conn, _ := db.connection.Acquire(ctx)
	defer conn.Release()

	selectString := "SELECT id FROM Users WHERE notification = True"

	rows, err := conn.Query(ctx, selectString)
	if err != nil {
		log.Printf("Get all notification error: %v", err)
		return []model.User{}
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		user := model.User{}
		if err := rows.Scan(&user.Id); err != nil {
			log.Printf("Parse notification user id error: %v", err)
			return []model.User{}
		}

		//courier.TimeFrom = time.Unix(timeFrom, 0)
		//courier.TimeTo = time.Unix(timeTo, 0)
		users = append(users, user)
	}
	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		log.Printf("Get notification error: %v", err)
		return []model.User{}
	}
	return users
}
