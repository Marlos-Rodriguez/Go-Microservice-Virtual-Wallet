package storage

import (
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/internal/utils"
	grpc "github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/grpc/client"
	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
)

//UserStorageService struct
type UserStorageService struct {
	db  *gorm.DB
	rdb *redis.Client
}

//NewUserStorageService Create a new storage user service
func NewUserStorageService(newDB *gorm.DB, newRDB *redis.Client) UserStorageService {
	go grpc.StartClient()
	go startTransactionClient()

	newDB.AutoMigrate(&models.User{}, &models.Profile{}, &models.Relation{})

	newService := UserStorageService{db: newDB, rdb: newRDB}

	return newService
}

//CloseDB Close DB for GRPC
func (s *UserStorageService) CloseDB() {
	s.db.Close()
	s.rdb.Close()
	grpc.CloseClient()
	closeTransactionClient()
}

//GetUser Get basic user info
func (s *UserStorageService) GetUser(ID string) (*models.UserResponse, error) {
	//Get info from Cache
	cacheResponse, err := s.GetUserCache(ID)

	if cacheResponse != nil && err == nil {
		return cacheResponse, nil
	}

	//Get info from DB
	var userDB *models.User = new(models.User)

	s.db.Where("user_id = ?", ID).First(&userDB)

	if err := s.db.Error; err != nil {
		return nil, err
	}

	//Here have to get the last transactions from the transaction service
	transactions, success, err := GetTransactions(ID)

	if !success || err != nil {
		return nil, err
	}

	if !userDB.IsActive {
		return nil, errors.New("User is not active")
	}

	//Assing the info for response
	userResponse := &models.UserResponse{
		UserID:           userDB.UserID,
		UserName:         userDB.UserName,
		Balance:          userDB.Balance,
		Avatar:           userDB.Profile.Avatar,
		LastTransactions: transactions,
	}

	//Set in Redis for Cache
	if err := s.SetUser(userDB); err != nil {
		return userResponse, err
	}

	return userResponse, nil
}

//GetProfileUser Get the profile info
func (s *UserStorageService) GetProfileUser(ID string) (*models.UserProfileResponse, error) {
	//Get info from Cache
	cacheProfile, err := s.GetProfileCache(ID)

	if cacheProfile != nil && err == nil {
		return cacheProfile, nil
	}

	//Get from DB
	var profileDB *models.Profile = new(models.Profile)

	s.db.Where("user_id = ?", ID).First(&profileDB)

	if err := s.db.Error; err != nil {
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

	//Set in Redis Cache
	if err := s.SetProfile(profileDB); err != nil {
		return profileResponse, err
	}

	return profileResponse, nil
}

//ModifyUser This modify the Complete User, this must not modify the Username or Email
func (s *UserStorageService) ModifyUser(m *models.User, ID string, newUsername string) (bool, error) {
	var change string

	//encrypt Password
	if len(m.Profile.Password) > 0 || m.Profile.Password != "" {
		m.Profile.Password, _ = utils.EncryptPassword(m.Profile.Password)

		change += "User change password & "
	}

	if newUsername != "" || len(newUsername) > 0 {
		if sucess, err := s.ModifyUsername(ID, m.UserName, newUsername); err != nil || sucess == false {
			return false, err
		}
		m.UserName = ""
	}

	if m.Profile.Email != "" || len(m.Profile.Email) > 0 {
		if sucess, err := s.ModifyEmail(ID, m.Profile.Email); err != nil || sucess == false {
			return false, err
		}
		m.Profile.Email = ""
	}

	m.UpdatedAt = time.Now()

	//Modify User in DB
	go s.db.Model(&models.User{}).Where("user_id = ?", ID).Update(&m)

	if err := s.db.Error; err != nil {
		return false, err
	}

	//Modify in Profile DB
	if err := s.db.Model(&models.Profile{}).Where("user_id = ?", ID).Update(&m.Profile).Error; err != nil {
		return false, err
	}

	//Make movement
	change += "Modify user " + m.UserID.String()

	succes, err := grpc.CreateMovement("User & Profile", change, "User Service")

	if err != nil {
		log.Println("Error in Create a movement: " + err.Error())
	}

	if succes == false {
		log.Println("Error in Create a movement")
	}

	s.UpdateUserCache(ID)

	return true, nil
}

//ModifyUsername Change the username if that not already exits
func (s *UserStorageService) ModifyUsername(ID string, currentUsername string, newUsername string) (bool, error) {
	//Check if exits a record with that username
	if err := s.db.Where("user_name = ?", newUsername).First(&models.User{}).Error; err != nil {
		//If not exits update the username
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//change username
			UserChange := &models.User{UserName: newUsername, UpdatedAt: time.Now()}

			err = s.db.Model(&models.User{}).Where("user_id = ?", ID).Updates(&UserChange).Error

			if err != nil {
				return false, err
			}

			go s.UpdateUserCache(ID)

			success, err := grpc.UpdateAuthCache(currentUsername, newUsername, "", "")

			if err != nil || success == false {
				log.Println("Error in Update the Auth Cache " + err.Error())
			}

			//Movement of change of Username
			var change string = "Modify UserName from " + currentUsername + " to " + newUsername

			success, err = grpc.CreateMovement("User", change, "User Service")

			if err != nil {
				log.Println("Error in Create a movement: " + err.Error())
			}

			if success == false {
				log.Println("Error in Create a movement")
			}

			//Modify username in from relations
			fromRelationChange := &models.Relation{FromName: newUsername, UpdatedAt: time.Now()}

			err = s.db.Model(&models.Relation{}).Where("from_name = ?", currentUsername).Updates(&fromRelationChange).Error

			if err != nil {
				return false, err
			}

			//Modify Username in to Relations
			toRelationChange := &models.Relation{ToName: newUsername, UpdatedAt: time.Now()}

			err = s.db.Model(&models.Relation{}).Where("to_name = ?", currentUsername).Updates(&toRelationChange).Error

			if err != nil {
				return false, err
			}

			go s.UpdateRelations(ID)

			success, err = grpc.CreateMovement("Relations", "Modify UserName in relations: "+currentUsername, "User Service")

			if err != nil {
				log.Println("Error in Create a movement: " + err.Error())
			}

			if success == false {
				log.Println("Error in Create a movement")
			}

			return true, nil
		}
	}

	//Not Error so record found and username exits
	return false, errors.New("Username already exists")
}

