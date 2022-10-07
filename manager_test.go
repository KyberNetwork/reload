package reload

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type priorityMockReloader struct {
	priority int
	m        *MockReloader
}

func TestManagerWithNotifierFunc(t *testing.T) {
	type TestCase struct {
		name        string
		prepare     func(ctrl *gomock.Controller) []priorityMockReloader
		notifierID  string
		notifierErr error
		expErr      bool
	}

	tests := []TestCase{
		{
			name: "If notifier fails it should end the execution with an error.",
			prepare: func(ctrl *gomock.Controller) []priorityMockReloader {
				reloader := NewMockReloader(ctrl)

				return []priorityMockReloader{
					{0, reloader},
				}
			},
			notifierID:  "test-id",
			notifierErr: fmt.Errorf("something"),
			expErr:      true,
		},
		{
			name: "Single reloader should be called with the expected trigger ID.",
			prepare: func(ctrl *gomock.Controller) []priorityMockReloader {
				reloader := NewMockReloader(ctrl)
				reloader.EXPECT().Reload(gomock.Any(), "test-id").Return(nil)

				return []priorityMockReloader{
					{0, reloader},
				}
			},
			notifierID: "test-id",
		},
		{
			name: "Single reloader error should get the error.",
			prepare: func(ctrl *gomock.Controller) []priorityMockReloader {
				reloader := NewMockReloader(ctrl)
				reloader.EXPECT().Reload(gomock.Any(), gomock.Any()).Return(fmt.Errorf("something"))

				return []priorityMockReloader{
					{0, reloader},
				}
			},
			notifierID: "test-id",
			expErr:     true,
		},
		{
			name: "Multiple reloaders should be called with the expected trigger ID.",
			prepare: func(ctrl *gomock.Controller) []priorityMockReloader {
				reloader1 := NewMockReloader(ctrl)
				reloader1.EXPECT().Reload(gomock.Any(), "test-id").Return(nil)

				reloader2 := NewMockReloader(ctrl)
				reloader2.EXPECT().Reload(gomock.Any(), "test-id").Return(nil)

				reloader3 := NewMockReloader(ctrl)
				reloader3.EXPECT().Reload(gomock.Any(), "test-id").Return(nil)

				return []priorityMockReloader{
					{0, reloader1},
					{0, reloader2},
					{0, reloader3},
				}
			},
			notifierID: "test-id",
		},
		{
			name: "Multiple reloaders with different priority should be called with the expected trigger ID.",
			prepare: func(ctrl *gomock.Controller) []priorityMockReloader {
				reloader1 := NewMockReloader(ctrl)
				reloader1.EXPECT().Reload(gomock.Any(), "test-id").Return(nil)

				reloader2 := NewMockReloader(ctrl)
				reloader2.EXPECT().Reload(gomock.Any(), "test-id").Return(nil)

				reloader3 := NewMockReloader(ctrl)
				reloader3.EXPECT().Reload(gomock.Any(), "test-id").Return(nil)

				return []priorityMockReloader{
					{2, reloader1},
					{0, reloader2},
					{1, reloader3},
				}
			},
			notifierID: "test-id",
		},
		{
			name: "Having multiple reloaders with different priority, if a lower priority errors, shouldn't call the next ones.",
			prepare: func(ctrl *gomock.Controller) []priorityMockReloader {
				reloader1 := NewMockReloader(ctrl)
				reloader1.EXPECT().Reload(gomock.Any(), "test-id").Return(fmt.Errorf("something"))

				reloader2 := NewMockReloader(ctrl)
				reloader2.EXPECT().Reload(gomock.Any(), "test-id").Return(nil)

				reloader3 := NewMockReloader(ctrl)

				reloader4 := NewMockReloader(ctrl)

				reloader5 := NewMockReloader(ctrl)

				return []priorityMockReloader{
					{10, reloader1},
					{4, reloader2},
					{25, reloader3},
					{20, reloader4},
					{25, reloader5},
				}
			},
			notifierID: "test-id",
			expErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			reloaders := tc.prepare(ctrl)

			// Prepare.
			m := NewManager()
			for _, r := range reloaders {
				m.RegisterReloader(r.priority, r.m)
			}
			notifierC := make(chan string)
			m.RegisterNotifier(NotifierFunc(func(context.Context) (string, error) {
				notifierID := <-notifierC
				return notifierID, tc.notifierErr
			}))

			// Execute.
			ctx, cancel := context.WithCancel(context.Background())
			checksFinished := make(chan struct{})
			go func() {
				err := m.Run(ctx)

				// Check.
				if tc.expErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}

				close(checksFinished)
			}()

			// Release the trigger to start the execution and checks.
			notifierC <- tc.notifierID

			// Wait for until the reloaders handle the trigger.
			// Then cancel the context in case the reloaders didn't throw
			// error.
			time.Sleep(10 * time.Millisecond)
			cancel()

			// Wait until everything has been checked.
			<-checksFinished
		})
	}
}

func TestManagerWithNotifierChan(t *testing.T) {
	type TestCase struct {
		name       string
		prepare    func(ctrl *gomock.Controller) []priorityMockReloader
		notifierID string
		expErr     bool
	}

	tests := []TestCase{
		{
			name: "Single reloader should be called with the expected trigger ID.",
			prepare: func(ctrl *gomock.Controller) []priorityMockReloader {
				reloader := NewMockReloader(ctrl)
				reloader.EXPECT().Reload(gomock.Any(), "test-id").Return(nil)

				return []priorityMockReloader{
					{0, reloader},
				}
			},
			notifierID: "test-id",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			reloaders := tc.prepare(ctrl)

			// Prepare.
			m := NewManager()
			for _, r := range reloaders {
				m.RegisterReloader(r.priority, r.m)
			}
			notifierC := make(chan string)
			m.RegisterNotifier(NotifierChan(notifierC))

			// Execute.
			ctx, cancel := context.WithCancel(context.Background())
			checksFinished := make(chan struct{})
			go func() {
				err := m.Run(ctx)

				// Check.
				if tc.expErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}

				close(checksFinished)
			}()

			// Release the trigger to start the execution and checks.
			notifierC <- tc.notifierID

			// Wait for until the reloaders handle the trigger.
			// Then cancel the context in case the reloaders didn't
			// error.
			time.Sleep(10 * time.Millisecond)
			cancel()

			// Wait until everything has been checked.
			<-checksFinished
		})
	}
}
