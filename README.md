## Entendendo OpenFGA (benchmarks)

### Requisitos

- **Go** (v1.21+ recomendado)
- **Docker & Docker Compose** (Para executar o servidor OpenFGA)

### Estrutura do Projeto

```
go/
├── cmd/                              # Entry points dos benchmarks
│   ├── benchmark_operators/          # Benchmark de operadores
│   └── benchmark_model_complexity/   # Benchmark de complexidade
├── pkg/
│   ├── benchmarks/                   # Funções de benchmark
│   ├── client/                       # Cliente OpenFGA
│   ├── config/                       # Configurações
│   ├── scenarios/                    # Cenários de teste
│   │   ├── operators/                # Cenários de operadores
│   │   └── model_complexity/         # Cenários de complexidade
│   └── utils/                        # Utilitários (seeder, charts)
```

### Executando os Benchmarks

Você pode usar o `Makefile` fornecido para executar os benchmarks. Isso iniciará automaticamente o servidor OpenFGA no Docker.

**Executar todos os benchmarks**
```bash
make benchmark
```

### Executar benchmarks específicos
- **Benchmark de Operadores:**
  ```bash
  make benchmark-operators
  ```
- **Benchmark de Complexidade do Modelo:**
  ```bash
  make benchmark-model-complexity
  ```

### Outros comandos úteis
- **Build apenas:** `make build`
- **Iniciar OpenFGA:** `make up`
- **Parar OpenFGA:** `make down`
- **Limpar binários e outputs:** `make clean`
