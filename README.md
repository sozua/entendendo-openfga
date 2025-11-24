## Entendendo OpenFGA (benchmarks)

### Requisitos

- **Node.js** (v20+ recomendado)
- **pnpm** (Gerenciador de pacotes)
- **Docker & Docker Compose** (Para executar o servidor OpenFGA)

### Instalação

1. Instale as dependências:
   ```bash
   pnpm install
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

**Gerenciar a infraestrutura manualmente**
- Iniciar OpenFGA: `make up`
- Parar OpenFGA: `make down`
