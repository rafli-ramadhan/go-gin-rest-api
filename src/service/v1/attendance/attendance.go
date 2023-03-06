package attendance

import (
	"fmt"
	"time"

	"go-rest-api/src/constant"
	"go-rest-api/src/http"
	"go-rest-api/src/model"
	"go-rest-api/src/pkg/pagination"
	"go-rest-api/src/repository/v1/attendance"
	"go-rest-api/src/service/v1/account"
	"go-rest-api/src/service/v1/location"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
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
	FindAttendanceHistory(accountID int, pgn pagination.Pagination, filter string) (responses []http.GetAttendance, err error)
	FindByLocation(accountID int, pgn pagination.Pagination) (responses []http.GetAttendanceByLocation, err error)
	Add(accountID int, request http.AddAttendance) (err error)
}

func (svc *Service) FindAttendanceHistory(accountID int, pgn pagination.Pagination, filter string) (responses []http.GetAttendance, err error) {
	accountExist, err := svc.account.CheckAccountByID(accountID)
	if err != nil {
		err = errors.Wrap(err, "check account by id")
		return
	}
	if !accountExist {
		err = constant.ErrAccountNotRegistered
		return
	}

	attendanceDatas, err := svc.repo.Find(accountID, pgn)
	if err != nil {
		err = errors.Wrap(err, "find attendance datas")
		return
	}

	if len(attendanceDatas) != 0 {
		now := time.Now()
		for i := range attendanceDatas {
			if filter == constant.FilterByDay {
				if attendanceDatas[i].CreatedAt.Year() != now.Year() || attendanceDatas[i].CreatedAt.Month() != now.Month() || attendanceDatas[i].CreatedAt.Day() != now.Day() {
					continue
				}
			} else if filter == constant.FilterByWeek {
				if attendanceDatas[i].CreatedAt.Year() != now.Year() || attendanceDatas[i].CreatedAt.Month() != now.Month() || now.Day() - attendanceDatas[i].CreatedAt.Day() > 7 {
					continue
				}
			} else if filter == constant.FilterByMonth {
				if attendanceDatas[i].CreatedAt.Year() != now.Year() || attendanceDatas[i].CreatedAt.Month() != now.Month() {
					continue
				}
			}

			attendance := http.GetAttendance{}
			location, err := svc.location.TakeLocationByID(attendanceDatas[i].LocationID)
			if err != nil {
				err = errors.Wrap(err, "check location by id")
				return nil, err
			}
			attendance.LocationID = attendanceDatas[i].LocationID
			attendance.LocationName = location.LocationName
			attendance.Status = attendanceDatas[i].Status
			attendance.Time = attendanceDatas[i].CreatedAt.Format("03:04 PM")
			if attendanceDatas[i].Status == constant.StatusCheckIn {
				attendance.Description = fmt.Sprintf("Check In - %s - %v", location.LocationName, attendanceDatas[i].CreatedAt.Format("03:04 PM"))
			} else {
				attendance.Description = fmt.Sprintf("Check Out - %s - %v", location.LocationName, attendanceDatas[i].CreatedAt.Format("03:04 PM"))
			}
			responses = append(responses, attendance)
		}
	}
	return
}

func (svc *Service) FindByLocation(accountID int, pgn pagination.Pagination) (responses []http.GetAttendanceByLocation, err error) {
	accountExist, err := svc.account.CheckAccountByID(accountID)
	if err != nil {
		err = errors.Wrap(err, "check account by id")
		return
	}
	if !accountExist {
		err = constant.ErrAccountNotRegistered
		return
	}

	attendanceDatas, err := svc.repo.Find(accountID, pgn)
	if err != nil {
		err = errors.Wrap(err, "find attendance datas")
		return
	}

	var locationNameCheck = make(map[int]bool)
	if len(attendanceDatas) != 0 {
		for i := range attendanceDatas {
			if locationNameCheck[attendanceDatas[i].LocationID] {
				continue
			} else {
				locationNameCheck[attendanceDatas[i].LocationID] = true
			}

			attendance := http.GetAttendanceByLocation{}
			location, err := svc.location.TakeLocationByID(attendanceDatas[i].LocationID)
			if err != nil {
				err = errors.Wrap(err, "check location by id")
				return nil, err
			}
			attendance.LocationID = attendanceDatas[i].LocationID
			attendance.LocationName = location.LocationName
			attendance.Address = location.Address
			attendance.Status = attendanceDatas[i].Status
			responses = append(responses, attendance)
		}
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
