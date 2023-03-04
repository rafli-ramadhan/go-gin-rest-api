package location

import (
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-rest-api/src/constant"
	"go-rest-api/src/http"
	"go-rest-api/src/model"
	"go-rest-api/src/repository/v1/location"
	"gorm.io/gorm"
)

type Service struct {
	repo location.Repositorier
}

func NewService(
	repositorier location.Repositorier,
) *Service {
	return &Service{
		repo: repositorier,
	}
}

type Servicer interface {
	TakeLocationByID(locationID int) (locations http.GetLocation, err error)
	TakeLocationByName(locationName string) (location model.Location, err error)
	Find(locationIDs []int) (locations []http.GetLocation, err error)
	CheckLocationByID(locationID int) (exist bool, err error)
	CheckLocationByName(locationName string) (exist bool, err error)
	Create(request http.CreateLocation) (err error)
	Update(locationID int, request http.UpdateLocation) (err error)
	Delete(locationID int) (err error)
}

func (svc *Service) TakeLocationByID(locationID int) (location http.GetLocation, err error) {
	takeLocation, err := svc.repo.TakeLocationByID(locationID)
	if err == gorm.ErrRecordNotFound {
		err = constant.ErrLocationNotExist
		return
	} else if err != nil {
		err = errors.Wrap(err, "take location")
		return
	}

	location = http.GetLocation{}
	copier.Copy(&location, &takeLocation)
	location.ID = int(location.ID)
	return
}

func (svc *Service) TakeLocationByName(locationName string) (location model.Location, err error) {
	location, err = svc.repo.TakeLocationByName(locationName)
	if err == gorm.ErrRecordNotFound {
		err = constant.ErrLocationNotExist
		return
	} else if err != nil {
		err = errors.Wrap(err, "take location by location name")
		return
	}
	return
}

func (svc *Service) Find(locationIDs []int) (locations []http.GetLocation, err error) {
	locationsData, err := svc.repo.Find(locationIDs)
	if err != nil {
		err = errors.Wrap(err, "find locations")
		return
	}
	for i := range locationsData {
		location := http.GetLocation{}
		copier.Copy(&location, &locationsData[i])
		location.ID = int(location.ID)
		locations = append(locations, location)
	} 
	return
}

func (svc *Service) CheckLocationByID(locationID int) (exist bool, err error) {
	exist = false
	_, err = svc.repo.TakeLocationByID(locationID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			exist = false
			err = nil
			return
		} else {
			err = errors.Wrap(err, "check location by email")
			return
		}
	}
	exist = true
	return
}

func (svc *Service) CheckLocationByName(locationName string) (exist bool, err error) {
	exist = false
	_, err = svc.repo.TakeLocationByName(locationName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			exist = false
			err = nil
			return
		} else {
			err = errors.Wrap(err, "check location by locationName")
			return
		}
	}
	exist = true
	return
}

func (svc *Service) Create(request http.CreateLocation) (err error) {
	if request.LocationName == "" {
		err = constant.ErrInvalidLocationName
		return
	}

	if request.Address == "" {
		err = constant.ErrInvalidAddress
		return
	}

	exist, err := svc.CheckLocationByName(request.LocationName)
	if err != nil {
		return
	}
	if exist {
		err = constant.ErrLocationAlreadyExist
		return
	} else {
		newLocation := model.Location{}
		copier.Copy(&newLocation, &request)
		if newLocation.PhotoURL == "" {
			newLocation.PhotoURL = "https://th.bing.com/th/id/OIP.gBRzG71aa1f6dy_MuGUwOAHaEo?pid=ImgDet&rs=1"
		}

		err = svc.repo.Create(newLocation)
		if err != nil {
			err = errors.Wrap(err, "create new location")
			return err
		}
	}
	return
}

func (svc *Service) Update(locationID int, request http.UpdateLocation) (err error) {
	exist, _ := svc.CheckLocationByID(locationID)
	if !exist {
		err = constant.ErrLocationNotExist
		return
	}

	if request.LocationName != "" {
	    locationName, _ := svc.CheckLocationByName(request.LocationName)
	    if locationName {
		    err = constant.ErrLocationNameAlreadyExist
		    return
	    }
	}

	location := model.Location{}
	copier.Copy(&location, &request)

	err = svc.repo.Update(locationID, location)
	if err != nil {
		err = errors.Wrap(err, "update location")
		return
	}
	return
}

func (svc *Service) Delete(locationID int) (err error) {
	_, err = svc.TakeLocationByID(locationID)
	if err != nil {
		err = errors.Wrap(err, "location is not exist")
		return
	}

	err = svc.repo.Delete(locationID)
	if err != nil {
		err = errors.Wrap(err, "delete location")
		return
	}
	return
}
