package config

import (
	"os"
	"testing"
	"time"
)

func TestAppConfig(t *testing.T) {
	_ = os.Setenv("PROJECT_ID", "123")
	_ = os.Setenv("DEBUG", "true")
	_ = os.Setenv("PORT", "3000")
	_ = os.Setenv("SERVER_TIMEOUT_READ", "10s")
	_ = os.Setenv("SERVER_TIMEOUT_WRITE", "10s")
	_ = os.Setenv("SERVER_TIMEOUT_IDLE", "10s")
	_ = os.Setenv("PUBLISHER_TOPIC_ID", "merchant-create")
	_ = os.Setenv("KYM_BUCKET_NAME", "bucket-dev")
	_ = os.Setenv("ORGANISATION_SERVICE_URL", "organisation-dev-url")
	_ = os.Setenv("RECIPIENT_SERVICE_URL", "recipient-dev-url")

	cfg, err := AppConfig()
	if err != nil {
		t.Errorf("failed to read config from env: %v", err)
	}

	cmp(t, "PROJECT_ID", cfg.ProjectID, "123")
	cmp(t, "DEBUG", cfg.Debug, true)
	cmp(t, "PORT", cfg.Server.Port, "3000")
	cmp(t, "SERVER_TIMEOUT_READ", cfg.Server.TimeoutRead, time.Second*10)
	cmp(t, "SERVER_TIMEOUT_WRITE", cfg.Server.TimeoutWrite, time.Second*10)
	cmp(t, "SERVER_TIMEOUT_IDLE", cfg.Server.TimeoutIdle, time.Second*10)
	cmp(t, "PUBLISHER_TOPIC_ID", cfg.MerchantCreatePublisher.TopicID, "merchant-create")
	cmp(t, "KYM_BUCKET_NAME", cfg.KymBucketName, "bucket-dev")
	cmp(t, "ORGANISATION_SERVICE_URL", cfg.OrganisationServiceURL, "organisation-dev-url")
	cmp(t, "RECIPIENT_SERVICE_URL", cfg.RecipientServiceURL, "recipient-dev-url")
}

func cmp(t *testing.T, field, got, want interface{}) {
	if got != want {
		t.Errorf("unexpected %s got %s; want %s", field, got, want)
	}
}
