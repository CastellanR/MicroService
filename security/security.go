package security

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CastellanR/UserFeedback-Microservice/tools/env"
	"github.com/CastellanR/UserFeedback-Microservice/tools/errors"
	gocache "github.com/patrickmn/go-cache"
	validator "gopkg.in/go-playground/validator.v9"
)

var cache = gocache.New(60*time.Minute, 10*time.Minute)

// User
type User struct {
	ID          string   `json:"id"  validate:"required"`
	Name        string   `json:"name"  validate:"required"`
	Permissions []string `json:"permissions"`
	Login       string   `json:"login"  validate:"required"`
}

// Validate token
/**
 * @apiDefine AuthHeader
 *
 * @apiExample {String} Auth Header
 *    Authorization=bearer {token}
 *
 * @apiErrorExample 401 Unauthorized
 *    HTTP/1.1 401 Unauthorized
 */
func Validate(token string) (*User, error) {
	// Si esta en cache, retornamos el cache
	if found, ok := cache.Get(token); ok {
		if user, ok := found.(*User); ok {
			return user, nil
		}
	}

	user, err := getRemote(token)
	if err != nil {
		return nil, errors.Unauthorized
	}

	// the token is added to cache and returned
	cache.Set(token, user, gocache.DefaultExpiration)

	return user, nil
}

func getRemote(token string) (*User, error) {
	// Search remote user
	req, err := http.NewRequest("GET", env.Get().SecurityServerURL+"/v1/users/current", nil)
	if err != nil {
		return nil, errors.Unauthorized
	}
	req.Header.Add("Authorization", "bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.Unauthorized
	}
	defer resp.Body.Close()

	user := &User{}
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Invalidate cache token
func Invalidate(token string) {
	cache.Delete(token[7:])
	log.Output(1, fmt.Sprintf("Token invalid: %s", token))
}
