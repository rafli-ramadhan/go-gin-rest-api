package attendance

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-rest-api/src/constant"
	"go-rest-api/src/http"
	"go-rest-api/src/model"
	"go-rest-api/src/repository/v1/attendance"
	"go-rest-api/src/service/v1/account"
	"go-rest-api/src/service/v1/location"
)

type Service struct {
	repo     attendance.Repositorier
	account  account.Servicer
	location location.Servicer
}

func NewService(
	repositorier attendance.Repositorier,
	accountSvc   account.Servicer,
	locationSvc  location.Servicer,
) *Service {
	return &Service{
		repo:     repositorier,
		account:  accountSvc,
		location: locationSvc,
	}
}

type Servicer interface {
	Find(accountID int) (responses []http.GetAttendance, err error)
	Add(accountID int, request http.AddAttendance) (err error)
}

func (svc *Service) Find(accountID int) (responses []http.GetAttendance, err error) {
	accountExist, err := svc.account.CheckAccountByID(accountID)
	if err != nil {
		err = errors.Wrap(err, "check account by id")
		return
	}
	if !accountExist {
		err = constant.ErrAccountNotRegistered
		return
	}

	attendaceDatas, err := svc.repo.Find(accountID)
	if err != nil {
		err = errors.Wrap(err, "find attendance datas")
		return
	}

	for i := range attendaceDatas {
		attendance := http.GetAttendance{}
		copier.Copy(&attendance, &attendaceDatas[i])
		responses = append(responses, attendance)
	}
	return
}

func (svc *Service) Add(accountID int, request http.AddAttendance) (err error) {
	accountExist, err := svc.account.CheckAccountByID(accountID)
	if err != nil {
		err = errors.Wrap(err, "check account by id")
		return
	}
	if !accountExist {
		err = constant.ErrAccountNotRegistered
		return
	}

	locationExist, err := svc.location.CheckLocationByID(request.LocationID)
	if err != nil {
		err = errors.Wrap(err, "check location by id")
		return
	}
	if !locationExist {
		err = constant.ErrLocationNotExist
		return
	}

	newAttendance := model.Attendance{}
	copier.Copy(&newAttendance, &request)
	newAttendance.AccountID = accountID
	newAttendance.CreatedAt = time.Now().UTC()
	newAttendance.UpdatedAt = time.Now().UTC()

	err = svc.repo.Create(accountID, newAttendance)
	if err != nil {
		err = errors.Wrap(err, "create new attendance")
		return err
	}
	return
}
