package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	MONGO_DATABASE      string `mapstructure:"MONGO_DATABASE"`
	MONGO_HOST          string `mapstructure:"MONGO_HOST"`
	MONGO_PORT          string `mapstructure:"MONGO_PORT"`
	MONGO_USER          string `mapstructure:"MONGO_USER"`
	MONGO_PASSWORD      string `mapstructure:"MONGO_PASSWORD"`
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
	Router  *fiber.App
	Config  ServiceConfigurations
	MongoDB *mongo.Client
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

type User struct {
	ID            primitive.ObjectID     `json:"id" bson:"_id" description:"User ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	Name          string                 `json:"name" bson:"name" validate:"required" description:"User Name" enabled:"true" label:"Name" editable:"true"`
	AccountNumber int                    `json:"accountNumber" bson:"accountNumber" length:"10" description:"Account Number" editable:"false" enabled:"false" label:"Account Number" unique:"true"`
	Email         string                 `json:"email" bson:"email" validate:"required,email" description:"Email" enabled:"true" label:"Email" editable:"true"`
	Password      string                 `json:"password" bson:"password" validate:"required" description:"Password" enabled:"true" label:"Password" editable:"true"`
	SaltStored    string                 `json:"saltStored" bson:"saltStored" description:"Salt Stored" enabled:"false" label:"Salt Stored" editable:"false"`
	Image         map[string]interface{} `json:"image" bson:"image" description:"Image" enabled:"true" label:"Image" editable:"true"`
	Active        bool                   `json:"active" bson:"active" description:"Active" enabled:"false" label:"Active" editable:"false"`
	Deleted       bool                   `json:"deleted" bson:"deleted" description:"Deleted" enabled:"false" label:"Deleted" editable:"false"`
	Region        string                 `json:"region" bson:"region" validate:"required" description:"Region" enabled:"true" label:"Region" editable:"true"`
	UserType      string                 `json:"userType" bson:"userType" validate:"required" description:"User Type" enabled:"true" label:"User Type" editable:"true"`
	CloudProvider string                 `json:"cloudProvider" bson:"cloudProvider" description:"Cloud Provider" enabled:"true" label:"Cloud Provider" editable:"true"`
	LastChangedBy string                 `json:"lastChangedBy" bson:"lastChangedBy" description:"Last Changed By" enabled:"false" label:"Last Changed By" editable:"false"`
	LastChangedOn time.Time              `json:"lastChangedOn" bson:"lastChangedOn" description:"Last Changed On" enabled:"false" label:"Last Changed On" editable:"false"`
	CreatedOn     time.Time              `json:"createdOn" bson:"createdOn" description:"Created On" enabled:"false" label:"Created On" editable:"false"`
	Personal      Personal               `json:"personal" bson:"personal" description:"Personal" enabled:"true" label:"Personal" editable:"true"`
	Organization  Organization           `json:"organization" bson:"organization" description:"Organization" enabled:"true" label:"Organization" editable:"true"`
	Settings      Settings               `json:"settings" bson:"settings" description:"Settings" enabled:"true" label:"Settings" editable:"true"`
	Address       Address                `json:"address" bson:"address" description:"Address" enabled:"true" label:"Address" editable:"true"`
}

type Personal struct {
	PhoneNumber int       `json:"phoneNumber" bson:"phoneNumber" length:"10" description:"Phone Number" editable:"true" enabled:"true" label:"Phone Number"`
	DateOfBirth time.Time `json:"dateOfBirth" bson:"dateOfBirth" description:"Date Of Birth" editable:"true" enabled:"true" label:"Date Of Birth"`
}

type Settings struct {
	Language     string `json:"language" bson:"language" description:"Language" editable:"true" enabled:"true" label:"Language"`
	CurrencyCode string `json:"currencyCode" bson:"currencyCode" description:"Currency Code" editable:"true" enabled:"true" label:"Currency Code"`
	DateFormat   string `json:"dateFormat" bson:"dateFormat" description:"Date Format" editable:"true" enabled:"true" label:"Date Format"`
}

type Organization struct {
	OrganizationName   string `json:"organizationName" bson:"organizationName" description:"Organization Name" editable:"true" enabled:"true" label:"Organization Name"`
	Industry           string `json:"industry" bson:"industry" description:"Industry" editable:"true" enabled:"true" label:"Industry"`
	ContactPerson      string `json:"contactPerson" bson:"contactPerson" description:"Contact Person" editable:"true" enabled:"true" label:"Contact Person"`
	ContactEmail       string `json:"contactEmail" bson:"contactEmail" description:"Contact Email" editable:"true" enabled:"true" label:"Contact Email"`
	Domain             string `json:"domain" bson:"domain" description:"Domain" editable:"true" enabled:"true" label:"Domain"`
	ContactPhoneNumber string `json:"contactPhoneNumber" bson:"contactPhoneNumber" length:"10" description:"Contact Phone Number" editable:"true" enabled:"true" label:"Contact Phone Number"`
}

type Address struct {
	Latitude  float64 `json:"latitude" bson:"latitude" description:"Latitude" editable:"true" enabled:"true" label:"Latitude"`
	Longitude float64 `json:"longitude" bson:"longitude" description:"Longitude" editable:"true" enabled:"true" label:"Longitude"`
	Village   string  `json:"village" bson:"village" description:"Village" editable:"true" enabled:"true" label:"Village"`
	District  string  `json:"district" bson:"district" description:"District" editable:"true" enabled:"true" label:"District"`
	State     string  `json:"state" bson:"state" description:"State" editable:"true" enabled:"true" label:"State"`
	Country   string  `json:"country" bson:"country" description:"Country" editable:"true" enabled:"true" label:"Country"`
	PinCode   int     `json:"pinCode" bson:"pinCode" length:"6" description:"Pin Code" editable:"true" enabled:"true" label:"Pin Code"`
	DoorNo    string  `json:"doorNo" bson:"doorNo" description:"Door No" editable:"true" enabled:"true" label:"Door No"`
	Street    string  `json:"street" bson:"street" description:"Street" editable:"true" enabled:"true" label:"Street"`
}

type UserOTP struct {
	OTP    int                `json:"otp" bson:"otp" required:"true"`
	Email  string             `json:"email" bson:"email" required:"true"`
	UserId string             `json:"userId" bson:"userId" required:"true"`
	ID     primitive.ObjectID `json:"id" bson:"_id" required:"true"`
	Expiry time.Time          `json:"expiry" bson:"expiry" required:"true"`
	Times  int                `json:"times" bson:"times" required:"true"`
}

type VerifyOTP struct {
	Email string `json:"email" bson:"email" required:"true"`
	OTP   int    `json:"otp" bson:"otp" required:"true"`
}

type SendOTPVerificationMail struct {
	Email string `json:"email" bson:"email" required:"true"`
}

type LoginCredentials struct {
	Email        string `json:"email" bson:"email" required:"true"`
	Password     string `json:"password" bson:"password" required:"true"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
}

type Tokens struct {
	UserId       string    `json:"userId" bson:"userId" required:"true"`
	Email        string    `json:"email" bson:"email" required:"true"`
	AccessToken  string    `json:"accessToken" bson:"accessToken" required:"true"`
	RefreshToken string    `json:"refreshToken" bson:"refreshToken" required:"true"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt" required:"true"`
	Expiry       int64     `json:"expiry" bson:"expiry" required:"true"`
}

type Password struct {
	UserId   primitive.ObjectID `json:"userId" bson:"userId"`
	Password string             `json:"password" bson:"password" required:"true"`
}
