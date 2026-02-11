package handler

import (
	appport "github.com/naoki-higashi-28/mdprev/internal/application/port"
	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

type mockFileRepo struct {
	readFn func(path string) ([]byte, error)
	rawFn  func(path string) (string, error)
}

func (m *mockFileRepo) ReadMarkdown(path string) ([]byte, error) {
	if m.readFn != nil {
		return m.readFn(path)
	}
	return nil, nil
}

func (m *mockFileRepo) ServeRaw(path string) (string, error) {
	if m.rawFn != nil {
		return m.rawFn(path)
	}
	return "", nil
}

type mockTreeRepo struct {
	listFn   func(path string) ([]model.Entry, error)
	searchFn func(query string) ([]model.Entry, error)
}

func (m *mockTreeRepo) ListEntries(path string) ([]model.Entry, error) {
	if m.listFn != nil {
		return m.listFn(path)
	}
	return nil, nil
}

func (m *mockTreeRepo) SearchEntries(query string) ([]model.Entry, error) {
	if m.searchFn != nil {
		return m.searchFn(query)
	}
	return nil, nil
}

type mockFileChangeSubscriber struct {
	events chan appport.FileChangeEvent
	cancel func()
}

func (m *mockFileChangeSubscriber) Subscribe() (<-chan appport.FileChangeEvent, func()) {
	if m.cancel == nil {
		m.cancel = func() {}
	}
	return m.events, m.cancel
}

func (m *mockFileChangeSubscriber) Close() error {
	return nil
}
