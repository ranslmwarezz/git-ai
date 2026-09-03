# Git-AI

<p align="center">
  <strong>Gerador de mensagens de commit com Inteligência Artificial</strong>
</p>

<p align="center">
  Gere mensagens de commit seguindo o padrão Conventional Commits a partir das alterações realizadas no código.
</p>

<p align="center">
  <a href="https://go.dev/">
    <img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  </a>
  <a href="https://ai.google.dev/">
    <img src="https://img.shields.io/badge/Google%20Gemini-API-8E75B2?style=for-the-badge&logo=google" alt="Google Gemini">
  </a>
  <a href="https://github.com/ranslmwarezz/git-ai/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/ranslmwarezz/git-ai?style=for-the-badge" alt="License">
  </a>
  <a href="https://github.com/ranslmwarezz/git-ai">
    <img src="https://img.shields.io/github/last-commit/ranslmwarezz/git-ai?style=for-the-badge" alt="Last Commit">
  </a>
  <!-- TODO: se tiver GitHub Actions rodando `go test`, adicione a badge de build aqui -->
</p>

---

## Sobre o projeto

O **Git-AI** é uma ferramenta de linha de comando desenvolvida em **Go** para gerar mensagens de commit automaticamente a partir das alterações adicionadas ao Git.

O projeto utiliza a **API do Google Gemini** para analisar o `git diff --cached` e gerar uma mensagem seguindo o padrão **Conventional Commits**.

A ideia é tornar o processo de criação de commits mais simples, mantendo mensagens claras, descritivas e padronizadas.

<!--
  TODO: explique aqui o fluxo real — isso é a primeira dúvida de quem for testar:
  - O Git-AI comita automaticamente com a mensagem gerada?
  - Ou ele abre uma interface (Bubble Tea) pra você revisar, editar e aprovar antes do commit?
  Exemplo de frase pra usar:
  "Ao rodar o Git-AI, uma interface interativa no terminal exibe a mensagem sugerida,
  permitindo editar, regenerar ou confirmar antes do commit ser criado."
-->

### Desenvolvido com

* [Go](https://go.dev/)
* [Google Gemini API](https://ai.google.dev/)
* [Git](https://git-scm.com/)
* [Bubble Tea](https://github.com/charmbracelet/bubbletea)
* [Godotenv](https://github.com/joho/godotenv)

---

## Começando

Para executar o projeto localmente, siga os passos abaixo.

### Pré-requisitos

Certifique-se de possuir:

* Go 1.26 ou superior
* Git
* Uma API Key do Google Gemini ([obtenha aqui](https://aistudio.google.com/app/apikey))

### Instalação

**Opção 1 — via `go install` (recomendado)**

```bash
go install github.com/ranslmwarezz/git-ai@latest
```

Isso instala o binário `git-ai` em `$GOPATH/bin` (ou `$HOME/go/bin`). Certifique-se de que esse diretório está no seu `PATH`.

**Opção 2 — a partir do código-fonte**

```bash
git clone https://github.com/ranslmwarezz/git-ai.git
cd git-ai
go run ./cmd/git-ai
```

> `go run` já resolve as dependências automaticamente — não é necessário rodar `go mod download` antes.

### Configuração

Crie um arquivo `.env` na raiz do projeto (ou exporte a variável no seu shell):

```env
GEMINI_API_KEY=sua_api_key_aqui
```

<!-- TODO: se houver outras variáveis configuráveis (modelo do Gemini, idioma da mensagem, etc), documente aqui. Exemplo:
GEMINI_MODEL=gemini-2.0-flash
COMMIT_LANGUAGE=pt-BR
-->

---

## Utilização

1. Adicione as alterações que deseja incluir no commit:

```bash
git add .
```

2. Execute o Git-AI:

```bash
git-ai
```

<!-- Se ainda não tiver o binário instalado via `go install`, use: go run ./cmd/git-ai -->

3. O Git-AI analisa as alterações em staging, envia o diff para o Google Gemini e exibe a mensagem sugerida no terminal.

<!--
  TODO: descreva aqui o que acontece a partir daqui. Por exemplo:
  - Teclas de atalho pra aceitar/editar/regenerar
  - Se o commit é criado automaticamente ao aceitar
  - Um GIF ou screenshot da interface ajuda MUITO mais que texto nessa parte
-->

Exemplo de mensagem gerada:

```text
feat: adiciona integração com API do Gemini
```

---

## Conventional Commits

As mensagens geradas pelo Git-AI seguem o padrão **Conventional Commits**.

Alguns exemplos:

```text
feat: adiciona autenticação de usuários

fix: corrige conexão com banco de dados

test: adiciona testes para o cliente Gemini

refactor: simplifica tratamento de erros

docs: atualiza documentação do projeto
```

---

## Testes

Para executar os testes automatizados:

```bash
go test ./...
```

Para executar os testes exibindo os detalhes:

```bash
go test ./... -v
```

---

## Solução de problemas


| Problema                        | Possível causa                                                 | Solução                                                                                            |
| ------------------------------- | -------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `GEMINI_API_KEY não encontrada` | Arquivo `.env` ausente ou variável de ambiente não configurada | Verifique se o `.env` está na raiz do projeto ou configure a variável `GEMINI_API_KEY` manualmente |
| Erro de autenticação na API     | API Key inválida ou sem permissão                              | Gere ou configure uma nova chave no [Google AI Studio](https://aistudio.google.com/app/apikey)     |
| Nenhuma alteração em staging  | Nenhum arquivo foi adicionado ao staging                       | Execute `git add .` antes de executar o Git-AI                                                     |



<!-- TODO: ajuste essa tabela conforme os erros reais que o Git-AI trata/loga -->

---

## Licença

Distribuído sob a licença MIT. Consulte o arquivo `LICENSE` para mais informações.
