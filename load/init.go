package load

import (
	"log"
	"os"

	"github.com/spf13/viper"
)


func Env() {
	// loading the .env file
	masterEnv := os.Getenv("MONOPOLY_ENV")

	println("Start Game")
	println("MONOPOLY_ENV is ", masterEnv)
	viper.SetConfigFile(masterEnv)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.Fatalf("Error reading config file: ", err)
	}


}