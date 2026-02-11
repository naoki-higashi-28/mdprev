package watch

import (
	"errors"
	"testing"

	appport "github.com/naoki-higashi-28/mdprev/internal/application/port"
)

type mockSubscriber struct {
	events chan appport.FileChangeEvent
	cancel func()
}

func (m *mockSubscriber) Subscribe() (<-chan appport.FileChangeEvent, func()) {
	return m.events, m.cancel
}

func (m *mockSubscriber) Close() error {
	return errors.New("not used")
}

func TestSubscribeFileChangesUseCaseExecute(t *testing.T) {
	events := make(chan appport.FileChangeEvent, 1)
	cancelCalled := false
	sub := &mockSubscriber{
		events: events,
		cancel: func() { cancelCalled = true },
	}

	uc := NewSubscribeFileChangesUseCase(sub)
	ch, cancel := uc.Execute()

	events <- appport.FileChangeEvent{Path: "a.md", Type: appport.FileChangeTypeWrite}
	got := <-ch
	if got.Path != "a.md" || got.Type != appport.FileChangeTypeWrite {
		t.Fatalf("unexpected event: %#v", got)
	}

	cancel()
	if !cancelCalled {
		t.Fatalf("cancel function was not called")
	}
}
