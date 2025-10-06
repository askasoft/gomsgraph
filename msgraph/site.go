package msgraph

import (
	"context"
	"time"
)

type Site struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	WebURL          string    `json:"webUrl"`
	DisplayName     string    `json:"displayName"`
	IsPersonalSite  bool      `json:"isPersonalSite"`
	CreatedDateTime time.Time `json:"createdDateTime"`
}

func (s *Site) String() string {
	return toString(s)
}

func (gc *GraphClient) GetSite(ctx context.Context, siteID string) (*Site, error) {
	url := gc.Endpoint("/sites/%s", siteID)
	site := &Site{}
	err := gc.DoGet(ctx, url, site)
	return site, err
}

func (gc *GraphClient) GetSites(ctx context.Context) ([]*Site, string, error) {
	return DoGets[*Site](ctx, gc, gc.Endpoint("/sites"))
}

func (gc *GraphClient) ListSites(ctx context.Context) ([]*Site, error) {
	return DoList[*Site](ctx, gc, gc.Endpoint("/sites"))
}

func (gc *GraphClient) IterSites(ctx context.Context, itf func(*Site) error) error {
	return DoIter(ctx, gc, gc.Endpoint("/sites"), itf)
}

func (gc *GraphClient) GetSubSites(ctx context.Context, siteID string) ([]*Site, string, error) {
	return DoGets[*Site](ctx, gc, gc.Endpoint("/sites/%s/sites", siteID))
}

func (gc *GraphClient) ListSubSites(ctx context.Context, siteID string) (sites []*Site, err error) {
	return DoList[*Site](ctx, gc, gc.Endpoint("/sites/%s/sites", siteID))
}

func (gc *GraphClient) IterSubSites(ctx context.Context, siteID string, itf func(*Site) error) error {
	return DoIter(ctx, gc, gc.Endpoint("/sites/%s/sites", siteID), itf)
}
