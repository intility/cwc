package config

import (
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/intility/cwc/pkg/errors"
)

func DefaultValidator(cfg *Config) error {
	var validationErrors []string

	validationErrors = append(validationErrors, validateAPIKey(cfg)...)
	validationErrors = append(validationErrors, validateEndpoint(cfg)...)
	validationErrors = append(validationErrors, validateProvider(cfg)...)
	validationErrors = append(validationErrors, validateAPIVersion(cfg)...)

	if len(validationErrors) > 0 {
		return &errors.ConfigValidationError{Errors: validationErrors}
	}

	return nil
}

func validateAPIKey(cfg *Config) []string {
	var errors []string
	if cfg.APIKey() == "" {
		errors = append(errors, "apiKey must be provided and not be empty")
	}

	return errors
}

func validateEndpoint(cfg *Config) []string {
	var errors []string
	if cfg.Endpoint == "" {
		errors = append(errors, "endpoint must be provided and not be empty")
	}

	return errors
}

func validateProvider(cfg *Config) []string {
	var errors []string

	if cfg.Provider == "" {
		cfg.Provider = providerAzure
	}

	switch cfg.Provider {
	case providerAzure:
		if cfg.ModelDeployment == "" {
			cfg.ModelDeployment = openai.GPT4TurboPreview
		}
	case providerOpenai:
		if cfg.Model == "" {
			errors = append(errors, "model must be provided and not be empty")
		}
	default:
		errors = append(errors,
			"provider not supported. supported providers:"+strings.Join(SupportedProviders, " "))
	}

	return errors
}

func validateAPIVersion(cfg *Config) []string {
	var errors []string
	if cfg.APIVersion == "" {
		errors = append(errors, "apiversion must be provided and not be empty")
	}

	return errors
}
