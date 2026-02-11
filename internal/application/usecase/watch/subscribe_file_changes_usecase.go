package watch

import appport "github.com/naoki-higashi-28/mdprev/internal/application/port"

// SubscribeFileChangesUseCase handles file change subscriptions.
type SubscribeFileChangesUseCase struct {
	subscriber appport.FileChangeSubscriber
}

// NewSubscribeFileChangesUseCase creates a new SubscribeFileChangesUseCase.
func NewSubscribeFileChangesUseCase(subscriber appport.FileChangeSubscriber) *SubscribeFileChangesUseCase {
	return &SubscribeFileChangesUseCase{subscriber: subscriber}
}

// Execute subscribes to file change events.
func (uc *SubscribeFileChangesUseCase) Execute() (<-chan appport.FileChangeEvent, func()) {
	return uc.subscriber.Subscribe()
}
