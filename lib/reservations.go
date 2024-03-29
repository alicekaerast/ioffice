package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alicekaerast/ioffice/schema"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func (i *IOffice) ListReservations() {
	reservations := i.GetReservations()
	fmt.Printf("Upcoming reservations: %v\n", len(reservations))

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("ID", "Start", "Location Name", "Location ID", "Floor", "Building", "Checked In?")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, reservation := range reservations {
		unixTimeUTC := time.Unix(reservation.StartDate/1000, 0)
		tbl.AddRow(reservation.ID, unixTimeUTC.Format(time.RFC822), reservation.Room.Name, reservation.Room.ID, reservation.Room.Floor.Name, reservation.Room.Floor.Building.Name, reservation.CheckedIn)
	}
	tbl.Print()
}

func (i *IOffice) GetReservations() schema.Reservations {
	endpoint := "v2/reservations/?showOnlyMyReservations=true"
	body := i.Request("GET", endpoint, nil)
	reservations := schema.Reservations{}
	json.Unmarshal([]byte(body), &reservations)
	return reservations
}

func (i *IOffice) CheckIn(reservationID string) {
	endpoint := "v2/reservations/" + reservationID + "/checkIn"
	body := i.Request("PUT", endpoint, bytes.NewBuffer([]byte("")))
	checkinResponse := schema.CheckinResponse{}
	json.Unmarshal(body, &checkinResponse)
	if checkinResponse.Error != "" {
		fmt.Println(checkinResponse.ErrorDescription)
	} else {
		fmt.Printf("Checked In: %v\n", checkinResponse.CheckedIn)
	}
}

func (i *IOffice) CancelReservation(reservationID string) {
	endpoint := "v2/reservations/" + reservationID + "/cancel"
	body := i.Request("PUT", endpoint, bytes.NewBuffer([]byte("")))
	cancellationResponse := schema.CancellationResponse{}
	json.Unmarshal(body, &cancellationResponse)
	if cancellationResponse.Error != "" {
		fmt.Println(cancellationResponse.ErrorDescription)
	} else {
		fmt.Printf("Booking for %v %v\n", cancellationResponse.Room.Name, cancellationResponse.CancellationReason)
	}
}

func (i *IOffice) CreateReservation(user schema.User, roomID int, date time.Time) {
	endpoint := "v2/reservations"

	// {"guests":[],"notes":"","user":{"id":2409},"center":{"id":74},"room":{"id":672},"numberOfPeople":1,"startDate":1647011700000,"endDate":1647015300000,"allDay":false}
	reservationRequest := schema.ReservationRequest{
		Guests: nil,
		Notes:  "",
		User: struct {
			ID int `json:"id"`
		}{
			ID: user.ID,
		},
		Center: struct {
			ID int `json:"id"`
		}{
			ID: 74,
		},
		Room: struct {
			ID int `json:"id"`
		}{
			ID: roomID,
		},
		NumberOfPeople: 1,
		StartDate:      date.Unix() * 1000,
		EndDate:        date.Unix() * 1000,
		AllDay:         true,
	}

	jsonReservationRequest, _ := json.Marshal(reservationRequest)
	body := i.Request("POST", endpoint, bytes.NewBuffer(jsonReservationRequest))
	reservationResponse := schema.ReservationResponse{}
	json.Unmarshal(body, &reservationResponse)

	if reservationResponse.Error != "" {
		fmt.Println(reservationResponse.ErrorDescription)
	} else {
		fmt.Printf("Reserved: %v for %v\n", reservationResponse.Room.Name, reservationResponse.User.Name)
	}
}
