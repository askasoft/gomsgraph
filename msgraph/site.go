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

func (gc *GraphClient) GetSite(ctx context.Context, siteID string, options ...string) (*Site, error) {
	url := gc.Endpoint("/sites/%s", siteID) + optionsQuery(options...)

	site := &Site{}
	err := gc.DoGet(ctx, url, site)
	return site, err
}

func (gc *GraphClient) getSitesURL(options ...string) string {
	return gc.Endpoint("/sites") + optionsQuery(options...)
}

func (gc *GraphClient) GetSites(ctx context.Context, options ...string) ([]*Site, string, error) {
	return DoGets[*Site](ctx, gc, gc.getSitesURL(options...))
}

func (gc *GraphClient) ListSites(ctx context.Context, options ...string) ([]*Site, error) {
	return DoList[*Site](ctx, gc, gc.getSitesURL(options...))
}

func (gc *GraphClient) IterSites(ctx context.Context, itf func(*Site) error, options ...string) error {
	return DoIter(ctx, gc, gc.getSitesURL(options...), itf)
}

func (gc *GraphClient) getSubSitesURL(siteID string, options ...string) string {
	return gc.Endpoint("/sites/%s/sites", siteID) + optionsQuery(options...)
}

func (gc *GraphClient) GetSubSites(ctx context.Context, siteID string, options ...string) ([]*Site, string, error) {
	return DoGets[*Site](ctx, gc, gc.getSubSitesURL(siteID, options...))
}

func (gc *GraphClient) ListSubSites(ctx context.Context, siteID string, options ...string) (sites []*Site, err error) {
	return DoList[*Site](ctx, gc, gc.getSubSitesURL(siteID, options...))
}

func (gc *GraphClient) IterSubSites(ctx context.Context, siteID string, itf func(*Site) error, options ...string) error {
	return DoIter(ctx, gc, gc.getSubSitesURL(siteID, options...), itf)
}
