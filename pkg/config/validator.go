package config

import (
	"strings"

	"github.com/intility/cwc/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

func DefaultValidator(cfg *Config) error {
	var validationErrors []string

	if cfg.APIKey() == "" {
		validationErrors = append(validationErrors, "apiKey must be provided and not be empty")
	}

	if cfg.Endpoint == "" {
		validationErrors = append(validationErrors, "endpoint must be provided and not be empty")
	}

	if cfg.Provider == "" {
		cfg.Provider = "azure"
	}

	switch cfg.Provider {
	case "azure":
		if cfg.ModelDeployment == "" {
			cfg.ModelDeployment = openai.GPT4TurboPreview
		}
	case "openai":
		if cfg.Model == "" {
			validationErrors = append(validationErrors, "model must be provided and not be empty")
		}
	default:
		validationErrors = append(validationErrors, "provider not supported. supported providers:"+strings.Join(SupportedProviders, " "))
	}

	if cfg.ApiVersion == "" {
		validationErrors = append(validationErrors, "apiversion must be provided and not be empty")
	}

	if len(validationErrors) > 0 {
		return &errors.ConfigValidationError{Errors: validationErrors}
	}

	return nil
}
