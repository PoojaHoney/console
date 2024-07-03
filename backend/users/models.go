package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceConfigurations struct {
	SERVICE_PORT        string `mapstructure:"SERVICE_PORT"`
	SERVICE_NAME        string `mapstructure:"SERVICE_NAME"`
	SERVICE_DOMAIN      string `mapstructure:"SERVICE_DOMAIN"`
	SERVICE_HOST        string `mapstructure:"SERVICE_HOST"`
	SERVICE_BASEPATH    string `mapstructure:"SERVICE_BASEPATH"`
	SERVICE_VERSION     string `mapstructure:"SERVICE_VERSION"`
	TOKEN_COLLECTION    string `mapstructure:"TOKEN_COLLECTION"`
	USERS_COLLECTION    string `mapstructure:"USERS_COLLECTION"`
	USER_OTP_COLLECTION string `mapstructure:"USER_OTP_COLLECTION"`
	POSTGRES_DATABASE   string `mapstructure:"POSTGRES_DATABASE"`
	POSTGRES_HOST       string `mapstructure:"POSTGRES_HOST"`
	POSTGRES_PORT       string `mapstructure:"POSTGRES_PORT"`
	POSTGRES_USERNAME   string `mapstructure:"POSTGRES_USERNAME"`
	POSTGRES_PASSWORD   string `mapstructure:"POSTGRES_PASSWORD"`
	JAEGER_URL          string `mapstructure:"JAEGER_URL"`
	TRACER_ENABLED      bool   `mapstructure:"TRACER_ENABLED"`
	SMTP_SERVER         string `mapstructure:"SMTP_SERVER"`
	SMTP_BASEEMAIL      string `mapstructure:"SMTP_BASEEMAIL"`
	SMTP_PASSKEY        string `mapstructure:"SMTP_PASSKEY"`
	SMTP_PORT           string `mapstructure:"SMTP_PORT"`
	SUPPORT_EMAIL       string `mapstructure:"SUPPORT_EMAIL"`
	SUPPORT_PHONE       string `mapstructure:"SUPPORT_PHONE"`
}

type Service struct {
	Router     *fiber.App
	Config     ServiceConfigurations
	PostgresDB *gorm.DB
}

type FieldCatalogue struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Editable    bool   `json:"editable"`
	Enabled     bool   `json:"enabled"`
	Label       string `json:"label"`
	Unique      bool   `json:"unique"`
	Length      int    `json:"length"`
}

type AllFieldCatalogues struct {
	Fields      []FieldCatalogue       `json:"fields"`
	Table       string                 `json:"table"`
	ChildTables map[string]interface{} `json:"childTables"`
	Description string                 `json:"description"`
}

type Personal struct {
	PhoneNumber int       `json:"phoneNumber" gorm:"not null" length:"10" description:"Phone Number"`
	DateOfBirth time.Time `json:"dateOfBirth" gorm:"not null" description:"Date Of Birth"`
}

type Settings struct {
	Language     string `json:"language" gorm:"not null" description:"Language"`
	CurrencyCode string `json:"currencyCode" gorm:"not null" description:"Currency Code"`
	DateFormat   string `json:"dateFormat" gorm:"not null" description:"Date Format"`
}

type Organization struct {
	OrganizationName   string `json:"organizationName" gorm:"not null" description:"Organization Name"`
	Industry           string `json:"industry" gorm:"not null" description:"Industry"`
	ContactPerson      string `json:"contactPerson" gorm:"not null" description:"Contact Person"`
	ContactEmail       string `json:"contactEmail" gorm:"not null" description:"Contact Email"`
	Domain             string `json:"domain" gorm:"not null" description:"Domain"`
	ContactPhoneNumber string `json:"contactPhoneNumber" gorm:"not null" length:"10" description:"Contact Phone Number"`
}

type Address struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4()" description:"Address ID"`
	UserID    uuid.UUID      `gorm:"not null" json:"userId" description:"User ID"` // Foreign key for User
	Latitude  float64        `json:"latitude" gorm:"not null" description:"Latitude"`
	Longitude float64        `json:"longitude" gorm:"not null" description:"Longitude"`
	Village   string         `json:"village" gorm:"not null" description:"Village"`
	District  string         `json:"district" gorm:"not null" description:"District"`
	State     string         `json:"state" gorm:"not null" description:"State"`
	Country   string         `json:"country" gorm:"not null" description:"Country"`
	PinCode   int            `json:"pinCode" gorm:"not null" length:"6" description:"Pin Code"`
	DoorNo    string         `json:"doorNo" gorm:"not null" description:"Door No"`
	Street    string         `json:"street" gorm:"not null" description:"Street"`
}

type UserOTP struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4()" description:"UserOTP ID"`
	OTP       int            `json:"otp" gorm:"not null"`
	Email     string         `json:"email" gorm:"not null"`
	UserId    uuid.UUID      `json:"userId" gorm:"not null"`
	Expiry    time.Time      `json:"expiry" gorm:"not null"`
	Times     int            `json:"times" gorm:"not null"`
}

type VerifyOTP struct {
	Email string `json:"email" gorm:"not null"`
	OTP   int    `json:"otp" gorm:"not null"`
}

type SendOTPVerificationMail struct {
	Email string `json:"email" gorm:"not null"`
}

type LoginCredentials struct {
	Email        string `json:"email" gorm:"not null"`
	Password     string `json:"password" gorm:"not null"`
	RefreshToken string `json:"refreshToken"`
}

type Tokens struct {
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4()" description:"UserOTP ID"`
	UserId       uuid.UUID      `json:"userId" gorm:"not null"`
	Email        string         `json:"email" gorm:"not null"`
	AccessToken  string         `json:"accessToken" gorm:"not null"`
	RefreshToken string         `json:"refreshToken" gorm:"not null"`
	Expiry       int64          `json:"expiry" gorm:"not null"`
}

type Password struct {
	UserId   uuid.UUID `json:"userId" gorm:"not null"`
	Password string    `json:"password" gorm:"not null"`
}

type User struct {
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4()" description:"User ID"`
	Name          string         `json:"name" gorm:"not null" validate:"required" description:"User Name"`
	AccountNumber int            `json:"accountNumber" gorm:"unique;not null" length:"10" description:"Account Number"`
	Email         string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email" description:"Email"`
	Password      string         `json:"password" gorm:"not null" validate:"required" description:"Password"`
	SaltStored    string         `json:"saltStored" description:"Salt Stored"`
	Image         string         `json:"image" gorm:"-" description:"Image"` // Ignored by GORM
	Active        bool           `json:"active" gorm:"default:false" description:"Active"`
	Deleted       bool           `json:"deleted" gorm:"default:false" description:"Deleted"`
	Region        string         `json:"region" gorm:"not null" validate:"required" description:"Region"`
	UserType      string         `json:"userType" gorm:"not null" validate:"required" description:"User Type"`
	CloudProvider string         `json:"cloudProvider" description:"Cloud Provider"`
	LastChangedBy string         `json:"lastChangedBy" description:"Last Changed By"`
	CreatedOn     time.Time      `json:"createdOn" gorm:"autoCreateTime" description:"Created On"`
	Personal      Personal       `json:"personal" gorm:"embedded" description:"Personal"`
	Organization  Organization   `json:"organization" gorm:"embedded" description:"Organization"`
	Settings      Settings       `json:"settings" gorm:"embedded" description:"Settings"`
	Address       Address        `json:"address" gorm:"foreignKey:UserID" description:"Address"`
}
