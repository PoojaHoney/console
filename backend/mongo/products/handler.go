package main

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/ffjson/ffjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createProduct(srv *Service, product Product, createdBy string) (Product, error) {
	productCollection := getMongoCollection(srv, srv.Config.PRODUCT_COLLECTION)
	recordsExists, err := getMongoRecordsCount(
		productCollection,
		bson.D{{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "productID", Value: product.ProductID}},
			}}})
	if err != nil {
		return Product{}, err
	}
	if recordsExists > 0 {
		return Product{}, errors.New("product already exists")
	}
	product.Status = "draft"
	product.ID = primitive.NewObjectID()
	product.CreatedOn = time.Now()
	product.CreatedBy = createdBy
	_, err = productCollection.InsertOne(context.TODO(), product)
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func createVersion(srv *Service, version ProductVersion) (ProductVersion, error) {
	productVersionsCollection := getMongoCollection(srv, srv.Config.PRODUCT_VERSIONS_COLLECTION)
	productCollection := getMongoCollection(srv, srv.Config.PRODUCT_COLLECTION)
	productExists, err := getMongoRecordsCount(
		productCollection,
		bson.D{{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "productID", Value: version.ProductID}},
			}}})
	if err != nil {
		return ProductVersion{}, err
	}
	if productExists == 0 {
		return ProductVersion{}, errors.New("product not found")
	}
	recordsExists, err := getMongoRecordsCount(
		productVersionsCollection,
		bson.D{{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "version", Value: version.Version}},
			}}})
	if err != nil {
		return ProductVersion{}, err
	}
	if recordsExists > 0 {
		return ProductVersion{}, errors.New("version already exists")
	}
	version.Status = "draft"
	// version.ID = primitive.NewObjectID()
	version.CreatedOn = time.Now()
	_, err = productVersionsCollection.InsertOne(context.TODO(), version)
	if err != nil {
		return ProductVersion{}, err
	}
	return version, nil
}

func createMicroService(srv *Service, microservice ProductMicroService) (ProductMicroService, error) {
	microserviceCollection := getMongoCollection(srv, srv.Config.PRODUCT_MICROSERVICES_COLLECTION)
	productCollection := getMongoCollection(srv, srv.Config.PRODUCT_COLLECTION)
	product, err := getProduct(srv, microservice.ProductID, productCollection)
	if err != nil {
		return ProductMicroService{}, err
	}
	if product.ProductID == "" {
		return ProductMicroService{}, errors.New("product not found")
	}
	recordsExists, err := getMongoRecordsCount(
		microserviceCollection,
		bson.D{{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "name", Value: microservice.Name}},
				bson.D{{Key: "portNumber", Value: microservice.PortNumber}},
				bson.D{{Key: "basePath", Value: microservice.BasePath}},
			}}, {
			Key: "$and", Value: bson.A{
				bson.D{{Key: "version", Value: microservice.Version}},
			}}})
	if err != nil {
		return ProductMicroService{}, err
	}
	if recordsExists > 0 {
		return ProductMicroService{}, errors.New("microservice already exists")
	}
	microservice.Status = "draft"
	microservice.CreatedOn = time.Now()
	// microservice.ID = primitive.NewObjectID()
	_, err = microserviceCollection.InsertOne(context.TODO(), microservice)
	if err != nil {
		return ProductMicroService{}, err
	}
	// product.MicroServicesCount = product.MicroServicesCount + 1
	// product.DatabasesCount = product.DatabasesCount + len(microservice.Databases)
	// _, err = productCollection.UpdateOne(context.TODO(), bson.D{{Key: "productID", Value: microservice.ProductID}}, bson.D{{Key: "$set", Value: product}})
	// if err != nil {
	// 	return ProductMicroService{}, err
	// }
	return microservice, nil
}

