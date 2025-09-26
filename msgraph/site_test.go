package msgraph

import (
	"context"
	"fmt"
	"testing"

	"github.com/askasoft/pango/test/require"
)

func TestGetSite(t *testing.T) {
	gc := testNewGraphClient(t)

	s, err := gc.GetSite(context.TODO(), "root")

	require.NoError(t, err)

	fmt.Printf("* #%d %s: %s - %s\n", 0, s.ID, s.Name, s.WebURL)

	sss, err := gc.GetSubSites(context.TODO(), s.ID)
	require.NoError(t, err)
	for j, ss := range sss {
		fmt.Printf("    * #%d %s: %s - %s\n", j, ss.ID, ss.Name, ss.WebURL)
	}
}

func TestGetSites(t *testing.T) {
	gc := testNewGraphClient(t)

	sites, err := gc.GetSites(context.TODO())

	require.NoError(t, err)

	for i, s := range sites {
		fmt.Printf("* #%d %s: %s - %s\n", i, s.ID, s.Name, s.WebURL)

		sss, err := gc.GetSubSites(context.TODO(), s.ID)
		require.NoError(t, err)
		for j, ss := range sss {
			fmt.Printf("    * #%d %s: %s - %s\n", j, ss.ID, ss.Name, ss.WebURL)
		}
	}
}
