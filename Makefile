
watch:
	@air -v > /dev/null 2> /dev/null || go install github.com/cosmtrek/air@latest
	@air

run:
	go run server.go

open:
	open http://localhost:8080

test:
	go test ./...

e2e:
	e2e/node_modules/.bin/mocha testGolangHtmx.js --no-timeouts --reporter spec --browser=phantomjs
