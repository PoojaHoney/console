package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func createProduct(srv *Service, product Product, createdBy string) (Product, error) {
	productExists := srv.ProductDB.Where("product_id = ?", product.ProductID).First(&Product{})
	if productExists.Error == nil {
		return Product{}, errors.New("product already exists")
	}
	if productExists.Error != nil && productExists.Error != gorm.ErrRecordNotFound {
		return Product{}, productExists.Error
	}
	product.Status = "draft"
	product.CreatedBy = createdBy
	result := srv.ProductDB.Create(&product)
	if result.Error != nil {
		return Product{}, result.Error
	}
	return product, nil
}

func createVersion(srv *Service, version ProductVersion) (ProductVersion, error) {
	productExists := srv.ProductDB.Where("product_id = ?", version.ProductID).First(&Product{})
	if productExists.Error != nil {
		return ProductVersion{}, productExists.Error
	}
	var existingVersion ProductVersion
	versionExists := srv.ProductDB.Where("version = ?", version.Version).First(&existingVersion)
	if versionExists.Error == nil {
		return ProductVersion{}, errors.New("version already exists")
	}
	if versionExists.Error != nil && versionExists.Error != gorm.ErrRecordNotFound {
		return ProductVersion{}, versionExists.Error
	}
	version.Status = "draft"
	// version.ID = primitive.NewObjectID()
	result := srv.ProductDB.Create(&version)
	if result.Error != nil {
		return ProductVersion{}, result.Error
	}
	return version, nil
}

func createMicroService(srv *Service, microservice ProductMicroService) (ProductMicroService, error) {
	product, err := getProduct(srv.ProductDB, map[string]interface{}{"product_id": microservice.ProductID})
	if err != nil {
		return ProductMicroService{}, err
	}
	if product.ProductID == "" {
		return ProductMicroService{}, errors.New("product not found")
	}
	microserviceExists := srv.ProductDB.Where("name = ?", microservice.Name).Or("portNumber = ?", microservice.PortNumber).Or(
		"basePath = ?", microservice.BasePath).Where("version = ?", microservice.Version).First(&ProductMicroService{})
	if microserviceExists.Error == nil {
		return ProductMicroService{}, errors.New("microservice already exists")
	}
	if microserviceExists.Error != nil && microserviceExists.Error != gorm.ErrRecordNotFound {
		return ProductMicroService{}, microserviceExists.Error
	}
	microservice.Status = "draft"
	result := srv.ProductDB.Create(&microservice)
	if result.Error != nil {
		return ProductMicroService{}, result.Error
	}
	// product.MicroServicesCount = product.MicroServicesCount + 1
	// product.DatabasesCount = product.DatabasesCount + len(microservice.Databases)
	// _, err = productCollection.UpdateOne(context.TODO(), bson.D{{Key: "product_id", Value: microservice.ProductID}}, bson.D{{Key: "$set", Value: product}})
	// if err != nil {
	// 	return ProductMicroService{}, err
	// }
	return microservice, nil
}

func createPlan(srv *Service, plan ProductPlan) (ProductPlan, error) {
	product, err := getProduct(srv.ProductDB, map[string]interface{}{"product_id": plan.ProductID})
	if err != nil {
		return ProductPlan{}, err
	}
	if product.ProductID == "" {
		return ProductPlan{}, errors.New("product not found")
	}
	planExists := srv.ProductDB.Where("plan = ?", plan.Plan).Or("name = ?", plan.Name).First(&ProductPlan{})
	if planExists.Error == nil {
		return ProductPlan{}, errors.New("plan already exists")
	}
	if planExists.Error != nil && planExists.Error != gorm.ErrRecordNotFound {
		return ProductPlan{}, planExists.Error
	}
	// plan.ID = primitive.NewObjectID()
	plan.Status = "draft"
	result := srv.ProductDB.Create(&plan)
	if result.Error != nil {
		return ProductPlan{}, result.Error
	}
	return plan, nil
}

func createResource(srv *Service, resource ProductResource) (ProductResource, error) {
	product, err := getProduct(srv.ProductDB, map[string]interface{}{"product_id": resource.ProductID})
	if err != nil {
		return ProductResource{}, err
	}
	if product.ProductID == "" {
		return ProductResource{}, errors.New("product not found")
	}
	resourceExists := srv.ProductDB.Where("name = ?", resource.Name).Or("kind = ?", resource.Kind).Or(
		"taskType = ?", resource.TaskType).First(&ProductResource{})
	if resourceExists.Error == nil {
		return ProductResource{}, errors.New("resource already exists")
	}
	if resourceExists.Error != nil && resourceExists.Error != gorm.ErrRecordNotFound {
		return ProductResource{}, resourceExists.Error
	}
	// resource.ID = primitive.NewObjectID()
	resource.Status = "draft"
	result := srv.ProductDB.Create(&resource)
	if result.Error != nil {
		return ProductResource{}, result.Error
	}
	return resource, nil
}

