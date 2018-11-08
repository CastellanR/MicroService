package feedback

import (
	"context"
	"log"

	"github.com/CastellanR/UserFeedback-Microservice/tools/db"
	"github.com/CastellanR/UserFeedback-Microservice/tools/errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// New dao es interno a este modulo, nadie fuera del modulo tiene acceso
func getDao() (db.Collection, error) {
	database, err := db.Get()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("feedback")

	_, err = collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.NewDocument(
				bson.EC.String("idProduct", ""),
			),
			Options: bson.NewDocument(
				bson.EC.Boolean("unique", true),
			),
		},
	)
	if err != nil {
		log.Output(1, err.Error())
	}

	coll := db.WrapCollection(collection)
	return coll, nil
}

// Insert into Database
func Insert(feedback *Feedback, cartID string) (string, error) {

	if err := feedback.validateSchema(); err != nil {
		return "", err
	}

	if _, err := getDao().InsertOne(context.Background(), feed); err != nil {
		return nil, err
	}

	return feed.ID, nil
}

// Find  and return the feedbacks from database
func Find(productID string) ([]*Feedback, error) {

	filter := bson.NewDocument(bson.EC.String("productID", productID))
	cur, err := getDao().Find(context.Background(), filter, nil)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	feedbacks := []*Feedback{}
	for cur.Next(context.Background()) {
		feedback := &Feedback{}
		if err := cur.Decode(feedback); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, feedback)
	}

	return feedbacks, nil
}

// FindByIDAndUpdate  and update a feedback from database
func FindByIDAndUpdate(feedbackID string) (string, error) {

	_id, err := objectid.FromHex(feedbackID)

	if err != nil {
		return nil, errors.ErrID
	}

	filter := bson.NewDocument(bson.EC.ObjectID("_id", _id))
	change := bson.NewDocument(bson.EC.String("moderated", true))

	feedback, err := getDao().FindOneAndUpdate(context.Background(), filter, change)

	if err != nil {
		return "", err
	}

	return feedbackID, nil
}
