package main

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/ffjson/ffjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRegistration registers a new user.
// @Summary Register a new user
// @Description Register a new user with the provided user data.
// @Accept json
// @Produce json
// @Param user body User true "User object containing user data"
// @Router /api/user/v1/register [post]
func (srv *Service) UserRegistration(ctx *fiber.Ctx) error {
	var user User

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	if err := validate.Struct(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in validating body",
			"error":   err.Error()})
	}
	user, err := createUser(user, srv, "")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in creating user",
			"error":   err.Error()})
	}
	err = sendOTPVerificationMail(srv, user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in sending email",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "User created successfully & Send Verfication Email",
		"user":    user})
}

// Create Creates a new user.
// @Summary Creates a new user
// @Description Creates a new user with the provided user data.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body User true "User object containing user data"
// @Router /api/user/v1/create [post]
func (srv *Service) Create(ctx *fiber.Ctx) error {
	var user User

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	if err := validate.Struct(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in validating body",
			"error":   err.Error()})
	}
	user, err := createUser(user, srv, ctx.Context().UserValue("email").(string))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in creating user",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    user})
}

// Update Updates a user.
// @Summary Updates a user
// @Description Updates a user with the provided user data.
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Param user body User true "User object containing user data"
// @Router /api/user/v1/update/{id} [put]
func (srv *Service) Update(ctx *fiber.Ctx) error {
	var user User
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user id"})
	}
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	if err := validate.Struct(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in validating body",
			"error":   err.Error()})
	}
	user.ID, _ = primitive.ObjectIDFromHex(id)
	_, err := updateUser(user, srv, ctx.Context().UserValue("email").(string))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in updating user",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "User updated successfully"})
}

// UpdatePassword Updates a user password.
// @Summary Updates a user password
// @Description Updates a user password with the provided user password data.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param password body Password true "User object containing user password data"
// @Router /api/user/v1/updatePassword/{id} [put]
func (srv *Service) UpdatePassword(ctx *fiber.Ctx) error {
	var userPassword Password
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user id"})
	}
	if err := ctx.BodyParser(&userPassword); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	if err := validate.Struct(userPassword); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in validating body",
			"error":   err.Error()})
	}
	userId, _ := primitive.ObjectIDFromHex(id)
	userPassword.UserId = userId
	_, err := updateUserPassword(userPassword, srv, ctx.Context().UserValue("email").(string))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in updating user password",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "User Password updated successfully"})
}

// Delete Deletes a user.
// @Summary Deletes a user
// @Description Deletes a user with the provided user id.
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Param hardDelete path bool true "Hard delete or soft delete"
// @Router /api/user/v1/delete/{id}/{hardDelete} [delete]
func (srv *Service) Delete(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	if userId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user id"})
	}
	var forceDelete bool
	if ctx.Params("hardDelete") == "" {
		forceDelete = false
	}
	forceDelete, _ = strconv.ParseBool(ctx.Params("hardDelete"))
	err := deleteUser(userId, forceDelete, srv)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in deleting a user",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "User Deleted successfully"})
}

// GetFieldCatalogues Field catalogues of Users Service.
// @Summary Fields Catalogues of user service
// @Description Gets all fields catalogues of user service tables.
// @Produce json
// @Router /api/user/v1/fieldCatalogues [get]
func (srv *Service) FieldCatalogues(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":         "Field catalogues",
		"fieldCatalogues": fieldCatalogues(),
	})
}

// Get Get users
// @Summary Get or Read users from database based on filters if any
// @Description Get or Read users from database based on filters if any.
// @Produce json
// @Param filters query string false "Filters for user retrieval"
// @Security BearerAuth
// @Router /api/user/v1/get [get]
func (srv *Service) Get(ctx *fiber.Ctx) error {
	filtersStr := ctx.Query("filters")
	var filters []map[string]interface{}
	var mongoFilters bson.M
	if filtersStr != "" {
		if err := ffjson.Unmarshal([]byte(filtersStr), &filters); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid filters format"})
		}
		mongoFilters = prepareMongoFilters(filters)
	}
	users, err := getUsers(srv, mongoFilters)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in getting users",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Users",
		"users":   users,
	})
}

// sendOTPVerificationMail Send OTP Verification Mail To User.
// @Summary Send OTP Verification Mail To User
// @Description Send OTP Verification Mail To User.
// @Accept json
// @Produce json
// @Param email body SendOTPVerificationMail true "Object containing email"
// @Router /api/user/v1/sendOTPVerificationMail [post]
func (srv *Service) SendOTPVerificationMail(ctx *fiber.Ctx) error {
	var email SendOTPVerificationMail
	if err := ffjson.Unmarshal(ctx.Body(), &email); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	user, err := getUser(srv, bson.D{{Key: "email", Value: email.Email}})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User Does not exists",
			"error":   err.Error()})
	}
	err = sendOTPVerificationMail(srv, user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in sending email",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": " Send Verification Email"})
}

// VerifyOTP Verify the OTP from User Mail.
// @Summary Verify the OTP from User Mail
// @Description Verify the OTP from User Mail.
// @Accept json
// @Produce json
// @Param otp body VerifyOTP true "Object containing otp and email"
// @Router /api/user/v1/verifyOTP [post]
func (srv *Service) VerifyOTP(ctx *fiber.Ctx) error {
	var otpInput VerifyOTP
	if err := ffjson.Unmarshal(ctx.Body(), &otpInput); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	err := verifyOTP(srv, otpInput)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in OTP Verification",
			"error":   err.Error()})
	}
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "OTP Verified Successfully"})
}

// Login Verify the User by login api with email and password.
// @Summary Verify the User by login api with email and password
// @Description Verify the User by login api with email and password.
// @Accept json
// @Produce json
// @Param loginCrds body LoginCredentials true "Object containing email, password and token"
// @Router /api/user/v1/login [post]
func (srv *Service) Login(ctx *fiber.Ctx) error {
	var loginCrds LoginCredentials
	if err := ffjson.Unmarshal(ctx.Body(), &loginCrds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in parsing body",
			"error":   err.Error()})
	}
	if err := validate.Struct(loginCrds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing required fields",
			"error":   err.Error(),
		})
	}
	token, err := login(srv, loginCrds)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to login",
			"error":   err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(token)
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
				"error":   errors.New("invalid userId or User does not exists"),
			})
		}
		if user["userId"].(string) != "" || user["email"].(string) != "" {
			ctx.Context().SetUserValue("userId", user["userId"].(string))
			ctx.Context().SetUserValue("email", user["email"].(string))
			ctx.Next()
		}
		return nil
	}
}
