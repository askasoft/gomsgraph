package msgraph

import (
	"context"
	"fmt"
	"testing"

	"github.com/askasoft/pango/test/require"
)

func TestGetSites(t *testing.T) {
	gc := testNewGraphClient(t)

	sites, err := gc.GetSites(context.TODO())

	require.NoError(t, err)

	for i, s := range sites {
		fmt.Printf("* #%d %s: %s - %s\n", i, s.ID, s.Name, s.WebURL)
	}
}
