package service

import (
	"errors"

	"gotickets/internal/dto"
	"gotickets/internal/models"
	"gotickets/internal/repository"
)

type ReservationService interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*models.Reservation, error)
	GetMyReservations(userID uint) ([]dto.ReservationResponse, error)
	CancelReservation(userID uint, reservationID uint) error
	GetAllReservations() ([]dto.ReservationResponse, error)
}

type reservationService struct {
	resRepo repository.ReservationRepository
}

func NewReservationService(resRepo repository.ReservationRepository) ReservationService {
	return &reservationService{resRepo}
}

func (s *reservationService) CreateReservation(userID uint, req dto.CreateReservationRequest) (*models.Reservation, error) {
	reservation := &models.Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       "active",
	}

	err := s.resRepo.CreateReservationTx(reservation)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.ReservationResponse, error) {
	reservations, err := s.resRepo.GetReservationsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var res []dto.ReservationResponse
	for _, r := range reservations {
		item := dto.ReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
		}
		if r.Zone != nil {
			item.Zone = &dto.ReservationZoneResponse{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			}
		}
		res = append(res, item)
	}

	return res, nil
}

func (s *reservationService) CancelReservation(userID uint, reservationID uint) error {
	reservation, err := s.resRepo.GetReservationByID(reservationID)
	if err != nil {
		return err
	}

	if reservation.UserID != userID {
		return errors.New("forbidden")
	}

	return s.resRepo.UpdateReservationStatus(reservationID, "cancelled")
}

func (s *reservationService) GetAllReservations() ([]dto.ReservationResponse, error) {
	reservations, err := s.resRepo.GetAllReservations()
	if err != nil {
		return nil, err
	}

	var res []dto.ReservationResponse
	for _, r := range reservations {
		item := dto.ReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
		}
		if r.Zone != nil {
			item.Zone = &dto.ReservationZoneResponse{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			}
		}
		if r.User != nil {
			item.User = &dto.UserResponse{
				ID:        r.User.ID,
				Name:      r.User.Name,
				Email:     r.User.Email,
				Role:      r.User.Role,
				CreatedAt: r.User.CreatedAt,
				UpdatedAt: r.User.UpdatedAt,
			}
		}
		res = append(res, item)
	}

	return res, nil
}
