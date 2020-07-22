package domain

import "campus-delivery/domain/model"

type DataBase interface {
	Connect() error
	CloseConnection()
	AddUser(model.User) error
	DeleteUser(int64) error
	DeleteUserByTimer()
	AddUserWithOrder(*model.Courier) error
	DeleteUserWithOrder(int64) error
	GetUserWithOrder(model.User) ([]model.Courier, error)
	GetUserInfo(int64) (model.User, error)
	CheckCourier(int64) (bool, error)
	ChangeNotificationStatus(int64, bool) error
	ChangeLocation(int64, float64, float64) error
	GetNotificationUser(model.User) []model.User
	GetAllNotificationUser() []model.User
	AddRating(model.Rating) error
	GetAllCourier() []model.Courier
	GetAllUsers() []model.User
	//SubscribeUser(model.UserRequest) error
	//DeleteSubscriberUser(int64) error
}
