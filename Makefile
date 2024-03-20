.PHONY: test watch run open e2e


test:
	go test ./...

watch:
	@air -v > /dev/null 2> /dev/null || go install github.com/cosmtrek/air@latest
	@air -build.exclude_dir e2e

run:
	go run server.go

open:
	open http://localhost:8080

depgraph:
	godepgraph -s -onlyprefixes github.com/xpmatteo/todomvc-golang \
		github.com/xpmatteo/todomvc-golang \
		| dot -Tpng -o /tmp/godepgraph.png \
		&& open /tmp/godepgraph.png

# does not work
e2e:
	e2e/node_modules/.bin/mocha testGolangHtmx.js --no-timeouts --reporter spec --browser=phantomjs
