package clientApi

import (
	"campus-delivery/api"
)

type Controller interface {
	Registration(api.RegistrationRequest) api.CodeResponse
	DeleteUser(api.DeleteRequest) api.CodeResponse
	AddUserWithOrder(api.AddCourierRequest) api.NewCourier
	DeleteUserWithOrder(api.DeleteRequest) api.CodeResponse
	GetCourier(api.UserRequest) api.GetCourierResponse
	CheckCourier(api.UserRequest) api.CodeResponse
	ChangeNotificationStatus(api.RegistrationRequest) api.CodeResponse
	ChangeLocation(api.RegistrationRequest) api.CodeResponse
	GetUserInfo(api.UserRequest) api.RegistrationRequest
	ListenCourier()
	AddRating(api.RatingRequest) api.CodeResponse
	GetAllUser(api.UserRequest) api.Users
	GetAllNotificationUser(api.UserRequest) api.Users
}
