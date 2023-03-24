package main

import "os"

type MySQLConfig struct {
	EnvRef        string
	MySQLDatabase string
	MySQLUser     string
	MySQLPassword string
	MySQLHost     string
	MySQLPort     string
}

func Envs() MySQLConfig {
	envRef := os.Getenv("ENV_REF")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")

	return MySQLConfig{
		EnvRef:        envRef,
		MySQLDatabase: mysqlDatabase,
		MySQLUser:     mysqlUser,
		MySQLPassword: mysqlPassword,
		MySQLHost:     mysqlHost,
		MySQLPort:     mysqlPort,
	}
}

func (p MySQLConfig) MysqlString() string {
	return p.MySQLUser + ":" + p.MySQLPassword + "@tcp(" + p.MySQLHost + ":" + p.MySQLPort + ")/" + p.MySQLDatabase
}
