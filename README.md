# Test Task

### Implementation of HTTP API Server
1) Method to retrieve all records you have collected.
2) Method to retrieve records for a selected day.

## Sequence of Actions

1) Create the database (MySQL) `exchange_rates`
2) Perform migration  
`go install github.com/pressly/goose/v3/cmd/goose@latest`  
`goose -dir sql/schema mysql "mike:mike@tcp(127.0.0.1:3306)/exchange_rates" up`
3) Set the environment variable  
   `CONFIG_PATH=config/local.yaml`
4) Start the server
5) Populate the tables with data using PHP scripts 
   https://github.com/kostevmf/exchange-rates-collector

## Endpoints
1) Method to retrieve all records you have collected  
`GET /exrates/all`
2) Method to retrieve records for a selected day  
`GET /exrates/ondate/{date}`  
`{date}` - required parameter, for example _2025-01-15_