func createPlan(srv *Service, plan ProductPlan) (ProductPlan, error) {
	plansCollection := getMongoCollection(srv, srv.Config.PRODUCT_PLANS_COLLECTION)
	productCollection := getMongoCollection(srv, srv.Config.PRODUCT_COLLECTION)
	productExists, err := getMongoRecordsCount(
		productCollection,
		bson.D{{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "productID", Value: plan.ProductID}},
			}}})
	if err != nil {
		return ProductPlan{}, err
	}
	if productExists == 0 {
		return ProductPlan{}, errors.New("product not found")
	}
	recordsExists, err := getMongoRecordsCount(
		plansCollection,
		bson.D{{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "plan", Value: plan.Plan}},
				bson.D{{Key: "name", Value: plan.Name}},
			}}})
	if err != nil {
		return ProductPlan{}, err
	}
	if recordsExists > 0 {
		return ProductPlan{}, errors.New("plan already exists")
	}
	// plan.ID = primitive.NewObjectID()
	plan.Status = "draft"
	plan.CreatedOn = time.Now()
	_, err = plansCollection.InsertOne(context.TODO(), plan)
	if err != nil {
		return ProductPlan{}, err
	}
	return plan, nil
}

func createResource(srv *Service, resource ProductResource) (ProductResource, error) {
	resourceCollection := getMongoCollection(srv, srv.Config.PRODUCT_RESOURCES_COLLECTION)
	productCollection := getMongoCollection(srv, srv.Config.PRODUCT_COLLECTION)
	productExists, err := getMongoRecordsCount(
		productCollection,
		bson.D{{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "productID", Value: resource.ProductID}},
			}}})
	if err != nil {
		return ProductResource{}, err
	}
	if productExists == 0 {
		return ProductResource{}, errors.New("product not found")
	}
	recordsExists, err := getMongoRecordsCount(
		resourceCollection,
		bson.D{{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "kind", Value: resource.Kind}},
				bson.D{{Key: "taskType", Value: resource.TaskType}},
			}}, {
			Key: "$and", Value: bson.A{
				bson.D{{Key: "name", Value: resource.Name}},
			}}})
	if err != nil {
		return ProductResource{}, err
	}
	if recordsExists > 0 {
		return ProductResource{}, errors.New("resource already exists")
	}
	// resource.ID = primitive.NewObjectID()
	resource.Status = "draft"
	resource.CreatedOn = time.Now()
	_, err = resourceCollection.InsertOne(context.TODO(), resource)
	if err != nil {
		return ProductResource{}, err
	}
	return resource, nil
}

func createConfiguration(srv *Service, configuration ProductConfiguration) (ProductConfiguration, error) {
	configCollection := getMongoCollection(srv, srv.Config.PRODUCT_CONFIG_COLLECTION)
	productCollection := getMongoCollection(srv, srv.Config.PRODUCT_COLLECTION)
	productExists, err := getMongoRecordsCount(
		productCollection,
		bson.D{{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "productID", Value: configuration.ProductID}},
			}}})
	if err != nil {
		return ProductConfiguration{}, err
	}
	if productExists == 0 {
		return ProductConfiguration{}, errors.New("product not found")
	}
	configuration.ID = primitive.NewObjectID()
	configuration.Status = "draft"
	configuration.CreatedOn = time.Now()
	_, err = configCollection.InsertOne(context.TODO(), configuration)
	if err != nil {
		return ProductConfiguration{}, err
	}
	return configuration, nil
}

