package storage

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/net/context"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
)

//SetUser Set the User in redis Cache
func (s *UserStorageService) SetUser(userDB *models.User) error {
	redisUser, err := json.Marshal(userDB)

	if err != nil {
		return err
	}
	status := s.rdb.Set(context.Background(), "User:"+userDB.UserID.String(), redisUser, time.Hour*72)

	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

//GetUserCache Get info from redis if exits
func (s *UserStorageService) GetUserCache(ID string) (*models.UserResponse, error) {
	//Get from Redis
	val := s.rdb.Get(context.Background(), "User:"+ID)

	err := val.Err()

	if err != nil && err != redis.Nil {
		return nil, err
	}

	//Convert to response
	var userDB *models.User = new(models.User)

	if err != redis.Nil {
		userBytes, _ := val.Bytes()
		json.Unmarshal(userBytes, &userDB)

		if userDB.IsActive == false {
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

		return userResponse, nil
	}

	return nil, err
}

//GetProfileCache Get the profile info if exits
func (s *UserStorageService) GetProfileCache(ID string) (*models.UserProfileResponse, error) {
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

	return nil, err
}

//SetProfile set the profile in redis cache
func (s *UserStorageService) SetProfile(profileDB *models.Profile) error {
	redisUser, err := json.Marshal(profileDB)

	if err != nil {
		return err
	}
	status := s.rdb.Set(context.Background(), "Profile:"+profileDB.UserID.String(), redisUser, time.Hour*72)

	if status.Err() != nil {
		return err
	}

	return err
}

//UpdateUserCache Update User & Profile in Cache
func (s *UserStorageService) UpdateUserCache(ID string) {
	//Get info from DB
	var userDB *models.User = new(models.User)
	var profileDB *models.Profile = new(models.Profile)

	go s.db.Where("user_id = ?", &ID).First(&userDB)
	s.db.Where("user_id = ?", ID).First(&profileDB)

	err := s.db.Error

	if err != nil {
		log.Fatalln("Error in get the info from DB for cache " + err.Error())
	}

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

//SetRelationCache Set One page of 20 relations
func (s *UserStorageService) SetRelationCache(relations []*models.Relation, ID string) {
	relationsCache, _ := json.Marshal(relations)

	var key string = "Relations:" + ID

	status := s.rdb.Set(context.Background(), key, relationsCache, time.Hour*72)

	if status.Err() != nil {
		log.Println("Error in set in the cache " + status.Err().Error())
	}
}

//GetRelationsCache Get last Relations from Redis
func (s *UserStorageService) GetRelationsCache(ID string) ([]*models.RelationReponse, error) {
	//Get info from redis
	val := s.rdb.Get(context.Background(), "Relations:"+ID)

	err := val.Err()

	if err != nil && err != redis.Nil {
		log.Println("Error in get the cache " + err.Error())
	}

	//Convert for response
	var relationDB []*models.Relation

	if err != redis.Nil {
		userBytes, _ := val.Bytes()
		json.Unmarshal(userBytes, &relationDB)

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

	return nil, err
}

//UpdateRelations Update the first relations of User in Cache
func (s *UserStorageService) UpdateRelations(ID string) error {
	//Get info from DB
	var relationDB []*models.Relation = []*models.Relation{new(models.Relation)}

	parseID, err := uuid.Parse(ID)

	if err != nil {
		return err
	}

	s.db.Where("from_user = ?", parseID).Or("to_user = ? AND mutual = true", parseID).Find(&relationDB).Limit(30)

	if err := s.db.Error; err != nil {
		return err
	}

	relationsCache, _ := json.Marshal(relationDB)

	var key string = "Relations:" + ID

	status := s.rdb.Set(context.Background(), key, relationsCache, time.Hour*72)

	if status.Err() != nil {
		return status.Err()
	}

	return nil
}