func createConfiguration(srv *Service, configuration ProductConfiguration) (ProductConfiguration, error) {
	product, err := getProduct(srv.ProductDB, map[string]interface{}{"product_id": configuration.ProductID})
	if err != nil {
		return ProductConfiguration{}, err
	}
	if product.ProductID == "" {
		return ProductConfiguration{}, errors.New("product not found")
	}
	configuration.Status = "draft"
	result := srv.ProductDB.Create(&configuration)
	if result.Error != nil {
		return ProductConfiguration{}, result.Error
	}
	return configuration, nil
}

func activateProduct(srv *Service, product string) (FullProductDetails, error) {
	productInfo, err := getProduct(srv.ProductDB, map[string]interface{}{"product_id": product})
	if err != nil {
		return FullProductDetails{}, err
	}
	configuration, err := getProductConfigurations(srv.ProductDB, map[string]interface{}{"product_id": product})
	if err != nil {
		return FullProductDetails{}, err
	}
	versions, err := getProductVersions(srv.ProductDB, map[string]interface{}{"product_id": product})
	if err != nil {
		return FullProductDetails{}, err
	}
	microservices, err := getProductMicroServices(srv.ProductDB, map[string]interface{}{"product_id": product})
	if err != nil {
		return FullProductDetails{}, err
	}
	plans, err := getProductPlans(srv.ProductDB, map[string]interface{}{"product_id": product})
	if err != nil {
		return FullProductDetails{}, err
	}
	resources, err := getProductResources(srv.ProductDB, map[string]interface{}{"product_id": product})
	if err != nil {
		return FullProductDetails{}, err
	}

	result := srv.ProductDB.Unscoped().Delete(&Product{}, productInfo.ID)
	if result.Error != nil {
		return FullProductDetails{}, result.Error
	}

	result = srv.ProductDB.Unscoped().Delete(&ProductConfiguration{}, configuration)
	if result.Error != nil {
		return FullProductDetails{}, result.Error
	}

	result = srv.ProductDB.Unscoped().Delete(&ProductConfiguration{}, configuration.ID)
	if result.Error != nil {
		srv.ProductDB.Create(&productInfo)
		return FullProductDetails{}, result.Error
	}

	// result = srv.ProductDB.Unscoped().Delete(&versions)
	result = srv.ProductDB.Exec("DELETE FROM product_versions WHERE ", "product_id = ?", product)
	if result.Error != nil {
		srv.ProductDB.Create(&productInfo)
		srv.ProductDB.Create(&configuration)
		return FullProductDetails{}, result.Error
	}

	result = srv.ProductDB.Exec("DELETE FROM product_microservices WHERE ", "product_id = ?", product)
	if result.Error != nil {
		srv.ProductDB.Create(&productInfo)
		srv.ProductDB.Create(&configuration)
		srv.ProductDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&versions)
		return FullProductDetails{}, result.Error
	}

	result = srv.ProductDB.Exec("DELETE FROM product_plans WHERE ", "product_id = ?", product)
	if result.Error != nil {
		srv.ProductDB.Create(&productInfo)
		srv.ProductDB.Create(&configuration)
		srv.ProductDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&versions)
		srv.ProductDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&microservices)
		return FullProductDetails{}, result.Error
	}

	result = srv.ProductDB.Exec("DELETE FROM product_resources WHERE ", "product_id = ?", product)
	if result.Error != nil {
		srv.ProductDB.Create(&productInfo)
		srv.ProductDB.Create(&configuration)
		srv.ProductDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&versions)
		srv.ProductDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&microservices)
		srv.ProductDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&plans)
		return FullProductDetails{}, result.Error
	}

	for index, _ := range microservices {
		microservices[index].Status = "active"
	}
	for index, _ := range plans {
		plans[index].Status = "active"
	}
	for index, _ := range resources {
		resources[index].Status = "active"
	}
	for index, _ := range versions {
		versions[index].Status = "active"
	}
	productFullDetails := FullProductDetails{
		MicroServices:      microservices,
		Plans:              plans,
		Resources:          resources,
		Configuration:      *configuration,
		Versions:           versions,
		Name:               productInfo.Name,
		Description:        productInfo.Description,
		Status:             productInfo.Status,
		Type:               productInfo.Type,
		ProductID:          productInfo.ProductID,
		Image:              productInfo.Image,
		Providers:          productInfo.Providers,
		MicroServicesCount: len(microservices),
		DatabasesCount:     productInfo.DatabasesCount,
		CreatedBy:          productInfo.CreatedBy,
	}
	result = srv.ProductDB.Create(&productFullDetails)
	if result.Error != nil {
		srv.ProductDB.Create(&productInfo)
		srv.ProductDB.Create(&configuration)
		srv.ProductDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&versions)
		srv.ProductDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&microservices)
		srv.ProductDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&plans)
		srv.ProductDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&resources)
		return FullProductDetails{}, err
	}
	return productFullDetails, nil
}

