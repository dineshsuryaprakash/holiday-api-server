# holiday-api-server
## Overview
An example api server project that populates json data from a url and host endpoints based on it.
## Installation
1. Clone the repository
```sh
git clone  https://github.com/dineshsuryaprakash/holiday-api-server.git
```
2. Navigate to the directory
```sh
cd holiday-api-server
```
3. Build the code
```sh
make build
```
## Usage
### Start the server app
The apps takes in two command line flags as input
  1) listenAddress -- User can choose the listen address for the app server(default: ":8080").
  2) BankHolidayJsonEndpointURL -- User can give a custom URL to endpoint that hosts the bank holiday json data(default: "https://www.gov.uk/bank-holidays.json"). It is mandatory that the custom URL should have same json templet as "https://www.gov.uk/bank-holidays.json"
```sh
./myapp -listenAddress ":8080" -BankHolidayJsonEndpointURL "https://www.gov.uk/bank-holidays.json"
```
Alternatively users can also use "go run" to start the server app
```sh
go run main.go -listenAddress ":8081" -BankHolidayJsonEndpointURL "https://example.com/new-url"
```

### Endpoints
- '/get_year': This endpoint retrieves holiday events for a specific year. You can specify the year by providing the "year" query parameter in the URL. The endpoint returns a JSON response containing holiday events grouped by divisions for the specified year.
```http
GET http://localhost:8080/get_year?year=2023
```

- '/get_division_bunting':  This endpoint filters holiday events based on the division and bunting status. You can specify the division name and whether the event has bunting using the division and bunting query parameters. The endpoint returns a JSON response with holiday events that match the specified criteria.
```http
GET http://localhost:8080/get_division_bunting?division=england-and-wales&bunting=true
```

- '/get_year_less_event': This endpoint is similar to '/get_year', but it returns holiday events for a specific year with some event details removed. You can specify the year using the year query parameter. The endpoint omits bunting and notes fields from the event details and create the JSON response only with title and date.
```http
GET http://localhost:8080/get_year_less_event?year=2023
```
