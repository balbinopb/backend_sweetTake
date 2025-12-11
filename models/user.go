package models

import "time"

type User struct {
    UserID       uint      `gorm:"primaryKey" json:"user_id"`

    // PersonalID   string    `gorm:"unique;not null" json:"personal_id"`
    Username     string    `gorm:"unique;not null" json:"username"`
    Email        string    `gorm:"unique;not null" json:"email"`

    
    
    FullName    *string    `json:"full_name,omitempty"`
    Password    string     `json:"-" gorm:"size:255;not null"`


    Gender       string     `json:"gender,omitempty"`
    DateOfBirth  *time.Time `json:"date_of_birth,omitempty"`
    Height       *float64   `json:"height,omitempty"`
    Weight       *float64   `json:"weight,omitempty"`
    ContactInfo  string     `json:"contact_info,omitempty"`

    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
