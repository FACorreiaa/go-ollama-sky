package structs

import (
	"database/sql/driver"
	"encoding/json"
	"log"
	"time"

	"errors"

	"github.com/google/uuid"
)

type FlightStatus string

const (
	Scheduled FlightStatus = "scheduled"
	Active    FlightStatus = "active"
	Landed    FlightStatus = "landed"
	Canceled  FlightStatus = "canceled"
	Incident  FlightStatus = "incident"
	Diverted  FlightStatus = "diverted"
)

type LiveFlights struct {
	ID           uuid.UUID    `json:"id" db:"id"`
	FlightDate   string       `json:"flight_date,omitempty"`
	FlightStatus FlightStatus `json:"flight_status,omitempty"`
	Departure    struct {
		Airport         string      `json:"airport"`
		Timezone        string      `json:"timezone"`
		Iata            string      `json:"iata"`
		Icao            string      `json:"icao"`
		Terminal        string      `json:"terminal"`
		Gate            interface{} `json:"gate"`
		Delay           *int        `json:"delay"`
		Scheduled       time.Time   `json:"scheduled"`
		Estimated       time.Time   `json:"estimated"`
		Actual          time.Time   `json:"actual"`
		EstimatedRunway time.Time   `json:"estimated_runway"`
		ActualRunway    time.Time   `json:"actual_runway"`
	} `json:"departure,omitempty"`
	Arrival struct {
		Airport         string      `json:"airport"`
		Timezone        string      `json:"timezone"`
		Iata            string      `json:"iata"`
		Icao            string      `json:"icao"`
		Terminal        interface{} `json:"terminal"`
		Gate            interface{} `json:"gate"`
		Baggage         interface{} `json:"baggage"`
		Delay           *int        `json:"delay"`
		Scheduled       time.Time   `json:"scheduled"`
		Estimated       time.Time   `json:"estimated"`
		Actual          time.Time   `json:"actual"`
		EstimatedRunway time.Time   `json:"estimated_runway"`
		ActualRunway    time.Time   `json:"actual_runway"`
	} `json:"arrival,omitempty"`
	Airline struct {
		Name string `json:"name"`
		Iata string `json:"iata"`
		Icao string `json:"icao"`
	} `json:"airline,omitempty"`
	Flight struct {
		Number     string `json:"number"`
		Iata       string `json:"iata"`
		Icao       string `json:"icao"`
		Codeshared struct {
			AirlineName  string `json:"airline_name"`
			AirlineIata  string `json:"airline_iata"`
			AirlineIcao  string `json:"airline_icao"`
			FlightNumber string `json:"flight_number"`
			FlightIata   string `json:"flight_iata"`
			FlightIcao   string `json:"flight_icao"`
		} `json:"codeshared,omitempty"`
	} `json:"flight"`
	Aircraft struct {
		AircraftRegistration string `json:"registration"`
		AircraftIata         string `json:"iata"`
		AircraftIcao         string `json:"icao"`
		AircraftIcao24       string `json:"icao24"`
	} `json:"aircraft,omitempty"`
	Live struct {
		LiveUpdated         string  `json:"updated"`
		LiveLatitude        float32 `json:"latitude,omitempty"`
		LiveLongitude       float32 `json:"longitude,omitempty"`
		LiveAltitude        float32 `json:"altitude"`
		LiveDirection       float32 `json:"direction"`
		LiveSpeedHorizontal float32 `json:"speed_horizontal"`
		LiveSpeedVertical   float32 `json:"speed_vertical"`
		LiveIsGround        bool    `json:"is_ground"`
	} `json:"live,omitempty"`
	CreatedAt CustomTime `json:"created_at"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	var dateStr string
	err := json.Unmarshal(data, &dateStr)
	if err != nil {
		log.Println("Error parsing date: ", err)
		return err
	}

	// Check if the date is "0000-00-00" and set it to a default value
	if dateStr == "0000-00-00" {
		ct.Time = time.Time{} // Assign zero value of CustomTime
		return nil
	}

	// Check if the date string is empty and set it to a default value
	if dateStr == "" {
		ct.Time = time.Time{} // Assign zero value of CustomTime
		return nil
	}

	// Parse the date using the predefined time layout
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		log.Println("Error parsing date: ", err)
		return err
	}

	ct.Time = t
	return nil
}

//

// Implement driver.Valuer interface.

func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time.Format(time.RFC3339), nil
}

// Implement sql.Scanner interface.

func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		// Handle NULL values by setting the time to the zero value
		ct.Time = time.Time{}
		return nil
	}

	switch t := value.(type) {
	case time.Time:
		ct.Time = t
		return nil
	case []byte:
		parsedTime, err := time.Parse("2006-01-02", string(t))
		if err != nil {
			return err
		}
		ct.Time = parsedTime
		return nil
	case string:
		parsedTime, err := time.Parse("2006-01-02", t)
		if err != nil {
			return err
		}
		ct.Time = parsedTime
		return nil
	default:
		return errors.New("unsupported Scan value for CustomTime")
	}
}
