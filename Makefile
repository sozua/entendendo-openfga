.PHONY: up down build benchmark-operators benchmark-model-complexity benchmark all clean

up:
	@mkdir -p output bin
	@docker rm -f benchmark-openfga-postgres benchmark-openfga-openfga benchmark-openfga-migrate 2>/dev/null || true
	@docker compose up -d

down:
	docker compose down || true

build:
	cd go && go mod tidy && go build -o ../bin/benchmark_operators ./cmd/benchmark_operators && go build -o ../bin/benchmark_model_complexity ./cmd/benchmark_model_complexity

benchmark-operators: up build
	./bin/benchmark_operators

benchmark-model-complexity: up build
	./bin/benchmark_model_complexity

benchmark: benchmark-operators benchmark-model-complexity

clean:
	rm -rf bin output

all: benchmark
