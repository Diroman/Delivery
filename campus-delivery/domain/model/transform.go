package model

import (
	"log"
	"time"

	"campus-delivery/api"
)

func RegistrationRequestToUser(request api.RegistrationRequest) User {
	return User{
		Id:           request.Id,
		NickName:     request.Nickname,
		FirstName:    request.FirstName,
		SecondName:   request.SecondName,
		Latitude:     request.Location.Latitude,
		Longitude:    request.Location.Longitude,
		Notification: request.Notification,
		Rating:       0,
	}
}

func AddCourierRequestToCourier(request api.AddCourierRequest) Courier {
	return Courier{
		User:        User{Id: request.Id},
		Shop:        request.ShopName,
		Description: request.Description,
		TimeFrom:    time.Unix(request.TimeRange.TimeFrom, 0),
		TimeTo:      time.Unix(request.TimeRange.TimeTo, 0),
		Link:        request.Link,
		ChatId:      request.ChatId,
	}
}

func RatingRequestToRating(request api.RatingRequest) Rating {
	return Rating{
		Id:     request.Id,
		Rating: request.Rating,
	}
}

func NewCourierToNewCourierResponse(message ChanMessage) api.NewCourier {
	courier := message.Courier
	log.Printf("time: %v, unix time: %v <%T>", courier.TimeFrom, message.Courier.TimeFrom.Unix(), message.Courier.TimeFrom.Unix())
	return api.NewCourier{
		Ids: message.Ids,
		Courier: &api.Courier{
			Id:         courier.User.Id,
			FirstName:  courier.User.FirstName,
			SecondName: courier.User.SecondName,
			Rating:     0,
			TimeRange: &api.TimeRange{
				TimeFrom: message.Courier.TimeFrom.Unix(),
				TimeTo:   message.Courier.TimeTo.Unix(),
			},
			ShopName:    courier.Shop,
			Description: courier.Description,
			Link:        courier.Link,
			ChatId:      courier.ChatId,
		},
	}
}

func CourierToCourierResponse(couriers []Courier) api.GetCourierResponse {
	respCouriers := make([]*api.Courier, len(couriers))

	for i, courier := range couriers {
		respCourier := api.Courier{
			Id:         courier.User.Id,
			FirstName:  courier.User.FirstName,
			SecondName: courier.User.SecondName,
			Rating:     courier.User.Rating,
			TimeRange: &api.TimeRange{
				TimeFrom: courier.TimeFrom.Unix(),
				TimeTo:   courier.TimeTo.Unix(),
			},
			ShopName:    courier.Shop,
			Link:        courier.Link,
			ChatId:      courier.ChatId,
			Description: courier.Description,
		}
		respCouriers[i] = &respCourier
		log.Printf("time: %v, unix time: %v <%T>", courier.TimeFrom, respCouriers[i].TimeRange.TimeFrom, respCouriers[i].TimeRange.TimeFrom)
	}
	return api.GetCourierResponse{Courier: respCouriers}
}

func UserToRegistrationRequest(user User) api.RegistrationRequest {
	return api.RegistrationRequest{
		Id:         user.Id,
		Nickname:   user.NickName,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Location: &api.Location{
			Latitude:  user.Latitude,
			Longitude: user.Longitude,
		},
		Notification: user.Notification,
	}
}

func RequestToUserRequest(request api.UserRequest) UserRequest {
	return UserRequest{
		Id:       request.Id,
		TimeFrom: time.Unix(request.TimeRange.TimeFrom, 0),
		TimeTo:   time.Unix(request.TimeRange.TimeTo, 0),
	}
}

func CodeToCodeResponse(code Code) api.CodeResponse {
	return api.CodeResponse{
		Code:  code.Code,
		Error: code.Error,
	}
}

func UserToUsersResponse(users []User) api.Users {
	ids := []int64{}
	for _, user := range users {
		ids = append(ids, user.Id)
	}
	return api.Users{Ids: ids}
}
