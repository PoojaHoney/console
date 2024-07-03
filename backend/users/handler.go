package main

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	// user.ID = primitive.NewObjectID()
	user.SaltStored = salt
	user.Active = false
	user.LastChangedBy = createdBy
	user.Deleted = false
	userExists := srv.PostgresDB.Where("email = ?", user.Email).Or("account_number = ?", user.AccountNumber).First(&User{})
	if userExists.Error == nil {
		return User{}, errors.New("user already exists")
	}
	if userExists.Error != nil && userExists.Error != gorm.ErrRecordNotFound {
		return User{}, err
	}
	result := srv.PostgresDB.Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func getUser(db *gorm.DB, filters map[string]interface{}) (*User, error) {
	var user User
	for key, value := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}
	result := db.First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func getUsers(db *gorm.DB, filters map[string]interface{}) ([]User, error) {
	var users []User
	query := db.Model(&User{})
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
	// for _, filter := range filters {
	// 	condition := query
	// 	for key, value := range filter {
	// 		condition = condition.Or(fmt.Sprintf("%s = ?", key), value)
	// 	}
	// 	query = condition
	// }
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func getUserOTP(filters map[string]interface{}, db *gorm.DB) (UserOTP, error) {
	var userOTP UserOTP
	for key, value := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}
	result := db.First(&userOTP)
	if result.Error != nil {
		return UserOTP{}, result.Error
	}
	return userOTP, nil
}

func updateUser(user User, srv *Service, updatedBy string) (User, error) {
	userExists, err := getUser(srv.PostgresDB, map[string]interface{}{"id": user.ID})
	if err != nil {
		return User{}, err
	}
	if userExists == nil {
		return User{}, errors.New("user not found")
	}
	user.LastChangedBy = updatedBy
	result := srv.PostgresDB.Save(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func deleteUser(userId string, hardDelete bool, srv *Service) error {
	id, _ := uuid.Parse(userId)
	user, err := getUser(srv.PostgresDB, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}
	if user.Deleted && !hardDelete {
		return errors.New("user already deleted")
	}
	if hardDelete {
		result := srv.PostgresDB.Unscoped().Delete(&User{}, id)
		if result.Error != nil {
			return result.Error
		}
	} else {
		user.Deleted = true
		result := srv.PostgresDB.Save(&user)
		if result.Error != nil {
			return result.Error
		}
	}
	return err
}

func updateUserPassword(userPassword Password, srv *Service, updatedBy string) (User, error) {
	user, err := getUser(srv.PostgresDB, map[string]interface{}{"id": userPassword.UserId})
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
	result := srv.PostgresDB.Save(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return *user, nil
}

func sendOTPVerificationMail(srv *Service, user User) error {

	userOTP, err := getUserOTP(map[string]interface{}{"email": user.Email}, srv.PostgresDB)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	fmt.Println(userOTP)
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
		UserId: user.ID,
		Expiry: time.Now().Add(time.Minute * 30),
		Times:  userOTP.Times + 1,
	}

	if userOTP.OTP != 0 {
		otpVerification.ID = userOTP.ID
		result := srv.PostgresDB.Save(&otpVerification)
		if result.Error != nil {
			return err
		}
		return nil
	}
	result := srv.PostgresDB.Create(&otpVerification)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func verifyOTP(srv *Service, otpInput VerifyOTP) error {
	userOTP, err := getUserOTP(map[string]interface{}{"email": otpInput.Email}, srv.PostgresDB)
	if err != nil {
		return err
	}
	if userOTP.OTP != int(otpInput.OTP) {
		return errors.New("wrong OTP")
	}
	if userOTP.Expiry.Add(time.Hour * 5).Before(time.Now()) {
		return errors.New("OTP expired")
	}
	user, err := getUser(srv.PostgresDB, map[string]interface{}{"id": userOTP.UserId})
	if err != nil {
		return err
	}
	user.Active = true
	result := srv.PostgresDB.Save(&user)
	if result.Error != nil {
		return err
	}
	result = srv.PostgresDB.Unscoped().Delete(&UserOTP{}, userOTP.ID)
	if result.Error != nil {
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
		userID, _ := uuid.Parse(newAccessToken["userId"].(string))
		token := Tokens{
			UserId:       userID,
			AccessToken:  newAccessToken["access_token"].(string),
			Email:        loginCrds.Email,
			RefreshToken: loginCrds.RefreshToken,
			Expiry:       timeNow.Add(time.Hour * 24 * 7).Unix(),
		}
		result := srv.PostgresDB.Save(&token)
		if result.Error != nil {
			return Tokens{}, result.Error
		}
		return token, nil
	}
	user, err := getUser(srv.PostgresDB, map[string]interface{}{"email": loginCrds.Email, "active": true})
	if err != nil {
		err := fmt.Sprintf("%s or user is inactive", err.Error())
		return Tokens{}, errors.New(err)
	}
	valid := verifyHashPassword(loginCrds.Password, user.Password, user.SaltStored, 0)
	if !valid {
		return Tokens{}, errors.New("email or password is incorrect")
	}
	newTokens, err := generateTokens(user.ID.String(), user.Email)
	if err != nil {
		return Tokens{}, err
	}
	timeNow := time.Now()
	token := Tokens{
		UserId:       user.ID,
		Email:        loginCrds.Email,
		AccessToken:  newTokens["access_token"].(string),
		RefreshToken: newTokens["refreshToken"].(string),
		Expiry:       timeNow.Add(time.Hour * 24 * 7).Unix(),
	}
	result := srv.PostgresDB.Create(&token)
	if result.Error != nil {
		return Tokens{}, result.Error
	}
	return token, nil
}

// func getUsers(srv *Service, filters bson.M) ([]User, error) {
// 	if filters == nil {
// 		filters = bson.M{}
// 	}
// 	users, err := getMongoRecords(getMongoCollection(srv, srv.Config.USERS_COLLECTION), filters)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var resultUsers []User
// 	for _, user := range users {
// 		var userID primitive.ObjectID
// 		if oid, ok := user["_id"].(primitive.ObjectID); ok {
// 			userID = oid
// 		}
// 		resultUser := User{}
// 		userBytes, err := ffjson.Marshal(user)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if err := ffjson.Unmarshal(userBytes, &resultUser); err != nil {
// 			return nil, err
// 		}
// 		resultUser.ID = userID
// 		resultUsers = append(resultUsers, resultUser)
// 	}
// 	return resultUsers, nil
// }

// func getUser(srv *Service, filters bson.D) (User, error) {
// 	if filters == nil {
// 		filters = bson.D{}
// 	}
// 	user, err := getMongoRecord(getMongoCollection(srv, srv.Config.USERS_COLLECTION), filters)
// 	if err != nil {
// 		return User{}, err
// 	}
// 	var userID primitive.ObjectID
// 	if oid, ok := user["_id"].(primitive.ObjectID); ok {
// 		userID = oid
// 	}
// 	resultUser := User{}
// 	userBytes, err := ffjson.Marshal(user)
// 	if err != nil {
// 		return User{}, err
// 	}
// 	if err := ffjson.Unmarshal(userBytes, &resultUser); err != nil {
// 		return User{}, err
// 	}
// 	resultUser.ID = userID
// 	return resultUser, nil
// }

func getToken(srv *Service, filters map[string]interface{}) (Tokens, error) {
	var token Tokens
	db := srv.PostgresDB
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
