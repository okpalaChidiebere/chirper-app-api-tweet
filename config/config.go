package config

import (
	"log"
	"os"
)

type env struct {
	UserTable string
	TweetTable string
}

type aws struct {
	Aws_region       string
	Aws_profile      string
}

type Config struct {
	Dev env
	Aws aws
	Prod env
}

func NewConfig() *Config {
	return &Config{
	   Dev: env{
			UserTable: "chirper-app-users-dev",
			TweetTable: "chirper-app-tweets-dev",
	   },
		Aws: aws{
			Aws_region:       mustGetenv("AWS_REGION"),
			Aws_profile:      mustGetenv("AWS_PROFILE"),
		},
		Prod: env{
			UserTable: "",
	   },
	}
}

func (c *Config) IsLocal() bool {
	return c.Aws.Aws_profile != "DEPLOYED"
}

func mustGetenv (k string) string {
	v, ok := os.LookupEnv(k)
	if !ok {
		log.Fatalf("Warning: %s environment variable is not set.", k)
	}
	return v
}