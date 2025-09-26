package msgraph

import (
	"context"
	"fmt"
	"testing"

	"github.com/askasoft/pango/test/require"
)

var testSiteDisplayName = "SharePoint連携テスト - テストチャネル"

func TestGetRootSiteDrives(t *testing.T) {
	gc := testNewGraphClient(t)

	drives, err := gc.GetSiteDrives(context.TODO(), "root", "root")
	require.NoError(t, err)
	for j, d := range drives {
		fmt.Printf("    * #%d %s: %s - %s\n", j, d.ID, d.Name, d.WebURL)
		require.NotNil(t, d.Root)

		if d.Name == "ドキュメント" {
			testGetDriveItemChildren(t, gc, 2, d.ID, d.Root.ID)
		}
	}
}

func TestGetSiteDrives(t *testing.T) {
	gc := testNewGraphClient(t)

	sites, err := gc.GetSites(context.TODO())

	require.NoError(t, err)

	for i, s := range sites {
		if s.DisplayName == testSiteDisplayName {
			fmt.Printf("* #%d %s: %s - %s\n", i, s.ID, s.DisplayName, s.WebURL)

			drives, err := gc.GetSiteDrives(context.TODO(), s.ID, "root")
			require.NoError(t, err)
			for j, d := range drives {
				fmt.Printf("    * #%d %s: %s - %s\n", j, d.ID, d.Name, d.WebURL)
				require.NotNil(t, d.Root)

				testGetDriveItemChildren(t, gc, 2, d.ID, d.Root.ID)
			}
			return
		}
	}
}

func testGetDriveItemChildren(t *testing.T, gc *GraphClient, indent int, driveID, itemID string) {
	items, err := gc.GetDriveItemChildren(context.TODO(), driveID, itemID, "listItem($expand=fields)")
	require.NoError(t, err)

	for k, f := range items {
		fmt.Printf("%*s #%d %s: %s - %d\n", indent*4, "*", k, f.ID, f.Name, f.Size)
		// fmt.Println(f.String())

		if f.Folder != nil {
			testGetDriveItemChildren(t, gc, indent+1, driveID, f.ID)
		}
	}
}
