package api

// ToBooking performs a lossy conversion from driverCarpoolBooking to Booking
// objects
func (dcb DriverCarpoolBooking) ToBooking() *Booking {
	return &Booking{
		Id:                     dcb.Id,
		Driver:                 dcb.Driver,
		Passenger:              User{},
		PassengerPickupLat:     dcb.PassengerPickupLat,
		PassengerPickupLng:     dcb.PassengerPickupLng,
		PassengerDropLat:       dcb.PassengerDropLat,
		PassengerDropLng:       dcb.PassengerDropLng,
		PassengerPickupAddress: dcb.PassengerPickupAddress,
		PassengerDropAddress:   dcb.PassengerDropAddress,
		Status:                 BookingStatus(dcb.Status),
		Duration:               dcb.Duration,
		Distance:               dcb.Distance,
		WebUrl:                 &dcb.WebUrl,
		Car:                    dcb.Car,
		Price:                  dcb.Price,
	}
}

// ToBooking performs a lossy conversion from passengerCarpoolBooking to
// Booking objects
func (pcb PassengerCarpoolBooking) ToBooking() *Booking {
	return &Booking{
		Id:                     pcb.Id,
		Driver:                 User{},
		Passenger:              pcb.Passenger,
		PassengerPickupLat:     pcb.PassengerPickupLat,
		PassengerPickupLng:     pcb.PassengerPickupLng,
		PassengerDropLat:       pcb.PassengerDropLat,
		PassengerDropLng:       pcb.PassengerDropLng,
		PassengerPickupAddress: pcb.PassengerPickupAddress,
		PassengerDropAddress:   pcb.PassengerDropAddress,
		Status:                 BookingStatus(pcb.Status),
		Duration:               pcb.Duration,
		Distance:               pcb.Distance,
		WebUrl:                 &pcb.WebUrl,
	}
}

// ToDriverCarpoolBooking performs a lossy conversion from
// Booking to DriverCarpoolBooking objects
func (b Booking) ToDriverCarpoolBooking() *DriverCarpoolBooking {
	if b.WebUrl == nil {
		mockURL := ""
		b.WebUrl = &mockURL
	}
	return &DriverCarpoolBooking{
		Car:            b.Car,
		Driver:         b.Driver,
		Price:          b.Price,
		CarpoolBooking: b.toCarpoolBooking(),
	}
}

// ToPassengerCarpoolBooking performs a lossy conversion from
// Booking to PassengerCarpoolBooking objects
func (b Booking) ToPassengerCarpoolBooking() *PassengerCarpoolBooking {
	return &PassengerCarpoolBooking{
		Passenger:      b.Passenger,
		CarpoolBooking: b.toCarpoolBooking(),
	}
}

func (b Booking) toCarpoolBooking() CarpoolBooking {
	if b.WebUrl == nil {
		mockURL := ""
		b.WebUrl = &mockURL
	}
	return CarpoolBooking{
		Distance:               b.Distance,
		Duration:               b.Duration,
		Id:                     b.Id,
		PassengerDropAddress:   b.PassengerDropAddress,
		PassengerDropLat:       b.PassengerDropLat,
		PassengerDropLng:       b.PassengerDropLng,
		PassengerPickupAddress: b.PassengerPickupAddress,
		PassengerPickupDate:    b.PassengerPickupDate,
		PassengerPickupLat:     b.PassengerPickupLat,
		PassengerPickupLng:     b.PassengerPickupLng,
		Status:                 CarpoolBookingStatus(b.Status),
		WebUrl:                 *b.WebUrl,
	}
}
