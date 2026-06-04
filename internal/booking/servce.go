package booking

type service struct {
	bookingRepo Repository
}

func NewService(bookingRepo Repository) *service {
	return &service{
		bookingRepo: bookingRepo,
	}
}
