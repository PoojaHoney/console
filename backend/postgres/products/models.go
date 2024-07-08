package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceConfigurations struct {
	SERVICE_PORT                string `mapstructure:"SERVICE_PORT"`
	SERVICE_NAME                string `mapstructure:"SERVICE_NAME"`
	SERVICE_DOMAIN              string `mapstructure:"SERVICE_DOMAIN"`
	SERVICE_HOST                string `mapstructure:"SERVICE_HOST"`
	SERVICE_BASEPATH            string `mapstructure:"SERVICE_BASEPATH"`
	SERVICE_VERSION             string `mapstructure:"SERVICE_VERSION"`
	TOKEN_TABLE                 string `mapstructure:"TOKEN_TABLE"`
	USERS_DATABASE              string `mapstructure:"USERS_DATABASE"`
	POSTGRES_DEFAULT_DATABASE   string `mapstructure:"POSTGRES_DEFAULT_DATABASE"`
	PRODUCT_TABLE               string `mapstructure:"PRODUCT_TABLE"`
	PRODUCT_VERSIONS_TABLE      string `mapstructure:"PRODUCT_VERSIONS_TABLE"`
	PRODUCT_MICROSERVICES_TABLE string `mapstructure:"PRODUCT_MICROSERVICES_TABLE"`
	PRODUCT_RESOURCES_TABLE     string `mapstructure:"PRODUCT_RESOURCES_TABLE"`
	PRODUCT_PLANS_TABLE         string `mapstructure:"PRODUCT_PLANS_TABLE"`
	PRODUCT_CONFIG_TABLE        string `mapstructure:"PRODUCT_CONFIG_TABLE"`
	POSTGRES_DATABASE           string `mapstructure:"POSTGRES_DATABASE"`
	POSTGRES_HOST               string `mapstructure:"POSTGRES_HOST"`
	POSTGRES_PORT               string `mapstructure:"POSTGRES_PORT"`
	POSTGRES_USERNAME           string `mapstructure:"POSTGRES_USERNAME"`
	POSTGRES_PASSWORD           string `mapstructure:"POSTGRES_PASSWORD"`
	JAEGER_URL                  string `mapstructure:"JAEGER_URL"`
	TRACER_ENABLED              bool   `mapstructure:"TRACER_ENABLED"`
	SMTP_SERVER                 string `mapstructure:"SMTP_SERVER"`
	SMTP_BASEEMAIL              string `mapstructure:"SMTP_BASEEMAIL"`
	SMTP_PASSKEY                string `mapstructure:"SMTP_PASSKEY"`
	SMTP_PORT                   string `mapstructure:"SMTP_PORT"`
	SUPPORT_EMAIL               string `mapstructure:"SUPPORT_EMAIL"`
	SUPPORT_PHONE               string `mapstructure:"SUPPORT_PHONE"`
	PRODUCT_FULL_DETAILS        string `mapstructure:"PRODUCT_FULL_DETAILS"`
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

type Service struct {
	Router    *fiber.App
	Config    ServiceConfigurations
	ProductDB *gorm.DB
	UsersDB   *gorm.DB
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

type FullProductDetails struct {
	DeletedAt          gorm.DeletedAt        `gorm:"index" json:"deletedAt"`
	ID                 uuid.UUID             `json:"id" gorm:"type:uuid;primaryKey" description:"Product ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	Name               string                `json:"name" gorm:"not null" description:"Product Name" required:"true" editable:"true" enabled:"true" label:"Name" unique:"true"`
	Description        string                `json:"description" gorm:"not null" description:"Product Description" required:"true" editable:"true" enabled:"true" label:"Description"`
	Status             string                `json:"status" gorm:"not null" description:"Product Status" editable:"false" enabled:"true" label:"Status"`
	Type               string                `json:"type" gorm:"not null" description:"Product Type" editable:"true" enabled:"true" required:"true" label:"Type"`
	ProductID          string                `json:"productID" gorm:"uniqueIndex;not null"  description:"Product Product ID" editable:"true" required:"true" enabled:"true" label:"Product ID"`
	Image              string                `json:"image" gorm:"not null" description:"Product Image" editable:"true" enabled:"true" label:"Image"`
	Providers          []ProductProvider     `json:"providers" gorm:"type:text[]" description:"Cloud Providers" editable:"true" required:"true" enabled:"true" label:"Providers"`
	MicroServicesCount int                   `json:"microServicesCount" gorm:"not null" description:"Number of Microservices" editable:"false" enabled:"false" label:"Microservices Count"`
	CreatedAt          time.Time             `json:"createdAt" description:"Created On" editable:"false" enabled:"true" label:"Created On"`
	UpdatedAt          time.Time             `json:"updatedAt" description:"Updated On" editable:"false" enabled:"true" label:"Updated On"`
	DatabasesCount     int                   `json:"databasesCount" gorm:"not null" description:"Number of Databases" editable:"false" enabled:"true" label:"Databases Count"`
	CreatedBy          string                `json:"createdBy" gorm:"not null" description:"Created By" editable:"false" enabled:"true" label:"Created By"`
	MicroServices      []ProductMicroService `json:"microServices" gorm:"foreignKey:productID"`
	Resources          []ProductResource     `json:"resources" gorm:"foreignKey:productID"`
	Configuration      ProductConfiguration  `json:"configurations" gorm:"foreignKey:productID"`
	Plans              []ProductPlan         `json:"plans" gorm:"foreignKey:productID"`
	Versions           []ProductVersion      `json:"versions" gorm:"foreignKey:productID"`
}

type Product struct {
	ID                 uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" description:"Product ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	Name               string            `json:"name" gorm:"not null" description:"Product Name" required:"true" editable:"true" enabled:"true" label:"Name" unique:"true"`
	Description        string            `json:"description" gorm:"not null" description:"Product Description" required:"true" editable:"true" enabled:"true" label:"Description"`
	Status             string            `json:"status" gorm:"not null" description:"Product Status" editable:"false" enabled:"true" label:"Status"`
	Type               string            `json:"type" gorm:"not null" description:"Product Type" editable:"true" enabled:"true" required:"true" label:"Type"`
	ProductID          string            `json:"productID" gorm:"uniqueIndex;not null" description:"Product Product ID" editable:"true" required:"true" enabled:"true" label:"Product ID"`
	Image              string            `json:"image" gorm:"not null" description:"Product Image" editable:"true" enabled:"true" label:"Image"`
	Providers          []ProductProvider `gorm:"foreignKey:ProductID;not null" json:"providers" description:"Cloud Providers" editable:"true" required:"true" enabled:"true" label:"Providers"`
	MicroServicesCount int               `json:"microServicesCount" gorm:"not null" description:"Number of Microservices" editable:"false" enabled:"false" label:"Microservices Count"`
	CreatedAt          time.Time         `json:"createdAt" description:"Created On" editable:"false" enabled:"true" label:"Created On"`
	DeletedAt          gorm.DeletedAt    `gorm:"index" json:"deletedAt"`
	UpdatedAt          time.Time         `json:"updatedAt" description:"Updated On" editable:"false" enabled:"true" label:"Updated On"`
	DatabasesCount     int               `json:"databasesCount" gorm:"not null" description:"Number of Databases" editable:"false" enabled:"true" label:"Databases Count"`
	CreatedBy          string            `json:"createdBy" gorm:"not null" description:"Created By" editable:"false" enabled:"true" label:"Created By"`
}

type ProductProvider struct {
	gorm.Model
	ProductID string `json:"productID" gorm:"not null" `
	Provider  string `json:"provider" gorm:"not null" `
}

type ProductMicroServiceDatabase struct {
	MicroserviceID  uuid.UUID      `json:"microserviceID" gorm:"not null"`
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey" description:"ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	CreatedAt       time.Time      `json:"createdAt" description:"Created On" editable:"false" enabled:"true" label:"Created On"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	UpdatedAt       time.Time      `json:"updatedAt" description:"Updated On" editable:"false" enabled:"true" label:"Updated On"`
	Type            string         `json:"type" description:"Database Type" editable:"true" enabled:"true" label:"Type"`
	Database        string         `json:"database" description:"Database Name" editable:"true" enabled:"true" label:"Database"`
	Name            string         `json:"name" description:"Database Name" editable:"true" enabled:"true" label:"Name"`
	PortNumber      string         `json:"portNumber" description:"Database Port Number" editable:"true" enabled:"true" label:"Port Number"`
	DefaultUsername string         `json:"defaultUsername" description:"Database Default Username" editable:"true" enabled:"true" label:"Default Username"`
	DefaultPassword string         `json:"defaultPassword" description:"Database Default Password" editable:"true" enabled:"true" label:"Default Password"`
	Host            string         `json:"host" description:"Database Host" editable:"true" enabled:"true" label:"Host"`
	Version         string         `json:"version" description:"Database Version" editable:"true" enabled:"true" label:"Version"`
	Image           string         `json:"image" description:"Image" editable:"true" enabled:"true" label:"Image"`
	MockDataPath    string         `json:"mockDataPath" description:"Database Mock Data Path" editable:"true" enabled:"true" label:"Mock Data Path"`
}

type ProductMicroService struct {
	ID         uuid.UUID                     `json:"id" gorm:"type:uuid;primaryKey" description:"Microservice ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	Name       string                        `json:"name" description:"Microservice Name" required:"true" editable:"true" enabled:"true" label:"Name"`
	Type       string                        `json:"type" description:"Microservice Type" editable:"true" required:"true" enabled:"true" label:"Type"`
	ProductID  string                        `json:"productID" required:"true" description:"Product ID" editable:"true" enabled:"true" label:"Product ID"`
	PortNumber string                        `json:"portNumber" description:"Microservice Port Number" editable:"true" required:"true" enabled:"true" label:"Port Number"`
	Status     string                        `json:"status" description:"Product Status" editable:"true" enabled:"true" label:"Status"`
	Host       string                        `json:"host" description:"Microservice Host" editable:"true" enabled:"true" label:"Host"`
	BasePath   string                        `json:"basePath" required:"true" description:"Microservice Base Path" editable:"true" enabled:"true" label:"Base Path"`
	Version    string                        `json:"version" required:"true" description:"Microservice Version" editable:"true" enabled:"true" label:"Version"`
	DeletedAt  gorm.DeletedAt                `gorm:"index" json:"deletedAt"`
	CreatedAt  time.Time                     `json:"createdAt" description:"Created On" editable:"false" enabled:"false" label:"Created On"`
	UpdatedAt  time.Time                     `json:"updatedAt" description:"Updated On" editable:"false" enabled:"false" label:"Updated On"`
	Databases  []ProductMicroServiceDatabase `gorm:"foreignKey:MicroserviceID" json:"productDatabases" description:"Product Databases" editable:"true" enabled:"true" label:"Databases"`
}

type ProductResourceVersions struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey" description:"ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	ResourceID     uuid.UUID      `json:"resourceID" gorm:"not null"`
	ResourceName   string         `json:"resourceName" description:"Resource Name" editable:"true" enabled:"true" label:"Resource Name"`
	DevelopmentTag string         `json:"developmentTag" description:"Development Tag" editable:"true" enabled:"true" label:"Development Tag"`
	ProductTag     string         `json:"productTag" description:"Product Tag" editable:"true" enabled:"true" label:"Product Tag"`
	StagingTag     string         `json:"stagingTag" description:"Staging Tag" editable:"true" enabled:"true" label:"Staging Tag"`
	TestingTag     string         `json:"testingTag" description:"Testing Tag" editable:"true" enabled:"true" label:"Testing Tag"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	CreatedAt      string         `json:"createdAt" description:"Created On" editable:"true" enabled:"true" label:"Created On"`
	UpdatedAt      string         `json:"updatedAt" description:"Updated On" editable:"true" enabled:"true" label:"Updated On"`
	Latest         bool           `json:"latest" description:"Latest" editable:"true" enabled:"true" label:"Latest"`
}

type ProductResource struct {
	ID                         uuid.UUID                            `json:"id" gorm:"type:uuid;primaryKey" description:"Resource ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	Kind                       string                               `json:"kind" required:"true" description:"Resource Kind" editable:"true" enabled:"true" label:"Resource Kind"`
	TaskType                   string                               `json:"taskType" required:"true" description:"Task Type" editable:"true" enabled:"true" label:"Task Type"`
	DeploymentTemplateName     string                               `json:"deploymentTemplateName" description:"Deployment Template Name" editable:"true" enabled:"true" label:"Deployment Template Name"`
	DeploymentTemplateLocation string                               `json:"deploymentTemplateLocation" description:"Deployment Template Location" editable:"true" enabled:"true" label:"Deployment Template Location"`
	Status                     string                               `json:"status" description:"Product Status" editable:"true" enabled:"true" label:"Status"`
	Name                       string                               `json:"name" description:"Name" required:"true" editable:"true" enabled:"true" label:"Name"`
	ProductID                  string                               `json:"productID" required:"true" description:"Product ID" editable:"true" enabled:"true" label:"Product ID"`
	DeletedAt                  gorm.DeletedAt                       `gorm:"index" json:"deletedAt"`
	CreatedAt                  time.Time                            `json:"createdAt" description:"Created On" editable:"false" enabled:"false" label:"Created On"`
	UpdatedAt                  time.Time                            `json:"updatedAt" description:"Updated On" editable:"false" enabled:"false" label:"Updated On"`
	Versions                   []ProductResourceVersions            `json:"versions" gorm:"foreignKey:ResourceID;not null" description:"Resource Versions" editable:"true" enabled:"true" label:"Resource Versions"`
	EnvironmentVariablesFile   string                               `json:"environmentVariablesFile" description:"Environment Variables File" editable:"true" enabled:"true" label:"Environment Variables File"`
	ProductVersion             string                               `json:"productVersion" description:"Product Version" editable:"true" required:"true" enabled:"true" label:"Product Version"`
	ExposedEnvVariables        []ProductResourceExposedEnvVariables `json:"exposedEnvVariables" gorm:"foreignKey:ResourceID;not null" description:"Exposed Environment Variables" editable:"true" enabled:"true" label:"Exposed Environment Variables"`
}

type ProductResourceExposedEnvVariables struct {
	gorm.Model
	ResourceID         uuid.UUID `json:"ResourceID" gorm:"not null" `
	ExposedEnvVariable string    `json:"exposedEnvVariable" gorm:"not null" `
}

type ProductConfiguration struct {
	ID                    uuid.UUID                                 `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();uniqueIndex" description:"Configuration ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	ProductID             string                                    `json:"productID" required:"true" description:"Product ID" editable:"true" enabled:"true" label:"Product ID"`
	DefaultMemory         float32                                   `json:"defaultMemory" required:"true" description:"Default Memory" editable:"true" enabled:"true" label:"Default Memory"`
	DefaultRAM            float32                                   `json:"defaultRAM" required:"true" description:"Default RAM" editable:"true" enabled:"true" label:"Default RAM"`
	StartupScriptFilePath string                                    `json:"startupScriptFilePath" description:"Startup Script File Path" editable:"true" enabled:"true" label:"Startup Script File Path"`
	Status                string                                    `json:"status" description:"Product Status" editable:"true" enabled:"true" label:"Status"`
	ArtifactRegistryName  string                                    `json:"artifactRegistryName" required:"true" description:"Artifact Registry Name" editable:"true" enabled:"true" label:"Artifact Registry Name"`
	NetworkTags           []ProductConfigurationNetworkTags         `json:"networkTags" gorm:"foreignKey:ProductID;references:ProductID;not null" description:"Network Tags" editable:"true" enabled:"true" label:"Network Tags"`
	DeletedAt             gorm.DeletedAt                            `gorm:"index" json:"deletedAt"`
	CreatedAt             time.Time                                 `json:"createdAt" description:"Created On" editable:"false" enabled:"false" label:"Created On"`
	UpdatedAt             time.Time                                 `json:"updatedAt" description:"Updated On" editable:"false" enabled:"false" label:"Updated On"`
	ProviderPermissions   []ProductConfigurationProviderPermissions `json:"providerPermissions" gorm:"foreignKey:ProductID;references:ProductID;not null"`
	EnvironmentsSupport   []ProductEnvironmentSupport               `gorm:"foreignKey:ProductID;references:ProductID;not null" json:"environmentsSupport" description:"Environments Support" editable:"true" enabled:"true" label:"Environments Support"`
}

type ProductProviderPermissions struct {
	gorm.Model
	Provider   string `json:"provider" gorm:"not null" `
	Permission string `json:"permission" gorm:"not null" `
}

type ProductConfigurationProviderPermissions struct {
	ID          uuid.UUID                    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();uniqueIndex" description:"ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	CreatedAt   time.Time                    `json:"createdAt" description:"Created On" editable:"false" enabled:"true" label:"Created On"`
	DeletedAt   gorm.DeletedAt               `gorm:"index;not null" json:"deletedAt"`
	UpdatedAt   time.Time                    `json:"updatedAt" description:"Updated On" editable:"false" enabled:"true" label:"Updated On"`
	ProductID   string                       `json:"productID" gorm:"not null;uniqueIndex"`
	Permissions []ProductProviderPermissions `gorm:"foreignKey:Provider;not null" json:"permissions" description:"Permissions" editable:"true" enabled:"true" label:"Permissions"`
	Provider    string                       `json:"provider" gorm:"not null" description:"Provider" editable:"true" enabled:"true" label:"Provider"`
	Enabled     bool                         `json:"enabled" gorm:"not null" description:"Enabled" editable:"true" enabled:"true" label:"Enabled"`
}
type ProductConfigurationNetworkTags struct {
	gorm.Model
	ProductID  string `json:"productID" gorm:"not null" `
	NetworkTag string `json:"networkTag" gorm:"not null" `
}

type ProductEnvironmentSupport struct {
	ID          uint   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();uniqueIndex"`
	ProductID   string `json:"productID" gorm:"not null"`
	Provider    string `json:"provider" gorm:"not null" description:"Provider" editable:"true" enabled:"true" label:"Provider"`
	Enabled     bool   `json:"enabled" gorm:"not null" description:"Enabled" editable:"true" enabled:"true" label:"Enabled"`
	Environment string `json:"environment" gorm:"not null" description:"Environment Name" editable:"true" enabled:"true" label:"Environment Name"`
}

type ProductPlan struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey" description:"Plan ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	ProductID   string         `json:"productID" required:"true" description:"Product ID" editable:"true" enabled:"true" label:"Product ID"`
	Plan        string         `json:"plan" description:"Plan" required:"true" editable:"true" enabled:"true" label:"Plan"`
	Name        string         `json:"name" description:"Name" required:"true" editable:"true" enabled:"true" label:"Name"`
	Description string         `json:"description" required:"true" description:"Description" editable:"true" enabled:"true" label:"Description"`
	Status      string         `json:"status" description:"Product Status" editable:"true" enabled:"true" label:"Status"`
	Active      string         `json:"active" description:"Active" editable:"true" enabled:"true" label:"Active"`
	Image       string         `json:"image" description:"Image" editable:"true" enabled:"true" label:"Image"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	CreatedAt   time.Time      `json:"createdAt" description:"Created On" editable:"false" enabled:"false" label:"Created On"`
	UpdatedAt   time.Time      `json:"updatedAt" description:"Updated On" editable:"false" enabled:"false" label:"Updated On"`
}

type ProductVersion struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" description:"Plan ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	Name        string         `json:"name" description:"Name" editable:"true" required:"true" enabled:"true" label:"Name"`
	Version     string         `json:"version" description:"Version" editable:"true" required:"true" enabled:"true" label:"Version"`
	BuildNumber string         `json:"buildNumber" description:"Build Number" editable:"true" required:"true" enabled:"true" label:"Build Number"`
	ProductID   string         `json:"productID" description:"Product ID" editable:"true" required:"true" enabled:"true" label:"Product ID"`
	Status      string         `json:"status" description:"Product Status" editable:"false" enabled:"true" label:"Status"`
	Description string         `json:"description" description:"Description" required:"true" editable:"true" enabled:"true" label:"Description"`
	Type        string         `json:"type" description:"Type" editable:"true" required:"true" enabled:"true" label:"Type"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	CreatedAt   time.Time      `json:"createdAt" description:"Created On" editable:"false" enabled:"false" label:"Created On"`
	UpdatedAt   time.Time      `json:"updatedAt" description:"Updated On" editable:"false" enabled:"false" label:"Updated On"`
}
