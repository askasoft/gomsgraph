package msgraph

import (
	"context"
	"net/url"
	"time"

	"github.com/askasoft/pango/str"
)

type Drive struct {
	ID                   string     `json:"id"`
	Name                 string     `json:"name"`
	WebURL               string     `json:"webUrl"`
	Description          string     `json:"description"`
	DriveType            string     `json:"driveType"`
	Root                 *DriveItem `json:"root"`
	CreatedDateTime      time.Time  `json:"createdDateTime"`
	LastModifiedDateTime time.Time  `json:"lastModifiedDateTime"`
}

func (d *Drive) String() string {
	return toString(d)
}

type Drives struct {
	Values []*Drive `json:"value"`
}

func (gc *GraphClient) GetSiteDrives(ctx context.Context, siteID string, expand ...string) ([]*Drive, error) {
	u := gc.Endpoint("/sites/%s/drives", siteID)
	if len(expand) > 0 {
		u += "?$expand=" + url.QueryEscape(str.Join(expand, " "))
	}

	drives := &Drives{}
	err := gc.DoGet(ctx, u, drives)
	return drives.Values, err
}

type DriveFolder struct {
	ChildCount int `json:"childCount"`
}

type DriveFile struct {
	MimeType string `json:"mimeType"`
}

type DriveItem struct {
	ID                   string         `json:"id"`
	Name                 string         `json:"name"`
	ETag                 string         `json:"eTag"`
	WebURL               string         `json:"webUrl"`
	Size                 int64          `json:"size"`
	Folder               *DriveFolder   `json:"folder"`
	File                 *DriveFile     `json:"file"`
	ListItem             map[string]any `json:"listItem"`
	CreatedDateTime      time.Time      `json:"createdDateTime"`
	LastModifiedDateTime time.Time      `json:"lastModifiedDateTime"`
}

func (i *DriveItem) String() string {
	return toString(i)
}

type DriveItems struct {
	Values []*DriveItem `json:"value"`
}

func (gc *GraphClient) GetDriveRoot(ctx context.Context, driveID string) (*DriveItem, error) {
	u := gc.Endpoint("/drives/%s/root", driveID)
	item := &DriveItem{}
	err := gc.DoGet(ctx, u, item)
	return item, err
}

func (gc *GraphClient) GetDriveItemChildren(ctx context.Context, driveID, itemID string, expand ...string) ([]*DriveItem, error) {
	u := gc.Endpoint("/drives/%s/items/%s/children", driveID, itemID)
	if len(expand) > 0 {
		u += "?$expand=" + url.QueryEscape(str.Join(expand, " "))
	}

	items := &DriveItems{}
	err := gc.DoGet(ctx, u, items)
	return items.Values, err
}

func (gc *GraphClient) GetDriveItemContent(ctx context.Context, driveID, itemID string) ([]byte, error) {
	u := gc.Endpoint("/drives/%s/items/%s/content", driveID, itemID)

	return gc.DoDownload(ctx, u)
}

func (gc *GraphClient) SaveDriveItemContent(ctx context.Context, driveID, itemID string, savePath string) error {
	u := gc.Endpoint("/drives/%s/items/%s/content", driveID, itemID)

	return gc.DoSaveFile(ctx, u, savePath)
}
