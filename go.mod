module myapp

go 1.18

replace example.com/greetings => ../greetings

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)