func activateProduct(srv *Service, product string) (FullProductDetails, error) {
	productCollection := getMongoCollection(srv, srv.Config.PRODUCT_COLLECTION)
	productInfo, err := getProduct(srv, product, productCollection)
	if err != nil {
		return FullProductDetails{}, err
	}
	configurationCollection := getMongoCollection(srv, srv.Config.PRODUCT_CONFIG_COLLECTION)
	configuration, err := getProductConfigurations(srv, configurationCollection, product)
	if err != nil {
		return FullProductDetails{}, err
	}
	versionCollection := getMongoCollection(srv, srv.Config.PRODUCT_VERSIONS_COLLECTION)
	versions, err := getProductVersions(srv, versionCollection, product)
	if err != nil {
		return FullProductDetails{}, err
	}
	microserviceCollection := getMongoCollection(srv, srv.Config.PRODUCT_MICROSERVICES_COLLECTION)
	microservices, err := getProductMicroServices(srv, microserviceCollection, product)
	if err != nil {
		return FullProductDetails{}, err
	}
	plansCollection := getMongoCollection(srv, srv.Config.PRODUCT_PLANS_COLLECTION)
	plans, err := getProductPlans(srv, plansCollection, product)
	if err != nil {
		return FullProductDetails{}, err
	}
	resourcesCollection := getMongoCollection(srv, srv.Config.PRODUCT_RESOURCES_COLLECTION)
	resources, err := getProductResources(srv, resourcesCollection, product)
	if err != nil {
		return FullProductDetails{}, err
	}

	_, err = productCollection.DeleteOne(context.TODO(), bson.D{{Key: "productID", Value: product}})
	if err != nil {
		return FullProductDetails{}, err
	}

	_, err = configurationCollection.DeleteOne(context.TODO(), bson.D{{Key: "productID", Value: product}})
	if err != nil {
		productCollection.InsertOne(context.TODO(), productInfo)
		return FullProductDetails{}, err
	}

	_, err = versionCollection.DeleteMany(context.TODO(), bson.D{{Key: "productID", Value: product}})
	if err != nil {
		productCollection.InsertOne(context.TODO(), productInfo)
		configurationCollection.InsertOne(context.TODO(), configuration)
		return FullProductDetails{}, err
	}

	_, err = microserviceCollection.DeleteMany(context.TODO(), bson.D{{Key: "productID", Value: product}})
	if err != nil {
		productCollection.InsertOne(context.TODO(), productInfo)
		configurationCollection.InsertOne(context.TODO(), configuration)
		insertProductVersions(srv, versionCollection, versions)
		return FullProductDetails{}, err
	}

	_, err = plansCollection.DeleteMany(context.TODO(), bson.D{{Key: "productID", Value: product}})
	if err != nil {
		productCollection.InsertOne(context.TODO(), productInfo)
		configurationCollection.InsertOne(context.TODO(), configuration)
		insertProductVersions(srv, versionCollection, versions)
		insertProductMicroServices(srv, microserviceCollection, microservices)
		return FullProductDetails{}, err
	}

	_, err = resourcesCollection.DeleteMany(context.TODO(), bson.D{{Key: "productID", Value: product}})
	if err != nil {
		productCollection.InsertOne(context.TODO(), productInfo)
		configurationCollection.InsertOne(context.TODO(), configuration)
		insertProductVersions(srv, versionCollection, versions)
		insertProductMicroServices(srv, microserviceCollection, microservices)
		insertProductPlans(srv, plansCollection, plans)
		return FullProductDetails{}, err
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
		Configuration:      configuration,
		Versions:           versions,
		ID:                 primitive.NewObjectID(),
		Name:               productInfo.Name,
		Description:        productInfo.Description,
		Status:             productInfo.Status,
		Type:               productInfo.Type,
		ProductID:          productInfo.ProductID,
		Image:              productInfo.Image,
		Providers:          productInfo.Providers,
		MicroServicesCount: len(microservices),
		CreatedOn:          productInfo.CreatedOn,
		UpdatedOn:          time.Now(),
		DatabasesCount:     productInfo.DatabasesCount,
		CreatedBy:          productInfo.CreatedBy,
	}
	_, err = getMongoCollection(srv, srv.Config.PRODUCT_FULL_DETAILS).InsertOne(context.TODO(), productFullDetails)
	if err != nil {
		productCollection.InsertOne(context.TODO(), productInfo)
		configurationCollection.InsertOne(context.TODO(), configuration)
		insertProductVersions(srv, versionCollection, versions)
		insertProductMicroServices(srv, microserviceCollection, microservices)
		insertProductPlans(srv, plansCollection, plans)
		insertProductResources(srv, resourcesCollection, resources)
		return FullProductDetails{}, err
	}
	return productFullDetails, nil
}

