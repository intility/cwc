package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/intility/cwc/pkg/config"
	"github.com/intility/cwc/pkg/errors"
	cwcui "github.com/intility/cwc/pkg/ui"
)

var (
	providerFlag        string //nolint:gochecknoglobals
	apiKeyFlag          string //nolint:gochecknoglobals
	endpointFlag        string //nolint:gochecknoglobals
	modelFlag           string //nolint:gochecknoglobals
	modelDeploymentFlag string //nolint:gochecknoglobals
	apiVersionFlag      string //nolint:gochecknoglobals
)

func promptUserForConfig(ui cwcui.UI) { //nolint: varnamelen
	if providerFlag == "" {
		ui.PrintMessage("Enter Provider name (azure,openai): ", cwcui.MessageTypeInfo)
		providerFlag = config.SanitizeInput(ui.ReadUserInput())
	}

	if apiKeyFlag == "" {
		ui.PrintMessage("Enter the OpenAI API Key: ", cwcui.MessageTypeInfo)
		apiKeyFlag = config.SanitizeInput(ui.ReadUserInput())
	}

	if endpointFlag == "" {
		ui.PrintMessage("Enter the OpenAI API Endpoint: ", cwcui.MessageTypeInfo)
		endpointFlag = config.SanitizeInput(ui.ReadUserInput())
	}

	if providerFlag == "azure" {
		if modelDeploymentFlag == "" {
			ui.PrintMessage("Enter the Azure OpenAI Model Deployment: ", cwcui.MessageTypeInfo)
			modelDeploymentFlag = config.SanitizeInput(ui.ReadUserInput())
		}

		if apiVersionFlag == "" {
			apiVersionFlag = config.APIVersion
		}
	}

	if providerFlag == "openai" {
		if modelFlag == "" {
			ui.PrintMessage("Enter the Model name: ", cwcui.MessageTypeInfo)
			modelFlag = config.SanitizeInput(ui.ReadUserInput())
		}

		if apiVersionFlag == "" {
			ui.PrintMessage("Enter the API Version: ", cwcui.MessageTypeInfo)
			apiVersionFlag = config.SanitizeInput(ui.ReadUserInput())
		}
	}
}

func createLoginCmd() *cobra.Command {
	ui := cwcui.NewUI() //nolint:varnamelen
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Azure OpenAI",
		Long: "Login will prompt you to enter your Azure OpenAI API key " +
			"and other relevant information required for authentication.\n" +
			"Your credentials will be stored securely in your keyring and will never be exposed on the file system directly.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Prompt for other required authentication details (apiKey, endpoint, version, and deployment)
			promptUserForConfig(ui)

			cfg := config.NewConfig(config.NewConfigParams{
				Provider:        providerFlag,
				Endpoint:        endpointFlag,
				APIVersion:      apiVersionFlag,
				ModelDeployment: modelDeploymentFlag,
				Model:           modelFlag,
			})

			cfg.SetAPIKey(apiKeyFlag)

			provider := config.NewDefaultProvider()
			err := provider.SaveConfig(cfg)
			if err != nil {
				if validationErr, ok := errors.AsConfigValidationError(err); ok {
					for _, e := range validationErr.Errors {
						ui.PrintMessage(e+"\n", cwcui.MessageTypeError)
					}

					return nil // suppress the error
				}

				return fmt.Errorf("error saving configuration: %w", err)
			}

			ui.PrintMessage("config saved successfully\n", cwcui.MessageTypeSuccess)

			return nil
		},
	}

	cmd.Flags().StringVarP(&providerFlag, "provider", "p", "",
		"Provider name. Supported providers: "+strings.Join(config.SupportedProviders, " "))
	cmd.Flags().StringVarP(&modelFlag, "model", "m", "", "OpenAI (compatible) Model")
	cmd.Flags().StringVarP(&apiKeyFlag, "api-key", "k", "", "API Key")
	cmd.Flags().StringVarP(&endpointFlag, "endpoint", "e", "", "API Endpoint")
	cmd.Flags().StringVarP(&apiVersionFlag, "api-version", "v", "", "API Version")
	cmd.Flags().StringVarP(&modelDeploymentFlag, "model-deployment", "d", "", "Azure OpenAI Model Deployment")

	return cmd
}
