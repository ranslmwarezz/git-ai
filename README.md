# Git-AI

<p align="center">
  <strong>Gerador de mensagens de commit com Inteligência Artificial</strong>
</p>

<p align="center">
  Gere, revise e execute commits seguindo o padrão Conventional Commits diretamente pelo terminal.
</p>

<p align="center">
  <a href="https://go.dev/">
    <img src="https://img.shields.io/badge/Go-1.26.1-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  </a>
  <a href="https://ai.google.dev/">
    <img src="https://img.shields.io/badge/Google%20Gemini-API-8E75B2?style=for-the-badge&logo=google" alt="Google Gemini">
  </a>
  <a href="https://github.com/charmbracelet/bubbletea">
    <img src="https://img.shields.io/badge/Bubble%20Tea-TUI-00ADD8?style=for-the-badge" alt="Bubble Tea">
  </a>
  <a href="https://github.com/ranslmwarezz/git-ai">
    <img src="https://img.shields.io/github/last-commit/ranslmwarezz/git-ai?style=for-the-badge" alt="Last Commit">
  </a>
</p>

---

## Sobre o projeto

O **Git-AI** é uma ferramenta de linha de comando desenvolvida em **Go** para auxiliar na criação de commits seguindo o padrão **Conventional Commits**.

A ferramenta analisa as alterações adicionadas ao staging através do `git diff --cached`, envia o diff para a **API do Google Gemini** e gera uma mensagem de commit em português.

Depois de gerar a mensagem, uma TUI interativa permite que o usuário revise o resultado, edite a mensagem, realize o commit e, opcionalmente, faça o push para o repositório remoto.

A proposta é combinar **Inteligência Artificial + Git + uma interface interativa no terminal** para tornar a criação de commits mais simples, rápida e padronizada.

### Desenvolvido com

- [Go](https://go.dev/)
- [Google Gemini API](https://ai.google.dev/)
- [Git](https://git-scm.com/)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- [Godotenv](https://github.com/joho/godotenv)

---

## Como funciona

O fluxo principal do Git-AI é:

```text
Alterações no projeto
  │
  ▼
    git add
  │
  ▼
    git-ai
  │
  ▼
 git diff --cached
  │
  ▼
  Google Gemini
  │
  ▼
Mensagem sugerida
  │
  ▼
┌───────────────────────┐
│ Realizar commit       │
│ Editar mensagem       │
│ Cancelar              │
└───────────────────────┘
  │
  ▼
   git commit
  │
  ▼
┌───────────────────────┐
│ Fazer push            │
│ Sair                  │
└───────────────────────┘
  │
  ▼
     git push
```

O commit e o push nunca são executados automaticamente. O usuário precisa confirmar cada operação.

---

## Começando

Para executar o projeto localmente, siga os passos abaixo.

### Pré-requisitos

Certifique-se de possuir:

- Go 1.26.1 ou superior
- Git
- Uma API Key do Google Gemini ([obtenha aqui](https://aistudio.google.com/app/apikey))

### Instalação

**Opção 1 — via `go install`**

```bash
go install github.com/ranslmwarezz/git-ai/cmd/git-ai@latest
```

Isso instala o binário `git-ai` no diretório de binários do Go. Certifique-se de que esse diretório esteja no seu `PATH`.

**Opção 2 — a partir do código-fonte**

```bash
git clone https://github.com/ranslmwarezz/git-ai.git
cd git-ai
go run ./cmd/git-ai
```

> `go run` resolve as dependências automaticamente — não é necessário executar `go mod download` antes.

### Configuração

Crie um arquivo `.env` na raiz do projeto ou configure a variável de ambiente diretamente:

```env
GEMINI_API_KEY=sua_api_key_aqui
```

Para configurar a API Key globalmente, independente do repositório atual:

```bash
git-ai config
```

A configuração global é salva em `~/.config/git-ai/config` com permissões restritas ao usuário.

Ao procurar a API Key, o Git-AI utiliza esta ordem de prioridade:

1. Variável de ambiente `GEMINI_API_KEY`.
2. Arquivo `.env` do diretório atual.
3. Configuração global criada por `git-ai config`.

A chave é utilizada para autenticar as requisições realizadas à API do Google Gemini.

---

## Utilização

### 1. Adicione as alterações ao staging

```bash
git add .
```

Você também pode adicionar arquivos específicos:

```bash
git add arquivo.go
```

### 2. Execute o Git-AI

Se o binário estiver instalado:

```bash
git-ai
```

Ou, executando diretamente a partir do código-fonte:

```bash
go run ./cmd/git-ai
```

### 3. Revise a mensagem sugerida

O Git-AI apresenta uma interface interativa no terminal com a mensagem gerada pela IA.

Você pode escolher entre:

```text
❯ Realizar commit
  Editar mensagem
  Cancelar
```

### 4. Edite a mensagem, se necessário

A opção **Editar mensagem** permite alterar a sugestão gerada pelo Gemini antes de criar o commit.

Pressione:

```text
Enter → salvar a edição
Esc   → cancelar a edição
```

### 5. Realize o commit

Ao escolher **Realizar commit**, o Git-AI executa:

```bash
git commit -m "mensagem"
```

A operação só acontece após a confirmação do usuário.

### 6. Faça push, se desejar

Após um commit realizado com sucesso, o Git-AI oferece a opção de fazer push.

O push também exige uma confirmação explícita antes de executar:

```bash
git push
```

O Git-AI não executa `git push` automaticamente.

---

## Conventional Commits

As mensagens geradas pelo Git-AI seguem o padrão **Conventional Commits**.

Alguns tipos comuns:

```text
feat: adiciona autenticação de usuários

fix: corrige conexão com banco de dados

test: adiciona testes para o cliente Gemini

refactor: simplifica tratamento de erros

docs: atualiza documentação do projeto

chore: atualiza dependências
```

A mensagem gerada é curta, objetiva e escrita em português.

---

## Interface

A interface do Git-AI utiliza **Bubble Tea** e **Lip Gloss** para fornecer uma experiência interativa diretamente no terminal.

Principais atalhos:

| Tecla     | Ação                   |
| --------- | ---------------------- |
| `↑` / `↓` | Navegar entre opções   |
| `Enter`   | Selecionar / confirmar |
| `Esc`     | Cancelar / voltar      |
| `q`       | Sair                   |

Durante a edição da mensagem, `Enter` salva a alteração e `Esc` cancela a edição.

---

## Testes

Para executar todos os testes automatizados:

```bash
go test ./...
```

Para executar os testes com detalhes:

```bash
go test ./... -v
```

O projeto possui testes para os principais componentes, incluindo cliente Git, serviço de commit, cliente da API Gemini e TUI.

---

## Solução de problemas

| Problema                         | Possível causa                                              | Solução                                                                                        |
| -------------------------------- | ----------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| `GEMINI_API_KEY não configurada` | Variável de ambiente, `.env` e configuração global ausentes | Execute `git-ai config` ou configure `GEMINI_API_KEY` manualmente                              |
| Erro de autenticação na API      | API Key inválida ou sem permissão                           | Gere ou configure uma nova chave no [Google AI Studio](https://aistudio.google.com/app/apikey) |
| Nenhuma alteração em staging     | Nenhum arquivo foi adicionado ao staging                    | Execute `git add .` antes de executar o Git-AI                                                 |

<!-- TODO: ajuste essa tabela conforme os erros reais que o Git-AI trata/loga -->

---

## Licença

Este projeto é distribuído sob a licença **MIT**. Consulte o arquivo [`LICENSE`](LICENSE) para mais informações.
