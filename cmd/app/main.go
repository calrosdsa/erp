package main

import (
	"log"	
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	// "github.com/spf13/viper"
)


func init(){
	viper.SetConfigFile(`../../configs/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	// viper.SetCon

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// 94442fed-cf6f-4255-bc22-e4dffcbdec4e
	// core.Init()
}
