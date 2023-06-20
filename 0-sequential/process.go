package sequential

import (
	"context"
	"io"
	"time"
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

type Processor struct {
	Client ClientInterface
	Store  StoreInterface
}

func (p *Processor) Process(ctx context.Context) (t time.Time, _ error) {
	meetings, err := p.Client.ListMeetings(ctx, nil)
	if err != nil {
		return t, err
	}

	for _, meeting := range meetings {
		args, err := p.TransformToDatum(ctx, meeting)
		if err != nil {
			return t, err
		}

		if err := p.Store.CreateMeetingDatum(ctx, args); err != nil {
			return t, err
		}

		t = meeting.Start
	}

	return t, nil
}

func (p *Processor) TransformToDatum(ctx context.Context, m Meeting) (args CreateMeetingDatumArguments, _ error) {
	rc, err := p.Client.DownloadMeeting(ctx, m.DownloadURL)
	if err != nil {
		return args, err
	}

	participants, err := p.Client.GetMeetingParticipants(ctx, m.ID)
	if err != nil {
		return args, err
	}

	return CreateMeetingDatumArguments{
		Topic:        m.Topic,
		Content:      rc,
		Participants: participants,
	}, nil
}
