package msgraph

import (
	"context"
	"time"
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

func (gc *GraphClient) getSiteDrivesURL(siteID string, options ...string) string {
	return gc.Endpoint("/sites/%s/drives", siteID) + optionsQuery(options...)
}

func (gc *GraphClient) GetSiteDrives(ctx context.Context, siteID string, options ...string) ([]*Drive, string, error) {
	return DoGets[*Drive](ctx, gc, gc.getSiteDrivesURL(siteID, options...))
}

func (gc *GraphClient) ListSiteDrives(ctx context.Context, siteID string, options ...string) ([]*Drive, error) {
	return DoList[*Drive](ctx, gc, gc.getSiteDrivesURL(siteID, options...))
}

func (gc *GraphClient) IterSiteDrives(ctx context.Context, siteID string, itf func(*Drive) error, options ...string) error {
	return DoIter(ctx, gc, gc.getSiteDrivesURL(siteID, options...), itf)
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

func (di *DriveItem) String() string {
	return toString(di)
}

func (gc *GraphClient) GetDriveRoot(ctx context.Context, driveID string, options ...string) (*DriveItem, error) {
	url := gc.Endpoint("/drives/%s/root", driveID) + optionsQuery(options...)

	item := &DriveItem{}
	err := gc.DoGet(ctx, url, item)
	return item, err
}

func (gc *GraphClient) getDriveItemChildrenURL(driveID, itemID string, options ...string) string {
	return gc.Endpoint("/drives/%s/items/%s/children", driveID, itemID) + optionsQuery(options...)
}

func (gc *GraphClient) GetDriveItemChildren(ctx context.Context, driveID, itemID string, options ...string) ([]*DriveItem, string, error) {
	return DoGets[*DriveItem](ctx, gc, gc.getDriveItemChildrenURL(driveID, itemID, options...))
}

func (gc *GraphClient) ListDriveItemChildren(ctx context.Context, driveID, itemID string, options ...string) ([]*DriveItem, error) {
	return DoList[*DriveItem](ctx, gc, gc.getDriveItemChildrenURL(driveID, itemID, options...))
}

func (gc *GraphClient) IterDriveItemChildren(ctx context.Context, driveID, itemID string, itf func(*DriveItem) error, options ...string) error {
	return DoIter(ctx, gc, gc.getDriveItemChildrenURL(driveID, itemID, options...), itf)
}

func (gc *GraphClient) GetDriveItemContent(ctx context.Context, driveID, itemID string) ([]byte, error) {
	return gc.DoDownload(ctx, gc.Endpoint("/drives/%s/items/%s/content", driveID, itemID))
}

func (gc *GraphClient) SaveDriveItemContent(ctx context.Context, driveID, itemID string, savePath string) error {
	return gc.DoSaveFile(ctx, gc.Endpoint("/drives/%s/items/%s/content", driveID, itemID), savePath)
}
