package database

import "time"

type MySQLConfig struct {
	Username string `mapstructure:",omitempty"`
	Password string `mapstructure:",omitempty"`
	Host     string `mapstructure:",omitempty"`
	Port     string `mapstructure:",omitempty"`
	Database string `mapstructure:",omitempty"`

	ConnMaxLifetime time.Duration `mapstructure:"max_life_time,omitempty"`
	MaxIdleConns    int           `mapstructure:"max_idle,omitempty"`
	MaxOpenConns    int           `mapstructure:"max_open,omitempty"`
}
