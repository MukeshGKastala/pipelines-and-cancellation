package concurrent_test

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	concurrent "example.com/pipelines-and-cancellation/1-bounded-concurrency"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	const (
		// 200 milliseconds
		maxNetworkLatency = 200

		maxNumberOfMeetings = 250

		problematicMeetingID = "113"
	)

	p := concurrent.Processor{
		Client: &concurrent.ClientInterfaceMock{
			ListPaginatedMeetingsFunc: func(ctx context.Context, params *concurrent.ListPaginatedMeetingsParams) (concurrent.ListPaginatedMeetingsResponse, error) {
				// Simulate I/O
				time.Sleep(time.Duration(rand.Intn(maxNetworkLatency)) * time.Millisecond)

				var begin int
				if params.NextPageToken != nil {
					begin, _ = strconv.Atoi(*params.NextPageToken)
				}

				// Cap number of meetings
				if begin >= maxNumberOfMeetings {
					return concurrent.ListPaginatedMeetingsResponse{}, nil
				}

				end := begin + 10
				if params.PageSize != nil {
					end = begin + *params.PageSize
				}

				select {
				case <-ctx.Done():
					return concurrent.ListPaginatedMeetingsResponse{}, ctx.Err()
				default:
					return concurrent.ListPaginatedMeetingsResponse{
						NextPageToken: strconv.Itoa(end),
						Meetings:      generateMeetings(begin, end),
					}, nil
				}
			},
			DownloadMeetingFunc: func(ctx context.Context, url string) (io.ReadCloser, error) {
				// Simulate I/O
				time.Sleep(time.Duration(rand.Intn(maxNetworkLatency)) * time.Millisecond)

				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
					return os.Open(url)
				}
			},
			GetMeetingParticipantsFunc: func(ctx context.Context, meetingID string) ([]concurrent.Participant, error) {
				// Simulate I/O
				time.Sleep(time.Duration(rand.Intn(maxNetworkLatency)) * time.Millisecond)

				if meetingID == problematicMeetingID {
					// Simulate retries
					time.Sleep(time.Duration(rand.Intn(maxNetworkLatency*2)) * time.Millisecond)
					err := ctx.Err()
					if err == nil {
						err = fmt.Errorf("forced participants error for %s", meetingID)
					}
					return nil, err
				}

				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
					return []concurrent.Participant{
						{
							ID:    uuid.NewString(),
							Name:  "John Doe",
							Email: "john.doe@example.com",
						},
						{
							ID:    uuid.NewString(),
							Name:  "Jane Doe",
							Email: "jane.doe@example.com",
						},
					}, nil
				}
			},
		},
		Store: &concurrent.StoreInterfaceMock{
			CreateMeetingDatumFunc: func(ctx context.Context, args concurrent.CreateMeetingDatumArguments) error {
				// Simulate I/O
				time.Sleep(time.Duration(rand.Intn(maxNetworkLatency)) * time.Millisecond)
				io.Copy(io.Discard, args.Content)
				args.Content.Close()
				return ctx.Err()
			},
		},
		Cfg: concurrent.Config{
			TransformerConcurrency: 3,
			UploaderConcurrency:    5,
		},
	}

	got, gerr := p.Process(context.Background())

	assert.ErrorContains(t, gerr, problematicMeetingID)
	days, _ := strconv.Atoi(problematicMeetingID)
	assert.Less(t, got, time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, days)) // IDs start at 1
	assert.Equal(t, 2, runtime.NumGoroutine())
}

func generateMeetings(begin, end int) []concurrent.Meeting {
	l := end - begin
	meetings := make([]concurrent.Meeting, 0, l)

	start := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, begin)

	for i := begin + 1; i <= end; i++ {
		id := strconv.Itoa(i)
		meetings = append(meetings, concurrent.Meeting{
			ID:          id,
			Topic:       fmt.Sprintf("Meeting %s", id),
			Start:       start,
			DownloadURL: filepath.Join("../testdata", "test.mp4"),
		})
		start = start.AddDate(0, 0, 1)
	}

	return meetings
}
