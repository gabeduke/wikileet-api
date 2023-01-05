# Wikileet

## Create .env file

    ```sh
    touch .env && echo "CFG_DATABASE_USER=postgres
    CFG_DATABASE_PASSWORD=postgres
    CFG_SESSION_SECRET=asdf
    CFG_AUTH_INTERNAL=true
    CFG_PORT=8080" >> .env
    ```

## Database

PostgreSQL 15.1 running on localhost:5432

`docker compose up -d postgres`

## API / BACKEND

Compile the application locally

`go run main.go`

## FRONTEND

Change to the /web directory and build with npm

`cd web`
`npm i`
`npm run watch`

Now navigate to `localhost:8080` in your browser and login:
> admin/admin OR test/test
