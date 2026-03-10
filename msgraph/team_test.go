package msgraph

import (
	"context"
	"fmt"
	"testing"

	"github.com/askasoft/pango/test/require"
)

func TestGetTeams(t *testing.T) {
	gc := testNewGraphClient(t)

	teams, err := gc.ListTeams(context.TODO())

	require.NoError(t, err)

	for i, t := range teams {
		fmt.Printf("* #%d %s: %s\n", i, t.ID, t.DisplayName)
	}
}
