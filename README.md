# Avito test task 2023 (https://github.com/avito-tech/backend-trainee-assignment-2023)

## Getting started
1. Create `.env` file in the root of the project with a similar structure:
```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=root
POSTGRES_DB=avito_db

CONFIG_PATH=config/local.yml
```

or simply rename [`template.env`](template.env) to `.env` with your configuration

2. Start postgres container: \
`docker-compose up postgres -d`
3. Start server: \
`go run cmd/avito-slug/main.go`


### Swagger endpoint: /swagger/index.html#/
![swagger.png](attachments%2Fswagger.png)


