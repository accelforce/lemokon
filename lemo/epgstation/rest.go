package epgstation

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
)

type (
	Client struct {
		c *resty.Client
	}

	Time time.Time

	Record struct {
		ID        int    `json:"id"`
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
