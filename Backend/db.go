package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-central-1"))

// https://github.com/aws/aws-sdk-go/issues/810
// GetAll returns all the items in the table.
// Note: if the table is big (over a few MB of data) only partial results will be returned with this call
// func (s *DynamoSvc) GetAll(refItem interface{}) (interface{}, error) {
//     params := &dynamodb.ScanInput{
//         TableName: s.table,
//     }
//     resp, err := s.db.Scan(params)
//     if err != nil {
//         return nil, fmt.Errorf("failed to get all items from dynamo: %s", err)
//     }

//     refType := reflect.TypeOf(refItem)
//     results := reflect.MakeSlice(reflect.SliceOf(refType), 0, int(*resp.Count))
//     for _, item := range resp.Items {
//         newItem := reflect.Indirect(reflect.New(refType.Elem()))
//         if err = unmarshalDynamoObject(item, newItem.Addr().Interface()); err != nil {
//             return nil, fmt.Errorf("failed to unmarshal the response from dynamo: %s", err)
//         }
//         results = reflect.Append(results, newItem.Addr())
//     }

//     return results.Interface(), nil
// }


func getItem(id string) (*todo, error) {
    input := &dynamodb.GetItemInput{
        TableName: aws.String("todo"),
        Key: map[string]*dynamodb.AttributeValue{
            "id": {
                S: aws.String(id),
            },
        },
    }

    result, err := db.GetItem(input)
    if err != nil {
        return nil, err
    }
    if result.Item == nil {
        return nil, nil
    }

    td := new(todo)
    err = dynamodbattribute.UnmarshalMap(result.Item, td)
    if err != nil {
        return nil, err
    }

    return td, nil
}

// Add a book record to DynamoDB.
func putItem(td *todo) error {
    input := &dynamodb.PutItemInput{
        TableName: aws.String("todo"),
        Item: map[string]*dynamodb.AttributeValue{
            "id": {
                S: aws.String(td.id),
            },
            "Title": {
                S: aws.String(td.Title),
            },
            "Completed": {
               BOOL: aws.Bool(td.Completed),
            },
            "CreatedAt": {
                S: aws.String(td.CreatedAt),
            },
        },
    }

    _, err := db.PutItem(input)
    return err
}
