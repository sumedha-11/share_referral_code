package config

import (
	"encoding/json"
	"fmt"
	"github.com/sumedha-11/share_referral_code/pkg/sql"
	"io/ioutil"
	"log"
	"os"
)

var Config *AppConfig

type AppConfig struct {
	Common     *CommonConfig `json:"CommonConfig"`
	DBConfig   *sql.DBConfig `json:"DBConfig"`
	GoogleCred *Credentials  `json:"GoogleCred"`
}

// Credentials which stores google ids.
type Credentials struct {
	Cid      string `json:"cid"`
	Csecret  string `json:"csecret"`
	Redirect string `json:"redirect"`
}

type CommonConfig struct {
	ServerPort string `json:"serverPort"`
}

func ReadConfig(configFile string) {
	Config = &AppConfig{}
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(file, &Config); err != nil {
		log.Fatalf("unable to marshal config data")
		return
	}
	fmt.Println("config loaded ", Config)
}
