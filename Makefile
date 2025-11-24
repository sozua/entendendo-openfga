.PHONY: up down benchmark benchmark-operators benchmark-model-complexity all

up:
	mkdir output && docker compose up -d || true

down:
	docker compose down || true

benchmark-operators: up
	node src/benchmark_operators.js

benchmark-model-complexity: up
	node src/benchmark_model_complexity.js

benchmark: benchmark-operators benchmark-model-complexity

all: benchmark
