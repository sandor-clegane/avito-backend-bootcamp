package flat

import (
	"avito-backend-bootcamp/internal/model"
	"avito-backend-bootcamp/pkg/utils/sl"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type FlatRepository interface {
	GetFlat(ctx context.Context, ID int64) (*model.Flat, error)
	SaveFlat(ctx context.Context, houseID, price, fooms int64) (*model.Flat, error)
	UpdateFlat(ctx context.Context, flat *model.Flat) (*model.Flat, error)
	FlatListByHouseID(ctx context.Context, houseID int64) ([]*model.Flat, error)
}

type EventRepository interface {
	PublishEvent(ctx context.Context, eventType model.EventType, payload string) error
}

type Cache[K comparable, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Remove(key K)
}

type Service struct {
	log             *slog.Logger
	flatRepository  FlatRepository
	eventRepository EventRepository
	cache           Cache[int64, string]
	trManager       *manager.Manager
}

func New(
	log *slog.Logger,
	flatRepository FlatRepository,
	eventRepository EventRepository,
	cache Cache[int64, string],
	trManager *manager.Manager,
) *Service {
	return &Service{
		log:             log,
		flatRepository:  flatRepository,
		eventRepository: eventRepository,
		cache:           cache,
		trManager:       trManager,
	}
}

func (s *Service) CreateFlat(ctx context.Context, houseID, price, rooms int64) (*model.Flat, error) {
	const op = "flat.UpdateCreateFlatFlat"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("houseID", houseID),
		slog.Int64("price", price),
		slog.Int64("rooms", rooms),
	)

	flat, err := s.flatRepository.SaveFlat(ctx, houseID, price, rooms)
	if err != nil {
		log.Error("failed to save flat", sl.Err(err))
		return nil, err
	}

	// it is necessary to invalidate the cache
	// since new flat was created
	s.cache.Remove(flat.HouseID)

	return flat, nil
}

func (s *Service) UpdateFlat(ctx context.Context, ID int64, status model.FlatStatus) (flat *model.Flat, err error) {
	const op = "flat.UpdateFlat"

	log := s.log.With(
		slog.String("op", op),
		slog.String("status", string(status)),
	)

	flat, err = s.flatRepository.GetFlat(ctx, ID)
	if err != nil {
		log.Error("failed to find flat", sl.Err(err))
		return nil, err
	}

	switch status {
	case model.StatusApproved:
		err = flat.Approve()
	case model.StatusDeclined:
		err = flat.Decline()
	case model.StatusOnModeration:
		err = flat.StartModeration()
	default:
		return nil, model.ErrImpossibleTransition
	}
	if err != nil {
		log.Error("failed to change status", sl.Err(err))
		return nil, err
	}

	s.trManager.Do(ctx, func(ctx context.Context) error {
		flat, err = s.flatRepository.UpdateFlat(ctx, flat)
		if err != nil {
			log.Error("failed to update flat in db", sl.Err(err))
			return err
		}

		// it is necessary to invalidate the cache
		// since flat was updated
		s.cache.Remove(flat.HouseID)

		if flat.Status == model.StatusApproved {
			eventPayload := fmt.Sprintf(
				`{"house_id": %d}`,
				flat.HouseID,
			)

			err = s.eventRepository.PublishEvent(ctx, model.FlatApproved, eventPayload)
			if err != nil {
				log.Error("fauled to publish event", sl.Err(err))
				return err
			}
		}

		return nil
	})

	return flat, nil
}

// GetFlatListByHouseID retrieves a list of flats for a given house ID,
// applying visibility rules based on the user role.
func (s *Service) GetFlatListByHouseID(ctx context.Context, houseID int64, userRole model.UserType) ([]*model.Flat, error) {
	const op = "flat.GetFlatListByHouseID"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_type", string(userRole)),
		slog.Int64("house_id", houseID),
	)

	// Cache hit
	cachedList, ok := s.cache.Get(houseID)
	if ok {
		var flats []*model.Flat
		err := json.Unmarshal([]byte(cachedList), &flats)
		if err != nil {
			log.Error("invalid json in cache", sl.Err(err))
			s.cache.Remove(houseID)
		} else {
			return flats, nil
		}
	}

	// Fetch from database
	flatList, err := s.flatRepository.FlatListByHouseID(ctx, houseID)
	if err != nil {
		log.Error("failed to get flat list", sl.Err(err))
		return nil, err
	}

	// Cache result
	flatsJSON, err := json.Marshal(flatList)
	if err != nil {
		log.Error("failed to marshal flat list", sl.Err(err))
	} else {
		s.cache.Set(houseID, string(flatsJSON))
	}

	return filterFlats(flatList, userRole), nil
}

// filterFlats filters the flat list based on the user role.
func filterFlats(flatList []*model.Flat, userRole model.UserType) []*model.Flat {
	var visibleFlats []*model.Flat

	for _, flat := range flatList {
		switch userRole {
		case model.Client:
			if flat.Status == model.StatusApproved {
				visibleFlats = append(visibleFlats, flat)
			}
		case model.Moderator:
			visibleFlats = append(visibleFlats, flat)
		}
	}

	return visibleFlats
}
