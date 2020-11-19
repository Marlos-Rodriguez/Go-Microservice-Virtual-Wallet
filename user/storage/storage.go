package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
)

//UserStorageService struct
type UserStorageService struct {
	db *gorm.DB
}

type relationChannel struct {
	Err error
	ID  string
}

//NewUserStorageService Create a new storage user service
func NewUserStorageService(newDB *gorm.DB) *UserStorageService {
	return &UserStorageService{db: newDB}
}

//GetUser Get basic user info
func (u *UserStorageService) GetUser(ID string) (*models.UserResponse, error) {
	//Get info from DB
	var userDB *models.User

	u.db.Where("user_id = ?", ID).First(&userDB)

	if err := u.db.Error; err != nil {
		userResponse := &models.UserResponse{}
		return userResponse, err
	}

	//Here have to get the last transactions from the transaction service

	//Assing the info for response
	userResponse := &models.UserResponse{
		UserID:   userDB.UserID,
		UserName: userDB.UserName,
		Balance:  userDB.Balance,
		Avatar:   userDB.Profile.Avatar,
	}

	return userResponse, nil
}

//GetProfileUser Get the profile info
func (u *UserStorageService) GetProfileUser(ID string) (*models.UserProfileResponse, error) {
	//Get info from DB
	var profileDB *models.Profile

	u.db.Where("user_id = ?", ID).First(&profileDB)

	if err := u.db.Error; err != nil {
		profileResponse := &models.UserProfileResponse{}
		return profileResponse, err
	}

	//Assing the info for response
	profileResponse := &models.UserProfileResponse{
		UserID:    profileDB.UserID,
		FirstName: profileDB.FirstName,
		LastName:  profileDB.LastName,
		Email:     profileDB.Email,
		Birthday:  profileDB.Birthday,
		Biography: profileDB.Biography,
		CreatedAt: profileDB.CreatedAt,
		UpdatedAt: profileDB.UpdatedAt,
	}

	return profileResponse, nil
}

//ModifyUser This modify the Complete User, this must not modify the Username or Email
func (u *UserStorageService) ModifyUser(m *models.User) (bool, error) {
	//encrypt Password
	if len(m.Profile.Password) > 0 || m.Profile.Password != "" {
		m.Profile.Password, _ = EncryptPassword(m.Profile.Password)
	}

	//Modify in DB
	if err := u.db.Save(&m).Error; err != nil {
		return false, err
	}

	return true, nil
}

//ModifyUsername Change the username if that not already exits
func (u *UserStorageService) ModifyUsername(ID string, newUsername string) (bool, error) {
	var userDB *models.User
	//Check if exits a record with that username
	if err := u.db.Where("user_name = ?", newUsername).First(&userDB).Error; err != nil {
		//If not exits update the username
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = u.db.Model(&models.User{}).Where("user_id = ?", ID).Update("user_name", newUsername).Error

			if err != nil {
				return false, err
			}

			return true, nil
		}
	}

	//Not Error so record found and username exits
	return false, errors.New("Username already exists")
}

//ModifyEmail Change the username if that not already exits
func (u *UserStorageService) ModifyEmail(ID string, newEmail string) (bool, error) {
	//Check if exits a record with that email
	if err := u.db.Where("user_name = ?", newEmail).First(&models.User{}).Error; err != nil {
		//If not exits update the username
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = u.db.Model(&models.User{}).Where("user_id = ?", ID).Update("email", newEmail).Error

			if err != nil {
				return false, err
			}

			return true, nil
		}
	}

	//Not Error so record found and email exits
	return false, errors.New("Email already exists")
}

