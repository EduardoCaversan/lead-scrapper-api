# lead-scrapper-api

## Visão Geral do Projeto

O `lead-scrapper-api` é um projeto desenvolvido em Go (Golang) que funciona como uma API para realizar web scraping. Seu principal objetivo é buscar e extrair informações de contato (como e-mails, telefones e links de mídias sociais) de páginas da web com base em palavras-chave fornecidas. Este projeto foi criado com foco em aprendizado e exploração de conceitos de desenvolvimento de APIs, concorrência em Go e técnicas de web scraping.

## Funcionalidades

*   **API RESTful:** Expõe um endpoint `/scrape` que aceita requisições POST contendo palavras-chave para iniciar o processo de scraping.
*   **Busca Inteligente:** Utiliza o DuckDuckGo HTML como motor de busca inicial para encontrar URLs relevantes com base nas palavras-chave.
*   **Extração de Contatos:** Navega pelas URLs encontradas e tenta extrair e-mails, números de telefone e links de perfis de mídias sociais usando expressões regulares.
*   **Processamento Paralelo:** Otimiza o tempo de execução ao processar múltiplas requisições de scraping e analisar diversas páginas em paralelo, aproveitando a concorrência nativa do Go (goroutines e channels).
*   **Tratamento de Erros e Timeout:** Inclui mecanismos para lidar com falhas de rede, erros de parsing de HTML e timeouts, garantindo a resiliência da aplicação.

## Tecnologias Utilizadas

*   **Go (Golang):** Linguagem de programação principal, escolhida por sua performance, concorrência e tipagem forte.
*   **Fiber:** Um framework web rápido e minimalista para Go, utilizado para construir a API de forma eficiente.
*   **`golang.org/x/net/html`:** Biblioteca padrão do Go para parsing e manipulação de documentos HTML.
*   **`regexp`:** Pacote para trabalhar com expressões regulares, essencial para a extração de padrões de contato.
*   **`context`:** Utilizado para gerenciar contextos de requisição, incluindo cancelamento e timeouts.
*   **`sync`:** Pacote para sincronização de goroutines (WaitGroup, Mutex).

## Estrutura do Projeto

O projeto segue uma estrutura modular para facilitar a organização e manutenção:

```
. 
├── cmd
│   └── server
│       └── main.go         # Ponto de entrada da aplicação, inicializa o servidor Fiber.
├── internal
│   ├── handler
│   │   └── scrape_handler.go # Lógica de manipulação da requisição de scraping.
│   ├── model
│   │   └── lead.go           # Definições das estruturas de dados (ScrapeRequest, LeadResult).
│   ├── service
│   │   └── scrapper.go       # Lógica principal de web scraping e orquestração.
│   └── util
│       ├── http_client.go    # Funções utilitárias para requisições HTTP.
│       └── parser.go         # Funções utilitárias para parsing e extração de dados com regex.
├── go.mod                  # Módulos e dependências do Go.
└── go.sum                  # Checksums das dependências.
```

## Como Rodar o Projeto

Para rodar este projeto localmente, siga os passos abaixo:

1.  **Pré-requisitos:** Certifique-se de ter o Go instalado em sua máquina (versão 1.23.0 ou superior, conforme `go.mod`).

2.  **Clonar o Repositório:**
    ```bash
    git clone https://github.com/EduardoCaversan/lead-scrapper-api.git
    cd lead-scrapper-api
    ```

3.  **Instalar Dependências:**
    ```bash
    go mod tidy
    ```

4.  **Executar a Aplicação:**
    ```bash
    go run cmd/server/main.go
    ```

    O servidor será iniciado na porta `8080`. Você verá uma mensagem no console:
    `Servidor rodando em http://localhost:8080`

## Exemplo de Uso (com `curl`)

Com o servidor rodando, você pode enviar uma requisição POST para o endpoint `/scrape`:

```bash
curl -X POST -H "Content-Type: application/json" \
-d '{"keywords": ["desenvolvedor go", "empresas de tecnologia"]}' \
http://localhost:8080/scrape
```

O resultado será um JSON contendo os leads encontrados, com informações como palavra-chave, título da página e URL.

## Pontos de Aprendizado

Este projeto foi uma excelente oportunidade para:

*   Aprofundar o conhecimento em Go, incluindo gerenciamento de módulos, estrutura de projetos e boas práticas.
*   Explorar o desenvolvimento de APIs RESTful com o framework Fiber.
*   Compreender os desafios e as abordagens para web scraping, incluindo parsing de HTML e extração de dados.
*   Aplicar e entender a importância da concorrência em Go para otimização de tarefas I/O-bound.
*   Desenvolver habilidades em tratamento de erros e construção de aplicações resilientes.

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues, propor melhorias ou enviar pull requests. Seu feedback e colaboração são muito importantes para o crescimento deste projeto e para o meu aprendizado.

---

## 👤 Autor

Desenvolvido por [Eduardo Caversan](mailto:educaversan.dev@gmail.com)
GitHub: [@EduardoCaversan](https://github.com/EduardoCaversan)

---

**Nota:** Este é um projeto de estudo e demonstração. O web scraping deve ser feito de forma ética e em conformidade com os termos de serviço dos sites visitados. Utilize-o com responsabilidade.