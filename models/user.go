package models

import "time"

type User struct {
    UserID       uint      `gorm:"primaryKey" json:"user_id"`

    PersonalID   string    `gorm:"unique;not null" json:"personal_id"`
    Username     string    `gorm:"unique;not null" json:"username"`
    PasswordHash string    `gorm:"not null" json:"password_hash"`
    Email        string    `gorm:"unique;not null" json:"email"`
    FullName     string    `gorm:"not null" json:"full_name"`

    Role         RoleType  `gorm:"type:role_type" json:"role"`

    Gender       string     `json:"gender,omitempty"`
    DateOfBirth  *time.Time `json:"date_of_birth,omitempty"`
    Height       *float64   `json:"height,omitempty"`
    Weight       *float64   `json:"weight,omitempty"`
    ContactInfo  string     `json:"contact_info,omitempty"`

    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
