BIN := pgnstats pgnstats-wasm

build:
	for target in $(BIN); do \
		go build -o bin/$$target ./cmd/$$target; \
	done

test:
	go run ./cmd/pgnstats/. -f ./cmd/pgnstats/testdata/pgn/carlsen.pgn -v -fp '' -o ./cmd/pgnstats/testdata/carlsen.json -i
	go run ./cmd/pgnstats/. -f ./cmd/pgnstats/testdata/pgn/carlsen.pgn -v -fp 'Carlsen,M' -o ./cmd/pgnstats/testdata/carlsen.json -i