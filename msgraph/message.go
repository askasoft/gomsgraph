package msgraph

import (
	"context"
	"time"
)

type Message struct {
	ID                   string          `json:"id"`
	ETag                 string          `json:"etag"`
	MessageType          string          `json:"messageType"`
	Subject              string          `json:"subject"`
	Summary              string          `json:"summary"`
	ReplyToID            string          `json:"replyToId"`
	ChatID               string          `json:"chatId"`
	Importance           string          `json:"importance"`
	Locale               string          `json:"locale"`
	WebURL               string          `json:"webUrl"`
	From                 FromInfo        `json:"from"`
	Body                 BodyInfo        `json:"body"`
	ChannelIdentity      ChannelIdentity `json:"channelIdentity"`
	Attachments          []any           `json:"attachments"`
	Mentions             []Mention       `json:"mentions"`
	Reactions            []any           `json:"reactions"`
	MessageHistory       []any           `json:"messageHistory"`
	CreatedDateTime      time.Time       `json:"createdDateTime"`
	LastModifiedDateTime time.Time       `json:"lastModifiedDateTime"`
	LastEditedDateTime   *time.Time      `json:"lastEditedDateTime"`
	DeletedDateTime      *time.Time      `json:"deletedDateTime"`
}

func (m *Message) String() string {
	return toString(m)
}

type FromInfo struct {
	Application any       `json:"application"`
	Device      any       `json:"device"`
	User        *UserInfo `json:"user"`
}

type UserInfo struct {
	ID               string `json:"id"`
	DisplayName      string `json:"displayName"`
	UserIdentityType string `json:"userIdentityType"`
}

type BodyInfo struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

type ChannelIdentity struct {
	TeamID    string `json:"teamId"`
	ChannelID string `json:"channelId"`
}

type Mention struct {
	ID          int       `json:"id"`
	MentionText string    `json:"mentionText"`
	Mentioned   Mentioned `json:"mentioned"`
}

type Mentioned struct {
	Application  any       `json:"application"`
	Device       any       `json:"device"`
	Conversation any       `json:"conversation"`
	User         *UserInfo `json:"user"`
}

func (gc *GraphClient) getChannelMessagesURL(teamID, channelID string, options ...string) string {
	return gc.Endpoint("/teams/%s/channels/%s/messages", teamID, channelID) + optionsQuery(options...)
}

func (gc *GraphClient) GetChannelMessages(ctx context.Context, teamID, channelID string, options ...string) ([]*Message, string, error) {
	return DoGets[*Message](ctx, gc, gc.getChannelMessagesURL(teamID, channelID, options...))
}

func (gc *GraphClient) ListChannelMessages(ctx context.Context, teamID, channelID string, options ...string) ([]*Message, error) {
	return DoList[*Message](ctx, gc, gc.getChannelMessagesURL(teamID, channelID, options...))
}

func (gc *GraphClient) IterChannelMessages(ctx context.Context, teamID, channelID string, itf func(*Message) error, options ...string) error {
	return DoIter(ctx, gc, gc.getChannelMessagesURL(teamID, channelID, options...), itf)
}

func (gc *GraphClient) getChannelMessageRepliesURL(teamID, channelID, msgID string, options ...string) string {
	return gc.Endpoint("/teams/%s/channels/%s/messages/%s/replies", teamID, channelID, msgID) + optionsQuery(options...)
}

func (gc *GraphClient) GetChannelMessageReplies(ctx context.Context, teamID, channelID, msgID string, options ...string) ([]*Message, string, error) {
	return DoGets[*Message](ctx, gc, gc.getChannelMessageRepliesURL(teamID, channelID, msgID, options...))
}

func (gc *GraphClient) ListChannelMessageReplies(ctx context.Context, teamID, channelID, msgID string, options ...string) ([]*Message, error) {
	return DoList[*Message](ctx, gc, gc.getChannelMessageRepliesURL(teamID, channelID, msgID, options...))
}

func (gc *GraphClient) IterChannelMessageReplies(ctx context.Context, teamID, channelID, msgID string, itf func(*Message) error, options ...string) error {
	return DoIter(ctx, gc, gc.getChannelMessageRepliesURL(teamID, channelID, msgID, options...), itf)
}
