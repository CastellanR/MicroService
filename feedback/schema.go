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
	IDUser    string            `bson:"idUser" validate:"required"`
	text      string            `bson:"text" validate:"required"`
	IDProduct string            `bson:"idProduct" validate:"required"`
	rate      int               `bson:"rate"  validate:"required"`
	moderated bool              `bson:"moderated" validate:"required"`
	created   time.Time         `bson:"created" validate:"required"`
	updated   time.Time         `bson:"updated" validate:"required"`
}

// New create feedback
func New() *Feedback {
	return &Feedback{
		ID:        objectid.New(),
		IDUser:    "",
		text:      "",
		IDProduct: "",
		rate:      0,
		moderated: false,
		created:   time.Now(),
		updated:   time.Now(),
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
