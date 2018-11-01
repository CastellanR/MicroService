package rabbit

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"github.com/CastellanR/UserFeedback-Microservice/security"
	"github.com/CastellanR/UserFeedback-Microservice/tools/env"
	"github.com/CastellanR/UserFeedback-Microservice/tools/errors"
	"github.com/streadway/amqp"
)

// ErrChannelNotInitialized Rabbit channel could not be initialized
var ErrChannelNotInitialized = errors.NewCustom(400, "Channel not initialized")

type message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Init se queda escuchando broadcasts de logout
func Init() {
	go func() {
		for {
			listenLogout()
			fmt.Println("RabbitMQ conectando en 5 segundos.")
			time.Sleep(5 * time.Second)
		}
	}()
}

/**
 * @api {direct} cart/article-exist Product Validation
 * @apiGroup RabbitMQ POST
 *
 * @apiDescription Sending a validation request for a product.
 *
 * @apiSuccessExample {json} Message
 *     {
			"type": "article-exist",
			"queue": "cart",
			"exchange": "cart",
			"message" : {
				"referenceId": "{cartId}",
 *             	"articleId": "{articleId}",
			}
		}
 */

func ProductValidation(productID,cartID) error {
	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	chn, err := conn.Channel()
	if err != nil {
		return err
	}
	defer chn.Close()	

	queue, err := chn.QueueDeclare(
		"cart", // name
		true,  // durable
		false,  // delete when unused
		false,   // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return err
	}
	
	msg, err = ch.Consume(
		queue.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
)

	if err != nil {
		return err
	}	

	return msg
}

/**
 * @api {fanout} feedback/ Send Feedback
 * @apiGroup RabbitMQ POST
 *
 * @apiDescription Sending new feedback.
 *
 * @apiSuccessExample {json} Message
 *     {
			"type": "article-exist",
			"queue": "feedback",
			"exchange": "feedback",
			"message" : {
				"articleId": "{articleId}",
			}
		}
 */

 func sendFeedback(feedback) error {
	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	chn, err := conn.Channel()
	if err != nil {
		return err
	}
	defer chn.Close()	

	queue, err := chn.ExchangeDeclare(
		"feedback_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
)
	if err != nil {
		return err
	}
	
	err = ch.Publish(
		"feedback_topic",           // exchange
		"feedback",       // routing key
		false,        // mandatory
		false,		// immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(feedback),
		}
	)

	if err != nil {
		return err
	}	

	return nil
}

/**
 * @api {direct} feedback/article-validation Product Validation
 * @apiGroup RabbitMQ GET
 *
 * @apiDescription Listen validation product messages from cart.
 *
 * @apiSuccessExample {json} Message
 * 		{
 *      	"type": "article-exist",
 *			"message" : 
 *				{
 *					"articleId": "{articleId}",
 *					"valid": True|False
 *				}
 *      }
 */

 func listenProductValidation() error {
	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	chn, err := conn.Channel()
	if err != nil {
		return err
	}
	defer chn.Close()	

	queue, err := chn.QueueDeclare(
		"cart", // name
		true,  // durable
		false,  // delete when unused
		false,   // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return err
	}
	
	err = ch.Publish(
		"",           // exchange
		queue.Name,       // routing key
		false,        // mandatory
		false,
		amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(cartID,productID),
		}
	)

	if err != nil {
		return err
	}	

	return nil
}

/**
 * @api {fanout} auth/logout Logout de Usuarios
 * @apiGroup RabbitMQ GET
 *
 * @apiDescription Escucha de Messages logout desde auth.
 *
 * @apiSuccessExample {json} Message
 *     {
 *        "type": "logout",
 *        "message": "{tokenId}"
 *     }
 */
func listenLogout() error {
	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	chn, err := conn.Channel()
	if err != nil {
		return err
	}
	defer chn.Close()

	err = chn.ExchangeDeclare(
		"auth",   // name
		"fanout", // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	queue, err := chn.QueueDeclare(
		"auth", // name
		false,  // durable
		false,  // delete when unused
		true,   // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return err
	}

	err = chn.QueueBind(
		queue.Name, // queue name
		"",         // routing key
		"auth",     // exchange
		false,
		nil)
	if err != nil {
		return err
	}

	mgs, err := chn.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return err
	}

	fmt.Println("RabbitMQ conectado")

	go func() {
		for d := range mgs {
			log.Output(1, "Mensage recibido")
			newMessage := &message{}
			err = json.Unmarshal(d.Body, newMessage)
			if err == nil {
				if newMessage.Type == "logout" {
					security.Invalidate(newMessage.Message)
				}
			}
		}
	}()

	fmt.Print("Closed connection: ", <-conn.NotifyClose(make(chan *amqp.Error)))

	return nil
}
