package routes

import (
	"github.com/CastellanR/UserFeedback-Microservice/feedback"
	"github.com/CastellanR/UserFeedback-Microservice/tools/errors"
	"github.com/gin-gonic/gin"
)

//NewFeedbackRequest structure
type NewFeedbackRequest struct {
	IDUser    string `json:"idUser" binding:"required"`
	text      string `json:"text" binding:"required"`
	IDProduct string `json:"idproduct" binding:"required"`
	rate      int    `json:"rate" binding:"required"`
	cartID    string
}

//GetFeedbacksRequest structure
type GetFeedbacksRequest struct {
	productID string `json:"idproduct" binding:"required"`
}

//ModerateFeedbackRequest structure
type ModerateFeedbackRequest struct {
	feedbackID string `json:"_id" binding:"required"`
}

// NewFeedback Create feedback
/**
 * @api {post} /v1/feedback Create Feedback
 * @apiName Create Feedback
 * @apiGroup Feedback
 *
 * @apiDescription Create a Feedback from a user.
 *
 * @apiExample {json} Body
 *    {
		"idUser" : "{ User Id }",
		"text" :  "{ Feedback Content }",
		"idProduct" : "{ Product Id }",
		"rate" : "{ Feedback Rate }",
		}

 *
 * @apiSuccessExample {json} Response
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "{ Feedback Id }"
 *     }
 *
 * @apiUse AuthHeader
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
*/
func NewFeedback(c *gin.Context) {

	body := NewFeedbackRequest{}

	if err := validateAuthentication(c); err != nil {
		errors.Handle(c, err)
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		errors.Handle(c, err)
		return
	}

	fdbk := feedback.New()
	fdbk.IDProduct = body.IDProduct
	fdbk.IDUser = body.IDUser
	//fdbk.text = body.text
	//fdbk.rate = body.rate

	dao, err := feedback.GetDao()
	if err != nil {
		errors.Handle(c, err)
		return
	}
	id, err := dao.Insert(fdbk, body.cartID)

	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(200, gin.H{
		"id": id,
	})
}

// GetFeedbacks Get feedback list of a product
/**
 * @api {get} /v1/feedback/:productID Get Feedbacks
 * @apiName Get Feedbacks
 * @apiGroup Feedback
 *
 * @apiDescription Get feedback list of a product
 *
 * @apiExample {json} Body
 *    {
		"idProduct" : "{ Product Id }",
		}
*
 * @apiSuccessExample {json} Response
* {
	{
	"id" : "{ Feedback Id }"
	"idUser" : "{ User Id }",
	"text" :  "{ Feedback Content }",
	"idProduct" : "{ Product Id }",
	"rate" : "{ Feedback Rate }",
	"moderated" : "{ Feedback Moderate Status Boolean }"
	"created" : "{ Creation Date }",
	"updated" : "{ Modification Date }",
	}
}

 *
 * @apiUse AuthHeader
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
*/
func GetFeedbacks(c *gin.Context) {
	if err := validateAuthentication(c); err != nil {
		errors.Handle(c, err)
		return
	}

	body := GetFeedbacksRequest{}

	if err := c.ShouldBindJSON(&body); err != nil {
		errors.Handle(c, err)
		return
	}

	productID := c.Param("productID")

	var data []*feedback.Feedback
	var err error

	dao, err := feedback.GetDao()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	data, err = dao.Find(productID)

	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(200, data)
}

// ModerateFeedback Moderate feedback
/**
 * @api {post} /v1/feedback/:id Moderate Feedback
 * @apiName Moderate Feedback
 * @apiGroup Feedback
 *
 * @apiDescription Moderate a  Feedback from a user.
 *
 * @apiExample {json} Body
 *    {
		"id": "{ Feedback Id }",
		}

 *
 * @apiSuccessExample {json} Response
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "{ Feedback Id }"
 *     }
 *
 * @apiUse AuthHeader
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
*/
func ModerateFeedback(c *gin.Context) {
	if err := validateAuthentication(c); err != nil {
		errors.Handle(c, err)
		return
	}

	body := ModerateFeedbackRequest{}

	if err := c.ShouldBindJSON(&body); err != nil {
		errors.Handle(c, err)
		return
	}

	feedbackID := c.Param("feedbackID")

	var err error

	dao, err := feedback.GetDao()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	id, err := dao.FindByIDAndUpdate(feedbackID)

	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(200, gin.H{
		"id": id,
	})
}
