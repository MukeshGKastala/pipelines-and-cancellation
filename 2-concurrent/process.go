package concurrent

import (
	"context"
	"io"
	"time"

	"github.com/sourcegraph/conc/stream"
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

// ListPaginatedMeetingsParams defines parameters for ListPaginatedMeetings.
type ListPaginatedMeetingsParams struct {
	From *time.Time
	To   *time.Time

	NextPageToken *string
	PageSize      *int
}

type ListPaginatedMeetingsResponse struct {
	NextPageToken string
	Meetings      []Meeting
}

//go:generate go run github.com/matryer/moq/... -out autogen_client.go . ClientInterface
type ClientInterface interface {
	// Meetings must be ordered by Start time
	ListPaginatedMeetings(ctx context.Context, params *ListPaginatedMeetingsParams) (ListPaginatedMeetingsResponse, error)
	DownloadMeeting(ctx context.Context, url string) (io.ReadCloser, error)
	GetMeetingParticipants(ctx context.Context, meetingID string) ([]Participant, error)
}

type CreateMeetingDatumArguments struct {
	Topic        string
	Start        time.Time
	Content      io.ReadCloser
	Participants []Participant
}

//go:generate go run github.com/matryer/moq/... -out autogen_store.go . StoreInterface
type StoreInterface interface {
	CreateMeetingDatum(ctx context.Context, args CreateMeetingDatumArguments) error
}

type Config struct {
	TransformerConcurrency int
	UploaderConcurrency    int
}

type Processor struct {
	Client ClientInterface
	Store  StoreInterface
	Cfg    Config
}

func (p *Processor) Process(ctx context.Context) (time.Time, error) {
	g, ctx := errgroup.WithContext(ctx)

	// Create channels to connect the stages
	meetings := make(chan Meeting)
	datums := make(chan CreateMeetingDatumArguments)
	tms := make(chan time.Time)

	// Source
	g.Go(func() error {
		return p.produce(ctx, meetings)
	})

	// Stage 2
	g.Go(func() error {
		return <-p.transform(ctx, meetings, datums)
	})

	// Stage 3
	g.Go(func() error {
		return <-p.upload(ctx, datums, tms)
	})

	// Sink
	// Track last successfully uploaded meeting's start time
	var t time.Time
	for tm := range tms {
		t = tm
	}

	return t, g.Wait()
}

func (p *Processor) produce(ctx context.Context, out chan<- Meeting) error {
	defer close(out)

	var nextPageToken string

	// Handle pagination
	for {
		resp, err := p.Client.ListPaginatedMeetings(ctx, &ListPaginatedMeetingsParams{
			NextPageToken: &nextPageToken,
		})
		if err != nil {
			return err
		}

		for _, meeting := range resp.Meetings {
			select {
			case out <- meeting:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if resp.NextPageToken == "" {
			break
		}

		nextPageToken = resp.NextPageToken
	}

	return nil
}

func (p *Processor) transform(ctx context.Context, in <-chan Meeting, out chan<- CreateMeetingDatumArguments) <-chan error {
	errorC := make(chan error)

	go func() {
		defer close(errorC)
		defer close(out)
		s := stream.New().WithMaxGoroutines(p.Cfg.TransformerConcurrency)
		for m := range in {
			m := m
			s.Go(func() stream.Callback {
				enriched, err := p.enrich(ctx, m)
				return func() {
					if err != nil {
						select {
						// No-op if we've already sent an error
						case errorC <- err:
							// Disable sends
							out = nil
						default:
						}
					} else {
						select {
						case out <- enriched:
						case <-ctx.Done():
							return
						}
					}
				}
			})
		}
		s.Wait()
		select {
		case errorC <- nil:
		default:
		}
	}()

	return errorC
}

func (p *Processor) enrich(ctx context.Context, m Meeting) (CreateMeetingDatumArguments, error) {
	args := CreateMeetingDatumArguments{
		Topic: m.Topic,
		Start: m.Start,
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

func (p *Processor) upload(ctx context.Context, in <-chan CreateMeetingDatumArguments, out chan<- time.Time) <-chan error {
	errorC := make(chan error)

	go func() {
		defer close(errorC)
		defer close(out)
		s := stream.New().WithMaxGoroutines(p.Cfg.UploaderConcurrency)
		for args := range in {
			args := args
			s.Go(func() stream.Callback {
				err := p.Store.CreateMeetingDatum(ctx, args)
				return func() {
					if err != nil {
						select {
						// No-op if we've already sent an error
						case errorC <- err:
							// Disable sends
							out = nil
						default:
						}
					} else {
						select {
						case out <- args.Start:
						case <-ctx.Done():
							return
						}
					}
				}
			})
		}
		s.Wait()
		select {
		case errorC <- nil:
		default:
		}
	}()

	return errorC
}
