package controller

import (
	"campus-delivery/api"
	"campus-delivery/clientApi"
	"campus-delivery/domain"
	"campus-delivery/domain/model"
	"log"
)

var n = 0

type Controller struct {
	usecase domain.UseCase
}

func NewController(usecase domain.UseCase) clientApi.Controller {
	return &Controller{
		usecase: usecase,
	}
}

func (ctrl Controller) Registration(request api.RegistrationRequest) api.CodeResponse {
	m := n
	n++

	user := model.RegistrationRequestToUser(request)
	log.Printf("(%v) Get registration request: <%v>", m, user)
	code := ctrl.usecase.Registration(user)
	log.Printf("(%v) Registration response: <%v>", m, code)
	return model.CodeToCodeResponse(code)
}

func (ctrl Controller) DeleteUser(request api.DeleteRequest) api.CodeResponse {
	m := n
	n++

	log.Printf("(%v) Get delete user request: <%v>", m, request.Id)
	code := ctrl.usecase.DeleteUser(request.Id)
	log.Printf("(%v) Delete user response: <%v>", m, code)
	return model.CodeToCodeResponse(code)
}

func (ctrl Controller) AddUserWithOrder(request api.AddCourierRequest) api.NewCourier {
	m := n
	n++

	courier := model.AddCourierRequestToCourier(request)
	log.Printf("(%v) Get add courier request: <%v>", m, courier)
	couriers := ctrl.usecase.AddUserWithOrder(courier)
	log.Printf("(%v) Add courier response: <%v>", m, couriers)
	return model.NewCourierToNewCourierResponse(couriers)
}

func (ctrl Controller) DeleteUserWithOrder(request api.DeleteRequest) api.CodeResponse {
	m := n
	n++

	log.Printf("(%v) Get delete courier request: <%v>", m, request.Id)
	code := ctrl.usecase.DeleteUserWithOrder(request.Id)
	log.Printf("(%v) Delete courier response: <%v>", m, code)
	return model.CodeToCodeResponse(code)
}

func (ctrl Controller) GetCourier(request api.UserRequest) api.GetCourierResponse {
	m := n
	n++

	log.Printf("(%v) Get courier request: <%v>", m, request.Id)
	couriers := ctrl.usecase.GetCourier(request.Id)
	log.Printf("(%v) Couiriers response: <%v>", m, couriers)
	return model.CourierToCourierResponse(couriers)
}

func (ctrl Controller) CheckCourier(request api.UserRequest) api.CodeResponse {
	m := n
	n++

	log.Printf("(%v) Check courier request: <%v>", m, request.Id)
	code := ctrl.usecase.CheckCourier(request.Id)
	log.Printf("(%v) Check courier response: <%v>", m, code)
	return model.CodeToCodeResponse(code)
}

func (ctrl Controller) ChangeNotificationStatus(request api.RegistrationRequest) api.CodeResponse {
	m := n
	n++

	log.Printf("(%v) Change notification status request: <%v>", m, request)
	code := ctrl.usecase.ChangeNotificationStatus(request.Id, request.Notification)
	log.Printf("(%v) Change notification status response: <%v>", m, code)
	return model.CodeToCodeResponse(code)
}

func (ctrl Controller) ChangeLocation(request api.RegistrationRequest) api.CodeResponse {
	m := n
	n++

	log.Printf("(%v) Change location request: <%v>", m, request)
	code := ctrl.usecase.CheckLocation(request.Id, float64(request.Location.Latitude), float64(request.Location.Longitude))
	log.Printf("(%v) Change location response: <%v>", m, code)
	return model.CodeToCodeResponse(code)
}

func (ctrl Controller) GetUserInfo(request api.UserRequest) api.RegistrationRequest {
	m := n
	n++

	log.Printf("(%v) Get user info request: <%v>", m, request)
	user := ctrl.usecase.GetUserInfo(request.Id)
	log.Printf("(%v) Get user info response: <%v>", m, user)
	return model.UserToRegistrationRequest(user)
}

func (ctrl Controller) GetAllUser(request api.UserRequest) api.Users {
	m := n
	n++

	log.Printf("(%v) Get all user request: <%v>", m, request)
	users := ctrl.usecase.GetAllUser()
	log.Printf("(%v) Get all user response: <%v>", m, users)
	return model.UserToUsersResponse(users)
}

func (ctrl Controller) GetAllNotificationUser(request api.UserRequest) api.Users {
	m := n
	n++

	log.Printf("(%v) Get all notification user request: <%v>", m, request)
	users := ctrl.usecase.GetAllNotificationUser()
	log.Printf("(%v) Get all notification user response: <%v>", m, users)
	return model.UserToUsersResponse(users)
}

func (ctrl Controller) ListenCourier() {
	//for {
	//	nc := <-ctrl.usecase.(*usecase.UseCase).NewCourier
	//	courierResp := model.NewCourierToNewCourierResponse(nc)
	//	log.Printf("Send courier: <%v>", courierResp)
	//	go func(){
	//		ctrl.NewCourier <- courierResp
	//	}()
	//}
}

func (ctrl Controller) AddRating(request api.RatingRequest) api.CodeResponse {
	m := n
	n++

	rating := model.RatingRequestToRating(request)
	log.Printf("(%v) Get add rating request: <%v>", m, rating)
	code := ctrl.usecase.AddRating(rating)
	log.Printf("(%v) Rating response: <%v>", m, code)
	return model.CodeToCodeResponse(code)
}
