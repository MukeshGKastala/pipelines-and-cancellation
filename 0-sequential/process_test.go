package sequential_test

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	sequential "example.com/pipelines-and-cancellation/0-sequential"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	// 200 milliseconds
	const maxNetworkLatency = 200

	p := sequential.Processor{
		Client: &sequential.ClientInterfaceMock{
			ListMeetingsFunc: func(_ context.Context, params *sequential.ListMeetingsParams) ([]sequential.Meeting, error) {
				time.Sleep(time.Duration(rand.Intn(maxNetworkLatency)) * time.Millisecond)
				return []sequential.Meeting{
					{
						ID:          "1",
						Topic:       "First Meeting",
						Start:       time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
						DownloadURL: filepath.Join("../testdata", "test.mp4"),
					},
					{
						ID:          "2",
						Topic:       "Second Meeting",
						Start:       time.Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC),
						DownloadURL: filepath.Join("../testdata", "test.mp4"),
					},
					{
						ID:          "3",
						Topic:       "Third Meeting",
						Start:       time.Date(2023, time.January, 3, 0, 0, 0, 0, time.UTC),
						DownloadURL: filepath.Join("../testdata", "test.mp4"),
					},
					{
						ID:          "4",
						Topic:       "Fourth Meeting",
						Start:       time.Date(2023, time.January, 4, 0, 0, 0, 0, time.UTC),
						DownloadURL: filepath.Join("../testdata", "test.mp4"),
					},
				}, nil
			},
			DownloadMeetingFunc: func(_ context.Context, url string) (io.ReadCloser, error) {
				// Simulate I/O
				time.Sleep(time.Duration(rand.Intn(maxNetworkLatency)) * time.Millisecond)
				return os.Open(url)
			},
			GetMeetingParticipantsFunc: func(_ context.Context, meetingID string) ([]sequential.Participant, error) {
				// Simulate I/O
				time.Sleep(time.Duration(rand.Intn(maxNetworkLatency)) * time.Millisecond)
				if meetingID == "3" {
					// Simulate retries
					time.Sleep(time.Duration(rand.Intn(maxNetworkLatency*2)) * time.Millisecond)
					return nil, fmt.Errorf("forced participants error for %s", meetingID)
				}

				return []sequential.Participant{
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
			},
		},
		Store: &sequential.StoreInterfaceMock{
			CreateMeetingDatumFunc: func(_ context.Context, args sequential.CreateMeetingDatumArguments) error {
				// Simulate I/O
				time.Sleep(time.Duration(rand.Intn(maxNetworkLatency)) * time.Millisecond)
				io.Copy(io.Discard, args.Content)
				args.Content.Close()
				return nil
			},
		},
	}

	got, gerr := p.Process(context.Background())

	assert.ErrorContains(t, gerr, "forced participants error")
	assert.Equal(t, time.Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC), got)
	calls := p.Store.(*sequential.StoreInterfaceMock).CreateMeetingDatumCalls()
	assert.Len(t, calls, 2)
}
