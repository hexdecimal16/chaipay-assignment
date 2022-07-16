# Chai Pay Stripe Api Integration Assignment

This container is used to test the Stripe API integration.

## Environment Variables

- `DB_HOST` - The hostname of the database.
- `DB_PORT` - The port of the database.
- `DB_USER` - The username of the database.
- `DB_PASSWORD` - The password of the database.
- `DB_NAME` - The name of the database.
- `STRIPE_KEY` - The Stripe API key.

## Dependencies

The following dependencies are required:

- [GOLANG](https://golang.org/): The Go language.
- [Docker](https://www.docker.com/): The Docker CLI.
- [Docker Compose](https://docs.docker.com/compose/install/): The Docker Compose CLI.

## Usage

You can run the docker container with the following command:

```bash
  docker-compose up --build
```

To start unitesting, you can run the following command:

```bash
  go test
```
