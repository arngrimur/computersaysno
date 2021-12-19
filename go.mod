module computersaysno

go 1.17

require RESTendpoints v0.0.0

require (
	datastructures v0.0.0 // indirect
	db v0.0.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
)

replace RESTendpoints => ./RESTendpoints

replace datastructures => ./datastructures

replace db => ./db
