
watch:
	air -v || go install github.com/cosmtrek/air@latest
	air

run:
	go run server.go

open:
	open http://localhost:8080

test:
	go test