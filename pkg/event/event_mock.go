package event

type MockEvent struct{}

var _ Event = (*MockEvent)(nil)
