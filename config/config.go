package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"log/slog"
	"os"
	"sync"
)

const (
	FileConfig  = "tradier-config.yml"
	APIEndpoint = "https://api.tradier.com"
)

var k = koanf.New(".")

type Config struct {
	Simulation bool `koanf:"simulation"`
	Sandbox    bool `koanf:"sandbox"`
	Debug      bool `koanf:"debug"`
}

type Dsn struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Dbname   string `koanf:"dbname"`
}

// singleton
var configMutex = &sync.Mutex{}
var configSingleton *Config

func Instance() *Config {
	if configSingleton != nil {
		return configSingleton
	}
	configMutex.Lock()
	defer configMutex.Unlock()
	if configSingleton == nil {
		configSingleton = &Config{}
		SetDefaults(configSingleton)
		loadFromFile(configSingleton)
	}
	return configSingleton
}

func SetDefaults(cfg *Config) {
}

func loadFromFile(cfg *Config) error {
	path, err := FindConfigFile()
	if err != nil {
		return err
	}
	slog.Debug("using config file", slog.String("path", path))
	if err = k.Load(file.Provider(path), yaml.Parser()); err != nil {
		LogError("error loading config", err)
		os.Exit(4)
	}
	//k.Print()
	err = k.Unmarshal("", cfg)
	if err != nil {
		LogError("unmarshal error", err)
		os.Exit(5)
	}
	return nil
}

func (c *Config) Log() {
	slog.Debug("config", slog.Bool("simulation", c.Simulation))
	slog.Debug("config", slog.Bool("sandbox", c.Sandbox))
}
