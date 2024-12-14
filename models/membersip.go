package models

import "time"

type Privilege struct {
	Name    string `json:"name"`
	Allowed bool   `json:"allowed"`
}

type Membership struct {
	ID             int         `json:"id"`
	UserID         int         `json:"user_id"`         // ID пользователя, которому принадлежит абонемент
	StartDate      time.Time   `json:"start_date"`      // Дата начала абонемента
	EndDate        time.Time   `json:"end_date"`        // Дата окончания абонемента
	MembershipName string      `json:"membership_name"` // Название абонемента
	Price          float64     `json:"price"`           // Цена абонемента
	Status         string      `json:"status"`
	Description    string      `json:"description"`   // Описание абонемента
	BookingHours   string      `json:"booking_hours"` // Часы бронирования
	Privileges     []Privilege `json:"privileges"`    // Привилегии абонемента
}
