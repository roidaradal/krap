package app

import (
	"fmt"
	"os"
)

type Env = string

const (
	EnvDev  Env = "dev"
	EnvProd Env = "prod"
)

// Initializer function with name
type Initializer struct {
	Fn   func() error
	Name string
}

// Check required env keys
func CheckRequiredEnvKeys(keys []string) error {
	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		if value == "" || !ok {
			return fmt.Errorf("missing env variable: %s", key)
		}
	}
	return nil
}
