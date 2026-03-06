package utils

import (
	"fmt"
	"os"
	"strings"
)

func StrToPtr(s string) *string {
	return &s
}

func GetRequiredEnvVar(envVarName string) (string, error) {
	envVarValue := strings.TrimSpace(os.Getenv(envVarName))
	if envVarValue == "" {
		return "", fmt.Errorf("environment variable %s is not set", envVarName)
	}
	return envVarValue, nil
}

func GetOptionalEnvVar(envVarName string, defaultValue string) string {
	envVarValue := strings.TrimSpace(os.Getenv(envVarName))
	if envVarValue == "" {
		return defaultValue
	}
	return envVarValue
}
