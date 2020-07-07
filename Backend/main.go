package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "context"
    "time"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    
//    "github.com/gofrs/uuid"
)

var ddb = dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-central-1"))

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

type todo struct {
		id        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt string    `json:"created_at"`
	}

type ListTodosResponse struct {
	Todos		[]todo  `json:"todos"`
}

func ListTodos(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("ListTodos")

	var (
		tableName = aws.String("todo")
	)

	// Read from DynamoDB
	input := &dynamodb.ScanInput{
		TableName: tableName,
	}
	result, _ := ddb.Scan(input)

	// Construct todos from response
	var todoList []todo
	for _, i := range result.Items {
		todo := todo{}
		if err := dynamodbattribute.UnmarshalMap(i, &todo); err != nil {
			fmt.Println("Failed to unmarshal")
			fmt.Println(err)
		}
		todoList = append(todoList, todo)
	}

	// Success HTTP response
	body, _ := json.Marshal(&ListTodosResponse{
		Todos: todoList,
	})
	return events.APIGatewayProxyResponse{
		Body: string(body),
		StatusCode: 200,
	}, nil
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    switch req.HTTPMethod {
    case "GET":
        return show(req)
    case "POST":
        return AddTodo(req)
    default:
        return clientError(http.StatusMethodNotAllowed)
    }
}


func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Get the `id` query string parameter from the request and
    // validate it.
    id := req.QueryStringParameters["id"]


    // Fetch the book record from the database based on the isbn value.
    todo, err := getItem(id)
    if err != nil {
        return serverError(err)
    }
    if todo == nil {
        return clientError(http.StatusNotFound)
    }

    // The APIGatewayProxyResponse.Body field needs to be a string, so
    // we marshal the book record into JSON.
    js, err := json.Marshal(todo)
    if err != nil {
        return serverError(err)
    }

    // Return a response with a 200 OK status and the JSON book record
    // as the body.
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { 
    if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
        return clientError(http.StatusNotAcceptable)
    }

    td := new(todo)
//    err := json.Unmarshal([]byte(req.Body), td)
//    if err != nil {
//        return clientError(http.StatusUnprocessableEntity)
//    }
    tableName := "todo"
    av, err := dynamodbattribute.MarshalMap(td)
    if err != nil {
        fmt.Println("Got error marshalling new movie item:")
        fmt.Println(err.Error())
        os.Exit(1)
    }

    input := &dynamodb.PutItemInput{
        Item:      av,
        TableName: aws.String(tableName),
    }

    _, err = ddb.PutItem(input)
    if err != nil {
        fmt.Println("Got error calling PutItem:")
        fmt.Println(err.Error())
        os.Exit(1)
    }


//    if !isbnRegexp.MatchString(bk.ISBN) {
//       return clientError(http.StatusBadRequest)
//    }
//    if td.Title == "" || td.CreatedAt == "" {
//        return clientError(http.StatusBadRequest)
//    }

//    err = putItem(td)
//    if err != nil {
//        return serverError(err)
//    }

    return events.APIGatewayProxyResponse{
        StatusCode: 201,
        Headers:    map[string]string{"Location": fmt.Sprintf("/todo?id=%s", td.id)},
    }, nil
}

func AddTodo(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    fmt.Println("AddTodo")

    var (
 //       uuid = uuid.Must(uuid.NewV4()).String()
        tableName = aws.String("todo")
    )


    // Initialize todo
    otodo := &todo{
 //       id:                 uuid,
        Completed:          false,
        CreatedAt:          time.Now().String(),
    }

    // Parse request body
    json.Unmarshal([]byte(request.Body), otodo)

    // Write to DynamoDB
    item, _ := dynamodbattribute.MarshalMap(otodo)
    input := &dynamodb.PutItemInput{
        Item: item,
        TableName: tableName,
    }
    if _, err := ddb.PutItem(input); err != nil {
        return events.APIGatewayProxyResponse{ // Error HTTP response
            Body: err.Error(),
            StatusCode: 500,
        }, nil
    } else {
        body, _ := json.Marshal(otodo)
        return events.APIGatewayProxyResponse{ // Success HTTP response
            Body: string(body),
            StatusCode: 200,
        }, nil
    }
}


// https://github.com/thedevsaddam/todoapp/blob/master/main.go
// func fetchTodos(w http.ResponseWriter, r *http.Request) {
// 	todos := []todoModel{}

// 	if err := db.C(collectionName).
// 		Find(bson.M{}).
// 		All(&todos); err != nil {
// 		rnd.JSON(w, http.StatusProcessing, renderer.M{
// 			"message": "Failed to fetch todo",
// 			"error":   err,
// 		})
// 		return
// 	}

// 	todoList := []todo{}
// 	for _, t := range todos {s
// 		todoList = append(todoList, todo{
// 			ID:        t.ID.Hex(),
// 			Title:     t.Title,
// 			Completed: t.Completed,
// 			CreatedAt: t.CreatedAt,
// 		})
// 	}

// 	rnd.JSON(w, http.StatusOK, renderer.M{
// 		"data": todoList,
// 	})
// }

// Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
    errorLogger.Println(err.Error())

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body:       http.StatusText(http.StatusInternalServerError),
    }, nil
}

// Similarly add a helper for send responses relating to client errors.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: status,
        Body:       http.StatusText(status),
    }, nil
}

func main() {
    lambda.Start(router)
}


