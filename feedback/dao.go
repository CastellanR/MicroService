package feedback

import (
	"github.com/CastellanR/UserFeedback-Microservice/tools/env"
	"github.com/CastellanR/UserFeedback-Microservice/tools/errors"
	"github.com/go-redis/redis"
)

// Insert feedback into Database
func Insert(feedback *Feedback) (string, error) {
	if err := feedback.validateSchema(); err != nil {
		return "", err
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

	return feedback.ID, nil
}

// Find  and return the feedbacks from database
func Find(productID string) (*Feedback, error) {
	client := client()
	data, err := client.Get(productID).Result()
	if err != nil {
		return nil, errors.NotFound
	}

	result := Feedback{
		ID:       productID,
		Feedback: data,
	}
	return &result, nil
}

func FindByIDAndUpdate(feedbackID string) (*Feedback, error) {
	client := client()
	data, err := client.Get(feedbackID).Result()
	if err != nil {
		return nil, errors.NotFound
	}

	result := Feedback{
		ID: feedbackID,
	}
	return &result, nil
}

func client() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     env.Get().RedisURL,
		Password: "",
		DB:       0,
	})
}
