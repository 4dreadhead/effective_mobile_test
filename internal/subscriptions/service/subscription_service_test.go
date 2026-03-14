package service_test

import (
	"context"
	"effective_mobile_test/internal/subscriptions/repository"
	"errors"
	"testing"
	"time"

	apperrors "effective_mobile_test/internal/platform/errors"
	"effective_mobile_test/internal/subscriptions/model"
	"effective_mobile_test/internal/subscriptions/service"
)

// --- Mock ---

type mockRepo struct {
	createFn func(ctx context.Context, sub model.Subscription) (*model.Subscription, error)
	getFn    func(ctx context.Context, id uint64) (*model.Subscription, error)
	updateFn func(ctx context.Context, sub *model.Subscription) (*model.Subscription, error)
	deleteFn func(ctx context.Context, id uint64) error
	listFn   func(ctx context.Context, filter repository.ListFilter) ([]model.Subscription, error)
	totalFn  func(ctx context.Context, filter repository.TotalFilter) (int64, bool, error)
}

func (m *mockRepo) CreateSubscription(ctx context.Context, sub model.Subscription) (*model.Subscription, error) {
	return m.createFn(ctx, sub)
}
func (m *mockRepo) GetSubscription(ctx context.Context, id uint64) (*model.Subscription, error) {
	return m.getFn(ctx, id)
}
func (m *mockRepo) UpdateSubscription(ctx context.Context, sub *model.Subscription) (*model.Subscription, error) {
	return m.updateFn(ctx, sub)
}
func (m *mockRepo) DeleteSubscription(ctx context.Context, id uint64) error {
	return m.deleteFn(ctx, id)
}
func (m *mockRepo) ListSubscriptions(ctx context.Context, filter repository.ListFilter) ([]model.Subscription, error) {
	return m.listFn(ctx, filter)
}
func (m *mockRepo) TotalCost(ctx context.Context, filter repository.TotalFilter) (int64, bool, error) {
	return m.totalFn(ctx, filter)
}

// --- Helpers ---

func ptr[T any](v T) *T { return &v }

// --- Create ---