func getProduct(db *gorm.DB, filters map[string]interface{}) (*Product, error) {
	var product Product
	for key, value := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}
	result := db.First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func getProductConfigurations(db *gorm.DB, filters map[string]interface{}) (*ProductConfiguration, error) {
	var configurations ProductConfiguration
	for key, value := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}
	result := db.First(&configurations)
	if result.Error != nil {
		return nil, result.Error
	}
	return &configurations, nil
}

func getProductMicroServices(db *gorm.DB, filters map[string]interface{}) ([]ProductMicroService, error) {
	var microservices []ProductMicroService
	for key, value := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}
	query := db.Model(&ProductMicroService{})
	for operation, filter := range filters {
		if operation == "or" {
			for key, value := range filter.(map[string]interface{}) {
				query = query.Or(fmt.Sprintf("%s = ?", key), value)
			}
		} else if operation == "and" {
			for key, value := range filter.(map[string]interface{}) {
				query = query.Where(fmt.Sprintf("%s = ?", key), value)
			}
		}
	}
	if err := query.Find(&microservices).Error; err != nil {
		return nil, err
	}
	if len(microservices) == 0 {
		return []ProductMicroService{}, errors.New("product microservices not found")
	}
	return microservices, nil
}

func getProductVersions(db *gorm.DB, filters map[string]interface{}) ([]ProductVersion, error) {
	var versions []ProductVersion
	for key, value := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}
	query := db.Model(&ProductVersion{})
	for operation, filter := range filters {
		if operation == "or" {
			for key, value := range filter.(map[string]interface{}) {
				query = query.Or(fmt.Sprintf("%s = ?", key), value)
			}
		} else if operation == "and" {
			for key, value := range filter.(map[string]interface{}) {
				query = query.Where(fmt.Sprintf("%s = ?", key), value)
			}
		}
	}
	if err := query.Find(&versions).Error; err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return []ProductVersion{}, errors.New("product versions not found")
	}
	return versions, nil
}

func getProductPlans(db *gorm.DB, filters map[string]interface{}) ([]ProductPlan, error) {
	var plans []ProductPlan
	for key, value := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}
	query := db.Model(&ProductPlan{})
	for operation, filter := range filters {
		if operation == "or" {
			for key, value := range filter.(map[string]interface{}) {
				query = query.Or(fmt.Sprintf("%s = ?", key), value)
			}
		} else if operation == "and" {
			for key, value := range filter.(map[string]interface{}) {
				query = query.Where(fmt.Sprintf("%s = ?", key), value)
			}
		}
	}
	if err := query.Find(&plans).Error; err != nil {
		return nil, err
	}
	if len(plans) == 0 {
		return []ProductPlan{}, errors.New("product plans not found")
	}
	return plans, nil
}

func getProductResources(db *gorm.DB, filters map[string]interface{}) ([]ProductResource, error) {
	var resources []ProductResource
	for key, value := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}
	query := db.Model(&ProductResource{})
	for operation, filter := range filters {
		if operation == "or" {
			for key, value := range filter.(map[string]interface{}) {
				query = query.Or(fmt.Sprintf("%s = ?", key), value)
			}
		} else if operation == "and" {
			for key, value := range filter.(map[string]interface{}) {
				query = query.Where(fmt.Sprintf("%s = ?", key), value)
			}
		}
	}
	if err := query.Find(&resources).Error; err != nil {
		return nil, err
	}
	if len(resources) == 0 {
		return []ProductResource{}, errors.New("product resources not found")
	}
	return resources, nil
}

