package account

import (
	"log"

	"github.com/forkyid/go-utils/v1/aes"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-rest-api/src/constant"
	"go-rest-api/src/http"
	"go-rest-api/src/model"
	"go-rest-api/src/pkg/bcrypt"
	"go-rest-api/src/repository/v1/account"
	"gorm.io/gorm"
)

type Service struct {
	repo account.Repositorier
}

func NewService(
	repositorier account.Repositorier,
) *Service {
	return &Service{
		repo: repositorier,
	}
}

type Servicer interface {
	TakeAccountByID(accountID int) (accounts http.GetUser, err error)
	TakeAccountByUsername(username string) (account model.User, err error)
	Find(accountIDs []int) (accounts []http.GetUser, err error)
	CheckAccountByID(accountID int) (exist bool, err error)
	CheckAccountByEmail(email string) (exist bool, err error)
	CheckAccountByKTPNumber(ktpNumber string) (exist bool, err error)
	CheckAccountByPhoneNumber(phoneNumber string) (exist bool, err error)
	CheckAccountByUsername(username string) (exist bool, err error)
	Create(request http.RegisterUser) (err error)
	Update(accountID int, request http.UpdateUser) (err error)
	Delete(accountID int) (err error)
}

func (svc *Service) TakeAccountByID(accountID int) (account http.GetUser, err error) {
	takeUser, err := svc.repo.TakeAccountByID(accountID)
	if err == gorm.ErrRecordNotFound {
		err = constant.ErrAccountNotRegistered
		return
	} else if err != nil {
		err = errors.Wrap(err, "take account")
		return
	}

	account = http.GetUser{}
	copier.Copy(&account, &takeUser)
	account.ID = aes.Encrypt(int(takeUser.ID))
	return
}

func (svc *Service) TakeAccountByUsername(username string) (account model.User, err error) {
	account, err = svc.repo.TakeAccountByUsername(username)
	if err == gorm.ErrRecordNotFound {
		err = constant.ErrAccountNotRegistered
		return
	} else if err != nil {
		err = errors.Wrap(err, "take account")
		return
	}
	return
}

func (svc *Service) Find(accountIDs []int) (accounts []http.GetUser, err error) {
	users, err := svc.repo.Find(accountIDs)
	if err != nil {
		err = errors.Wrap(err, "find accounts")
		return
	}
	for i := range users {
		account := http.GetUser{}
		copier.Copy(&account, &users[i])
		account.ID = aes.Encrypt(int(users[i].ID))
		log.Print(users[i].PhotoURL)
		accounts = append(accounts, account)
	} 
	return
}

func (svc *Service) CheckAccountByID(accountID int) (exist bool, err error) {
	exist = false
	_, err = svc.repo.TakeAccountByID(accountID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			exist = false
			err = nil
			return
		} else {
			err = errors.Wrap(err, "check account by email")
			return
		}
	}
	exist = true
	return
}

func (svc *Service) CheckAccountByEmail(email string) (exist bool, err error) {
	exist = false
	_, err = svc.repo.TakeAccountByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			exist = false
			err = nil
			return
		} else {
			err = errors.Wrap(err, "check account by email")
			return
		}
	}
	exist = true
	return
}

func (svc *Service) CheckAccountByKTPNumber(ktpNumber string) (exist bool, err error) {
	exist = false
	_, err = svc.repo.TakeAccountByKTPNumber(ktpNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			exist = false
			err = nil
			return
		} else {
			err = errors.Wrap(err, "check account by ktp number")
			return
		}
	}
	exist = true
	return
}

func (svc *Service) CheckAccountByPhoneNumber(phoneNumber string) (exist bool, err error) {
	exist = false
	_, err = svc.repo.TakeAccountByPhoneNumber(phoneNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			exist = false
			err = nil
			return
		} else {
			err = errors.Wrap(err, "check account by phone number")
			return
		}
	}
	exist = true
	return
}

func (svc *Service) CheckAccountByUsername(username string) (exist bool, err error) {
	exist = false
	_, err = svc.repo.TakeAccountByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			exist = false
			err = nil
			return
		} else {
			err = errors.Wrap(err, "check account by username")
			return
		}
	}
	exist = true
	return
}

func (svc *Service) Create(request http.RegisterUser) (err error) {
	exist, err := svc.CheckAccountByUsername(request.Username)
	if err != nil {
		return
	}

	if exist {
		err = constant.ErrAccountExist
		return
	} else {
		newAccount := model.User{}
		copier.Copy(&newAccount, &request)

		hashedPassword, err := bcrypt.GeneratePasswordHarsh(newAccount.Password)
		if err != nil {
			err = errors.Wrap(err, "hash password")
			return err
		}
		newAccount.Password = hashedPassword
		newAccount.PhotoURL = "https://thumbs.dreamstime.com/b/user-profile-avatar-solid-black-line-icon-simple-vector-filled-flat-pictogram-isolated-white-background-134042540.jpg"
		newAccount.Gender = "none"
		newAccount.IsVerified = false

		err = svc.repo.Create(newAccount)
		if err != nil {
			err = errors.Wrap(err, "create new account")
			return err
		}
	}
	return
}

func (svc *Service) Update(accountID int, request http.UpdateUser) (err error) {
	exist, _ := svc.CheckAccountByID(accountID)
	if !exist {
		err = constant.ErrAccountNotRegistered
		return
	}

	if request.Username != "" {
	    emailExist, _ := svc.CheckAccountByUsername(request.Username)
	    if emailExist {
		    err = constant.ErrUsernameAlreadyExist
		    return
	    }
	}

	if request.Email != "" {
	    emailExist, _ := svc.CheckAccountByEmail(request.Email)
	    if emailExist {
		    err = constant.ErrEmailAlreadyExist
		    return
	    }
	}

	if request.KTPNumber != "" {
	    emailExist, _ := svc.CheckAccountByKTPNumber(request.KTPNumber)
	    if emailExist {
		    err = constant.ErrKTPNumberAlreadyExist
		    return
	    }
	}

	if request.PhoneNumber != "" {
	    emailExist, _ := svc.CheckAccountByPhoneNumber(request.PhoneNumber)
	    if emailExist {
		    err = constant.ErrPhoneNumberAlreadyExist
		    return
	    }
	}

	err = svc.repo.Update(accountID, request)
	if err != nil {
		err = errors.Wrap(err, "update account")
		return
	}
	return
}

func (svc *Service) Delete(accountID int) (err error) {
	_, err = svc.TakeAccountByID(accountID)
	if err != nil {
		err = errors.Wrap(err, "account is not exist")
		return
	}

	err = svc.repo.Delete(accountID)
	if err != nil {
		err = errors.Wrap(err, "delete account")
		return
	}
	return
}
