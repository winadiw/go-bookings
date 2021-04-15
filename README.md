# Golang Bookings 
### based on Udemy course

This is the repository for my bookings and reservations project

## To run the app, use run.sh
`sudo chmod +x ./run.sh` (if not yet)
`./run.sh`

## To run the test
cd to directory  
`go test -v `
or from root
`go test -v ./...`

## To test the coverage  
`go test -coverprofile=coverage.out && go tool cover -html=coverage.out`


## Database
Using Soda for migrations [https://gobuffalo.io/en/docs/db/getting-started/]

### migration
make sure to copy `database.yml.example` to `database.yml`

new migration example
`soda g fizz CreateUserTable`  

run migration  
`soda migrate` 
