package main

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/ffjson/ffjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func fieldCatalogues(tables ...interface{}) map[string]interface{} {
	if tables == nil {
		tables = []interface{}{
			User{}, Personal{}, Organization{}, Address{}, Settings{}}
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
func createUser(user User, srv *Service, createdBy string) (User, error) {
	user.AccountNumber = generateAccountNumber()
	hashedPassword, salt, err := generateHashPassword(user.Password)
	if err != nil {
		return User{}, err
	}
	user.Password = hashedPassword
	user.ID = primitive.NewObjectID()
	user.SaltStored = salt
	user.Active = false
	user.Deleted = false
	user.LastChangedOn = time.Now()
	if createdBy == "" {
		createdBy = user.ID.Hex()
	}
	user.LastChangedBy = createdBy
	user.CreatedOn = time.Now()
	userCollection := getMongoCollection(srv, srv.Config.USERS_COLLECTION)
	recordsExists, err := getMongoRecordsCount(
		userCollection,
		bson.D{{
			Key: "$or", Value: bson.A{
				bson.D{{Key: "email", Value: user.Email}},
				bson.D{{Key: "accountNumber", Value: user.AccountNumber}},
			}}})
	if err != nil {
		return User{}, err
	}
	if recordsExists > 0 {
		return User{}, errors.New("user already exists")
	}
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func updateUser(user User, srv *Service, updatedBy string) (User, error) {
	count, err := getMongoRecordsCount(getMongoCollection(srv, srv.Config.USERS_COLLECTION), bson.D{{Key: "_id", Value: user.ID}})
	if err != nil {
		return User{}, err
	}
	if count == 0 {
		return User{}, errors.New("user not found")
	}
	user.LastChangedOn = time.Now()
	user.LastChangedBy = updatedBy

	_, err = getMongoCollection(srv, srv.Config.USERS_COLLECTION).UpdateOne(context.TODO(),
		bson.D{{Key: "_id", Value: user.ID}}, bson.D{{Key: "$set", Value: user}})
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func deleteUser(userId string, hardDelete bool, srv *Service) error {
	id, _ := primitive.ObjectIDFromHex(userId)
	var err error
	user, err := getUser(srv, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}
	if user.Deleted && !hardDelete {
		return errors.New("user already deleted")
	}
	if hardDelete {
		_, err = getMongoCollection(srv, srv.Config.USERS_COLLECTION).DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: id}})
	} else {
		_, err = getMongoCollection(srv, srv.Config.USERS_COLLECTION).UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: id}},
			bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: true}}}})
	}
	return err
}

func updateUserPassword(userPassword Password, srv *Service, updatedBy string) (User, error) {
	user, err := getUser(srv, bson.D{{Key: "_id", Value: userPassword.UserId}})
	if err != nil {
		return User{}, err
	}
	hashedPassword, salt, err := generateHashPassword(userPassword.Password)
	if err != nil {
		return User{}, err
	}
	user.Password = hashedPassword
	user.SaltStored = salt
	user.LastChangedBy = updatedBy
	user.LastChangedOn = time.Now()
	_, err = getMongoCollection(srv, srv.Config.USERS_COLLECTION).UpdateOne(context.TODO(),
		bson.D{{Key: "_id", Value: user.ID}}, bson.D{{Key: "$set", Value: user}})
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func sendOTPVerificationMail(srv *Service, user User) error {

	userOTP, err := getUserOTP(srv, bson.D{{Key: "email", Value: user.Email}})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	var otp int
	if userOTP.Expiry.Add(time.Hour * 5).After(time.Now()) {
		otp = userOTP.OTP
	} else {
		otp = generateOTP()
	}

	input := map[string]interface{}{
		"otp":      otp,
		"email":    user.Email,
		"expiry":   30,
		"subject":  "Email Verification for Console",
		"userName": user.Name,
	}
	templateContent, err := os.ReadFile("emailTemplates/verification.html")
	if err != nil {
		return err
	}
	emailBody, err := generateEmailBody(string(templateContent), input)
	if err != nil {
		return err
	}

	smtpServer := srv.Config.SMTP_SERVER
	smtpPort := srv.Config.SMTP_PORT
	smtpBaseMail := srv.Config.SMTP_BASEEMAIL
	smtpPassKey := srv.Config.SMTP_PASSKEY

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html\r\n\r\n%s",
		smtpBaseMail, input["email"].(string), input["subject"].(string), emailBody)
	auth := smtp.PlainAuth("", smtpBaseMail, smtpPassKey, smtpServer)
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, smtpBaseMail, []string{input["email"].(string)}, []byte(message))
	if err != nil {
		return err
	}

	otpVerification := UserOTP{
		OTP:    input["otp"].(int),
		Email:  input["email"].(string),
		UserId: user.ID.Hex(),
		Expiry: time.Now().Add(time.Minute * 30),
		Times:  userOTP.Times + 1,
	}

	if userOTP.OTP != 0 {
		otpVerification.ID = userOTP.ID
		_, err = getMongoCollection(srv, srv.Config.USER_OTP_COLLECTION).UpdateOne(context.TODO(),
			bson.D{{Key: "_id", Value: userOTP.ID}}, bson.D{{Key: "$set", Value: otpVerification}})
		if err != nil {
			return err
		}
		return nil
	}
	otpVerification.ID = primitive.NewObjectID()
	_, err = getMongoCollection(srv, srv.Config.USER_OTP_COLLECTION).InsertOne(context.TODO(), otpVerification)
	if err != nil {
		return err
	}
	return nil
}