func insertProductVersions(srv *Service, collection *mongo.Collection, versions []ProductVersion) error {
	if collection == nil {
		collection = getMongoCollection(srv, srv.Config.PRODUCT_VERSIONS_COLLECTION)
	}
	bytesVersion, _ := ffjson.Marshal(versions)
	documentVersions := []interface{}{}
	_ = ffjson.Unmarshal(bytesVersion, &documentVersions)
	_, err := collection.InsertMany(context.TODO(), documentVersions)
	return err
}

func insertProductMicroServices(srv *Service, collection *mongo.Collection, microservices []ProductMicroService) error {
	if collection == nil {
		collection = getMongoCollection(srv, srv.Config.PRODUCT_MICROSERVICES_COLLECTION)
	}
	bytesMicroServices, _ := ffjson.Marshal(microservices)
	documentMicroServices := []interface{}{}
	_ = ffjson.Unmarshal(bytesMicroServices, &documentMicroServices)
	_, err := collection.InsertMany(context.TODO(), documentMicroServices)
	return err
}

func insertProductPlans(srv *Service, collection *mongo.Collection, plans []ProductPlan) error {
	if collection == nil {
		collection = getMongoCollection(srv, srv.Config.PRODUCT_PLANS_COLLECTION)
	}
	bytesPlans, _ := ffjson.Marshal(plans)
	documentPlans := []interface{}{}
	_ = ffjson.Unmarshal(bytesPlans, &documentPlans)
	_, err := collection.InsertMany(context.TODO(), documentPlans)
	return err
}

func insertProductResources(srv *Service, collection *mongo.Collection, resources []ProductResource) error {
	if collection == nil {
		collection = getMongoCollection(srv, srv.Config.PRODUCT_RESOURCES_COLLECTION)
	}
	bytesResources, _ := ffjson.Marshal(resources)
	documentResources := []interface{}{}
	_ = ffjson.Unmarshal(bytesResources, &documentResources)
	_, err := collection.InsertMany(context.TODO(), documentResources)
	return err
}

func getProduct(srv *Service, productID string, collection *mongo.Collection) (Product, error) {
	if collection == nil {
		collection = getMongoCollection(srv, srv.Config.PRODUCT_COLLECTION)
	}
	document, err := getMongoRecord(collection, bson.D{{Key: "productID", Value: productID}})
	if err != nil {
		return Product{}, err
	}
	var _id primitive.ObjectID
	if oid, ok := document["_id"].(primitive.ObjectID); ok {
		_id = oid
	}
	if document["productID"].(string) == "" {
		return Product{}, errors.New("product not found")
	}
	structProduct := Product{}
	byteProduct, _ := ffjson.Marshal(document)
	_ = ffjson.Unmarshal(byteProduct, &structProduct)
	structProduct.ID = _id
	return structProduct, nil
}

func getProductConfigurations(srv *Service, collection *mongo.Collection, productID string) (ProductConfiguration, error) {
	if collection == nil {
		collection = getMongoCollection(srv, srv.Config.PRODUCT_CONFIG_COLLECTION)
	}
	document, err := getMongoRecord(collection, bson.D{{Key: "productID", Value: productID}})
	if err != nil {
		return ProductConfiguration{}, err
	}
	var _id primitive.ObjectID
	if oid, ok := document["_id"].(primitive.ObjectID); ok {
		_id = oid
	}
	if document["productID"].(string) == "" {
		return ProductConfiguration{}, errors.New("product not found")
	}
	structProductConfig := ProductConfiguration{}
	byteProduct, _ := ffjson.Marshal(document)
	_ = ffjson.Unmarshal(byteProduct, &structProductConfig)
	structProductConfig.ID = _id
	return structProductConfig, nil
}

func getProductMicroServices(srv *Service, collection *mongo.Collection, productID string) ([]ProductMicroService, error) {
	if collection == nil {
		collection = getMongoCollection(srv, srv.Config.PRODUCT_MICROSERVICES_COLLECTION)
	}
	microservices, err := getMongoRecords(collection, bson.D{{Key: "productID", Value: productID}})
	if err != nil {
		return []ProductMicroService{}, err
	}
	if len(microservices) == 0 {
		return []ProductMicroService{}, errors.New("product microservices not found")
	}
	documentMicroservices := []ProductMicroService{}
	bytesMicroServices, _ := ffjson.Marshal(microservices)
	_ = ffjson.Unmarshal(bytesMicroServices, &documentMicroservices)
	return documentMicroservices, nil
}

