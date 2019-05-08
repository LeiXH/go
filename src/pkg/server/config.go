package server

import (
	"log"
	"path/filepath"
	"strings"

	"pkg/database"
	"pkg/logger"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Option struct {
	PrinterConfig struct {
		Enabled  bool
		Port     string
		Width    string
		Height   string
		Speed    string
		Density  string
		Sensor   string
		Vertical string
		Offset   string
	} `yaml:"printer"`
	MysqlConfig  database.MySQLConfig `yaml:"mysql"`
	ServerConfig struct {
		Port string
		Mode string
	} `yaml:"server"`
}

func (srv *Server) loadConfig(configPath string) error {
	suf := filepath.Ext(configPath)
	if suf != ".yaml" && suf != ".yml" {
		logger.Fatalf("config file %s not recognized", configPath)
	}
	config := viper.New()
	config.SetConfigType("yaml")

	config.SetConfigName(strings.TrimSuffix(filepath.Base(configPath), suf))
	config.AddConfigPath(filepath.Dir(configPath))

	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
	return unmarshalViperConfigToObj(config.AllSettings(), &srv.config)
}

func unmarshalViperConfigToObj(confs interface{}, dst interface{}) error {
	tmp, _ := yaml.Marshal(confs)
	return yaml.Unmarshal(tmp, dst)
}
