package feedback

import (
	"strings"

	"github.com/CastellanR/UserFeedback-Microservice/tools/errors"
	uuid "github.com/satori/go.uuid"
	validator "gopkg.in/go-playground/validator.v9"
)

// Feedback structure
type Feedback struct {
	ID        string `json:id validate:"required"`
	IDUser    string `json:idUser validate:"required"`
	text      string `json:message validate:"required"`
	IDProduct string `json:idProduct validate:"required"`
	rate      int    `json:rate  validate:"required"`
	moderated bool   `json:idUser validate:"required"`
	created   string `json:created validate:"required"`
	updated   string `json:updated validate:"required"`
}

// New create feedback
func New() *Feedback {
	return &Feedback{
		ID:        uuid.NewV4().String(),
		IDUser:    "",
		text:      "",
		IDProduct: "",
		rate:      "",
		moderated: "",
		created:   "",
		updated:   "",
	}
}

// ErrData feedback not valid
var ErrData = errors.NewValidationField("feedback", "invalid")

func (e *Feedback) validateSchema() error {
	validate := validator.New()
	if err := validate.Struct(e); err != nil {
		return err
	}
	if strings.Index(e.Feedback, "data:feedback/") < 0 {
		return ErrData
	}
	return nil
}
