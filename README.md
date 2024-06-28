# Real-Time Message Classification (RTMC)

## Overview

This project is a real-time monitoring and classification system coded in Go, using gRPC services for communication between the client and server.

## Features

- User authentication and registration using JWT
- gRPC services for authentication and registration (`/internal/user`)
- CRUD operations for rules via gRPC services (`/internal/rule`)
- Parsing and storing user-defined rules
- Real-time data stream processing with gRPC (`StreamData`)
- Unit tests for the rules module
- Integration tests for end-to-end validation
- Load testing using k6

## Client Communication

The client communication code is located in `/client/client.go`. The client establishes a gRPC connection, adds the JWT token in metadata, and sends requests to the gRPC server.

User related Module: `/internal/user`
Rule related Module: `/internal/rule`
Authentication:  `/internal/middleware`

## Authentication and Registration

The user are registered and stored in MongoDB database and are authenticated by assigning a JWT token to the user.

## DSL Format

- `NOTIFY <notification message> WHEN <condition>`
- Conditions can include operations: `Count`, `Length`, `Contains`, `MIN`, `MAX`, `AND`, `OR`, `NOT`
- Operators allowed: `>`, `>=`, `=`, `!=`, `<=`, `<`
- Example condition: `<operation> <operator> <positive integer>`
- Conditions can go upto any depth eg:- MIN (MAX (---- , ----) , MIN (--- , ---) ) >= 6 

## Rule Struct
<condition> string is parsed and stored as
```go
type Rule struct {
    ID       primitive.ObjectID 
    Name     string          // operation
    Matcher  string          // Count or Cantains operation argument
    Ineq     string	     // operator
    Limit    int             // positive integer
    Obj      []primitive.ObjectID        // list of rule id incase of compound rules eg: MIN(operation1, operation2) then IDs of operation1,operation2 are here
    Notify   string          // notification message
    When     string          // storing whole condition as string
    Owners   []primitive.ObjectID  // list of users created this rule
} (parsing is done in /internal/rule/parse.go)
```
The parsing logic is located in /internal/rule/parse.go. This code handles the transformation of DSL rules into the Rule struct.

## Classification

The classification logic is implemented in `/internal/rule/service.go`. When a classification request is made, the server extracts the user ID from the token, retrieves all rules defined by the user, and evaluates the rules against the input string. The server then returns a list of notifications for the rules that are satisfied.

## Stream Data

The `StreamData` gRPC service handles streaming data from the client. The client sends a continuous stream of data, and the server processes each data point in real-time, sending back notifications as necessary.

## Unit Tests

Unit tests are implemented rules module. Database operations are mocked using a mock repository with predefined outputs to ensure the smallest code units are working correctly.

## Integration Tests 

Code logic is implemented in `./integration_test.go`. Integration tests are performed using a duplicate database to store users and rules. These tests ensure that different services work together as expected.

## Load Testing

Load testing is performed using k6. The load test script for the HTTP endpoint (`/classify`) is located in `loadtest_http.js`.
Results: 
![Logo](https://github.com/Prapul1614/RTMC/blob/main/Loadtest%20results/vu_1.png)(For one user sending requests in one sec)
![](https://github.com/Prapul1614/RTMC/blob/main/Loadtest%20results/vu_200.png)(For 200 virtual users send requests in one sec)

## Running the Project

1. Clone the repository:

    ```sh
    git clone https://github.com/Prapul1614/RTMC.git
    cd RTMC
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

3. Compile Protocol Buffers:

    ```sh
    protoc --go_out=. --go-grpc_out=. proto/*.proto
    ```

4. Run the server:

    ```sh
    go run cmd/api/main.go
    ```

5. Run the client:

    ```sh
    go run client/client.go
    ```

