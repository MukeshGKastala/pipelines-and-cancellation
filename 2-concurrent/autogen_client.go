// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package concurrent

import (
	"context"
	"io"
	"sync"
)

// Ensure, that ClientInterfaceMock does implement ClientInterface.
// If this is not the case, regenerate this file with moq.
var _ ClientInterface = &ClientInterfaceMock{}

// ClientInterfaceMock is a mock implementation of ClientInterface.
//
//	func TestSomethingThatUsesClientInterface(t *testing.T) {
//
//		// make and configure a mocked ClientInterface
//		mockedClientInterface := &ClientInterfaceMock{
//			DownloadMeetingFunc: func(ctx context.Context, url string) (io.ReadCloser, error) {
//				panic("mock out the DownloadMeeting method")
//			},
//			GetMeetingParticipantsFunc: func(ctx context.Context, meetingID string) ([]Participant, error) {
//				panic("mock out the GetMeetingParticipants method")
//			},
//			ListPaginatedMeetingsFunc: func(ctx context.Context, params *ListPaginatedMeetingsParams) (ListPaginatedMeetingsResponse, error) {
//				panic("mock out the ListPaginatedMeetings method")
//			},
//		}
//
//		// use mockedClientInterface in code that requires ClientInterface
//		// and then make assertions.
//
//	}
type ClientInterfaceMock struct {
	// DownloadMeetingFunc mocks the DownloadMeeting method.
	DownloadMeetingFunc func(ctx context.Context, url string) (io.ReadCloser, error)

	// GetMeetingParticipantsFunc mocks the GetMeetingParticipants method.
	GetMeetingParticipantsFunc func(ctx context.Context, meetingID string) ([]Participant, error)

	// ListPaginatedMeetingsFunc mocks the ListPaginatedMeetings method.
	ListPaginatedMeetingsFunc func(ctx context.Context, params *ListPaginatedMeetingsParams) (ListPaginatedMeetingsResponse, error)

	// calls tracks calls to the methods.
	calls struct {
		// DownloadMeeting holds details about calls to the DownloadMeeting method.
		DownloadMeeting []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// URL is the url argument value.
			URL string
		}
		// GetMeetingParticipants holds details about calls to the GetMeetingParticipants method.
		GetMeetingParticipants []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// MeetingID is the meetingID argument value.
			MeetingID string
		}
		// ListPaginatedMeetings holds details about calls to the ListPaginatedMeetings method.
		ListPaginatedMeetings []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Params is the params argument value.
			Params *ListPaginatedMeetingsParams
		}
	}
	lockDownloadMeeting        sync.RWMutex
	lockGetMeetingParticipants sync.RWMutex
	lockListPaginatedMeetings  sync.RWMutex
}

// DownloadMeeting calls DownloadMeetingFunc.
func (mock *ClientInterfaceMock) DownloadMeeting(ctx context.Context, url string) (io.ReadCloser, error) {
	if mock.DownloadMeetingFunc == nil {
		panic("ClientInterfaceMock.DownloadMeetingFunc: method is nil but ClientInterface.DownloadMeeting was just called")
	}
	callInfo := struct {
		Ctx context.Context
		URL string
	}{
		Ctx: ctx,
		URL: url,
	}
	mock.lockDownloadMeeting.Lock()
	mock.calls.DownloadMeeting = append(mock.calls.DownloadMeeting, callInfo)
	mock.lockDownloadMeeting.Unlock()
	return mock.DownloadMeetingFunc(ctx, url)
}

// DownloadMeetingCalls gets all the calls that were made to DownloadMeeting.
// Check the length with:
//
//	len(mockedClientInterface.DownloadMeetingCalls())
func (mock *ClientInterfaceMock) DownloadMeetingCalls() []struct {
	Ctx context.Context
	URL string
} {
	var calls []struct {
		Ctx context.Context
		URL string
	}
	mock.lockDownloadMeeting.RLock()
	calls = mock.calls.DownloadMeeting
	mock.lockDownloadMeeting.RUnlock()
	return calls
}

// GetMeetingParticipants calls GetMeetingParticipantsFunc.
func (mock *ClientInterfaceMock) GetMeetingParticipants(ctx context.Context, meetingID string) ([]Participant, error) {
	if mock.GetMeetingParticipantsFunc == nil {
		panic("ClientInterfaceMock.GetMeetingParticipantsFunc: method is nil but ClientInterface.GetMeetingParticipants was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		MeetingID string
	}{
		Ctx:       ctx,
		MeetingID: meetingID,
	}
	mock.lockGetMeetingParticipants.Lock()
	mock.calls.GetMeetingParticipants = append(mock.calls.GetMeetingParticipants, callInfo)
	mock.lockGetMeetingParticipants.Unlock()
	return mock.GetMeetingParticipantsFunc(ctx, meetingID)
}

// GetMeetingParticipantsCalls gets all the calls that were made to GetMeetingParticipants.
// Check the length with:
//
//	len(mockedClientInterface.GetMeetingParticipantsCalls())
func (mock *ClientInterfaceMock) GetMeetingParticipantsCalls() []struct {
	Ctx       context.Context
	MeetingID string
} {
	var calls []struct {
		Ctx       context.Context
		MeetingID string
	}
	mock.lockGetMeetingParticipants.RLock()
	calls = mock.calls.GetMeetingParticipants
	mock.lockGetMeetingParticipants.RUnlock()
	return calls
}

// ListPaginatedMeetings calls ListPaginatedMeetingsFunc.
func (mock *ClientInterfaceMock) ListPaginatedMeetings(ctx context.Context, params *ListPaginatedMeetingsParams) (ListPaginatedMeetingsResponse, error) {
	if mock.ListPaginatedMeetingsFunc == nil {
		panic("ClientInterfaceMock.ListPaginatedMeetingsFunc: method is nil but ClientInterface.ListPaginatedMeetings was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Params *ListPaginatedMeetingsParams
	}{
		Ctx:    ctx,
		Params: params,
	}
	mock.lockListPaginatedMeetings.Lock()
	mock.calls.ListPaginatedMeetings = append(mock.calls.ListPaginatedMeetings, callInfo)
	mock.lockListPaginatedMeetings.Unlock()
	return mock.ListPaginatedMeetingsFunc(ctx, params)
}

// ListPaginatedMeetingsCalls gets all the calls that were made to ListPaginatedMeetings.
// Check the length with:
//
//	len(mockedClientInterface.ListPaginatedMeetingsCalls())
func (mock *ClientInterfaceMock) ListPaginatedMeetingsCalls() []struct {
	Ctx    context.Context
	Params *ListPaginatedMeetingsParams
} {
	var calls []struct {
		Ctx    context.Context
		Params *ListPaginatedMeetingsParams
	}
	mock.lockListPaginatedMeetings.RLock()
	calls = mock.calls.ListPaginatedMeetings
	mock.lockListPaginatedMeetings.RUnlock()
	return calls
}
