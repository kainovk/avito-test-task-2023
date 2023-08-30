# Avito test task 2023 (https://github.com/avito-tech/backend-trainee-assignment-2023)

## Getting started

### Run locally

1. Create `.env` file in the root of the project with a similar structure:

```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=root
POSTGRES_DB=avito_db
CONFIG_PATH=config/local.yml
```

or simply rename [`template.env`](template.env) to `.env` with your configuration

2. Start docker containers: \
   `docker-compose up -d`

### Swagger endpoint: http://\<HOST>:\<PORT>/swagger/

![swagger.png](attachments%2Fswagger.png)

## Sample queries

### Segments

**Create New Segment** \
Request \
`POST` http://localhost:8080/segments 
```json
{
"name": "AVITO" 
}
```

Response: 200 
```json
{
    "status": "OK"
}
```

Request \
`POST` http://localhost:8080/segments
```json
{
"name": "AVITO" 
}
```

Response: 400
```json
{
   "status": "Error",
   "error": "segment already exists"
}
```

**Get All Segments** \
Request \
`GET` http://localhost:8080/segments

Response: 200
```json
{
   "segments": [
      "AVITO_VOICE_MESSAGES",
      "AVITO_DISCOUNT"
   ]
}
```

**Delete Segment** \
Request \
`DELETE` http://localhost:8080/segments/AVITO_VOICE_MESSAGES

Response: 200
```json
{
   "status": "OK"
}
```

Request \
`DELETE` http://localhost:8080/segments/AVITO_UNKNOWN

Response: 404
```json
{
   "status": "Error",
   "error": "segment does not exist"
}
```

### Users

**Create New User** \
Request \
`POST` http://localhost:8080/users
```json
{
"name": "Kirill" 
}
```

Response: 200
```json
{
    "status": "OK"
}
```

Request \
`POST` http://localhost:8080/users
```json
{
"name": "Kirill" 
}
```

Response: 400
```json
{
   "status": "Error",
   "error": "user already exists"
}
```

**Configure User Segments** \
Request \
`POST` http://localhost:8080/users/1/configure-segments
```json
{
   "segments_to_add": [
      {
         "slug": "AVITO_DISCOUNT",
         "delete_at": "2023-08-29T14:05:00Z"
      },
      {
         "slug": "AVITO_VOICE_MESSAGES"
      }
   ],
   "segments_to_del": null
}
```

Response: 200
```json
{
    "status": "OK"
}
```

**Note**: add to user with id=1 segments - AVITO_DISCOUNT and AVITO_VOICE_MESSAGES. 
AVITO_DISCOUNT will be deleted at time `delete_at` (or after 1 minute if `delete_at` in the past).

**Get User Segments** \
Request \
`GET` http://localhost:8080/users/1/segments

Response: 200
```json
{
   "segments": [
      "AVITO_VOICE_MESSAGES"
   ]
}
```
