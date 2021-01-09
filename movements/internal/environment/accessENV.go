package environment

import (
	"os"
	"sync"

	//Autoload the env
	_ "github.com/joho/godotenv/autoload"
)

var (
	environment      = map[string]string{}
	environmentMutex = sync.RWMutex{}
)

//AccessENV Return the ENV if exits
func AccessENV(key string) string {
	environmentMutex.RLock()
	if environment[key] != "" {
		val := environment[key]
		environmentMutex.RUnlock()
		return val
	}
	environmentMutex.RUnlock()

	val := os.Getenv(key)

	if val == "" || len(val) <= 0 {
		return ""
	}
	environmentMutex.Lock()
	environment[key] = val
	environmentMutex.Unlock()

	return val
}
