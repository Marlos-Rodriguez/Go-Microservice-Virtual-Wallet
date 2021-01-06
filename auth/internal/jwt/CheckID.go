package jwt

import (
	"errors"
)

//CheckID Check if the Passed ID match with JWT User ID
func CheckID(ID string) (bool, error) {
	if ID == "" || len(ID) <= 0 {
		return false, errors.New("Must sent ID")
	}

	if ID != UserID {
		return false, errors.New("The ID not match with Token ID")
	}

	return true, nil
}
