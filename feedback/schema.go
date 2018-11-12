package feedback

import (
	"time"

	"github.com/CastellanR/UserFeedback-Microservice/tools/errors"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	validator "gopkg.in/go-playground/validator.v9"
)

// Feedback structure
type Feedback struct {
	ID        objectid.ObjectID `bson:"_id"`
	UserID    string            `bson:"userId" validate:"required"`
	Text      string            `bson:"text" validate:"required"`
	ProductID string            `bson:"productId" validate:"required"`
	Rate      int               `bson:"rate"  validate:"required"`
	Moderated bool              `bson:"moderated"`
	Created   time.Time         `bson:"created" validate:"required"`
	Updated   time.Time         `bson:"updated" validate:"required"`
}

// New create feedback
func New() *Feedback {
	return &Feedback{
		ID:        objectid.New(),
		UserID:    "",
		Text:      "",
		ProductID: "",
		Rate:      0,
		Moderated: false,
		Created:   time.Now(),
		Updated:   time.Now(),
	}
}

// ErrData feedback not valid
var ErrData = errors.NewValidationField("feedback", "invalid")

func (e *Feedback) validateSchema() error {
	validate := validator.New()
	if err := validate.Struct(e); err != nil {
		return err
	}
	return nil
}
