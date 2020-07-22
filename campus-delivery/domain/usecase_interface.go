package domain

import "campus-delivery/domain/model"

type UseCase interface {
	Registration(model.User) model.Code
	DeleteUser(int64) model.Code
	AddUserWithOrder(model.Courier) model.ChanMessage
	DeleteUserWithOrder(int64) model.Code
	GetCourier(int64) []model.Courier
	CheckCourier(int64) model.Code
	CheckLocation(int64, float64, float64) model.Code
	ChangeNotificationStatus(int64, bool) model.Code
	AddRating(model.Rating) model.Code
	GetUserInfo(int64) model.User
	GetAllUser() []model.User
	GetAllNotificationUser() []model.User
}
