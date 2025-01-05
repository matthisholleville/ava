package configuration

import (
	"fmt"
	"reflect"

	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/spf13/viper"
)

type Configuration struct {
	Knowledge Knowledge `yaml:"knowledge,omitempty"`
	Executors Executors `yaml:"executors,omitempty"`
	AI        AI        `yaml:"ai,omitempty"`
	API       API       `yaml:"api,omitempty"`
	Events    Events    `yaml:"events,omitempty"`
}

type Knowledge struct {
	Github GithubKnowledge `yaml:"github,omitempty"`
}

type GithubKnowledge struct {
	Token string `yaml:"token,omitempty"`
}

type Executors struct {
	Enabled bool            `yaml:"enabled,omitempty"`
	K8S     K8SExecutors    `yaml:"k8s,omitempty"`
	Common  CommonExecutors `yaml:"common,omitempty"`
	Web     WebExecutors    `yaml:"web,omitempty"`
}

type K8SExecutors struct {
	Read  bool `yaml:"read,omitempty"`
	Write bool `yaml:"write,omitempty"`
}

type CommonExecutors struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type WebExecutors struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type AI struct {
	Type   string `yaml:"type,omitempty" example:"openai"`
	OpenAI OpenAI `yaml:"openai,omitempty"`
}

type OpenAI struct {
	APIKey string `yaml:"apiKey,omitempty" example:""`
}

type API struct {
	Chat      ChatAPI      `yaml:"chat,omitempty"`
	Knowledge KnowledgeAPI `yaml:"knowledge,omitempty"`
	Events    EventsAPI    `yaml:"events,omitempty"`
	Swagger   Swagger      `yaml:"swagger,omitempty"`
}

type Swagger struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type ChatAPI struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type KnowledgeAPI struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type EventsAPI struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type Events struct {
	Type  string      `yaml:"type,omitempty" example:"slack"`
	Slack SlackEvents `yaml:"slack,omitempty"`
}

type SlackEvents struct {
	ValidationToken string `yaml:"validation_token,omitempty"`
	BotToken        string `yaml:"bot_token,omitempty"`
}

func WriteInitConfig(logger logger.ILogger) {
	// Define default values for the configuration
	viper.SetDefault("executors.enabled", true)
	viper.SetDefault("executors.k8s.write", false)
	viper.SetDefault("executors.k8s.read", true)
	viper.SetDefault("executors.common.enabled", true)
	viper.SetDefault("executors.web.enabled", true)

	viper.SetDefault("ai.type", "openai")
	viper.SetDefault("ai.openai.apiKey", "${OPENAI_API_KEY}")

	viper.SetDefault("api.chat.enabled", true)
	viper.SetDefault("api.knowledge.enabled", true)
	viper.SetDefault("api.events.enabled", true)
	viper.SetDefault("api.swagger.enabled", true)

	viper.SetDefault("events.type", "slack")
	viper.SetDefault("events.slack.validationToken", "${SLACK_VALIDATION_TOKEN}")
	viper.SetDefault("events.slack.botToken", "${SLACK_BOT_TOKEN}")

	// Write the default configuration to a file
	if err := viper.SafeWriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
			logger.Debug("Configuration file already exists. No changes made.")
		} else {
			logger.Fatal(fmt.Sprintf("Unable to write configuration file: %v", err))
		}
	} else {
		logger.Info("Default configuration written successfully")
	}
}

func LoadConfiguration(logger logger.ILogger) *Configuration {
	var config Configuration
	err := viper.Unmarshal(&config)
	if err != nil {
		logger.Fatal("error on fetching configuration file")
	}

	if reflect.DeepEqual(config, Configuration{}) {
		logger.Fatal("Configuration is empty")
	}

	return &config
}
