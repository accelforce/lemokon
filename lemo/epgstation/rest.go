package epgstation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type (
	Client struct {
		c *resty.Client
	}

	Time time.Time

	Channel struct {
		ID                 int    `json:"id"`
		Name               string `json:"halfWidthName"`
		HasLogoData        bool   `json:"hasLogoData"`
		RemoteControlKeyID int    `json:"remoteControlKeyId"`
	}

	Record struct {
		ID        int    `json:"id"`
		ChannelID int    `json:"channelId"`
		Name      string `json:"name"`
		StartedAt Time   `json:"startAt"`
		EndsAt    Time   `json:"endAt"`
	}

	Records struct {
		Total   int      `json:"total"`
		Records []Record `json:"records"`
	}
)

func NewClient(url string) *Client {
	return &Client{
		c: resty.New().SetHostURL(url),
	}
}

func (t *Time) UnmarshalJSON(b []byte) error {
	var r time.Duration
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}
	*t = Time(time.Unix(0, int64(r*time.Millisecond)))
	return nil
}

func (client *Client) GetChannel(ctx context.Context, ID int) (*Channel, error) {
	r, err := client.c.R().
		SetContext(ctx).
		SetResult([]Channel{}).
		Get("/channels")
	if err != nil {
		return nil, err
	}
	for _, c := range *r.Result().(*[]Channel) {
		if c.ID == ID {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("could not find channel: %d", ID)
}

func (client *Client) GetRecording(ctx context.Context) (*Records, error) {
	r, err := client.c.R().
		SetContext(ctx).
		SetResult(Records{}).
		SetQueryParams(map[string]string{
			"isHalfWidth": "true",
		}).
		Get("/recording")
	if err != nil {
		return nil, err
	}
	return r.Result().(*Records), nil
}
