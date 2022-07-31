package config

import (
	"fmt"
	"time"

	"github.com/Netflix/go-env"
)

// Config loads and stores all the environment variables for the program
type Config struct {
	ProjectID              string `env:"PROJECT_ID,required=true"`
	KymBucketName          string `env:"KYM_BUCKET_NAME,required=true"`
	Debug                  bool   `env:"DEBUG,default=false"`
	OrganisationServiceURL string `env:"ORGANISATION_SERVICE_URL,required=true"`
	RecipientServiceURL    string `env:"RECIPIENT_SERVICE_URL,required=true"`

	Server struct {
		Port         string        `env:"PORT,default=8080"`
		TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required=true"`
		TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required=true"`
		TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required=true"`
	}

	MerchantCreatePublisher struct {
		TopicID string `env:"PUBLISHER_TOPIC_ID,required=true"`
	}
}

// AppConfig initializes the environment variables
func AppConfig() (*Config, error) {

	// parses env variables to Config struct
	config := &Config{}
	if _, err := env.UnmarshalFromEnviron(config); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return config, nil
}