//ModifyEmail Change the username if that not already exits
func (s *UserStorageService) ModifyEmail(ID string, newEmail string) (bool, error) {
	//Check if exits a record with that email
	if err := s.db.Where("email = ?", newEmail).First(&models.Profile{}).Error; err != nil {
		//If not exits update the username
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//Update DB
			UserChange := &models.Profile{Email: newEmail, UpdatedAt: time.Now()}

			err = s.db.Model(&models.Profile{}).Where("user_id = ?", ID).Updates(&UserChange).Error

			if err != nil {
				return false, err
			}

			profile, _ := s.GetProfileCache(ID)

			if profile != nil {
				success, err := grpc.UpdateAuthCache("", "", profile.Email, newEmail)
				if success == false || err != nil {
					log.Println("Error in Set the Auth cache " + err.Error())
				}
				profile = nil
			} else {
				profileDB, err := s.GetProfileUser(ID)

				success, err := grpc.UpdateAuthCache("", "", profileDB.Email, newEmail)

				if success == false || err != nil {
					log.Println("Error in Set the Auth cache " + err.Error())
				}
				profileDB = nil
			}

			//Set in Cache
			go s.UpdateUserCache(ID)

			//Create Movement
			succes, err := grpc.CreateMovement("Profile", "Modify Email", "User Service")

			if err != nil {
				log.Println("Error in Create a movement: " + err.Error())
			}

			if succes == false {
				log.Println("Error in Create a movement")
			}

			return true, nil
		}
	}

	//Not Error so record found and email exits
	return false, errors.New("Email already exists")
}

//GetRelations Get relations from one User
func (s *UserStorageService) GetRelations(ID string, page int) ([]*models.RelationReponse, error) {

	if page > 1 {
		relationsCache, err := s.GetRelationsCache(ID)
		if err != nil {
			log.Println("Error in get the Cache " + err.Error())
		}

		if relationsCache != nil {
			return relationsCache, nil
		}
	}

	//Get info from DB
	var relationDB []*models.Relation = []*models.Relation{new(models.Relation)}

	limit := page * 30

	s.db.Where("from_user = ?", ID).Or("to_user = ? AND mutual = true", ID).Find(&relationDB).Limit(limit)

	if err := s.db.Error; err != nil {
		return nil, nil
	}

	//Set in Cache
	if page > 1 {
		s.SetRelationCache(relationDB, ID)
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
func (s *UserStorageService) AddRelation(r *models.RelationRequest) (bool, error) {

	//Check if exits a relation but is not mutual
	exits, err := s.CheckMutualRelation(r.FromName, r.ToName, r.FromID)

	//If there was an error but the relation exits
	if err != nil && exits {
		return false, errors.New("Relation already exits")
	}

	//If the User exits and change the relation to mutual
	if exits == true && err == nil {
		return true, nil
	}

	//Check if exits the relation in DB
	exits, err = s.CheckExistingRelation(r.FromName, r.ToName, true)

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

	toID, err = s.GetIDName(r.ToName, r.ToEmail)

	//If there was an error
	if err != nil {
		return false, err
	}

	if len(toID) < 0 || toID == "" {
		return false, errors.New("Error in Get the ID of To user")
	}

	//Check if exits another the new user
	var isActive bool

	exits, isActive, err = s.CheckExistingUser(toID)

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
	if err := s.db.Create(&newRelation).Error; err != nil {
		return false, err
	}

	//Update Cache
	go s.UpdateRelations(newRelation.FromUser.String())

	//Create Movement
	var change string = "Create a new Relation between " + newRelation.FromName + " & " + newRelation.ToName

	succes, err := grpc.CreateMovement("Relations", change, "User Service")

	if err != nil {
		log.Println("Error in Create a movement: " + err.Error())
	}

	if succes == false {
		log.Println("Error in Create a movement")
	}

	return true, nil
}

//DeactivateRelation Deactive the relation in DB
func (s *UserStorageService) DeactivateRelation(relationID string, FromID string, ToID string) (bool, error) {
	//Check values
	if len(FromID) < 0 || len(ToID) < 0 {
		return false, errors.New("Must send boths variables")
	}

	var relationDB *models.Relation = new(models.Relation)

	s.db.Where("relation_id = ? AND from_user = ? AND to_user = ?", relationID, FromID, ToID).
		Or("relation_id = ? AND from_user = ? AND to_user = ? AND mutual = true", relationID, ToID, FromID).First(&relationDB)

	if s.db.Error != nil {
		return false, s.db.Error
	}

	relationDB.IsActive = false
	relationDB.UpdatedAt = time.Now()

	if err := s.db.Save(&relationDB).Error; err != nil {
		return false, err
	}

	go s.UpdateRelations(relationDB.FromUser.String())
	go s.UpdateRelations(relationDB.ToUser.String())

	succes, err := grpc.CreateMovement("Relations", "Deactive Relation: "+relationDB.RelationID.String(), "User Service")

	if err != nil {
		log.Println("Error in Create a movement: " + err.Error())
	}

	if succes == false {
		log.Println("Error in Create a movement")
	}

	return true, nil
}
