package storage

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/models"
)

//CheckExistingUser Check if the user Exits for Username or Email
func (s *AuthStorageService) CheckExistingUser(username string, email string) (string, bool, error) {
	//Check using cache
	ID, exits, err := s.CheckExistingUserCache(username, email)

	if err != nil {
		log.Println("Error in get the Cache " + err.Error())
	}

	if exits == true {
		return ID, true, nil
	}

	//Check in DB
	var UserDB *models.User = new(models.User)
	var ProfileDB *models.Profile = new(models.Profile)

	//Check Username
	if err := s.db.Where(&models.User{UserName: username}).First(&UserDB).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//Check Email
			if err = s.db.Where(&models.Profile{Email: email}).First(&ProfileDB).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return "", false, nil
				}

				return "", false, err
			}

			if ProfileDB.IsActive {
				return ProfileDB.UserID.String(), true, nil
			}
			return "", true, errors.New("User is not active")
		}

		return "", false, err
	}

	if UserDB.IsActive {
		return UserDB.UserID.String(), true, nil
	}
	return "", true, errors.New("User is not active")
}
