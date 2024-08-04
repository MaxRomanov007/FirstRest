package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string           `yaml:"env"`
	DB         DbConnConfig     `yaml:"database_connection"`
	HTTPServer HTTPServerConfig `yaml:"http_server"`
	Cars       CarsConfig       `yaml:"cars"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
	User        string        `yaml:"user" env-required:"true"`
	Pass        string        `yaml:"password" env-required:"true"`
}

type DbConnConfig struct {
	Server string `yaml:"server" env-required:"true"`
	DB     string `yaml:"database" env-required:"true"`
	Port   int    `yaml:"port" env-required:"true"`
	User   string `yaml:"username" env-required:"true"`
	Pass   string `yaml:"password" env-required:"true"`
	Ssl    bool   `yaml:"ssl" env-default:"false"`
}

type CarsConfig struct {
	ImagesPath     string `yaml:"images_path" env-required:"true"`
	ImagesFormName string `yaml:"images_form_name" env-default:"images"`
	DefaultOrderBy string `yaml:"default_order_by" env-default:"producer"`
}

func MustLoad() *Config {
	path := MustGetPath()

	return MustLoadByPath(path)
}

func MustGetPath() string {
	path := getPath()
	if path == "" {
		log.Fatal("config path not set")
	}

	return path
}

func getPath() string {
	if path := getPathByEnv(); path != "" {
		return path
	}
	return getPathByFlag()
}

func getPathByEnv() string {
	path := os.Getenv("CONFIG_PATH")
	return path
}

func getPathByFlag() string {
	var path string

	flag.StringVar(&path, "config_path", "", "path to config file")
	flag.Parse()
	return path
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatal("failed to read config:" + err.Error())
	}

	return &cfg
}
