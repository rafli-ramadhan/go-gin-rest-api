package location

import (
	"go-rest-api/src/connection"
	"go-rest-api/src/constant"
	"go-rest-api/src/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)
 
type DB struct {
	Master *gorm.DB
}

type Repository struct {
	dbMaster *gorm.DB
}

func NewRepository(
	db connection.DB,
) *Repository {
	return &Repository{
		dbMaster: db.Master,
	}
}

type Repositorier interface {
	TakeLocationByID(locationID int) (location model.Location, err error)
	TakeLocationByName(name string) (location model.Location, err error)
	Find(locationIDs []int) (locations []model.Location, err error)
	Create(location model.Location) (err error)
	Update(locationID int, request model.Location) (err error)
	Delete(locationID int) (err error)
}

func (repo *Repository) TakeLocationByID(locationID int) (location model.Location, err error) {
	query := repo.dbMaster.Model(&model.Location{}).
		Where("id", locationID).
		Take(&location)
	err = query.Error
	return
}

func (repo *Repository) TakeLocationByName(name string) (location model.Location, err error) {
	query := repo.dbMaster.Model(&model.Location{}).
		Where("name", name).
		Take(&location)
	err = query.Error
	return
}

func (repo *Repository) Find(locationIDs []int) (locations []model.Location, err error) {
	query := repo.dbMaster.Model(&model.Location{}).
		Find(&locations, locationIDs)
	err = query.Error
	return
}

func (repo *Repository) Create(location model.Location) (err error) {
	query := repo.dbMaster.Model(&location).Begin().
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"address": location.Address,
				"deleted_at": nil,
			})}).
		Create(&location)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}

	err = query.Commit().Error
	return
}

func (repo *Repository) Update(locationID int, request model.Location) (err error) {
	location := &model.Location{}
	query := repo.dbMaster.Model(&location).Begin().
		Where("id", locationID).
		Updates(request)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}

	err = query.Commit().Error
	return
}

func (repo *Repository) Delete(locationID int) (err error) {
	location := &model.Location{}
	query := repo.dbMaster.Model(location).Begin().
		Where("id", locationID).
		Delete(location)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}
	if query.RowsAffected != 1 {
		query.Rollback()
		err = constant.ErrInvalidID
		return
	}

	err = query.Commit().Error
	return
}