func verifyOTP(srv *Service, otpInput VerifyOTP) error {
	userOTP, err := getUserOTP(srv, bson.D{{Key: "email", Value: otpInput.Email}})
	if err != nil {
		return err
	}
	if userOTP.OTP != int(otpInput.OTP) {
		return errors.New("wrong OTP")
	}
	if userOTP.Expiry.Add(time.Hour * 5).Before(time.Now()) {
		return errors.New("OTP expired")
	}
	userId, _ := primitive.ObjectIDFromHex(userOTP.UserId)
	user, err := getUser(srv, bson.D{{Key: "_id", Value: userId}})
	if err != nil {
		return err
	}
	user.Active = true
	_, err = getMongoCollection(srv, srv.Config.USERS_COLLECTION).UpdateOne(context.TODO(),
		bson.D{{Key: "_id", Value: userId}}, bson.D{{Key: "$set", Value: user}})
	if err != nil {
		return err
	}
	_, err = getMongoCollection(srv, srv.Config.USER_OTP_COLLECTION).DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: userOTP.ID}})
	if err != nil {
		return err
	}
	return nil
}

func login(srv *Service, loginCrds LoginCredentials) (Tokens, error) {
	if loginCrds.RefreshToken != "" {
		newAccessToken, err := refreshAccessToken(loginCrds.RefreshToken)
		if err != nil {
			return Tokens{}, err
		}
		timeNow := time.Now()
		token := Tokens{
			UserId:       newAccessToken["userId"].(string),
			AccessToken:  newAccessToken["accessToken"].(string),
			Email:        loginCrds.Email,
			RefreshToken: loginCrds.RefreshToken,
			Expiry:       timeNow.Add(time.Hour * 24 * 7).Unix(),
			CreatedAt:    timeNow,
		}
		tokenCollection := getMongoCollection(srv, srv.Config.TOKEN_COLLECTION)
		_, err = tokenCollection.UpdateOne(context.TODO(), bson.D{{Key: "userId", Value: newAccessToken["userId"].(string)}},
			bson.D{{Key: "$set", Value: token}})
		if err != nil {
			return Tokens{}, err
		}
		return token, nil
	}
	user, err := getUser(srv, bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "email", Value: loginCrds.Email}},
		bson.D{{Key: "active", Value: bson.M{"$ne": false}}},
	}}})
	if err != nil {
		err := fmt.Sprintf("%s or user is inactive", err.Error())
		return Tokens{}, errors.New(err)
	}
	valid := verifyHashPassword(loginCrds.Password, user.Password, user.SaltStored, 0)
	if !valid {
		return Tokens{}, errors.New("email or password is incorrect")
	}
	newTokens, err := generateTokens(user.ID.Hex(), user.Email)
	if err != nil {
		return Tokens{}, err
	}
	timeNow := time.Now()
	token := Tokens{
		UserId:       user.ID.Hex(),
		Email:        loginCrds.Email,
		AccessToken:  newTokens["accessToken"].(string),
		RefreshToken: newTokens["refreshToken"].(string),
		Expiry:       timeNow.Add(time.Hour * 24 * 7).Unix(),
		CreatedAt:    timeNow,
	}
	tokenCollection := getMongoCollection(srv, srv.Config.TOKEN_COLLECTION)
	_, err = tokenCollection.InsertOne(context.TODO(), token)
	if err != nil {
		return Tokens{}, err
	}
	return token, nil
}

func getUsers(srv *Service, filters bson.M) ([]User, error) {
	if filters == nil {
		filters = bson.M{}
	}
	users, err := getMongoRecords(getMongoCollection(srv, srv.Config.USERS_COLLECTION), filters)
	if err != nil {
		return nil, err
	}
	var resultUsers []User
	for _, user := range users {
		var userID primitive.ObjectID
		if oid, ok := user["_id"].(primitive.ObjectID); ok {
			userID = oid
		}
		resultUser := User{}
		userBytes, err := ffjson.Marshal(user)
		if err != nil {
			return nil, err
		}
		if err := ffjson.Unmarshal(userBytes, &resultUser); err != nil {
			return nil, err
		}
		resultUser.ID = userID
		resultUsers = append(resultUsers, resultUser)
	}
	return resultUsers, nil
}

func getUser(srv *Service, filters bson.D) (User, error) {
	if filters == nil {
		filters = bson.D{}
	}
	user, err := getMongoRecord(getMongoCollection(srv, srv.Config.USERS_COLLECTION), filters)
	if err != nil {
		return User{}, err
	}
	var userID primitive.ObjectID
	if oid, ok := user["_id"].(primitive.ObjectID); ok {
		userID = oid
	}
	resultUser := User{}
	userBytes, err := ffjson.Marshal(user)
	if err != nil {
		return User{}, err
	}
	if err := ffjson.Unmarshal(userBytes, &resultUser); err != nil {
		return User{}, err
	}
	resultUser.ID = userID
	return resultUser, nil
}

func getUserOTP(srv *Service, filters bson.D) (UserOTP, error) {
	if filters == nil {
		filters = bson.D{}
	}
	userOTP, err := getMongoRecord(getMongoCollection(srv, srv.Config.USER_OTP_COLLECTION), filters)
	if err != nil {
		return UserOTP{}, err
	}
	var userOTPID primitive.ObjectID
	if oid, ok := userOTP["_id"].(primitive.ObjectID); ok {
		userOTPID = oid
	}
	var resultUserOTP UserOTP
	userOTPBytes, err := ffjson.Marshal(userOTP)
	if err != nil {
		return UserOTP{}, err
	}
	if err := ffjson.Unmarshal(userOTPBytes, &resultUserOTP); err != nil {
		return UserOTP{}, err
	}
	resultUserOTP.ID = userOTPID
	return resultUserOTP, nil
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