//CheckExistingUser Check existing User
func (u *UserStorageService) CheckExistingUser(ID string) (bool, error) {
	if err := u.db.Where("user_id = ?", ID).First(&models.User{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}

	return false, errors.New("User already exists")
}

//GetRelations Get relations from one User
func (u *UserStorageService) GetRelations(ID string, page int) ([]*models.RelationReponse, error) {
	//Get info from DB
	var relationDB []*models.Relation

	limit := page * 20

	u.db.Where("from_user = ?", ID).Or("to_user = ? AND mutual = true", ID).Find(&relationDB).Limit(limit)

	if err := u.db.Error; err != nil {
		var relationResponse []*models.RelationReponse
		return relationResponse, nil
	}

	//Assing the info for response
	var relationResponse []*models.RelationReponse

	for _, relation := range relationDB {
		//new model for append
		loopRelation := &models.RelationReponse{
			RelationID: relation.RelationID,
			FromUser:   relation.FromUser,
			FromName:   relation.FromName,
			ToUser:     relation.ToUser,
			ToName:     relation.ToName,
			CreatedAt:  relation.CreatedAt,
			UpdatedAt:  relation.UpdatedAt,
		}

		relationResponse = append(relationResponse, loopRelation)
	}

	return relationResponse, nil
}

//AddRelation Create a new Relation
func (u *UserStorageService) AddRelation(r *models.RelationRequest, fromID string) (bool, error) {

	exits, err := u.CheckExistingRelation(r.FromName, r.ToName)

	//If there was an error but the relation exits
	if err != nil && exits {
		return false, errors.New("Relation already exits")
	}

	//If the User exits and change the relation to mutual
	if exits == true && err == nil {
		return true, nil
	}

	//Create the channels
	fromChan := make(chan relationChannel, 1)
	toChan := make(chan relationChannel, 1)

	//Check if ID is pass and execute the gorutine
	if len(fromID) < 0 || fromID == "" {
		go u.getID(r.FromName, r.FromEmail, fromChan)
	}

	//Execute the gorutine for toUser
	go u.getID(r.FromName, r.FromEmail, toChan)

	//Create the Users variables
	var fromByteID uuid.UUID
	var toByteID uuid.UUID

	//If the user pass the ID
	if len(fromID) > 0 || fromID != "" {
		fromByteID, err = uuid.Parse(fromID)
		if err != nil {
			return false, err
		}
	}

	//Create the new relation with the other info
	newRelation := &models.Relation{
		RelationID: uuid.New(),
		FromName:   r.FromName,
		ToName:     r.ToName,
		Mutual:     false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsActive:   true,
	}

	//Select for the channels
	select {
	case fromIDDB := <-fromChan:
		//If the fromID is empty and not errors
		if len(fromID) < 0 || fromID != "" || fromIDDB.Err != nil {
			fromByteID, err = uuid.Parse(fromIDDB.ID)
			if err != nil {
				return false, err
			}
		}
	case toID := <-toChan:
		if toID.Err != nil {
			return false, toID.Err
		}
		toByteID, err = uuid.Parse(toID.ID)

		if err != nil {
			return false, err
		}
	}

	//assing the new ID
	newRelation.FromUser = fromByteID
	newRelation.ToUser = toByteID

	//Create relation in DB
	if err := u.db.Create(&newRelation).Error; err != nil {
		return false, err
	}

	return true, nil
}

//CheckExistingRelation Check if exits a relation and if is not mutual, If is not mutual update it
func (u *UserStorageService) CheckExistingRelation(fromUser string, toUser string) (bool, error) {
	//Check values
	if len(fromUser) < 0 || len(toUser) < 0 {
		return false, errors.New("Must send boths variables")
	}

	var relationDB *models.Relation

	if err := u.db.Where("from_user = ? AND to_user = ? AND mutual = false", toUser, fromUser).First(&relationDB).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	relationDB.Mutual = true

	if err := u.db.Save(&relationDB); err != nil {
		return true, err.Error
	}

	return true, nil
}

func (u *UserStorageService) getID(username string, email string, relationChan chan relationChannel) {
	//Get ID
	ID, sucess, err := u.GetUserFromNameEmail(username, email)

	//Create internal Channel
	internalChan := new(relationChannel)

	if sucess != true || err != nil {
		internalChan.Err = err
	}

	internalChan.Err = nil
	internalChan.ID = ID

	relationChan <- *internalChan
}

//GetUserFromNameEmail Get the ID of User from Username and email
func (u *UserStorageService) GetUserFromNameEmail(username string, email string) (string, bool, error) {
	//Check values
	if len(username) < 0 || len(email) < 0 {
		return "", false, errors.New("Must send boths variables")
	}

	//Get info from
	var userDB *models.User
	if err := u.db.Where("user_name = ?", username).First(&userDB).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", false, errors.New("User not found")
		}
		return "", false, errors.New("Error in Get the User")
	}

	if userDB.Profile.Email != email {
		return "", false, errors.New("the username and email not mach")
	}

	return userDB.UserID.String(), true, nil
}

//Here must create a funtion to deactive a relation
