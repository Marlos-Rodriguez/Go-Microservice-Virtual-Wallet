package storage

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"

	UserModels "github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
)

//CheckExistingUser Check if the user Exits for Username or Email
func (s *AuthStorageService) CheckExistingUser(username string, email string) (bool, error) {
	//Check using cache
	exits, err := s.CheckExistingUserCache(username, email)

	if err != nil {
		log.Println("Error in get the Cache " + err.Error())
	}

	if exits == true {
		return true, nil
	}

	//Check in DB
	var UserDB *UserModels.User = new(UserModels.User)
	var ProfileDB *UserModels.Profile = new(UserModels.Profile)

	//Check Username
	if err := s.db.Where(&UserModels.User{UserName: username}).First(&UserDB).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//Check Email
			if err = s.db.Where(&UserModels.Profile{Email: email}).First(&ProfileDB).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return false, nil
				}

				return false, err
			}

			return true, nil
		}

		return false, err
	}

	return true, nil
}
