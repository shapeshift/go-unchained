TEST_PATH ?= ./...
BENCH_TIME ?= 1s
TIMEOUT ?= 10m

.PHONY:  dev-cosmos dev-osmosis dev-thorchain

dev-cosmos:
	docker-compose -f build/cosmos/docker-compose.yml up --build

dev-osmosis:
	docker-compose -f build/osmosis/docker-compose.yml up --build

dev-thorchain:
	docker-compose -f build/thorchain/docker-compose.yml up --build


.PHONY: build test test-unit test-integration test-bench clean godoc generate

build:
	go build ./...

test: test-unit test-integration

test-unit:
	go test -v -cover -count=1 -tags=unit ${TEST_PATH}

test-integration:
	go test -v -cover -count=1 -p=1 -tags=integration ${TEST_PATH}

test-bench:
	go test -timeout=${TIMEOUT} -run=NONE -bench=. -benchtime=${BENCH_TIME} -tags=benchmark ${TEST_PATH}

clean:
	go clean ./...