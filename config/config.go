package config

import (
	"os"
)

const Layout = "20060102"

type TodoEnvironment struct {
	Port     string
	DBFile   string
	Password string
}

func checkEnvironment(env, baseValue string) string {
	if value, ok := os.LookupEnv(env); ok {
		return value
	}
	return baseValue
}

func GetEnv() *TodoEnvironment {
	port := checkEnvironment("TODO_PORT", "7540")
	dbfile := checkEnvironment("TODO_DBFILE", "")
	password := checkEnvironment("TODO_PASSWORD", "")

	return &TodoEnvironment{
		Port:     port,
		DBFile:   dbfile,
		Password: password,
	}
}