func TestCreate_Success(t *testing.T) {
	repo := &mockRepo{
		createFn: func(_ context.Context, sub model.Subscription) (*model.Subscription, error) {
			sub.ID = 1
			return &sub, nil
		},
	}
	uc := service.NewSubscriptionService(repo)

	result, err := uc.Create(context.Background(), service.CreateSubscriptionInput{
		ServiceName: "Yandex Plus",
		MonthlyCost: 400,
		UserID:      "user-1",
		From:        "01.2025",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID=1, got %d", result.ID)
	}
}

func TestCreate_EmptyServiceName(t *testing.T) {
	uc := service.NewSubscriptionService(&mockRepo{})

	_, err := uc.Create(context.Background(), service.CreateSubscriptionInput{
		ServiceName: "   ",
		MonthlyCost: 400,
		UserID:      "user-1",
		From:        "01.2025",
	})

	if !errors.Is(err, apperrors.ErrInvalidFields) {
		t.Errorf("expected ErrInvalidFields, got %v", err)
	}
}

func TestCreate_NegativeCost(t *testing.T) {
	uc := service.NewSubscriptionService(&mockRepo{})

	_, err := uc.Create(context.Background(), service.CreateSubscriptionInput{
		ServiceName: "Netflix",
		MonthlyCost: -1,
		UserID:      "user-1",
		From:        "01.2025",
	})

	if !errors.Is(err, apperrors.ErrInvalidFields) {
		t.Errorf("expected ErrInvalidFields, got %v", err)
	}
}

func TestCreate_ToBeforeFrom(t *testing.T) {
	uc := service.NewSubscriptionService(&mockRepo{})

	_, err := uc.Create(context.Background(), service.CreateSubscriptionInput{
		ServiceName: "Netflix",
		MonthlyCost: 500,
		UserID:      "user-1",
		From:        "06.2025",
		To:          "01.2025",
	})

	if !errors.Is(err, apperrors.ErrInvalidFields) {
		t.Errorf("expected ErrInvalidFields, got %v", err)
	}
}

func TestCreate_InvalidFromFormat(t *testing.T) {
	uc := service.NewSubscriptionService(&mockRepo{})

	_, err := uc.Create(context.Background(), service.CreateSubscriptionInput{
		ServiceName: "Netflix",
		MonthlyCost: 500,
		UserID:      "user-1",
		From:        "2025-01",
	})

	if !errors.Is(err, apperrors.ErrInvalidFields) {
		t.Errorf("expected ErrInvalidFields, got %v", err)
	}
}

// --- Get ---

func TestGet_NotFound(t *testing.T) {
	repo := &mockRepo{
		getFn: func(_ context.Context, id uint64) (*model.Subscription, error) {
			return nil, apperrors.ErrRecordNotFound
		},
	}
	uc := service.NewSubscriptionService(repo)

	_, err := uc.Get(context.Background(), 999)

	if !errors.Is(err, apperrors.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound, got %v", err)
	}
}

func TestGet_Success(t *testing.T) {
	expected := &model.Subscription{ID: 42, ServiceName: "Spotify"}
	repo := &mockRepo{
		getFn: func(_ context.Context, id uint64) (*model.Subscription, error) {
			return expected, nil
		},
	}
	uc := service.NewSubscriptionService(repo)

	result, err := uc.Get(context.Background(), 42)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 42 {
		t.Errorf("expected ID=42, got %d", result.ID)
	}
}

// --- Delete ---

func TestDelete_NotFound(t *testing.T) {
	repo := &mockRepo{
		deleteFn: func(_ context.Context, id uint64) error {
			return apperrors.ErrRecordNotFound
		},
	}
	uc := service.NewSubscriptionService(repo)

	err := uc.Delete(context.Background(), 999)

	if !errors.Is(err, apperrors.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound, got %v", err)
	}
}

func TestDelete_Success(t *testing.T) {
	repo := &mockRepo{
		deleteFn: func(_ context.Context, id uint64) error {
			return nil
		},
	}
	uc := service.NewSubscriptionService(repo)

	if err := uc.Delete(context.Background(), 1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- Update ---

func TestUpdate_NotFound(t *testing.T) {
	repo := &mockRepo{
		getFn: func(_ context.Context, id uint64) (*model.Subscription, error) {
			return nil, apperrors.ErrRecordNotFound
		},
	}
	uc := service.NewSubscriptionService(repo)

	_, err := uc.Update(context.Background(), 999, service.UpdateSubscriptionInput{})

	if !errors.Is(err, apperrors.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound, got %v", err)
	}
}

func TestUpdate_ToBeforeFrom(t *testing.T) {
	from := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	repo := &mockRepo{
		getFn: func(_ context.Context, id uint64) (*model.Subscription, error) {
			return &model.Subscription{
				ID:          1,
				ServiceName: "Netflix",
				MonthlyCost: 500,
				FromDate:    from,
			}, nil
		},
	}
	uc := service.NewSubscriptionService(repo)

	_, err := uc.Update(context.Background(), 1, service.UpdateSubscriptionInput{
		To: ptr("01.2025"),
	})

	if !errors.Is(err, apperrors.ErrInvalidFields) {
		t.Errorf("expected ErrInvalidFields, got %v", err)
	}
}

func TestUpdate_Success(t *testing.T) {
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	existing := &model.Subscription{
		ID:          1,
		ServiceName: "Netflix",
		MonthlyCost: 500,
		FromDate:    from,
	}
	repo := &mockRepo{
		getFn: func(_ context.Context, id uint64) (*model.Subscription, error) {
			return existing, nil
		},
		updateFn: func(_ context.Context, sub *model.Subscription) (*model.Subscription, error) {
			return sub, nil
		},
	}
	uc := service.NewSubscriptionService(repo)

	result, err := uc.Update(context.Background(), 1, service.UpdateSubscriptionInput{
		MonthlyCost: ptr(600),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MonthlyCost != 600 {
		t.Errorf("expected MonthlyCost=600, got %d", result.MonthlyCost)
	}
}

// --- TotalCost ---

func TestTotalCost_InvalidFrom(t *testing.T) {
	uc := service.NewSubscriptionService(&mockRepo{})

	_, _, err := uc.TotalCost(context.Background(), repository.TotalFilter{
		From: "bad",
		To:   "12.2025",
	})

	if !errors.Is(err, apperrors.ErrInvalidFields) {
		t.Errorf("expected ErrInvalidFields, got %v", err)
	}
}

func TestTotalCost_NoData(t *testing.T) {
	repo := &mockRepo{
		totalFn: func(_ context.Context, filter repository.TotalFilter) (int64, bool, error) {
			return 0, false, nil
		},
	}
	uc := service.NewSubscriptionService(repo)

	total, hasData, err := uc.TotalCost(context.Background(), repository.TotalFilter{
		UserID:      "user-1",
		ServiceName: "Netflix",
		From:        "01.2025",
		To:          "06.2025",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hasData {
		t.Error("expected hasData=false")
	}
	if total != 0 {
		t.Errorf("expected total=0, got %d", total)
	}
}

func TestTotalCost_WithData(t *testing.T) {
	repo := &mockRepo{
		totalFn: func(_ context.Context, filter repository.TotalFilter) (int64, bool, error) {
			return 2400, true, nil
		},
	}
	uc := service.NewSubscriptionService(repo)

	total, hasData, err := uc.TotalCost(context.Background(), repository.TotalFilter{
		UserID:      "user-1",
		ServiceName: "Netflix",
		From:        "01.2025",
		To:          "06.2025",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hasData {
		t.Error("expected hasData=true")
	}
	if total != 2400 {
		t.Errorf("expected total=2400, got %d", total)
	}
}
