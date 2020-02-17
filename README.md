# Challenge-Stone

API criada para fornecer a possibilidade de criação de contas e de transferência entre elas

### Instalação

Para instalar e deixar a API rodando é muito simples :) Você so precisa dedois comandos...

```
git clone https://github.com/luannevesb/challenge-stone-accounts.git
cd challenge-stone-accounts
go run cmd/main.go
```

Pronto,
Sua API estará rodando em http://localhost:8000/

## Rodando os testes

Essa API está equipada com mais de 80% de cobertura de testes unitários, para rodar eles é so digitar o seguinte comando

```
make test-coverage
```

Depois disso será criado um arquivo chamado "coverage.html", nele você pode ver a cobertura dos testes.

## Rotas presentes na API

No readme do desafio estava pedindo uma rota chamada ***"balance"*** com um "L" só, mas para manter a coerência com o nome do atributo de account, a rota ficou com o nome de ***ballance*** com os dois "L"

| ROTAS                   | VERBOS | PARÂMETROS                                                                                                      | OBJETIVO                                                         |   |
|-------------------------|--------|-----------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------|---|
| /accounts               | POST   | { "ballance": - Obrigatório - float "cpf": - Obrigatório - string - CPF válido "name": - Obrigatório - string } | Rota usada para criação de uma nova conta                        |   |
| /accounts               | GET    | -                                                                                                               | Rota usada para buscar as informações de todas as accounts       |   |
| /accounts/{id}          | GET    | -                                                                                                               | Rota usada para buscar as informações de uma account             |   |
| /accounts/{id}/ballance | GET    | -                                                                                                               | Rota usada para buscar as informações de ballance de uma account |   |


| ROTAS      | VERBOS | PARÂMETROS                                                                                                                      | OBJETIVO                                                    |
|------------|--------|---------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------|
| /transfers | GET    | -                                                                                                                               | Rota usada para buscar as informações de todas as transfers |
| /transfers | POST   | {	 "account_destination_id":- Obrigatório - string	 "account_origin_id": - Obrigatório - string "amount": - Obrigatório - float } | Rota usada para criação de uma nova transfer                |

## Tecnologias usadas

* [GO](https://golang.org) - A linguagem usada
* [Scribble](https://github.com/nanobox-io/golang-scribble) - DB JSON
* [MUX](github.com/gorilla/mux) - Usado para auxilixar na construção da API
* [Govalidator](github.com/thedevsaddam/govalidator) - Usado para auxilixar na validação das Requests
* [uuid](github.com/google/uuid) - Usado para gerar os ID'S das transfers

## Authors

* **Luan Neves** - *Trabalho completo* - [luannevesb](https://github.com/luannevesb)
