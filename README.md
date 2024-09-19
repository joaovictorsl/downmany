## Comunicação cliente-servidor

- Métodos
    - Join()
        - Informar conexão ao entrar na rede
        - Reafirmar conexão periodicamente
    - GetIPs()
        - Retorna IPs conhecidos pelo servidor

## Comunicação cliente-cliente

- Métodos
    - HasFile(hash uint64)
        - Retorna se o cliente tem ou não um arquivo
    - Download(start, end uint64)
        - Começa a fazer o download a partir de offset
     
## Lógica do projeto

  - Servidor armazena a lista de todos os IPS presentes na rede
  - Cada cliente possui um cache local para armazenar a lista de quem já é conhecido por ele na rede e a informação de quem também possui os arquivos que ele possui
  - Comunicação com o servidor:
      - Quando o cliente entra na rede e periodicamente para confirmar a presença
      - Quando o cliente não consegue o arquivo com a rede que ele já conhece
   
  - Um arquivo não existe no sistema quando todos presentes na rede informam que não possuem esse arquivo
  - Caso cliente não encontre arquivos na rede atual, realizamos um backoff exponencial

## Repositório

- `core`
  - lógica da aplicação
- `network`
  - lógica de comunicação
- `main.go`
  - script para iniciar servidor ou cliente

## Como rodar a demo

```
./run_demo.sh
```

Em seguida veja os logs dos clientes 1~4 e o servidor

## Docker

Buildar imagem
`docker build -t downmany:1.0 .`

Executar imagem
`docker run -p 3000:3000 -it downmany:1.0 ./bin/downmany -file_hash 123 -server_addr 127.0.0.1:8000 -port 1234`

TODO: docker-compose para subir 3 clientes e um servidor em portas diferentes. Adicionar volume para mapear a pasta /dataset
