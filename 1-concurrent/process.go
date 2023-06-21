package concurrent

import (
	"context"
	"io"
	"time"

	"golang.org/x/sync/errgroup"
)

type Meeting struct {
	ID          string
	Topic       string
	Start       time.Time
	DownloadURL string
}

type Participant struct {
	ID    string
	Name  string
	Email string
}

// ListMeetingsParams defines parameters for ListMeetings.
type ListMeetingsParams struct {
	From *time.Time
	To   *time.Time
}

//go:generate go run github.com/matryer/moq/... -out autogen_client.go . ClientInterface
type ClientInterface interface {
	// Meetings must be ordered by Start time
	ListMeetings(ctx context.Context, params *ListMeetingsParams) ([]Meeting, error)
	DownloadMeeting(ctx context.Context, url string) (io.ReadCloser, error)
	GetMeetingParticipants(ctx context.Context, meetingID string) ([]Participant, error)
}

type CreateMeetingDatumArguments struct {
	Topic        string
	Content      io.ReadCloser
	Participants []Participant
}

//go:generate go run github.com/matryer/moq/... -out autogen_store.go . StoreInterface
type StoreInterface interface {
	CreateMeetingDatum(ctx context.Context, args CreateMeetingDatumArguments) error
}

type Config struct {
	MeetingConcurrency int
}

type Processor struct {
	Client ClientInterface
	Store  StoreInterface
	Cfg    Config
}

func (p *Processor) Process(ctx context.Context) (time.Time, error) {
	// Source
	meetings, err := p.Client.ListMeetings(ctx, nil)
	if err != nil {
		return time.Time{}, err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(p.Cfg.MeetingConcurrency)

	tms := make(chan time.Time)

	// Stage 1
	go func() {
		defer close(tms)
		for _, meeting := range meetings {
			meeting := meeting
			g.Go(func() error {
				args, err := p.TransformToDatum(ctx, meeting)
				if err != nil {
					return err
				}

				if err := p.Store.CreateMeetingDatum(ctx, args); err != nil {
					return err
				}

				select {
				case tms <- meeting.Start:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			})
		}
		g.Wait()
	}()

	// Sink
	var t time.Time
	for tm := range tms {
		if tm.After(t) {
			t = tm
		}
	}

	return t, g.Wait()
}

func (p *Processor) TransformToDatum(ctx context.Context, m Meeting) (CreateMeetingDatumArguments, error) {
	args := CreateMeetingDatumArguments{
		Topic: m.Topic,
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		args.Content, err = p.Client.DownloadMeeting(ctx, m.DownloadURL)
		return err
	})

	g.Go(func() error {
		var err error
		args.Participants, err = p.Client.GetMeetingParticipants(ctx, m.ID)
		return err
	})

	return args, g.Wait()
}
