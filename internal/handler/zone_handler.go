package handler

import (
	"net/http"
	"strconv"

	"gotickets/internal/dto"
	"gotickets/internal/service"
	"gotickets/internal/utils"

	"github.com/labstack/echo/v5"
)

type ZoneHandler struct {
	zoneService service.ZoneService
}

func NewZoneHandler(zoneService service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService}
}

func (h *ZoneHandler) CreateZone(c *echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}

	zone, err := h.zoneService.CreateZone(req)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to create parking zone", err.Error())
	}

	res := dto.ZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}

	return utils.SendSuccess(c, http.StatusCreated, "Parking zone created successfully", res)
}

func (h *ZoneHandler) GetAllZones(c *echo.Context) error {
	zones, err := h.zoneService.GetAllZones()
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve parking zones", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "Parking zones retrieved successfully", zones)
}

func (h *ZoneHandler) GetZoneByID(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid zone ID", nil)
	}

	zone, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		return utils.SendError(c, http.StatusNotFound, "Parking zone not found", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "Parking zone retrieved successfully", zone)
}
