package storage

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/auth/models"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

//CheckExistingUserCache Check if the user Exits for Username or Email in Cache
func (s *AuthStorageService) CheckExistingUserCache(username string, email string) (string, bool, error) {
	if username == "" || len(username) <= 0 || email == "" || len(email) <= 0 {
		return "", false, errors.New("Review your Input")
	}

	val := s.rdb.Get(context.Background(), "RegisterUsername:"+username)

	if val.Err() != redis.Nil {
		return val.Val(), true, nil
	}

	val = s.rdb.Get(context.Background(), "RegisterEmail:"+email)

	if val.Err() != redis.Nil {
		return val.Val(), true, nil
	}

	return "", false, val.Err()
}

//SetRegisterCache Set in cache the user info
func (s *AuthStorageService) SetRegisterCache(username string, email string, user *models.User) {
	if username == "" || len(username) <= 0 || email == "" || len(email) <= 0 {
		log.Println("Review your Input")
		return
	}

	if err := s.rdb.Set(context.Background(), "RegisterUsername:"+username, user.UserID.String(), 0); err.Err() != nil {
		log.Println("Error in Auth Cache " + err.Err().Error())
	}
	if err := s.rdb.Set(context.Background(), "RegisterEmail:"+email, user.UserID.String(), 0); err.Err() != nil {
		log.Println("Error in Auth Cache " + err.Err().Error())
	}

	//Here must use User Service for set the user in cache
	go s.SetUser(user)
	s.SetProfile(&user.Profile)
}

//SetProfile set the profile in redis cache
func (s *AuthStorageService) SetProfile(profileDB *models.Profile) {
	redisUser, err := json.Marshal(profileDB)

	if err != nil {
		log.Println("Error in Marshal the user" + err.Error())
	}
	status := s.rdb.Set(context.Background(), "Profile:"+profileDB.UserID.String(), redisUser, time.Hour*72)

	if status.Err() != nil {
		log.Println("Error in set in the cache " + status.Err().Error())
	}
}

//SetUser Set the User in redis Cache
func (s *AuthStorageService) SetUser(userDB *models.User) {
	redisUser, err := json.Marshal(userDB)

	if err != nil {
		log.Println("Error in Marshal the user" + err.Error())
	}
	status := s.rdb.Set(context.Background(), "User:"+userDB.UserID.String(), redisUser, time.Hour*72)

	if status.Err() != nil {
		log.Println("Error in set in the cache " + status.Err().Error())
	}
}

//ChangeRegisterCache Change in cache the username and email
func (s *AuthStorageService) ChangeRegisterCache(oldUsername string, newUsername string, oldEmail string, newEmail string) (bool, error) {
	if len(oldUsername) > 0 && len(newUsername) > 0 {
		status := s.rdb.Rename(context.Background(), "RegisterUsername:"+oldUsername, "RegisterUsername:"+newUsername)
		if status.Err() != nil {
			return false, status.Err()
		}
		return true, nil
	}

	if len(oldEmail) > 0 && len(newEmail) > 0 {
		status := s.rdb.Rename(context.Background(), "RegisterEmail:"+oldEmail, "RegisterEmail:"+newEmail)
		if status.Err() != nil {
			return false, status.Err()
		}
		return true, nil
	}

	return false, errors.New("Invalid Input")
}

//GetProfileCache Get the profile info if exits
func (s *AuthStorageService) GetProfileCache(ID string) (*models.Profile, error) {
	//Get info from redis
	val := s.rdb.Get(context.Background(), "Profile:"+ID)

	err := val.Err()

	if err != nil && err != redis.Nil {
		log.Println("Error in get the cache " + err.Error())
	}

	//Convert for response
	var profileDB *models.Profile = new(models.Profile)

	if err != redis.Nil {
		userBytes, _ := val.Bytes()
		json.Unmarshal(userBytes, &profileDB)

		return profileDB, nil
	}

	return nil, err
}

//DeleteUserCache Update User & Profile in Cache
func (s *AuthStorageService) DeleteUserCache(ID string) {
	//Get info from DB
	var userDB *models.User = new(models.User)
	var profileDB *models.Profile = new(models.Profile)

	go s.db.Where("user_id = ?", &ID).First(&userDB)
	s.db.Where("user_id = ?", ID).First(&profileDB)

	err := s.db.Error

	if err != nil {
		log.Fatalln("Error in get the info from DB for cache " + err.Error())
	}

	go s.rdb.Del(context.Background(), "RegisterUsername:"+userDB.UserID.String())
	go s.rdb.Del(context.Background(), "RegisterEmail:"+userDB.UserID.String())

	//Convert to save
	redisUser, err := json.Marshal(userDB)

	if err != nil {
		log.Println("Error in Marshal the user" + err.Error())
	}

	redisProfile, err := json.Marshal(profileDB)

	if err != nil {
		log.Println("Error in Marshal the user" + err.Error())
	}

	//Save in redis
	status := s.rdb.Set(context.Background(), "User:"+userDB.UserID.String(), redisUser, time.Hour*72)

	if status.Err() != nil {
		log.Println("Error in set in the cache " + status.Err().Error())
	}

	status = s.rdb.Set(context.Background(), "Profile:"+profileDB.UserID.String(), redisProfile, time.Hour*72)

	if status.Err() != nil {
		log.Println("Error in set in the cache " + status.Err().Error())
	}
}