func fieldCatalogues(tables ...interface{}) map[string]interface{} {
	if tables == nil {
		tables = []interface{}{
			Product{}, ProductConfiguration{}, ProductMicroService{},
			ProductPlan{}, ProductResource{}, ProductMicroServiceDatabase{},
			ProductProviderPermissions{}, ProductResourceVersions{}, FullProductDetails{}}
	}
	allFieldCatalogues := make(map[string]interface{})
	for i := 0; i < len(tables); i++ {
		fields := []FieldCatalogue{}
		var table interface{} = tables[i]
		childFields := make(map[string]interface{})
		tableType := reflect.TypeOf(table)
		for i := 0; i < tableType.NumField(); i++ {
			temp := tableType.Field(i)
			length, _ := strconv.Atoi(temp.Tag.Get("length"))
			kind := temp.Type.Kind().String()
			if temp.Type.String() == "time.Time" {
				kind = "date"
			}
			field := FieldCatalogue{
				Name:        temp.Name,
				Type:        kind,
				Label:       temp.Tag.Get("label"),
				Description: temp.Tag.Get("description"),
				Required:    temp.Tag.Get("required") == "true",
				Unique:      temp.Tag.Get("unique") == "true",
				Length:      length,
				Editable:    temp.Tag.Get("editable") == "true",
				Enabled:     temp.Tag.Get("enabled") == "true",
			}
			fields = append(fields, field)
			if field.Type == "struct" {
				child := fieldCatalogues(reflect.New(temp.Type).Elem().Interface())
				for key, value := range child {
					childFields[key] = value
				}
			}
		}
		allFieldCatalogues[strings.ToLower(tableType.Name()[:1])+tableType.Name()[1:]] = AllFieldCatalogues{
			Fields:      fields,
			Table:       tableType.Name(),
			ChildTables: childFields,
		}
	}
	return allFieldCatalogues
}
func getToken(srv *Service, filters map[string]interface{}) (Tokens, error) {
	var token Tokens
	db := srv.ProductDB
	for key, value := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}
	resultToken := db.First(&token)
	if resultToken.Error != nil {
		return Tokens{}, resultToken.Error
	}
	return token, nil
}

func validateToken(srv *Service, authHeader string) (map[string]interface{}, error) {
	// tokenString := authHeader[len("Bearer "):]
	tokenString := authHeader
	var expiryTime int
	token, err := getToken(srv, map[string]interface{}{"access_token": tokenString})
	if err != nil {
		return nil, err
	}
	expiryTime = int(time.Since(token.CreatedAt.Add(time.Hour * 5)).Minutes())
	if token.UserId != uuid.Nil && (expiryTime < 30) {
		return map[string]interface{}{
			"userId": token.UserId,
			"email":  token.Email,
		}, nil
	} else {
		return nil, errors.New("authorization token has expired")
	}
}

// func createVersions(srv *Service, versions []ProductVersion) ([]ProductVersion, error) {
// 	productVersionsCollection := getMongoCollection(srv, srv.Config.PRODUCT_VERSIONS_TABLE)
// 	var createdVersions []ProductVersion
// 	for _, version := range versions {
// 		recordsExists, err := getMongoRecordsCount(
// 			productVersionsCollection,
// 			bson.D{{
// 				Key: "$or", Value: bson.A{
// 					bson.D{{Key: "Version", Value: version.Version}},
// 				}}})
// 		if err != nil {
// 			return nil, err
// 		}
// 		if recordsExists > 0 {
// 			continue
// 		}
// 		version.Status = "draft"
// 		version.CreatedOn = time.Now()
// 		createdVersions = append(createdVersions, version)
// 	}
// 	if len(createdVersions) == 0 {
// 		return createdVersions, nil
// 	}
// 	documents := []interface{}{}
// 	docBytes, _ := ffjson.Marshal(createdVersions)
// 	ffjson.Unmarshal(docBytes, &documents)
// 	_, err := productVersionsCollection.InsertMany(context.TODO(), documents)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return createdVersions, nil
// }

// func createMicroServices(srv *Service, microservices []ProductMicroService) ([]ProductMicroService, error) {
// 	microserviceCollection := getMongoCollection(srv, srv.Config.PRODUCT_MICROSERVICES_TABLE)
// 	var createdMicroServices []ProductMicroService
// 	for _, microservice := range microservices {
// 		recordsExists, err := getMongoRecordsCount(
// 			microserviceCollection,
// 			bson.D{{
// 				Key: "$or", Value: bson.A{
// 					bson.D{{Key: "Version", Value: microservice.Version}},
// 				}}})
// 		if err != nil {
// 			return nil, err
// 		}
// 		if recordsExists > 0 {
// 			continue
// 		}
// 		microservice.Status = "draft"
// 		microservice.CreatedOn = time.Now()
// 		createdMicroServices = append(createdMicroServices, microservice)
// 	}
// 	if len(createdMicroServices) == 0 {
// 		return createdMicroServices, nil
// 	}
// 	documents := []interface{}{}
// 	docBytes, _ := ffjson.Marshal(createdMicroServices)
// 	ffjson.Unmarshal(docBytes, &documents)
// 	_, err := microserviceCollection.InsertMany(context.TODO(), documents)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return createdMicroServices, nil
// }
