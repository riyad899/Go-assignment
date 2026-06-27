package service

import (
	"gotickets/internal/dto"
	"gotickets/internal/models"
	"gotickets/internal/repository"
)

type ZoneService interface {
	CreateZone(req dto.CreateZoneRequest) (*models.ParkingZone, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
	resRepo  repository.ReservationRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository, resRepo repository.ReservationRepository) ZoneService {
	return &zoneService{zoneRepo, resRepo}
}

func (s *zoneService) CreateZone(req dto.CreateZoneRequest) (*models.ParkingZone, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	err := s.zoneRepo.CreateZone(zone)
	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, err := s.zoneRepo.GetAllZones()
	if err != nil {
		return nil, err
	}

	var res []dto.ZoneResponse
	for _, zone := range zones {
		activeCount, err := s.resRepo.GetActiveReservationCountByZone(zone.ID)
		if err != nil {
			return nil, err
		}
		res = append(res, dto.ZoneResponse{
			ID:             zone.ID,
			Name:           zone.Name,
			Type:           zone.Type,
			TotalCapacity:  zone.TotalCapacity,
			AvailableSpots: zone.TotalCapacity - int(activeCount),
			PricePerHour:   zone.PricePerHour,
			CreatedAt:      zone.CreatedAt,
			UpdatedAt:      zone.UpdatedAt,
		})
	}

	return res, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.GetZoneByID(id)
	if err != nil {
		return nil, err
	}

	activeCount, err := s.resRepo.GetActiveReservationCountByZone(zone.ID)
	if err != nil {
		return nil, err
	}

	res := &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.TotalCapacity - int(activeCount),
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}

	return res, nil
}
