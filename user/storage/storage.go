package storage

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	grpc "github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/grpc/client"
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
	newDB.AutoMigrate(&models.User{}, &models.Profile{}, &models.Relation{})

	return &UserStorageService{db: newDB}
}

//CloseDB Close DB
func (u *UserStorageService) CloseDB() {
	u.db.Close()
}

//GetUser Get basic user info
func (u *UserStorageService) GetUser(ID string) (*models.UserResponse, error) {

	//Get info from DB
	var userDB *models.User = new(models.User)

	if u.db == nil {
		log.Println("DB is nil")
	}

	u.db.Where("user_id = ?", &ID).First(&userDB)

	if err := u.db.Error; err != nil {
		return nil, err
	}

	//Here have to get the last transactions from the transaction service

	//Assing the info for response
	userResponse := &models.UserResponse{
		UserID:   userDB.UserID,
		UserName: userDB.UserName,
		Balance:  userDB.Balance,
		Avatar:   userDB.Profile.Avatar,
	}

	var change string = "Get info of: " + userDB.UserName

	succes, err := grpc.CreateMovement("User & Transactions", change, "User Service")

	if err != nil {
		return nil, err
	}

	if succes == false {
		log.Fatalln("Error in Create a movement")
	}

	return userResponse, nil
}

//GetProfileUser Get the profile info
func (u *UserStorageService) GetProfileUser(ID string) (*models.UserProfileResponse, error) {
	//Get info from DB
	var profileDB *models.Profile = new(models.Profile)

	u.db.Where("user_id = ?", ID).First(&profileDB)

	if err := u.db.Error; err != nil {
		return nil, err
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

	var change string = "Get info Profile of: " + profileDB.UserID.String()

	//Create the movement in DB
	succes, err := grpc.CreateMovement("User & Profile", change, "User Service")

	if err != nil {
		return nil, err
	}

	if succes == false {
		log.Fatalln("Error in Create a movement")
	}

	return profileResponse, nil
}

//ModifyUser This modify the Complete User, this must not modify the Username or Email
func (u *UserStorageService) ModifyUser(m *models.User, ID string, newUsername string) (bool, error) {
	var change string

	//encrypt Password
	if len(m.Profile.Password) > 0 || m.Profile.Password != "" {
		m.Profile.Password, _ = EncryptPassword(m.Profile.Password)

		change += "User change password "
	}

	if newUsername != "" || len(newUsername) > 0 {
		if sucess, err := u.ModifyUsername(ID, m.UserName, newUsername); err != nil || sucess == false {
			return false, err
		}
		m.UserName = ""
	}

	if m.Profile.Email != "" || len(m.Profile.Email) > 0 {
		if sucess, err := u.ModifyEmail(ID, m.Profile.Email); err != nil || sucess == false {
			return false, err
		}
		m.Profile.Email = ""
	}

	//Modify User in DB
	if err := u.db.Model(&models.User{}).Where("user_id = ?", ID).Update(&m).Error; err != nil {
		return false, err
	}

	//Modify in Profile DB
	if err := u.db.Model(&models.Profile{}).Where("user_id = ?", ID).Update(&m.Profile).Error; err != nil {
		return false, err
	}

	change += "& Modify user " + m.UserID.String()

	succes, err := grpc.CreateMovement("User & Profile", change, "User Service")

	if err != nil {
		return false, err
	}

	if succes == false {
		log.Fatalln("Error in Create a movement")
	}

	return true, nil
}

//ModifyUsername Change the username if that not already exits
func (u *UserStorageService) ModifyUsername(ID string, currentUsername string, newUsername string) (bool, error) {
	var userDB *models.User = new(models.User)
	//Check if exits a record with that username
	if err := u.db.Where("user_name = ?", newUsername).First(&userDB).Error; err != nil {
		//If not exits update the username
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//change username
			err = u.db.Model(&models.User{}).Where("user_id = ?", ID).Update("user_name", newUsername).Error

			if err != nil {
				return false, err
			}

			var change string = "Modify UserName from " + currentUsername + " to " + newUsername

			succes, err := grpc.CreateMovement("User", change, "User Service")

			if err != nil {
				return false, err
			}

			if succes == false {
				log.Fatalln("Error in Create a movement")
			}

			err = u.db.Model(&models.Relation{}).Where("from_name = ?", currentUsername).Update("from_name", newUsername).Error

			if err != nil {
				return false, err
			}

			err = u.db.Model(&models.Relation{}).Where("to_name = ?", currentUsername).Update("to_name", newUsername).Error

			if err != nil {
				return false, err
			}

			succes, err = grpc.CreateMovement("Relations", "Modify UserName in relations: "+currentUsername, "User Service")

			if err != nil {
				return false, err
			}

			if succes == false {
				log.Fatalln("Error in Create a movement")
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
	if err := u.db.Where("email = ?", newEmail).First(&models.Profile{}).Error; err != nil {
		//If not exits update the username
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = u.db.Model(&models.Profile{}).Where("user_id = ?", ID).Update("email", newEmail).Error

			if err != nil {
				return false, err
			}

			succes, err := grpc.CreateMovement("Profile", "Modify Email", "User Service")

			if err != nil {
				return false, err
			}

			if succes == false {
				log.Fatalln("Error in Create a movement")
			}

			return true, nil
		}
	}

	//Not Error so record found and email exits
	return false, errors.New("Email already exists")
}

//CheckExistingUser Check existing User
func (u *UserStorageService) CheckExistingUser(ID string) (bool, bool, error) {

	var userDB *models.User = new(models.User)

	var (
		exits    bool
		isActive bool
	)

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

//GetRelations Get relations from one User
func (u *UserStorageService) GetRelations(ID string, page int) ([]*models.RelationReponse, error) {
	//Get info from DB
	var relationDB []*models.Relation = []*models.Relation{new(models.Relation)}

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
func (u *UserStorageService) AddRelation(r *models.RelationRequest) (bool, error) {

	//Check if exits a relation but is not mutual
	exits, err := u.CheckMutualRelation(r.FromName, r.ToName)

	//If there was an error but the relation exits
	if err != nil && exits {
		return false, errors.New("Relation already exits")
	}

	//If the User exits and change the relation to mutual
	if exits == true && err == nil {
		return true, nil
	}

	//Check if exits the relation in DB
	exits, err = u.CheckExistingRelation(r.FromName, r.ToName, true)

	//If there was an error
	if err != nil {
		if errors.Is(err, errors.New("The relation was reactived")) {
			return true, nil
		}
		return false, err
	}

	//If the User exits
	if exits == true && err == nil {
		return false, errors.New("Relations already exits")
	}

	//Get the other user ID
	var toID string

	toID, err = u.GetIDName(r.ToName, r.ToEmail)

	//If there was an error
	if err != nil {
		return false, err
	}

	if len(toID) < 0 || toID == "" {
		return false, errors.New("Error in Get the ID of To user")
	}

	//Check if exits another the new user
	var isActive bool

	exits, isActive, err = u.CheckExistingUser(toID)

	if err != nil {
		return false, err
	}

	if !exits || !isActive {
		return false, errors.New("User no exits or is not active")
	}

	//Convert the ID

	fromID, err := uuid.Parse(r.FromID)
	newtoID, err := uuid.Parse(toID)

	if err != nil {
		return false, errors.New("Error converting the ID in DB")
	}

	//Create the new relation with the other info
	newRelation := &models.Relation{
		RelationID: uuid.New(),
		FromUser:   fromID,
		FromName:   r.FromName,
		ToUser:     newtoID,
		ToName:     r.ToName,
		Mutual:     false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsActive:   true,
	}

	//Create relation in DB
	if err := u.db.Create(&newRelation).Error; err != nil {
		return false, err
	}

	var change string = "Create a new Relation between " + newRelation.FromName + " & " + newRelation.ToName

	succes, err := grpc.CreateMovement("Relations", change, "User Service")

	if err != nil {
		return false, err
	}

	if succes == false {
		log.Fatalln("Error in Create a movement")
	}

	return true, nil
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
		err = u.db.Model(&models.Relation{}).Where(&models.Relation{RelationID: relationDB.RelationID, IsActive: false}).Update("is_active", true).Error
		if err != nil {
			return false, err
		}

		return true, errors.New("The relation was reactived")
	}

	return true, nil
}

//CheckMutualRelation Check if exits a relation and if is not mutual, If is not mutual update it
func (u *UserStorageService) CheckMutualRelation(fromUser string, toUser string) (bool, error) {
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

	succes, err := grpc.CreateMovement("Relations", "Update relation to mutual", "User Service")

	if err != nil {
		return false, err
	}

	if succes == false {
		log.Fatalln("Error in Create a movement")
	}

	return true, nil
}

//DeactivateRelation Deactive the relation in DB
func (u *UserStorageService) DeactivateRelation(FromID string, ToID string) (bool, error) {
	//Check values
	if len(FromID) < 0 || len(ToID) < 0 {
		return false, errors.New("Must send boths variables")
	}

	var relationDB *models.Relation = new(models.Relation)

	u.db.Where("from_user = ? AND to_user = ?", FromID, ToID).Or("from_user = ? AND to_user = ? AND mutual = true", ToID, FromID).First(&relationDB)

	if u.db.Error != nil {
		return false, u.db.Error
	}

	relationDB.IsActive = false

	if err := u.db.Save(&relationDB).Error; err != nil {
		return false, err
	}

	succes, err := grpc.CreateMovement("Relations", "Deactive Relation: "+relationDB.RelationID.String(), "User Service")

	if err != nil {
		return false, err
	}

	if succes == false {
		log.Fatalln("Error in Create a movement")
	}

	return true, nil
}
