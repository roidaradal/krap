package app

import (
	"errors"
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

// Check if valid env ('dev' or 'prod')
func IsValidEnv(env Env) error {
	if env != EnvDev && env != EnvProd {
		return errors.New("invalid app env")
	}
	return nil
}

// Check if app env is 'prod'
func IsProdEnv(env Env) bool {
	return env == EnvProd
}

// Run initializers list
func RunInitializers(initializers []Initializer) error {
	for _, initializer := range initializers {
		err := initializer.Fn()
		if err != nil {
			return fmt.Errorf("%s: failed to initialize: %w", initializer.Name, err)
		}
	}
	return nil
}
