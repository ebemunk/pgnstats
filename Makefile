BIN := pgnstats pgnstats-wasm

build:
	for target in $(BIN); do \
		go build -o bin/$$target ./cmd/$$target; \
	done