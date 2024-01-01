package lib

import (
	"fmt"
	"os"
)

const (
	dbHostEnvVar = "DBHOST"
	DBHost       = "Host"

	dbUserEnvVar = "DBUSER"
	DBUsername   = "User"

	dbPassEnvVar = "DBPASSWORD"
	DBPass       = "Pass"

	dbPortEnvVar = "DBPORT"
	DBPort       = "Port"

	devEnvVar   = "DEV"
	Development = "Development"

	lookupErr = "ENVVAR for %q not found"
)

type Config struct {
	cfg map[string]string
}

func NewConfig(requiredConfigs ...string) (*Config, error) {
	c := &Config{
		cfg: map[string]string{},
	}

	if err := c.Populate(requiredConfigs); err != nil {
		return c, err
	}

	return c, nil
}

func (c *Config) Get(key string) (string, error) {
	value, ok := c.cfg[key]
	if !ok {
		return "", fmt.Errorf("config key %q not found", key)
	}
	return value, nil
}

func (c *Config) Set(key, value string) {
	c.cfg[key] = value
}

var (
	lookupMap = map[string]string{
		dbHostEnvVar: DBHost,
		dbUserEnvVar: DBUsername,
		dbPassEnvVar: DBPass,
		dbPortEnvVar: DBPort,
	}
)

func (c *Config) Populate(requiredConfigs []string) error {
	for envVarKey, cfgKey := range lookupMap {
		val, ok := os.LookupEnv(envVarKey)
		if !ok && contains(requiredConfigs, cfgKey) {
			return fmt.Errorf(lookupErr, envVarKey)
		}
		c.cfg[cfgKey] = val
	}

	if _, ok := os.LookupEnv(devEnvVar); !ok {
		c.cfg[Development] = "false"
	} else {
		c.cfg[Development] = "true"
	}

	return nil
}

func contains(slice []string, target string) bool {
	for _, elem := range slice {
		if elem == target {
			return true
		}
	}

	return false
}
