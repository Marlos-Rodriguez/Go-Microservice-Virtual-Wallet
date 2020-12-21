package environment

import (
	"os"
	//Autoload the env
	_ "github.com/joho/godotenv/autoload"
)

var environment = make(map[string]string)

//AccessENV Return the ENV if exits
func AccessENV(key string) (string, bool) {
	if environment[key] != "" {
		return environment[key], true
	}

	val := os.Getenv(key)

	if val == "" || len(val) <= 0 {
		return "", false
	}

	environment[key] = val

	return val, true
}
