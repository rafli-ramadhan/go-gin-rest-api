package account

import (
	"go-rest-api/src/connection"
	"go-rest-api/src/constant"
	entity "go-rest-api/src/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"go-rest-api/src/http"
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
	TakeAccountByID(accountID int) (account entity.User, err error)
	TakeAccountByEmail(email string) (account entity.User, err error)
	TakeAccountByKTPNumber(ktpNumber string) (account entity.User, err error)
	TakeAccountByPhoneNumber(phoneNumber string) (account entity.User, err error)
	TakeAccountByUsername(username string) (account entity.User, err error)
	Find(accountIDs []int) (accounts []entity.User, err error)
	Create(account entity.User) (err error)
	Update(accountID int, request http.UpdateUser) (err error)
	Delete(accountID int) (err error)
}

func (repo *Repository) TakeAccountByID(accountID int) (account entity.User, err error) {
	query := repo.dbMaster.Model(&entity.User{}).
		Where("id", accountID).
		Take(&account)
	err = query.Error
	return
}

func (repo *Repository) TakeAccountByEmail(email string) (account entity.User, err error) {
	query := repo.dbMaster.Model(&entity.User{}).
		Where("email", email).
		Take(&account)
	err = query.Error
	return
}

func (repo *Repository) TakeAccountByKTPNumber(ktpNumber string) (account entity.User, err error) {
	query := repo.dbMaster.Model(&entity.User{}).
		Where("ktp_number", ktpNumber).
		Take(&account)
	err = query.Error
	return
}

func (repo *Repository) TakeAccountByPhoneNumber(phoneNumber string) (account entity.User, err error) {
	query := repo.dbMaster.Model(&entity.User{}).
		Where("phone_number", phoneNumber).
		Take(&account)
	err = query.Error
	return
}

func (repo *Repository) TakeAccountByUsername(username string) (account entity.User, err error) {
	query := repo.dbMaster.Model(&entity.User{}).
		Where("username", username).
		Take(&account)
	err = query.Error
	return
}

func (repo *Repository) Find(accountIDs []int) (accounts []entity.User, err error) {
	query := repo.dbMaster.Model(&entity.User{}).
		Find(&accounts, accountIDs)
	err = query.Error
	return
}

func (repo *Repository) Create(account entity.User) (err error) {
	query := repo.dbMaster.Model(&account ).Begin().
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "username"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"username": account.Username,
				"full_name": account.FullName,
				"password": account.Password,
				"deleted_at": nil,
			})}).
		Create(&account)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}

	err = query.Commit().Error
	return
}

func (repo *Repository) Update(accountID int, request http.UpdateUser) (err error) {
	account := &entity.User{}
	query := repo.dbMaster.Model(&account ).Begin().
		Where("id", accountID).
		Updates(request)
	err = query.Error
	if err != nil {
		query.Rollback()
		return
	}

	err = query.Commit().Error
	return
}


func (repo *Repository) Delete(accountID int) (err error) {
	account := &entity.User{}
	query := repo.dbMaster.Model(account).Begin().
		Where("id", accountID).
		Delete(account )
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
