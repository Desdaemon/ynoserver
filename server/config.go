package server

import (
	"io/ioutil"
	"fmt"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	IP              string `yaml:"ip"`
	Port            int    `yaml:"port"`
	IndexPath       string `yaml:"index_path"`
	NumRooms        int    `yaml:"num_rooms"`
	BadRooms        string `yaml:"bad_rooms"`
	BadSounds       string `yaml:"bad_sounds"`
	PictureNames    string `yaml:"picture_names"`
	PicturePrefixes string `yaml:"picture_prefixes"`
	GameName        string `yaml:"game_name"`
	SignKey		 	string `yaml:"sign_key"`
	IPHubKey		string `yaml:"iphub_key"`
	Database		struct {
		User		string `yaml:"user"`
		Pass		string `yaml:"pass"`
		Host		string `yaml:"host"`
		Name		string `yaml:"name"`
	} `yaml:"database"`
	Logging         struct {
		File       string `yaml:"file"`
		MaxSize    int    `yaml:"max_size"`
		MaxBackups int    `yaml:"max_backups"`
		MaxAge     int    `yaml:"max_age"`
	} `yaml:"logging"`
}

func ParseConfig(file string) ServerConfig {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var config ServerConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	if config.IndexPath == "" {
		config.IndexPath = "games/default/index.json"
	}
	if config.IP == "" {
		config.IP = "127.0.0.1"
	}
	if config.Port == 0 {
		config.Port = 8080
	}
	if config.NumRooms == 0 {
		config.NumRooms = 100
	}
	if config.Logging.File == "" {
		config.Logging.File = "server.log"
	}
	if config.Logging.MaxSize == 0 {
		config.Logging.MaxSize = 50 // MB
	}
	if config.Logging.MaxBackups == 0 {
		config.Logging.MaxBackups = 6
	}
	if config.Logging.MaxAge == 0 {
		config.Logging.MaxAge = 28 // Days
	}

	fmt.Printf("%+v\n", config)

	return config
}

type Config struct {
	spriteNames []string
	systemNames []string
	soundNames []string
	ignoredSoundNames []string
	pictureNames []string
	picturePrefixes []string

	gameName string

	signKey string
	ipHubKey string

	dbUser string
	dbPass string
	dbHost string
	dbName string
}

func GetConfig(spriteNames []string, systemNames []string, soundNames []string, ignoredSoundNames []string, pictureNames []string, picturePrefixes []string, gameName string, signKey string, ipHubKey string, dbUser string, dbPass string, dbHost string, dbName string) (Config) {
	c := Config{
		spriteNames: spriteNames,
		systemNames: systemNames,
		soundNames: soundNames,
		ignoredSoundNames: ignoredSoundNames,
		pictureNames: pictureNames,
		picturePrefixes: picturePrefixes,
		gameName: gameName,

		signKey: signKey,
		ipHubKey: ipHubKey,

		dbUser: dbUser,
		dbPass: dbPass,
		dbHost: dbHost,
		dbName: dbName,
	}

	return c
}
