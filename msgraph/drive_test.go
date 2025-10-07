package msgraph

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/askasoft/pango/test/require"
)

func TestGetRootSiteDrives(t *testing.T) {
	gc := testNewGraphClient(t)

	drives, err := gc.ListSiteDrives(context.TODO(), "root", "$expand", "root")
	require.NoError(t, err)
	for j, d := range drives {
		fmt.Printf("    * #%d %s: %s - %s\n", j, d.ID, d.Name, d.WebURL)
		require.NotNil(t, d.Root)
	}
}

func TestGetSiteDrives(t *testing.T) {
	gc := testNewGraphClient(t)

	spts := os.Getenv("MSG_TEST_SITE")
	if spts == "" {
		t.Skip("MSG_TEST_SITE not set")
		return
	}

	s, err := gc.GetSite(context.TODO(), spts)

	require.NoError(t, err)

	fmt.Printf("* %s: %s - %s\n", s.ID, s.DisplayName, s.WebURL)

	drives, err := gc.ListSiteDrives(context.TODO(), s.ID, "$expand", "root")
	require.NoError(t, err)
	for j, d := range drives {
		fmt.Printf("    * #%d %s: %s - %s\n", j, d.ID, d.Name, d.WebURL)
		require.NotNil(t, d.Root)

		fmt.Printf("    * #%d ROOT: %s\n", j, d.Root.ID)
		testGetDriveItemChildren(t, gc, 2, d.ID, d.Root.ID)
	}
}

func testGetDriveItemChildren(t *testing.T, gc *GraphClient, indent int, driveID, itemID string) {
	items, err := gc.ListDriveItemChildren(context.TODO(), driveID, itemID, "$expand", "listItem($expand=fields)")
	require.NoError(t, err)

	for k, f := range items {
		fmt.Printf("%*s #%d %s: %s - %d\n", indent*4, "*", k, f.ID, f.Name, f.Size)
		// fmt.Println(f.String())

		if f.Folder != nil {
			testGetDriveItemChildren(t, gc, indent+1, driveID, f.ID)
		}
	}
}
