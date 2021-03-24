// +build epgstation

package lemo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/accelforce/lemokon/lemo/epgstation"
	"github.com/accelforce/lemokon/util"
)

func Lemo(ctx context.Context) error {
	host, err := util.CheckEnv("EPGSTATION_HOST")
	if err != nil {
		return err
	}
	secure := util.BoolEnv("EPGSTATION_SECURE")

	port := 80
	if secure {
		port = 443
	}

	if strings.Contains(host, ":") {
		arr := strings.Split(host, ":")
		host = arr[0]
		p, err := strconv.ParseInt(arr[1], 10, 32)
		port = int(p)
		if err != nil {
			return err
		}
	}

	schema := "http"
	if secure {
		schema = "https"
	}
	client := epgstation.NewClient(fmt.Sprintf("%s://%s:%d/api", schema, host, port))

	go observe(ctx, client)

	epgstation.Streaming(ctx, host, port, secure)
	return nil
}

func updateStatus(ctx context.Context, client *epgstation.Client) error {
	r, err := client.GetRecording(ctx)
	if err != nil {
		return err
	}
	ids := make([]int, len(r.Records))
	for i, r := range r.Records {
		ids[i] = r.ID
		if exists, err := epgstation.RecordingExists(r.ID); err != nil {
			return fmt.Errorf("error processing records: %w\n", err)
		} else if exists {
			continue
		}
		r, err := epgstation.NewRecording(r)
		if err != nil {
			return fmt.Errorf("error saving new record: %w\n", err)
		}
		fmt.Printf("New recording found: %d\n", r.ID)
		newRecording <- r
	}

	ended, err := epgstation.FinishRunning(ids)
	if err != nil {
		return err
	}
	for _, r := range ended {
		fmt.Printf("Ended recording found: %d\n", r.ID)
		endedRecording <- &r
	}

	return nil
}

func observe(ctx context.Context, client *epgstation.Client) {
	for {
		select {
		case <-epgstation.Connected:
			if err := updateStatus(ctx, client); err != nil {
				fmt.Printf("Failed to fetch status: %s\n", err)
			}
			break
		case <-epgstation.StatusUpdated:
			if err := updateStatus(ctx, client); err != nil {
				fmt.Printf("Failed to fetch status: %s\n", err)
			}
			break
		case <-epgstation.EncodeUpdated:
			break
		}
	}
}
