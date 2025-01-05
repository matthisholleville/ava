package configuration

import (
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
	Type   string `json:"type,omitempty" example:"openai"`
	OpenAI OpenAI `json:"openai,omitempty"`
}

type OpenAI struct {
	APIKey string `json:"apiKey,omitempty" example:""`
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
	Type  string      `json:"type,omitempty" example:"slack"`
	Slack SlackEvents `json:"slack,omitempty"`
}

type SlackEvents struct {
	ValidationToken string `json:"validation_token,omitempty"`
	BotToken        string `json:"bot_token,omitempty"`
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
