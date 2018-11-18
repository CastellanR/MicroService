package feedback

import (
	"context"
	"encoding/json"
	"log"

	"github.com/CastellanR/UserFeedback-Microservice/rabbit"
	"github.com/CastellanR/UserFeedback-Microservice/tools/db"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type daoStruct struct {
	collection db.Collection
}

//Dao s
type Dao interface {
	Insert(feedback *Feedback) error
	FindByIDAndUpdate(feedbackID string) error
	Find(productID string) ([]*Feedback, error)
}

//GetDao sda
func GetDao() (Dao, error) {
	database, err := db.Get()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("feedback")

	_, err = collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.NewDocument(
				bson.EC.String("productId", ""),
			),
			Options: bson.NewDocument(
				bson.EC.Boolean("unique", false),
			),
		},
	)
	if err != nil {
		log.Output(1, err.Error())
	}

	coll := db.WrapCollection(collection)
	return daoStruct{
		collection: coll,
	}, nil
}

// Insert into Database
func (d daoStruct) Insert(feedback *Feedback) error {

	if err := feedback.validateSchema(); err != nil {
		return err
	}

	if err := rabbit.ProductValidation(feedback.ProductID, feedback.ID); err != nil {
		return err
	}

	if _, err := d.collection.InsertOne(context.Background(), feedback); err != nil {
		return err
	}

	feed, err := json.Marshal(feedback)
	if err != nil {
		return err
	}

	rabbit.SendFeedback(string(feed[:]))
	return nil
}

// Find  and return the feedbacks from database
func (d daoStruct) Find(productID string) ([]*Feedback, error) {

	filter := bson.NewDocument(bson.EC.String("productId", productID), bson.EC.Boolean("moderated", false))
	cur, err := d.collection.Find(context.Background(), filter, nil)
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
func (d daoStruct) FindByIDAndUpdate(feedbackID string) error {

	_id, err := objectid.FromHex(feedbackID)

	if err != nil {
		return err
	}

	_, err = d.collection.UpdateOne(context.Background(),
		bson.NewDocument(bson.EC.ObjectID("_id", _id)),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.Boolean("moderated", true),
			),
		))

	if err != nil {
		return err
	}

	return nil
}
