# Curso GO Expert - Desafio Client-Server-API

<div>
    <img alt="Criado por Alcir Junior [Caju]" src="https://img.shields.io/badge/criado%20por-Alcir Junior [Caju]-%23f08700">
    <img alt="License" src="https://img.shields.io/badge/license-MIT-%23f08700">
</div>

---

## Observação
Como eu uso o Mac M1 estava com um erro de compilação do sqlite3 então usei o `Planet Scale` para não ter de subir um container com o `mysql` assim basta executar o exemplo, mas os tempos dos `contexts` ficaram maiores por se tratar de um banco externo.

___
## Descrição

O Curso GO Expert é uma formação completa para fazer com que pessoas desenvolvedoras sejam capazes de trabalhar em projetos expressivos sendo capazes de desenvolver aplicações de grande porte utilizando de boas práticas de desenvolvimento.

---

## Repositório Pai
https://github.com/alcir-junior-caju/study-go-expert

---

## Visualizar o projeto na IDE:

Para quem quiser visualizar o projeto na IDE clique no teclado a tecla `ponto`, esse recurso do GitHub é bem bacana

---
### Desafio Client-Server-API
Entregar dois sistemas em Go:
- client.go
- server.go

#### Os requisitos para cumprir este desafio são:

O `client.go` deverá realizar uma requisição HTTP no `server.go` solicitando a cotação do dólar.

O `server.go` deverá consumir a API contendo o câmbio de Dólar e Real no endereço: `https://economia.awesomeapi.com.br/json/last/USD-BRL` e em seguida deverá retornar no formato `JSON` o resultado para o cliente.

Usando o package `"context"`, o `server.go` deverá registrar no banco de dados `SQLite` cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de `200ms` e o timeout máximo para conseguir persistir os dados no banco deverá ser de `10ms`.

O `client.go` precisará receber do `server.go` apenas o valor atual do câmbio `(campo "bid" do JSON)`. Utilizando o package `"context"`, o `client.go` terá um timeout máximo de `300ms` para receber o resultado do `server.go`.

Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.

O `client.go` terá que salvar a cotação atual em um arquivo `"cotacao.txt"` no formato: `Dólar: {valor}`

O endpoint necessário gerado pelo `server.go` para este desafio será: `/cotacao` e a porta a ser utilizada pelo servidor HTTP será a `8080`.
