
---
# Background 
------

This is a Golang based todo application - It can be deployed on AWS as a lambda fuction and exposed through API gateway, the Database used in this example is DynamoDB 

# Backend 
------


## Install required packages

```
go get github.com/go-chi/chi
go get gopkg.in/mgo.v2
go get github.com/thedevsaddam/renderer
```

## Local test with DynamoDB

https://medium.com/@vschroeder/install-a-local-dynamodb-development-database-on-your-machine-82dc38d59503

```
aws dynamodb list-tables --endpoint-url http://localhost:8000
```


## Create table

```
aws dynamodb create-table --table-name todo --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5  --endpoint-url http://localhost:8000

aws2 dynamodb put-item --table-name todo --item '{"id": {"S": "978-0486298238"}, "Title": {"S": "finish dynamodb with go"},  "Completed":  {"BOOL": false}, "CreatedAt":  {"S": "05/02/2019"}}' --endpoint-url http://localhost:8000

# Test with Api-Gateway
"id"="123124234"&"title"="hello"&"completed"="False"&"created_at"="12/12/2020"

# GO Lambda function packaging
env GOOS=linux GOARCH=amd64 go build -o /tmp/main todo
zip -j /tmp/main.zip /tmp/<main></main>
cp /tmp/main.zip /home/equa/Cloud/aws-terraform-application/infrastucture-aws-passpes/lambda/files/go/booksapi.zip
```


# Frontend
------

## Project setup
```
yarn install
```

## Compiles and hot-reloads for development
```
yarn serve
```

## Compiles and minifies for production
```
yarn build
```

## Lints and fixes files
```
yarn lint
```

## Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).
