package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	appport "github.com/naoki-higashi-28/mdprev/internal/application/port"
	watchuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/watch"
)

func TestWatchHandlerHandleWatch(t *testing.T) {
	events := make(chan appport.FileChangeEvent, 1)
	cancelCalled := false
	connectCalled := false
	disconnectCalled := false

	sub := &mockFileChangeSubscriber{
		events: events,
		cancel: func() { cancelCalled = true },
	}

	h := NewWatchHandler(
		watchuc.NewSubscribeFileChangesUseCase(sub),
		func() { connectCalled = true },
		func() { disconnectCalled = true },
	)

	req := httptest.NewRequest(http.MethodGet, "/api/watch", nil)
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		h.HandleWatch(rr, req)
		close(done)
	}()

	events <- appport.FileChangeEvent{Path: "docs/a.md", Type: appport.FileChangeTypeCreate}
	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("handler did not return")
	}

	if !connectCalled {
		t.Fatal("onConnect was not called")
	}
	if !disconnectCalled {
		t.Fatal("onDisconnect was not called")
	}
	if !cancelCalled {
		t.Fatal("subscriber cancel was not called")
	}

	body := rr.Body.String()
	if !strings.Contains(body, "tree_change") {
		t.Fatalf("SSE body missing tree_change: %q", body)
	}
	if !strings.Contains(body, "file_change") {
		t.Fatalf("SSE body missing file_change: %q", body)
	}
}
