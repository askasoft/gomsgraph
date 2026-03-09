package msgraph

import (
	"os"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/log/httplog"
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
	logs.SetLevel(log.LevelDebug)
	logger := logs.GetLogger("MSG")

	gc := &GraphClient{
		TenantID:     tenantID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Transport:    httplog.LoggingRoundTripper(logger),
		Retryer:      NewRetryer(logger, 1, time.Second*3),
	}

	return gc
}
