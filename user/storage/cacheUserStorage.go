package storage

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"

	"github.com/Marlos-Rodriguez/go-postgres-wallet-back/user/models"
)

//SetUser Set the User in redis Cache
func (u *UserStorageService) SetUser(userDB *models.User) {
	redisUser, err := json.Marshal(userDB)

	if err != nil {
		log.Println("Error in Marshal the user" + err.Error())
	}
	status := u.rdb.Set(context.Background(), "User:"+userDB.UserID.String(), redisUser, time.Hour*72)

	if status.Err() != nil {
		log.Println("Error in set in the cache " + status.Err().Error())
	}
}

//GetUserCache Get info from redis if exits
func (u *UserStorageService) GetUserCache(ID string) (*models.UserResponse, error) {
	//Get from Redis
	val := u.rdb.Get(context.Background(), "User:"+ID)

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
func (u *UserStorageService) GetProfileCache(ID string) (*models.UserProfileResponse, error) {
	//Get info from redis
	val := u.rdb.Get(context.Background(), "Profile:"+ID)

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
func (u *UserStorageService) SetProfile(profileDB *models.Profile) {
	redisUser, err := json.Marshal(profileDB)

	if err != nil {
		log.Println("Error in Marshal the user" + err.Error())
	}
	status := u.rdb.Set(context.Background(), "Profile:"+profileDB.UserID.String(), redisUser, time.Hour*72)

	if status.Err() != nil {
		log.Println("Error in set in the cache " + status.Err().Error())
	}
}

//UpdateUserCache Update User & Profile in Cache
func (u *UserStorageService) UpdateUserCache(ID string) {
	//Get info from DB
	var userDB *models.User = new(models.User)
	var profileDB *models.Profile = new(models.Profile)

	go u.db.Where("user_id = ?", &ID).First(&userDB)
	u.db.Where("user_id = ?", ID).First(&profileDB)

	err := u.db.Error

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
	status := u.rdb.Set(context.Background(), "User:"+userDB.UserID.String(), redisUser, time.Hour*72)

	if status.Err() != nil {
		log.Println("Error in set in the cache " + status.Err().Error())
	}

	status = u.rdb.Set(context.Background(), "Profile:"+profileDB.UserID.String(), redisProfile, time.Hour*72)

	if status.Err() != nil {
		log.Println("Error in set in the cache " + status.Err().Error())
	}

	log.Println("User cache Updated")
}

//SetRelationCache Set One page of 20 relations
func (u *UserStorageService) SetRelationCache(relations []*models.Relation, ID string) {
	relationsCache, _ := json.Marshal(relations)

	var key string = "Relations:" + ID

	status := u.rdb.Set(context.Background(), key, relationsCache, time.Hour*72)

	if status.Err() != nil {
		log.Println("Error in set in the cache " + status.Err().Error())
	}
}

//GetRelationsCache Get last Relations from Redis
func (u *UserStorageService) GetRelationsCache(ID string) ([]*models.RelationReponse, error) {
	//Get info from redis
	val := u.rdb.Get(context.Background(), "Relations:"+ID)

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
func (u *UserStorageService) UpdateRelations(ID string) {
	//Get info from DB
	var relationDB []*models.Relation = []*models.Relation{new(models.Relation)}

	u.db.Where("from_user = ?", ID).Or("to_user = ? AND mutual = true", ID).Find(&relationDB).Limit(30)

	if err := u.db.Error; err != nil {
		log.Println("Error in DB in Cache " + err.Error())
	}

	relationsCache, _ := json.Marshal(relationDB)

	var key string = "Relations:" + ID

	status := u.rdb.Set(context.Background(), key, relationsCache, time.Hour*72)

	if status.Err() != nil {
		log.Println("Error in set in the cache " + status.Err().Error())
	}
}
