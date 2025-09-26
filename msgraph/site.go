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

type Sites struct {
	Values []*Site `json:"value"`
}

func (gc *GraphClient) GetSite(ctx context.Context, siteID string) (*Site, error) {
	url := gc.Endpoint("/sites/%s", siteID)
	site := &Site{}
	err := gc.DoGet(ctx, url, site)
	return site, err
}

func (gc *GraphClient) GetSites(ctx context.Context) ([]*Site, error) {
	url := gc.Endpoint("/sites")
	sites := &Sites{}
	err := gc.DoGet(ctx, url, sites)
	return sites.Values, err
}

func (gc *GraphClient) GetSubSites(ctx context.Context, siteID string) ([]*Site, error) {
	url := gc.Endpoint("/sites/%s/sites", siteID)
	sites := &Sites{}
	err := gc.DoGet(ctx, url, sites)
	return sites.Values, err
}
