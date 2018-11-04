package feedback

import (
	"github.com/CastellanR/UserFeedback-Microservice/tools/env"
	"github.com/CastellanR/UserFeedback-Microservice/tools/errors"
	"github.com/CastellanR/UserFeedback-Microservice/rabbit"
	"github.com/go-redis/redis"
)

// Insert feedback into Database
func Insert(feedback *Feedback, cartID) (string, error) {
	if err := feedback.validateSchema(); err != nil {
		return "", err
	}
	err := productValidation(feedback.IDProduct, cartID)

	if err != nil {
		return "Producto no valido", err
	}

	client := client()
	err := client.Set(
		feedback.ID,
		feedback.IDProduct,
		feedback.IDUser,
		feedback.moderated,
		feedback.rate,
		feedback.text
		,0).Err()
	if err != nil {
		return "", err
	}

	sendFeedback(feedback);
	return feedback.ID, nil
}

// Find  and return the feedbacks from database
func Find(productID string) (*Feedback, error) {
	client := client()
	data, err := client.Get(productID).Result()
	if err != nil {
		return nil, errors.NotFound
	}

	[]Feedback result := data
	return &result, nil
}

func FindByIDAndUpdate(feedbackID string) (ID, error) {
	client := client()
	data, err := client.Get(feedbackID).Result()
	if err != nil {
		return nil, errors.NotFound
	}
	moderated = true

	err := client.Set(
		data.ID,
		data.IDProduct,
		data.IDUser,
		moderated,
		data.rate,
		data.text
		,0).Err()
		
	if err != nil {
		return "", err
	}
	return data.ID, nil
}

func client() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     env.Get().RedisURL,
		Password: "",
		DB:       0,
	})
}
