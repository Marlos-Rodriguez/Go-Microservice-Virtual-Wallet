package storage

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"

	grpc "github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/grpc/client"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
)

//CheckExistingUser Check existing User
func (u *UserStorageService) CheckExistingUser(ID string) (bool, bool, error) {

	var userDB *models.User = new(models.User)

	var (
		exits    bool
		isActive bool
	)

	user, _ := u.GetUserCache(ID)

	if user != nil {
		exits = true
		isActive = true

		return exits, isActive, nil
	}

	if err := u.db.Where("user_id = ?", ID).First(userDB).Error; err != nil {
		exits = false
		isActive = false

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exits, isActive, errors.New("User Not exists")
		}
		return exits, isActive, err
	}

	exits = true

	if !userDB.IsActive {
		isActive = false
		return exits, isActive, nil
	}

	isActive = true

	return exits, isActive, nil
}

//GetIDName Get the ID from the username
func (u *UserStorageService) GetIDName(username string, email string) (string, error) {
	var userDB *models.User = new(models.User)

	if err := u.db.Where("username = ?", username).First(&userDB).Error; err != nil {
		return "", err
	}

	if len(userDB.UserID.String()) > 0 || userDB.UserID.String() != "" {
		return userDB.UserID.String(), nil
	}

	var profileDB *models.Profile = new(models.Profile)

	if err := u.db.Where("email = ?", email).First(&profileDB).Error; err != nil {
		return "", err
	}

	return profileDB.UserID.String(), nil
}

//CheckExistingRelation Check if exits any relations before create
func (u *UserStorageService) CheckExistingRelation(fromUser string, toUser string, active bool) (bool, error) {
	//Check values
	if len(fromUser) < 0 || len(toUser) < 0 {
		return false, errors.New("Must send boths variables")
	}

	var relationDB *models.Relation = new(models.Relation)

	err := u.db.Where(&models.Relation{FromName: fromUser, ToName: toUser}).Or(&models.Relation{FromName: toUser, ToName: fromUser, Mutual: true}).First(&relationDB).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	//If pass the variable and the relation is not active, reactive the relation
	if active == true && relationDB.IsActive == false {
		//Update the DB
		err = u.db.Model(&models.Relation{}).Where(&models.Relation{RelationID: relationDB.RelationID, IsActive: false}).Update("is_active", true).Error
		if err != nil {
			return false, err
		}

		//Update the Cache
		go u.UpdateRelations(relationDB.FromUser.String())
		go u.UpdateRelations(relationDB.ToUser.String())

		//Create the movement
		succes, err := grpc.CreateMovement("Relations", "Update relation to mutual", "User Service")

		if err != nil {
			log.Println("Error in Create a movement: " + err.Error())
		}

		if succes == false {
			log.Println("Error in Create a movement")
		}

		return true, errors.New("The relation was reactived")
	}

	return true, nil
}

//CheckMutualRelation Check if exits a relation and if is not mutual, If is not mutual update it
func (u *UserStorageService) CheckMutualRelation(fromUser string, fromID string, toUser string) (bool, error) {
	//Check values
	if len(fromUser) < 0 || len(toUser) < 0 {
		return false, errors.New("Must send boths variables")
	}

	//If the relations already exits with other user updated to mutual
	err := u.db.Model(&models.Relation{}).Where(&models.Relation{FromName: toUser, ToName: fromUser, Mutual: false}).Update("mutual", true).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return true, err
	}

	go u.UpdateRelations(fromID)

	succes, err := grpc.CreateMovement("Relations", "Update relation to mutual", "User Service")

	if err != nil {
		log.Println("Error in Create a movement: " + err.Error())
	}

	if succes == false {
		log.Println("Error in Create a movement")
	}

	toID, err := u.GetIDName(toUser, "")

	if err != nil {
		log.Println("Error in get the ID in cache " + err.Error())
	}

	if toID != "" {
		u.UpdateRelations(toID)
	}

	return true, nil
}
