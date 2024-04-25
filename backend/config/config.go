// config/config.go
package config

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

var (
	once           sync.Once
	configInstance *Config
	secretsMap     map[string]string
)

type Config struct {
	data           map[string]interface{}
	secretsManager *secretsmanager.SecretsManager
	awsSession     *session.Session
}

// GetConfig returns the global configuration instance.
func GetConfig() *Config {
	return configInstance
}

// LoadConfig loads the YAML configuration file.
func LoadConfig(env string, awsSession *session.Session) (*Config, error) {
	once.Do(func() {
		var filePath string
		if env == "local" {
			filePath = "config-local.yml"
		} else {
			filePath = "config-prod.yml"
		}
		configData, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Error reading configuration file: %v", err)
		}

		configMap := make(map[string]interface{})
		if err := yaml.Unmarshal(configData, &configMap); err != nil {
			log.Fatalf("Error parsing configuration file: %v", err)
		}

		configInstance = &Config{
			data:           configMap,
			awsSession:     awsSession,
			secretsManager: secretsmanager.New(awsSession),
		}

		configInstance.retrieveSecrets(env)
		for k, v := range secretsMap {
			configMap[k] = v
		}
	})

	return configInstance, nil
}

// retrieveSecrets automatically retrieves secrets from AWS Secrets Manager.
func (c *Config) retrieveSecrets(env string) {
	if secretsMap == nil {
		secretsMap = make(map[string]string)
	}

	secretName := "backend/" + env

	secretValue, err := c.getSecretFromAWS(secretName)
	if err != nil {
		log.Fatalf("Error retrieving secret from AWS Secrets Manager: %v", err)
	}

	var secretData map[string]string
	if err := json.Unmarshal([]byte(secretValue), &secretData); err != nil {
		log.Fatalf("Error parsing secret JSON: %v", err)
	}

	for key, value := range secretData {
		secretsMap[key] = value
	}
}

// getSecretFromAWS retrieves a secret from AWS Secrets Manager.
func (c *Config) getSecretFromAWS(key string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &key, // Use your secret's ARN or name here
	}

	result, err := c.secretsManager.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}

// Get retrieves a configuration value by key.
func (c *Config) Get(key string) interface{} {
	return c.data[key]
}

// GetString retrieves a configuration value as a string.
func (c *Config) GetString(key string) string {
	val, ok := c.data[key].(string)
	if !ok {
		return ""
	}
	return val
}

// GetSecret retrieves a secret from the configuration.
func (c *Config) GetSecret(key string) (string, error) {
	if secretValue, exists := secretsMap[key]; exists {
		return secretValue, nil
	}

	return "", nil
}
