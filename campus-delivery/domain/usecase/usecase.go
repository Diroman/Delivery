package usecase

import (
	"campus-delivery/domain"
	"campus-delivery/domain/model"
	"fmt"
	"log"
)

type UseCase struct {
	database domain.DataBase
}

func NewUseCase(database domain.DataBase) domain.UseCase {
	return &UseCase{
		database: database,
	}
}

func (uc UseCase) Registration(user model.User) model.Code {
	if err := uc.database.AddUser(user); err != nil {
		errString := fmt.Sprintf("User <%v> not add: %v", user.Id, err)
		log.Print(errString, err)

		return model.Code{
			Code:  errString,
			Error: true,
		}
	}
	return model.Code{
		Code:  "",
		Error: false,
	}
}

func (uc UseCase) DeleteUser(id int64) model.Code {
	if err := uc.database.DeleteUser(id); err != nil {
		errString := fmt.Sprintf("Can`t delete user <%v>: %v", id, err)
		log.Print(errString, err)

		return model.Code{
			Code:  errString,
			Error: true,
		}
	}
	return model.Code{
		Code:  "",
		Error: false,
	}
}

func (uc UseCase) AddUserWithOrder(courier model.Courier) model.ChanMessage {
	if err := uc.database.AddUserWithOrder(&courier); err != nil {
		errString := fmt.Sprintf("Error to add <%v> in <UserWithOrder>", courier.User.Id)
		log.Print(errString, err)

		return model.ChanMessage{
			Ids:     []int64{},
			Courier: model.Courier{},
		}
	}
	user, err := uc.database.GetUserInfo(courier.User.Id)
	if err != nil {
		errString := fmt.Sprintf("Error to get info about <%v> in <Users>", courier.User.Id)
		log.Print(errString, err)

		return model.ChanMessage{
			Ids:     []int64{},
			Courier: model.Courier{},
		}
	}
	courier.User.FirstName = user.FirstName
	courier.User.SecondName = user.SecondName

	couriersNear := uc.database.GetNotificationUser(user)

	ids := []int64{}
	for _, cour := range couriersNear {
		ids = append(ids, cour.Id)
	}

	return model.ChanMessage{
		Ids:     ids,
		Courier: courier,
	}
}

func (uc UseCase) DeleteUserWithOrder(id int64) model.Code {
	if err := uc.database.DeleteUserWithOrder(id); err != nil {
		errString := fmt.Sprintf("Error to get info about <%v> in <UserWithOrder>", id)
		log.Print(errString, err)

		return model.Code{
			Code:  errString,
			Error: true,
		}
	}

	return model.Code{
		Code:  "",
		Error: false,
	}
}

func (uc UseCase) GetCourier(id int64) []model.Courier {
	user, err := uc.database.GetUserInfo(id)
	if err != nil {
		return []model.Courier{}
	}

	couriersNear, err := uc.database.GetUserWithOrder(user)
	if err != nil {
		errString := fmt.Sprintf("Error to get near user for <%v> in <UserWithOrder>", id)
		log.Print(errString, err)
		return []model.Courier{}
	}
	return couriersNear
}

func (uc UseCase) CheckCourier(id int64) model.Code {
	isCourier, err := uc.database.CheckCourier(id)
	if err != nil {
		return model.Code{
			Code:  "Error",
			Error: false,
		}
	}

	return model.Code{
		Code:  "",
		Error: isCourier,
	}
}

func (uc UseCase) ChangeNotificationStatus(id int64, status bool) model.Code {
	if err := uc.database.ChangeNotificationStatus(id, status); err != nil {
		log.Print(err)
		return model.Code{
			Code:  "",
			Error: true,
		}
	}

	return model.Code{
		Code:  "",
		Error: false,
	}
}

func (uc UseCase) CheckLocation(id int64, lat, long float64) model.Code {
	if err := uc.database.ChangeLocation(id, lat, long); err != nil {
		log.Print(err)
		return model.Code{
			Code:  "",
			Error: true,
		}
	}

	return model.Code{
		Code:  "",
		Error: false,
	}
}

func (uc UseCase) AddRating(rating model.Rating) model.Code {
	if err := uc.database.AddRating(rating); err != nil {
		errString := fmt.Sprintf("Error to update rating for user <%v>: %v", rating.Id, err)
		log.Print(errString, err)

		return model.Code{
			Code:  errString,
			Error: true,
		}
	}

	return model.Code{
		Code:  "",
		Error: false,
	}
}

func (uc UseCase) GetUserInfo(id int64) model.User {
	user, err := uc.database.GetUserInfo(id)
	if err != nil {
		return model.User{Id: -1}
	}
	return user
}

func (uc UseCase) GetAllUser() []model.User {
	users := uc.database.GetAllUsers()
	return users
}

func (uc UseCase) GetAllNotificationUser() []model.User {
	users := uc.database.GetAllNotificationUser()
	return users
}
