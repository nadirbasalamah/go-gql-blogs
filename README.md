# go-gql-blogs

The blog application implemented in GraphQL with `gqlgen`.

## How To Use

There are two branches in this repository:

- `local-storage` for implementation with local storage.
- `main` for implementation with MongoDB database.

1. Clone the repository.

2. If `local-storage` branch is chosen, run the application with this command.

```
go run server.go
```

3. If `main` branch is chosen, copy the `.env` file.

```
cp .env.example .env
```

4. Configure the MongoDB database configuration inside the `.env` file.

5. Run the MongoDB database server or make sure the MongoDB database is online.

6. Run the application with this command.

```
go run server.go
```

7. Test the application with this command.

```
go test
```
