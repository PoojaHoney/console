package main

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/ffjson/ffjson"
)

// CreateProduct Creates a new product.
// @Summary Creates a new product
// @Description Creates a new product with the provided product data.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body Product true "Product object containing product data"
// @Router /api/product/v1/createProduct [post]
func (srv *Service) CreateProduct(ctx *fiber.Ctx) error {
	var product Product

	if err := ffjson.Unmarshal(ctx.Body(), &product); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	if err := validate.Struct(product); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in validating body",
			"error":   err.Error()})
	}
	// createdBy := ctx.Context().UserValue("email").(string)
	createdBy := ""
	createdProduct, err := createProduct(srv, product, createdBy)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in creating product",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "product created successfully",
		"data":    createdProduct})
}

// CreateMicroService Creates new microservices for a product.
// @Summary Creates new microservices for a product
// @Description  Creates new microservices for a product.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param microservices body ProductMicroService true "Microservices array containing microservice data"
// @Router /api/product/v1/createMicroService [post]
func (srv *Service) CreateMicroService(ctx *fiber.Ctx) error {
	var microservice ProductMicroService

	if err := ffjson.Unmarshal(ctx.Body(), &microservice); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	if err := validate.Struct(microservice); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in validating body",
			"error":   err.Error()})
	}
	createdMicroservice, err := createMicroService(srv, microservice)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in creating product microservice",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "product's microservice created successfully",
		"data":    createdMicroservice})
}

// CreateResource Creates a product resources.
// @Summary Creates a product resources
// @Description Creates a product resources with the provided resource data.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param resources body ProductResource true "Resources array containing resources data of a product"
// @Router /api/product/v1/createResource [post]
func (srv *Service) CreateResource(ctx *fiber.Ctx) error {
	var resource ProductResource

	if err := ffjson.Unmarshal(ctx.Body(), &resource); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	if err := validate.Struct(resource); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in validating body",
			"error":   err.Error()})
	}
	createdResource, err := createResource(srv, resource)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in creating resource",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "product's resource created successfully",
		"data":    createdResource})
}

// CreateConfiguration Creates a new product Configuration.
// @Summary Creates product Configuration
// @Description Creates product Configuration with the provided Configuration data.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param configuration body ProductConfiguration true "Configuration object containing Configuration data"
// @Router /api/product/v1/createConfiguration [post]
func (srv *Service) CreateConfiguration(ctx *fiber.Ctx) error {
	var configuration ProductConfiguration

	if err := ffjson.Unmarshal(ctx.Body(), &configuration); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	if err := validate.Struct(configuration); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in validating body",
			"error":   err.Error()})
	}
	createdConfiguration, err := createConfiguration(srv, configuration)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in creating product configuration",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "product's configuration created successfully",
		"data":    createdConfiguration})
}

// CreatePlan Creates a product Plans.
// @Summary Creates product Plans
// @Description Creates product Plans with the provided Plans data.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Plans body ProductPlan true "Plans array containing Product Plans data"
// @Router /api/product/v1/createPlan [post]
func (srv *Service) CreatePlan(ctx *fiber.Ctx) error {
	var plan ProductPlan

	if err := ffjson.Unmarshal(ctx.Body(), &plan); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	if err := validate.Struct(plan); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in validating body",
			"error":   err.Error()})
	}
	createdPlan, err := createPlan(srv, plan)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in creating product plan",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "product's plan created successfully",
		"data":    createdPlan})
}

// CreateVersion Creates a product versions.
// @Summary Creates product versions
// @Description Creates product versions with the provided versions data.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param version body ProductVersion true "Versions array containing Product versions data"
// @Router /api/product/v1/createVersion [post]
func (srv *Service) CreateVersion(ctx *fiber.Ctx) error {
	var version ProductVersion

	if err := ffjson.Unmarshal(ctx.Body(), &version); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	if err := validate.Struct(version); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in validating body",
			"error":   err.Error()})
	}
	createdVersion, err := createVersion(srv, version)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in creating product version",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "product version created successfully",
		"data":    createdVersion})
}

// ActivateProduct Activate the product, makes the product available to the user/customers.
// @Summary Activates the product
// @Description Activate the product, makes the product available to the user/customers.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body object true "Object Containing Product id"
// @Router /api/product/v1/activateProduct [post]
func (srv *Service) ActivateProduct(ctx *fiber.Ctx) error {
	product := map[string]interface{}{}

	if err := ffjson.Unmarshal(ctx.Body(), &product); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
		})
	}
	if product["productID"].(string) == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
		})

	}
	fullProductDetails, err := activateProduct(srv, product["productID"].(string))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in activating product",
		})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "product activated successfully",
		"data":    fullProductDetails,
	})
}

func (srv *Service) Update(ctx *fiber.Ctx) error {
	return nil
}

func (srv *Service) Delete(ctx *fiber.Ctx) error {
	return nil
}

func (srv *Service) GetFullProductDetails(ctx *fiber.Ctx) error {
	return nil
}

// GetFieldCatalogues Field catalogues of products Service.
// @Summary Fields Catalogues of product service
// @Description Gets all fields catalogues of product service tables.
// @Produce json
// @Security BearerAuth
// @Router /api/product/v1/fieldCatalogues [get]
func (srv *Service) FieldCatalogues(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Field catalogues",
		"data":    fieldCatalogues(),
	})
}

func ValidateToken(srv *Service) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization failed",
				"error":   errors.New("invalid authorization token"),
			})
		}
		user, err := validateToken(srv, authHeader)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization failed",
				"error":   err.Error(),
			})
		}
		if user == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization failed",
				"error":   errors.New("invalid productId or product does not exists"),
			})
		}
		if user["userId"].(string) != "" || user["email"].(string) != "" {
			ctx.Context().SetUserValue("email", user["email"].(string))
			ctx.Next()
		}
		return nil
	}
}
