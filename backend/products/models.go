package main

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ServiceConfigurations struct {
	SERVICE_PORT                     string `mapstructure:"SERVICE_PORT"`
	SERVICE_NAME                     string `mapstructure:"SERVICE_NAME"`
	SERVICE_DOMAIN                   string `mapstructure:"SERVICE_DOMAIN"`
	SERVICE_HOST                     string `mapstructure:"SERVICE_HOST"`
	SERVICE_BASEPATH                 string `mapstructure:"SERVICE_BASEPATH"`
	SERVICE_VERSION                  string `mapstructure:"SERVICE_VERSION"`
	TOKEN_COLLECTION                 string `mapstructure:"TOKEN_COLLECTION"`
	USERS_COLLECTION                 string `mapstructure:"USERS_COLLECTION"`
	PRODUCT_COLLECTION               string `mapstructure:"PRODUCT_COLLECTION"`
	PRODUCT_VERSIONS_COLLECTION      string `mapstructure:"PRODUCT_VERSIONS_COLLECTION"`
	PRODUCT_MICROSERVICES_COLLECTION string `mapstructure:"PRODUCT_MICROSERVICES_COLLECTION"`
	PRODUCT_RESOURCES_COLLECTION     string `mapstructure:"PRODUCT_RESOURCES_COLLECTION"`
	PRODUCT_PLANS_COLLECTION         string `mapstructure:"PRODUCT_PLANS_COLLECTION"`
	PRODUCT_CONFIG_COLLECTION        string `mapstructure:"PRODUCT_CONFIG_COLLECTION"`
	MONGO_DATABASE                   string `mapstructure:"MONGO_DATABASE"`
	MONGO_HOST                       string `mapstructure:"MONGO_HOST"`
	MONGO_PORT                       string `mapstructure:"MONGO_PORT"`
	MONGO_USER                       string `mapstructure:"MONGO_USER"`
	MONGO_PASSWORD                   string `mapstructure:"MONGO_PASSWORD"`
	JAEGER_URL                       string `mapstructure:"JAEGER_URL"`
	TRACER_ENABLED                   bool   `mapstructure:"TRACER_ENABLED"`
	SMTP_SERVER                      string `mapstructure:"SMTP_SERVER"`
	SMTP_BASEEMAIL                   string `mapstructure:"SMTP_BASEEMAIL"`
	SMTP_PASSKEY                     string `mapstructure:"SMTP_PASSKEY"`
	SMTP_PORT                        string `mapstructure:"SMTP_PORT"`
	SUPPORT_EMAIL                    string `mapstructure:"SUPPORT_EMAIL"`
	SUPPORT_PHONE                    string `mapstructure:"SUPPORT_PHONE"`
	PRODUCT_FULL_DETAILS             string `mapstructure:"PRODUCT_FULL_DETAILS"`
}

