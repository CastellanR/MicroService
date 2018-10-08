<a name="top"></a>
# UserFeedback  MicroService

UserFeedback Microservice

- [Feedback](#feedback)
	- [Create Feedback](#create-feedback)
	- [Get  Feedbacks](#get-feedbacks)
	- [Modify Feedback](#modify-feedback)
  - [Delete Feedback](#delete-feedback)
	
- [RabbitMQ_POST](#rabbitmq_post)
	- [Product Validation](#product-validation)
  - [Sending Feedback](#sending-feedback)
	
- [RabbitMQ_GET](#rabbitmq_get)
	- [Product Validation](#product-validation)
  - [User Logout](#user-logout)

# <a name='feedback'></a> Feedback

## <a name='create-feedback'></a> Create Feedback
[Back to top](#top)

<p>Add feedback to the server</p>

	POST /v1/feedback



### Examples

Body

```
{
  "idUser" : "{ User Id }",
  "text" :  "{ Feedback Content }",
  "idProduct" : "{ Product Id }",
  "rate" : "{ Feedback Rate }"
}
```
Auth Header

```
Authorization=bearer {token}
```

### Success Response

Response

```
HTTP/1.1 200 OK
{
  "id": "{ Feedback Id }"
}
```
400 Bad Request

```
HTTP/1.1 400 Bad Request
{
   "messages" : [
     {
       "path" : "{Property's Name}",
       "message" : "{ Error cause}"
     },
     ...
  ]
}
```
500 Server Error

```
HTTP/1.1 500 Internal Server Error
{
   "error" : "Not Found"
}
```


### Error Response

401 Unauthorized

```
HTTP/1.1 401 Unauthorized
```
## <a name='get-feedbacks'></a> Get Feedbacks
[Back to top](#top)

<p>Get feedbacks of a product</p>

	GET /v1/feedback/:productid



### Examples

Body

```
{
  "idProduct" : "{ Product Id }",
}
```
Auth Header

```
Authorization=bearer {token}
```

### Success Response

Response

```
{
  {
    "id" : "{ Feedback Id }"
    "idUser" : "{ User Id }",
    "text" :  "{ Feedback Content }",
    "idProduct" : "{ Product Id }",
    "rate" : "{ Feedback Rate }"
  }
}
```
400 Bad Request

```
HTTP/1.1 400 Bad Request
{
   "messages" : [
     {
       "path" : "{Property's Name}",
       "message" : "{ Error cause}"
     },
     ...
  ]
}
```
500 Server Error

```
HTTP/1.1 500 Internal Server Error
{
   "error" : "Not Found"
}
```


### Error Response

401 Unauthorized

```
HTTP/1.1 401 Unauthorized
```
## <a name='modify-feedback'></a> Modify Feedback
[Back to top](#top)

<p>Modify a Feedback</p>

	POST /v1/feedback/:id



### Examples

Body

```
{
  "text" :  "{ Feedback Content }",
  "rate" : "{ Feedback Rate }"
}

```
Auth Header

```
Authorization=bearer {token}
```

### Success Response

Response

```
HTTP/1.1 200 OK
{
  "id": "{ Feedback Id }"
}
```
400 Bad Request

```
HTTP/1.1 400 Bad Request
{
   "messages" : [
     {
       "path" : "{Property's Name}",
       "message" : "{ Error cause}"
     },
     ...
  ]
}
```
500 Server Error

```
HTTP/1.1 500 Internal Server Error
{
   "error" : "Not Found"
}
```


### Error Response

401 Unauthorized

```
HTTP/1.1 401 Unauthorized
```

## <a name='delete-feedback'></a> Delete Feedback
[Back to top](#top)

<p>Delete a Feedback</p>

	DELETE /v1/feedback/:id



### Examples

Body

```
{
  "id": "{ Feedback Id }"
}

```
Auth Header

```
Authorization=bearer {token}
```

### Success Response

Response

```
HTTP/1.1 200 OK
{
  "id": "{ Feedback Id }"
}
```
400 Bad Request

```
HTTP/1.1 400 Bad Request
{
   "messages" : [
     {
       "path" : "{Property's Name}",
       "message" : "{ Error cause}"
     },
     ...
  ]
}
```
500 Server Error

```
HTTP/1.1 500 Internal Server Error
{
   "error" : "Not Found"
}
```


### Error Response

401 Unauthorized

```
HTTP/1.1 401 Unauthorized
```

# <a name='rabbitmq_get'></a> RabbitMQ_POST

## <a name='product-validation'></a> Product Validation
[Back to top](#top)

<p>Sending a validation request for a product</p>

	DIRECT catalog/article-exist


### Example

Message

```
{
  "type": "article-exist",
  "queue": "feedback",
  "exchange": "feedback",
  "message" : {
      "articleId": "{articleId}",
  }
}
```

## <a name='sending-feedback'></a> Sending Feedback
[Back to top](#top)

<p>Send new feedback message</p>

	TOPIC  feedback/new-feedback




### Example

Message

```
{
   "type": "new-feedback",
   "queue": "feedback"
   "message": {
        "id" : "{ Feedback Id }"
        "idUser" : "{ User Id }",
        "text" :  "{ Feedback Content }",
        "idProduct" : "{ Product Id }",
        "rate" : "{ Feedback Rate }"
     }
}
```

# <a name='rabbitmq_get'></a> RabbitMQ_GET

## <a name='product-validation'></a> Product Validation
[Back to top](#top)

<p>Listening Validation messages</p>

	DIRECT feedback/article-exist




### Success Response

Message

```
{
  "type": "article-exist",
  "message" : {
      "articleId": "{articleId}",
      "valid": True|False
  }
}
```

## <a name='users-logout'></a> Users Logout
[Back to top](#top)

<p>Listening Logout messages from auth</p>

	FANOUT auth/logout




### Success Response

Message

```
{
   "type": "logout",
   "message": "{tokenId}"
}
```


