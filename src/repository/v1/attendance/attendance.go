package attendance

import (
	"go-rest-api/src/connection"
	"go-rest-api/src/model"
	"go-rest-api/src/pkg/pagination"

	"gorm.io/gorm"
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
	Find(accountID int, pgn pagination.Pagination) (attendaceDatas []model.Attendance, err error)
	Create(accountID int, account model.Attendance) (err error)
}

func (repo *Repository) Find(accountID int, pgn pagination.Pagination) (attendaceDatas []model.Attendance, err error) {
	query := repo.dbMaster.Model(&model.Attendance{}).
		Where("account_id", accountID).
		Order("created_at desc").
		Limit(pgn.Limit).
		Offset(pgn.Offset).
		Find(&attendaceDatas)
	err = query.Error
	return
}

func (repo *Repository) Create(accountID int, attendance model.Attendance) (err error) {
	query := repo.dbMaster.Model(&attendance).Begin().
		Create(&attendance)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}

	err = query.Commit().Error
	return
}
