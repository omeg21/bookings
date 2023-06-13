package models

import "time"

//User is the user model
type User struct{
	Id int
	FirstName string
	LastName string
	Email string
	Password string
	AccessLevel int
	CreatedAt  time.Time
	UpdatedAt time.Time
}

//Room is the room model 
type Room struct{
	Id int
	RoomName string
	CreatedAt  time.Time
	UpdatedAt time.Time
}

//Restriction is the Restrictions model
type Restriction struct{
	Id int
	RestrictionsName string
	CreatedAt  time.Time
	UpdatedAt time.Time
}
//Reservation is the Reservations model
type Reservation struct{
	Id int
	FirstName string
	LastName string
	Email string
	Phone string
	StartDate time.Time
	EndDate time.Time
	RoomId int
	CreatedAt  time.Time
	UpdatedAt time.Time
	Room Room
}

//RoomRestriction is the room restriction models
type RoomRestriction struct{
	Id int
	CreatedAt  time.Time
	UpdatedAt time.Time
	StartDate time.Time
	EndDate time.Time
	RoomId int
	Room Room
	ReservationID int
	Reservation Reservation
	RestrictionID int
	Restriction Restriction
}

//

