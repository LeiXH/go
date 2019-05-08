package database

import "time"

func getOrDefault(v, d string) string {
	if len(v) == 0 {
		return d
	}
	return v
}
func getOrDefaultInt(v, d int) int {
	if v <= 0 {
		return d
	}
	return v
}
func getOrDefaultTime(v, d time.Duration) time.Duration {
	if v.Seconds() < 0 {
		return d
	}
	return v
}

func mergeConfig(config *MySQLConfig) *MySQLConfig {
	mergeConfig := MySQLConfig{
		Username:        getOrDefault(config.Username, "root"),
		Password:        getOrDefault(config.Password, ""),
		Host:            getOrDefault(config.Host, "localhost"),
		Port:            getOrDefault(config.Port, "3306"),
		Database:        getOrDefault(config.Database, ""),
		ConnMaxLifetime: getOrDefaultTime(config.ConnMaxLifetime, 4*time.Hour),
		MaxIdleConns:    getOrDefaultInt(config.MaxIdleConns, 0),
		MaxOpenConns:    getOrDefaultInt(config.MaxIdleConns, 10),
	}

	return &mergeConfig
}