type Tokens struct {
	UserId       string    `json:"userId" bson:"userId" required:"true"`
	Email        string    `json:"email" bson:"email" required:"true"`
	AccessToken  string    `json:"accessToken" bson:"accessToken" required:"true"`
	RefreshToken string    `json:"refreshToken" bson:"refreshToken" required:"true"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt" required:"true"`
	Expiry       int64     `json:"expiry" bson:"expiry" required:"true"`
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

type FullProductDetails struct {
	ID                 primitive.ObjectID    `json:"id" bson:"_id" description:"Product ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	Name               string                `json:"name" bson:"name" description:"Product Name" required:"true" editable:"true" enabled:"true" label:"Name" unique:"true"`
	Description        string                `json:"description" bson:"description" description:"Product Description" required:"true" editable:"true" enabled:"true" label:"Description"`
	Status             string                `json:"status" bson:"status" description:"Product Status" editable:"false" enabled:"true" label:"Status"`
	Type               string                `json:"type" bson:"type" description:"Product Type" editable:"true" enabled:"true" required:"true" label:"Type"`
	ProductID          string                `json:"productID" bson:"productID" description:"Product Product ID" editable:"true" required:"true" enabled:"true" label:"Product ID"`
	Image              string                `json:"image" bson:"image" description:"Product Image" editable:"true" enabled:"true" label:"Image"`
	Providers          []string              `json:"providers" bson:"providers" description:"Cloud Providers" editable:"true" required:"true" enabled:"true" label:"Providers"`
	MicroServicesCount int                   `json:"microServicesCount" bson:"microServicesCount" description:"Number of Microservices" editable:"false" enabled:"false" label:"Microservices Count"`
	CreatedOn          time.Time             `json:"createdOn" bson:"createdOn" description:"Created On" editable:"false" enabled:"true" label:"Created On"`
	UpdatedOn          time.Time             `json:"updatedOn" bson:"updatedOn" description:"Updated On" editable:"false" enabled:"true" label:"Updated On"`
	DatabasesCount     int                   `json:"databasesCount" bson:"databasesCount" description:"Number of Databases" editable:"false" enabled:"true" label:"Databases Count"`
	CreatedBy          string                `json:"createdBy" bson:"createdBy" description:"Created By" editable:"false" enabled:"true" label:"Created By"`
	MicroServices      []ProductMicroService `json:"microServices"`
	Resources          []ProductResource     `json:"resources"`
	Configuration      ProductConfiguration  `json:"configurations"`
	Plans              []ProductPlan         `json:"plans"`
	Versions           []ProductVersion      `json:"versions"`
}

type Product struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id" description:"Product ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	Name               string             `json:"name" bson:"name" description:"Product Name" required:"true" editable:"true" enabled:"true" label:"Name" unique:"true"`
	Description        string             `json:"description" bson:"description" description:"Product Description" required:"true" editable:"true" enabled:"true" label:"Description"`
	Status             string             `json:"status" bson:"status" description:"Product Status" editable:"false" enabled:"true" label:"Status"`
	Type               string             `json:"type" bson:"type" description:"Product Type" editable:"true" enabled:"true" required:"true" label:"Type"`
	ProductID          string             `json:"productID" bson:"productID" description:"Product Product ID" editable:"true" required:"true" enabled:"true" label:"Product ID"`
	Image              string             `json:"image" bson:"image" description:"Product Image" editable:"true" enabled:"true" label:"Image"`
	Providers          []string           `json:"providers" bson:"providers" description:"Cloud Providers" editable:"true" required:"true" enabled:"true" label:"Providers"`
	MicroServicesCount int                `json:"microServicesCount" bson:"microServicesCount" description:"Number of Microservices" editable:"false" enabled:"false" label:"Microservices Count"`
	CreatedOn          time.Time          `json:"createdOn" bson:"createdOn" description:"Created On" editable:"false" enabled:"true" label:"Created On"`
	UpdatedOn          time.Time          `json:"updatedOn" bson:"updatedOn" description:"Updated On" editable:"false" enabled:"true" label:"Updated On"`
	DatabasesCount     int                `json:"databasesCount" bson:"databasesCount" description:"Number of Databases" editable:"false" enabled:"true" label:"Databases Count"`
	CreatedBy          string             `json:"createdBy" bson:"createdBy" description:"Created By" editable:"false" enabled:"true" label:"Created By"`
}

type ProductMicroServiceDatabase struct {
	Type            string `json:"type" bson:"type" description:"Database Type" editable:"true" enabled:"true" label:"Type"`
	Database        string `json:"database" bson:"database" description:"Database Name" editable:"true" enabled:"true" label:"Database"`
	Name            string `json:"name" bson:"name" description:"Database Name" editable:"true" enabled:"true" label:"Name"`
	PortNumber      string `json:"portNumber" bson:"portNumber" description:"Database Port Number" editable:"true" enabled:"true" label:"Port Number"`
	DefaultUsername string `json:"defaultUsername" bson:"defaultUsername" description:"Database Default Username" editable:"true" enabled:"true" label:"Default Username"`
	DefaultPassword string `json:"defaultPassword" bson:"defaultPassword" description:"Database Default Password" editable:"true" enabled:"true" label:"Default Password"`
	Host            string `json:"host" bson:"host" description:"Database Host" editable:"true" enabled:"true" label:"Host"`
	Version         string `json:"version" bson:"version" description:"Database Version" editable:"true" enabled:"true" label:"Version"`
	Image           string `json:"image" bson:"image" description:"Image" editable:"true" enabled:"true" label:"Image"`
	MockDataPath    string `json:"mockDataPath" bson:"mockDataPath" description:"Database Mock Data Path" editable:"true" enabled:"true" label:"Mock Data Path"`
}

type ProductMicroService struct {
	// ID         primitive.ObjectID            `json:"id" bson:"_id" description:"Microservice ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	Name       string                        `json:"name" bson:"name" description:"Microservice Name" required:"true" editable:"true" enabled:"true" label:"Name"`
	Type       string                        `json:"type" bson:"type" description:"Microservice Type" editable:"true" required:"true" enabled:"true" label:"Type"`
	ProductID  string                        `json:"productID" required:"true" bson:"productID" description:"Product ID" editable:"true" enabled:"true" label:"Product ID"`
	PortNumber string                        `json:"portNumber" bson:"portNumber" description:"Microservice Port Number" editable:"true" required:"true" enabled:"true" label:"Port Number"`
	Status     string                        `json:"status" bson:"status" description:"Product Status" editable:"true" enabled:"true" label:"Status"`
	Host       string                        `json:"host" bson:"host" description:"Microservice Host" editable:"true" enabled:"true" label:"Host"`
	BasePath   string                        `json:"basePath" bson:"basePath" required:"true" description:"Microservice Base Path" editable:"true" enabled:"true" label:"Base Path"`
	Version    string                        `json:"version" bson:"version" required:"true" description:"Microservice Version" editable:"true" enabled:"true" label:"Version"`
	CreatedOn  time.Time                     `json:"createdOn" bson:"createdOn" description:"Created On" editable:"false" enabled:"false" label:"Created On"`
	UpdatedOn  time.Time                     `json:"updatedOn" bson:"updatedOn" description:"Updated On" editable:"false" enabled:"false" label:"Updated On"`
	Databases  []ProductMicroServiceDatabase `json:"productDatabases" bson:"productDatabases" description:"Product Databases" editable:"true" enabled:"true" label:"Databases"`
}

type ProductProviderPermissions struct {
	Permissions []string `json:"permissions" bson:"permissions" description:"Permissions" editable:"true" enabled:"true" label:"Permissions"`
	Provider    string   `json:"provider" bson:"provider" description:"Provider" editable:"true" enabled:"true" label:"Provider"`
	Enabled     bool     `json:"enabled" bson:"enabled" description:"Enabled" editable:"true" enabled:"true" label:"Enabled"`
}

type ProductResourceVersions struct {
	ResourceName   string `json:"resourceName" bson:"resourceName" description:"Resource Name" editable:"true" enabled:"true" label:"Resource Name"`
	DevelopmentTag string `json:"developmentTag" bson:"developmentTag" description:"Development Tag" editable:"true" enabled:"true" label:"Development Tag"`
	ProductTag     string `json:"productTag" bson:"productTag" description:"Product Tag" editable:"true" enabled:"true" label:"Product Tag"`
	StagingTag     string `json:"stagingTag" bson:"stagingTag" description:"Staging Tag" editable:"true" enabled:"true" label:"Staging Tag"`
	TestingTag     string `json:"testingTag" bson:"testingTag" description:"Testing Tag" editable:"true" enabled:"true" label:"Testing Tag"`
	CreatedOn      string `json:"createdOn" bson:"createdOn" description:"Created On" editable:"true" enabled:"true" label:"Created On"`
	UpdatedOn      string `json:"updatedOn" bson:"updatedOn" description:"Updated On" editable:"true" enabled:"true" label:"Updated On"`
	Latest         bool   `json:"latest" bson:"latest" description:"Latest" editable:"true" enabled:"true" label:"Latest"`
}

type ProductResource struct {
	// ID                         primitive.ObjectID        `json:"id" bson:"_id" description:"Resource ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	Kind                       string                    `json:"kind" required:"true" bson:"kind" description:"Resource Kind" editable:"true" enabled:"true" label:"Resource Kind"`
	TaskType                   string                    `json:"taskType" required:"true" bson:"taskType" description:"Task Type" editable:"true" enabled:"true" label:"Task Type"`
	DeploymentTemplateName     string                    `json:"deploymentTemplateName" bson:"deploymentTemplateName" description:"Deployment Template Name" editable:"true" enabled:"true" label:"Deployment Template Name"`
	DeploymentTemplateLocation string                    `json:"deploymentTemplateLocation" bson:"deploymentTemplateLocation" description:"Deployment Template Location" editable:"true" enabled:"true" label:"Deployment Template Location"`
	Status                     string                    `json:"status" bson:"status" description:"Product Status" editable:"true" enabled:"true" label:"Status"`
	Name                       string                    `json:"name" bson:"name" description:"Name" required:"true" editable:"true" enabled:"true" label:"Name"`
	ProductID                  string                    `json:"productID" required:"true" bson:"productID" description:"Product ID" editable:"true" enabled:"true" label:"Product ID"`
	CreatedOn                  time.Time                 `json:"createdOn" bson:"createdOn" description:"Created On" editable:"false" enabled:"false" label:"Created On"`
	UpdatedOn                  time.Time                 `json:"updatedOn" bson:"updatedOn" description:"Updated On" editable:"false" enabled:"false" label:"Updated On"`
	Versions                   []ProductResourceVersions `json:"versions" bson:"versions" description:"Resource Versions" editable:"true" enabled:"true" label:"Resource Versions"`
	EnvironmentVariablesFile   string                    `json:"environmentVariablesFile" bson:"environmentVariablesFile" description:"Environment Variables File" editable:"true" enabled:"true" label:"Environment Variables File"`
	ProductVersion             string                    `json:"productVersion" bson:"productVersion" description:"Product Version" editable:"true" required:"true" enabled:"true" label:"Product Version"`
	ExposedEnvVariables        []string                  `json:"exposedEnvVariables" bson:"exposedEnvVariables" description:"Exposed Environment Variables" editable:"true" enabled:"true" label:"Exposed Environment Variables"`
}

type ProductConfiguration struct {
	ID                    primitive.ObjectID           `json:"id" bson:"_id" description:"Configuration ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	ProductID             string                       `json:"productID" required:"true" bson:"productID" description:"Product ID" editable:"true" enabled:"true" label:"Product ID"`
	DefaultMemory         float32                      `json:"defaultMemory" required:"true" bson:"defaultMemory" description:"Default Memory" editable:"true" enabled:"true" label:"Default Memory"`
	DefaultRAM            float32                      `json:"defaultRAM" required:"true" bson:"defaultRAM" description:"Default RAM" editable:"true" enabled:"true" label:"Default RAM"`
	ProviderPermissions   []ProductProviderPermissions `json:"providerPermissions" required:"true" bson:"providerPermissions" description:"Provider Permissions" editable:"true" enabled:"true" label:"Provider Permissions"`
	StartupScriptFilePath string                       `json:"startupScriptFilePath" bson:"startupScriptFilePath" description:"Startup Script File Path" editable:"true" enabled:"true" label:"Startup Script File Path"`
	Status                string                       `json:"status" bson:"status" description:"Product Status" editable:"true" enabled:"true" label:"Status"`
	ArtifactRegistryName  string                       `json:"artifactRegistryName" required:"true" bson:"artifactRegistryName" description:"Artifact Registry Name" editable:"true" enabled:"true" label:"Artifact Registry Name"`
	NetworkTags           []string                     `json:"networkTags" bson:"networkTags" description:"Network Tags" editable:"true" enabled:"true" label:"Network Tags"`
	CreatedOn             time.Time                    `json:"createdOn" bson:"createdOn" description:"Created On" editable:"false" enabled:"false" label:"Created On"`
	UpdatedOn             time.Time                    `json:"updatedOn" bson:"updatedOn" description:"Updated On" editable:"false" enabled:"false" label:"Updated On"`
	EnvironmentsSupport   []ProductEnvironmentSupport  `json:"environmentsSupport" bson:"environmentsSupport" description:"Environments Support" editable:"true" enabled:"true" label:"Environments Support"`
}

type ProductEnvironmentSupport struct {
	Provider    string `json:"provider" bson:"provider" description:"Provider" editable:"true" enabled:"true" label:"Provider"`
	Enabled     bool   `json:"enabled" bson:"enabled" description:"Enabled" editable:"true" enabled:"true" label:"Enabled"`
	Environment string `json:"environment" bson:"environment" description:"Environment Name" editable:"true" enabled:"true" label:"Environment Name"`
}

type ProductPlan struct {
	// ID          primitive.ObjectID `json:"id" bson:"_id" description:"Plan ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	ProductID   string    `json:"productID" bson:"productID" required:"true" description:"Product ID" editable:"true" enabled:"true" label:"Product ID"`
	Plan        string    `json:"plan" bson:"plan" description:"Plan" required:"true" editable:"true" enabled:"true" label:"Plan"`
	Name        string    `json:"name" bson:"name" description:"Name" required:"true" editable:"true" enabled:"true" label:"Name"`
	Description string    `json:"description" bson:"description" required:"true" description:"Description" editable:"true" enabled:"true" label:"Description"`
	Status      string    `json:"status" bson:"status" description:"Product Status" editable:"true" enabled:"true" label:"Status"`
	Active      string    `json:"active" bson:"active" description:"Active" editable:"true" enabled:"true" label:"Active"`
	Image       string    `json:"image" bson:"image" description:"Image" editable:"true" enabled:"true" label:"Image"`
	CreatedOn   time.Time `json:"createdOn" bson:"createdOn" description:"Created On" editable:"false" enabled:"false" label:"Created On"`
	UpdatedOn   time.Time `json:"updatedOn" bson:"updatedOn" description:"Updated On" editable:"false" enabled:"false" label:"Updated On"`
}

type ProductVersion struct {
	// ID          primitive.ObjectID `json:"id" bson:"_id" description:"Plan ID" editable:"false" enabled:"false" label:"ID" unique:"true"`
	Name        string    `json:"name" bson:"name" description:"Name" editable:"true" required:"true" enabled:"true" label:"Name"`
	Version     string    `json:"version" bson:"version" description:"Version" editable:"true" required:"true" enabled:"true" label:"Version"`
	BuildNumber string    `json:"buildNumber" bson:"buildNumber" description:"Build Number" editable:"true" required:"true" enabled:"true" label:"Build Number"`
	ProductID   string    `json:"productID" bson:"productID" description:"Product ID" editable:"true" required:"true" enabled:"true" label:"Product ID"`
	Status      string    `json:"status" bson:"status" description:"Product Status" editable:"false" enabled:"true" label:"Status"`
	Description string    `json:"description" bson:"description" description:"Description" required:"true" editable:"true" enabled:"true" label:"Description"`
	Type        string    `json:"type" bson:"type" description:"Type" editable:"true" required:"true" enabled:"true" label:"Type"`
	CreatedOn   time.Time `json:"createdOn" bson:"createdOn" description:"Created On" editable:"false" enabled:"false" label:"Created On"`
	UpdatedOn   time.Time `json:"updatedOn" bson:"updatedOn" description:"Updated On" editable:"false" enabled:"false" label:"Updated On"`
}