func getProductVersions(srv *Service, collection *mongo.Collection, productID string) ([]ProductVersion, error) {
	if collection == nil {
		collection = getMongoCollection(srv, srv.Config.PRODUCT_MICROSERVICES_COLLECTION)
	}
	versions, err := getMongoRecords(collection, bson.D{{Key: "productID", Value: productID}})
	if err != nil {
		return []ProductVersion{}, err
	}
	if len(versions) == 0 {
		return []ProductVersion{}, errors.New("product versions not found")
	}
	for index, version := range versions {
		if version["productID"].(string) == "" {
			continue
		}
		versions[index]["status"] = "active"
	}
	documentVersions := []ProductVersion{}
	bytesVersions, _ := ffjson.Marshal(versions)
	_ = ffjson.Unmarshal(bytesVersions, &documentVersions)
	return documentVersions, nil
}

func getProductPlans(srv *Service, collection *mongo.Collection, productID string) ([]ProductPlan, error) {
	if collection == nil {
		collection = getMongoCollection(srv, srv.Config.PRODUCT_MICROSERVICES_COLLECTION)
	}
	plans, err := getMongoRecords(collection, bson.D{{Key: "productID", Value: productID}})
	if err != nil {
		return []ProductPlan{}, err
	}
	if len(plans) == 0 {
		return []ProductPlan{}, errors.New("product plans not found")
	}
	documentPlans := []ProductPlan{}
	bytesPlans, _ := ffjson.Marshal(plans)
	_ = ffjson.Unmarshal(bytesPlans, &documentPlans)
	return documentPlans, nil
}

func getProductResources(srv *Service, collection *mongo.Collection, productID string) ([]ProductResource, error) {
	if collection == nil {
		collection = getMongoCollection(srv, srv.Config.PRODUCT_MICROSERVICES_COLLECTION)
	}
	resources, err := getMongoRecords(collection, bson.D{{Key: "productID", Value: productID}})
	if err != nil {
		return []ProductResource{}, err
	}
	if len(resources) == 0 {
		return []ProductResource{}, errors.New("product resources not found")
	}
	documentResources := []ProductResource{}
	bytesResources, _ := ffjson.Marshal(resources)
	_ = ffjson.Unmarshal(bytesResources, &documentResources)
	return documentResources, nil
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

func validateToken(srv *Service, authHeader string) (map[string]interface{}, error) {
	// tokenString := authHeader[len("Bearer "):]
	tokenString := authHeader
	var expiryTime int
	token, err := getToken(srv, bson.D{{Key: "accessToken", Value: tokenString}})
	if err != nil {
		return nil, err
	}
	expiryTime = int(time.Since(token.CreatedAt.Add(time.Hour * 5)).Minutes())
	if token.UserId != "" && (expiryTime < 30) {
		return map[string]interface{}{
			"userId": token.UserId,
			"email":  token.Email,
		}, nil
	} else {
		return nil, errors.New("authorization token has expired")
	}
}

func getToken(srv *Service, filters bson.D) (Tokens, error) {
	if filters == nil {
		filters = bson.D{}
	}
	token, err := getMongoRecord(getMongoCollection(srv, srv.Config.TOKEN_COLLECTION), filters)
	if err != nil {
		return Tokens{}, err
	}
	var resultToken Tokens
	tokenBytes, err := ffjson.Marshal(token)
	if err != nil {
		return Tokens{}, err
	}
	if err := ffjson.Unmarshal(tokenBytes, &resultToken); err != nil {
		return Tokens{}, err
	}
	return resultToken, nil
}

// func createVersions(srv *Service, versions []ProductVersion) ([]ProductVersion, error) {
// 	productVersionsCollection := getMongoCollection(srv, srv.Config.PRODUCT_VERSIONS_COLLECTION)
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
// 	microserviceCollection := getMongoCollection(srv, srv.Config.PRODUCT_MICROSERVICES_COLLECTION)
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
