package account

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
	TakeAccountByID(accountID int) (account model.Account, err error)
	TakeAccountByEmail(email string) (account model.Account, err error)
	TakeAccountByKTPNumber(ktpNumber string) (account model.Account, err error)
	TakeAccountByPhoneNumber(phoneNumber string) (account model.Account, err error)
	TakeAccountByUsername(username string) (account model.Account, err error)
	Find(accountIDs []int) (accounts []model.Account, err error)
	Create(account model.Account) (err error)
	Update(accountID int, request model.Account) (err error)
	Delete(accountID int) (err error)
}

func (repo *Repository) TakeAccountByID(accountID int) (account model.Account, err error) {
	query := repo.dbMaster.Model(&model.Account{}).
		Where("id", accountID).
		Take(&account)
	err = query.Error
	return
}

func (repo *Repository) TakeAccountByEmail(email string) (account model.Account, err error) {
	query := repo.dbMaster.Model(&model.Account{}).
		Where("email", email).
		Take(&account)
	err = query.Error
	return
}

func (repo *Repository) TakeAccountByKTPNumber(ktpNumber string) (account model.Account, err error) {
	query := repo.dbMaster.Model(&model.Account{}).
		Where("ktp_number", ktpNumber).
		Take(&account)
	err = query.Error
	return
}

func (repo *Repository) TakeAccountByPhoneNumber(phoneNumber string) (account model.Account, err error) {
	query := repo.dbMaster.Model(&model.Account{}).
		Where("phone_number", phoneNumber).
		Take(&account)
	err = query.Error
	return
}

func (repo *Repository) TakeAccountByUsername(username string) (account model.Account, err error) {
	query := repo.dbMaster.Model(&model.Account{}).
		Where("username", username).
		Take(&account)
	err = query.Error
	return
}

func (repo *Repository) Find(accountIDs []int) (accounts []model.Account, err error) {
	query := repo.dbMaster.Model(&model.Account{}).
		Find(&accounts, accountIDs)
	err = query.Error
	return
}

func (repo *Repository) Create(account model.Account) (err error) {
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

func (repo *Repository) Update(accountID int, request model.Account) (err error) {
	account := &model.Account{}
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
	account := &model.Account{}
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
