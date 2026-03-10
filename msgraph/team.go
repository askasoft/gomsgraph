package msgraph

import (
	"context"
	"time"
)

type Team struct {
	ID                          string         `json:"id"`
	DisplayName                 string         `json:"displayName"`
	Description                 string         `json:"description"`
	Visibility                  string         `json:"visibility"`
	WebURL                      string         `json:"webUrl"`
	IsArchived                  bool           `json:"isArchived"`
	IsMembershipLimitedToOwners bool           `json:"isMembershipLimitedToOwners"`
	MemberSettings              map[string]any `json:"memberSettings"`
	GuestSettings               map[string]any `json:"guestSettings"`
	MessagingSettings           map[string]any `json:"messagingSettings"`
	FunSettings                 map[string]any `json:"funSettings"`
	DiscoverySettings           map[string]any `json:"discoverySettings"`
	TagSettings                 map[string]any `json:"tagSettings"`
	Summary                     map[string]any `json:"summary"`
}

func (t *Team) String() string {
	return toString(t)
}

func (gc *GraphClient) GetTeam(ctx context.Context, teamID string, options ...string) (*Team, error) {
	url := gc.Endpoint("/teams/%s", teamID) + optionsQuery(options...)
	team := &Team{}
	err := gc.DoGet(ctx, url, team)
	return team, err
}

func (gc *GraphClient) getTeamsURL(options ...string) string {
	return gc.Endpoint("/teams") + optionsQuery(options...)
}

func (gc *GraphClient) GetTeams(ctx context.Context, options ...string) ([]*Team, string, error) {
	return DoGets[*Team](ctx, gc, gc.getTeamsURL(options...))
}

func (gc *GraphClient) ListTeams(ctx context.Context, options ...string) ([]*Team, error) {
	return DoList[*Team](ctx, gc, gc.getTeamsURL(options...))
}

func (gc *GraphClient) IterTeams(ctx context.Context, itf func(*Team) error, options ...string) error {
	return DoIter(ctx, gc, gc.getTeamsURL(options...), itf)
}

type Channel struct {
	ID              string    `json:"id"`
	DisplayName     string    `json:"displayName"`
	Description     string    `json:"description"`
	WebURL          string    `json:"webUrl"`
	IsArchived      bool      `json:"isArchived"`
	CreatedDateTime time.Time `json:"createdDateTime"`
}

func (c *Channel) String() string {
	return toString(c)
}

func (gc *GraphClient) getChannelsURL(teamID string, options ...string) string {
	return gc.Endpoint("/teams/%s/channels", teamID) + optionsQuery(options...)
}

func (gc *GraphClient) GetChannels(ctx context.Context, teamID string, options ...string) ([]*Channel, string, error) {
	return DoGets[*Channel](ctx, gc, gc.getChannelsURL(teamID, options...))
}

func (gc *GraphClient) ListChannels(ctx context.Context, teamID string, options ...string) ([]*Channel, error) {
	return DoList[*Channel](ctx, gc, gc.getChannelsURL(teamID, options...))
}

func (gc *GraphClient) IterChannels(ctx context.Context, teamID string, itf func(*Channel) error, options ...string) error {
	return DoIter(ctx, gc, gc.getChannelsURL(teamID, options...), itf)
}

func (gc *GraphClient) getAllChannelsURL(teamID string, options ...string) string {
	return gc.Endpoint("/teams/%s/allChannels", teamID) + optionsQuery(options...)
}

func (gc *GraphClient) GetAllChannels(ctx context.Context, teamID string, options ...string) ([]*Channel, string, error) {
	return DoGets[*Channel](ctx, gc, gc.getAllChannelsURL(teamID, options...))
}

func (gc *GraphClient) ListAllChannels(ctx context.Context, teamID string, options ...string) ([]*Channel, error) {
	return DoList[*Channel](ctx, gc, gc.getAllChannelsURL(teamID, options...))
}

func (gc *GraphClient) IterAllChannels(ctx context.Context, teamID string, itf func(*Channel) error, options ...string) error {
	return DoIter(ctx, gc, gc.getAllChannelsURL(teamID, options...), itf)
}

func (gc *GraphClient) getIncomingChannelsURL(teamID string, options ...string) string {
	return gc.Endpoint("/teams/%s/incomingChannels", teamID) + optionsQuery(options...)
}

func (gc *GraphClient) GetIncomingChannels(ctx context.Context, teamID string, options ...string) ([]*Channel, string, error) {
	return DoGets[*Channel](ctx, gc, gc.getIncomingChannelsURL(teamID, options...))
}

func (gc *GraphClient) ListIncomingChannels(ctx context.Context, teamID string, options ...string) ([]*Channel, error) {
	return DoList[*Channel](ctx, gc, gc.getIncomingChannelsURL(teamID, options...))
}

func (gc *GraphClient) IterIncomingChannels(ctx context.Context, teamID string, itf func(*Channel) error, options ...string) error {
	return DoIter(ctx, gc, gc.getIncomingChannelsURL(teamID, options...), itf)
}
