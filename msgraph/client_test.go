package msgraph

import (
	"os"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
)

func testNewGraphClient(t *testing.T) *GraphClient {
	tenantID := os.Getenv("MSG_TENANT_ID")
	if tenantID == "" {
		t.Skip("MSG_TENANT_ID not set")
		return nil
	}

	clientID := os.Getenv("MSG_CLIENT_ID")
	if tenantID == "" {
		t.Skip("MSG_CLIENT_ID not set")
		return nil
	}

	clientSecret := os.Getenv("MSG_CLIENT_SECRET")
	if clientSecret == "" {
		t.Skip("MSG_CLIENT_SECRET not set")
		return nil
	}

	logs := log.NewLog()
	logs.SetLevel(log.LevelInfo)
	gc := &GraphClient{
		TenantID:     tenantID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Logger:       logs.GetLogger("MSG"),
		MaxRetries:   1,
		RetryAfter:   time.Second * 3,
	}

	return gc
}
