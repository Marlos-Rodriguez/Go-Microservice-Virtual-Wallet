package environment

import (
	"os"
	//Autoload the env
	_ "github.com/joho/godotenv/autoload"
)

var environment = make(map[string]string)

//AccessENV Return the ENV if exits
func AccessENV(key string) string {
	if environment[key] != "" {
		return environment[key]
	}

	val := os.Getenv(key)

	if val == "" || len(val) <= 0 {
		return ""
	}

	environment[key] = val

	return val
}
