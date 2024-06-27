package main

import (
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"math/rand"
	textTemplate "text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/xdg-go/pbkdf2"
)

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
	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	access_tokenString, err := access_token.SignedString([]byte("r9tWqmF#Y2P$%78G"))
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
		"access_token": access_tokenString,
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
		"access_token": tokens["access_token"].(string),
		"email":        claims["email"].(string),
		"userId":       claims["userId"].(string),
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

func preparePostgresFilters(filters []map[string]interface{}) map[string]interface{} {
	if filters == nil {
		return nil
	}
	allFilters := map[string]interface{}{}
	andConditions := map[string]interface{}{}
	for _, filter := range filters {
		key := filter["key"].(string)
		value := filter["value"]

		orConditions := map[string]interface{}{}

		switch v := value.(type) {
		case string:
			if key == "id" || key == "_id" {
				temp, _ := uuid.Parse(value.(string))
				andConditions[key] = temp
			} else {
				andConditions[key] = v
			}
			// orConditions = append(orConditions, bson.M{key: bson.M{"$eq": v}})
		case []string:
			for _, val := range v {
				orConditions[key] = val
			}
		}

		if len(orConditions) > 0 || orConditions != nil {
			allFilters["or"] = orConditions
		}
	}

	if len(andConditions) > 0 {
		allFilters["and"] = andConditions
		return allFilters
	}
	return nil
}
