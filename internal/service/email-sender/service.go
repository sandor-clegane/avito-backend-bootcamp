package emailsender

import (
	"avito-backend-bootcamp/internal/infra/repository"
	"avito-backend-bootcamp/internal/model"
	r "avito-backend-bootcamp/pkg/utils/retry"
	"avito-backend-bootcamp/pkg/utils/sl"
	"errors"
	"fmt"

	"context"
	"encoding/json"
	"log/slog"
	"time"
)

const (
	retryAttempts = 5
	retryTimeout  = 1 * time.Second
)

type EmailSender interface {
	SendEmail(ctx context.Context, recipient string, message string) error
}

type SubscriptionRepository interface {
	SubsciptionListByHouseID(ctx context.Context, houseID int64) ([]*model.Subscription, error)
}

type EventRepository interface {
	GetNewEvent(ctx context.Context) (*model.Event, error)
	SetDone(ctx context.Context, eventID int64) error
}

type HouseRepository interface {
	GetHouse(ctx context.Context, id int64) (*model.House, error)
}

type Service struct {
	log                   *slog.Logger
	sender                EmailSender
	subscitpionRepository SubscriptionRepository
	eventRepository       EventRepository
	houseRepository       HouseRepository
	retrier               *r.Retrier
}

func New(
	log *slog.Logger,
	emailSender EmailSender,
	subscitpionRepository SubscriptionRepository,
	eventRepository EventRepository,
	houseRepository HouseRepository,
) *Service {
	return &Service{
		log:                   log,
		sender:                emailSender,
		subscitpionRepository: subscitpionRepository,
		eventRepository:       eventRepository,
		houseRepository:       houseRepository,
		retrier:               r.NewRetrier(retryAttempts, retryTimeout),
	}
}

type Payload struct {
	HouseID int64 `json:"house_id"`
}

// StartProcessEvents starts a goroutine that processes events periodically.
func (s *Service) StartProcessEvents(ctx context.Context, handlePeriod time.Duration) {
	const op = "email-sender.StartProcessEvents"

	log := s.log.With(slog.String("op", op))

	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("stopping event processing")
				return
			case <-ticker.C:
				// Process event
				err := s.processEvent(ctx)
				if err != nil {
					log.Error("failed to process event", sl.Err(err))
				}
			}
		}
	}()
}

// processEvents processes new events from the repository.
func (s *Service) processEvent(ctx context.Context) error {
	event, err := s.eventRepository.GetNewEvent(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			s.log.Info("no events to send")
			return nil
		}
		return fmt.Errorf("failed to get new event: %w", err)
	}

	// Unmarshal the event payload
	var payload Payload
	err = json.Unmarshal([]byte(event.Payload), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event payload: %w", err)
	}

	// Fetch subscribers and house information
	subscribers, err := s.subscitpionRepository.SubsciptionListByHouseID(ctx, payload.HouseID)
	if err != nil {
		return fmt.Errorf("failed to get subscribers list: %w", err)
	}

	house, err := s.houseRepository.GetHouse(ctx, payload.HouseID)
	if err != nil {
		return fmt.Errorf("failed to get house by ID: %w", err)
	}

	// Send emails to subscribers
	for _, sub := range subscribers {
		s.retrier.Retry(ctx, func() error {
			return s.sender.SendEmail(ctx, sub.Email, composeEmailMessage(house))
		})
	}

	// Mark the event as done
	if err := s.eventRepository.SetDone(ctx, event.ID); err != nil {
		return fmt.Errorf("failed to set event done: %w", err)
	}

	return nil
}

// composeEmailMessage composes the email message.
func composeEmailMessage(house *model.House) string {
	return fmt.Sprintf("В доме по адресу %s появилось новое объявление", house.Address)
}
