# Iniciar

docker-compose up -d --build --force-recreate

`go run send.go`

Para recibir mensajes:

`go run receiver_console.go`
o
`go run receiver_db.go`


# crear tabla
psql postgresql://postgres:postgres@localhost:5432/postgres 


CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL
);