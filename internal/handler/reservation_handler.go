package handler

import (
	"net/http"
	"strconv"

	"gotickets/internal/dto"
	"gotickets/internal/service"
	"gotickets/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type ReservationHandler struct {
	reservationService service.ReservationService
}

func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService}
}

func getUserID(c *echo.Context) uint {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	return uint(claims["id"].(float64))
}

func (h *ReservationHandler) CreateReservation(c *echo.Context) error {
	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}

	userID := getUserID(c)

	reservation, err := h.reservationService.CreateReservation(userID, req)
	if err != nil {
		if err.Error() == "parking zone is full" {
			return utils.SendError(c, http.StatusConflict, "Parking zone is full", err.Error())
		}
		return utils.SendError(c, http.StatusInternalServerError, "Failed to create reservation", err.Error())
	}

	res := dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}

	return utils.SendSuccess(c, http.StatusCreated, "Reservation confirmed successfully", res)
}

func (h *ReservationHandler) GetMyReservations(c *echo.Context) error {
	userID := getUserID(c)

	reservations, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve reservations", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "My reservations retrieved successfully", reservations)
}

func (h *ReservationHandler) CancelReservation(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid reservation ID", nil)
	}

	userID := getUserID(c)

	err = h.reservationService.CancelReservation(userID, uint(id))
	if err != nil {
		if err.Error() == "forbidden" {
			return utils.SendError(c, http.StatusForbidden, "Cannot cancel someone else's reservation", nil)
		}
		return utils.SendError(c, http.StatusInternalServerError, "Failed to cancel reservation", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "Reservation cancelled successfully", nil)
}

func (h *ReservationHandler) GetAllReservations(c *echo.Context) error {
	reservations, err := h.reservationService.GetAllReservations()
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve all reservations", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "All reservations retrieved successfully", reservations)
}
