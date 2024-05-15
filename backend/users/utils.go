package main

import (
	"bytes"
	"context"
	cryptoRand "crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"math/rand"
	textTemplate "text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/xdg-go/pbkdf2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getMongoCollection(srv *Service, collectioName string) *mongo.Collection {
	return srv.MongoDB.Database(srv.Config.MONGO_DATABASE).Collection(collectioName)
}

func getMongoRecordsCount(collection *mongo.Collection, filter interface{}) (int64, error) {
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getMongoRecord(collection *mongo.Collection, filter interface{}) (map[string]interface{}, error) {
	var record map[string]interface{}
	err := collection.FindOne(context.TODO(), filter).Decode(&record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func getMongoRecords(collection *mongo.Collection, filter interface{}) ([]map[string]interface{}, error) {
	var records []map[string]interface{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var record map[string]interface{}
		err := cursor.Decode(&record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func generateHashPassword(password string) (string, string, error) {
	key := make([]byte, 0x10)
	_, err := cryptoRand.Read(key)
	if err != nil {
		return "", "", err
	}
	hashedPassword := pbkdf2.Key([]byte(password), key, 1000, 49, sha1.New)
	return base64.StdEncoding.EncodeToString(hashedPassword), base64.StdEncoding.EncodeToString(key), nil
}

func generateAccountNumber() int {
	accountNumber := rand.Intn(9000000000) + 1000000000
	return accountNumber
}

func verifyHashPassword(password string, hashedPassword string, saltStored string, version int) bool {
	if version == 1 {
		hashedPassword, err := base64.StdEncoding.DecodeString(hashedPassword)
		if err != nil {
			return false
		}
		saltDst := make([]byte, 0x10)
		hashedKey := make([]byte, 0x20)
		copy(saltDst, hashedPassword[1:17])
		copy(hashedKey, hashedPassword[17:49])
		password_hash := pbkdf2.Key([]byte(password), saltDst, 1000, 49, sha1.New)
		return bytes.Equal(password_hash, hashedKey)
	}
	saltDst, _ := base64.StdEncoding.DecodeString(saltStored)
	hashPasswordStr, _ := base64.StdEncoding.DecodeString(hashedPassword)
	password_hash := pbkdf2.Key([]byte(password), saltDst, 1000, 49, sha1.New)
	return bytes.Equal(password_hash, hashPasswordStr)
}

func generateTokens(userId string, email string) (map[string]interface{}, error) {
	accessClaims := jwt.MapClaims{
		"userId": userId,
		"email":  email,
		"exp":    time.Now().Add(time.Minute * 15).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte("r9tWqmF#Y2P$%78G"))
	if err != nil {
		return nil, err
	}
	refreshClaims := jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte("6X!yT@Zu$2F&vP#N"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"accessToken":  accessTokenString,
		"refreshToken": refreshTokenString,
	}, nil
}

func refreshAccessToken(refreshToken string) (map[string]interface{}, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("6X!yT@Zu$2F&vP#N"), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	tokens, err := generateTokens(claims["userId"].(string), claims["email"].(string))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"accessToken": tokens["accessToken"].(string),
		"email":       claims["email"].(string),
		"userId":      claims["userId"].(string),
	}, nil
}

func generateOTP() int {
	otp := rand.Intn(900000) + 100000
	return otp
}

func generateEmailBody(emailContent string, input map[string]interface{}) (string, error) {
	template, err := textTemplate.New("emailTemplate").Parse(emailContent)
	if err != nil {
		return "", err
	}

	var templateOutput bytes.Buffer
	if err := template.Execute(&templateOutput, input); err != nil {
		return "", err
	}
	return templateOutput.String(), nil
}

func prepareMongoFilters(filters []map[string]interface{}) bson.M {
	if filters == nil {
		return bson.M{}
	}
	var andConditions []bson.M
	for _, filter := range filters {
		key := filter["key"].(string)
		value := filter["value"]

		var orConditions []bson.M

		switch v := value.(type) {
		case string:
			if key == "id" || key == "_id" {
				temp, _ := primitive.ObjectIDFromHex(value.(string))
				andConditions = append(andConditions, bson.M{key: bson.M{"$eq": temp}})
			} else {
				andConditions = append(andConditions, bson.M{key: bson.M{"$eq": v}})
			}
			// orConditions = append(orConditions, bson.M{key: bson.M{"$eq": v}})
		case []string:
			for _, val := range v {
				orConditions = append(orConditions, bson.M{key: bson.M{"$eq": val}})
			}
		}

		if len(orConditions) > 0 {
			andConditions = append(andConditions, bson.M{"$or": orConditions})
		}
	}

	if len(andConditions) > 0 {
		return bson.M{"$and": andConditions}
	}
	return bson.M{}
}
