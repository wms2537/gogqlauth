# Golang Authentication Backend with gqlgen

This project demonstrates a simple authentication backend using Golang and gqlgen. It provides a GraphQL API for user authentication, including features like login, password change, and token refresh.

## Table of Contents

- [Golang Authentication Backend with gqlgen](#golang-authentication-backend-with-gqlgen)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Project Structure](#project-structure)
  - [Getting Started](#getting-started)
  - [How It Works](#how-it-works)
  - [Authentication Flow](#authentication-flow)
  - [Running the Project](#running-the-project)
  - [Development with Air](#development-with-air)
  - [GraphQL Schema](#graphql-schema)
  - [Key Components](#key-components)
  - [Security Considerations](#security-considerations)
  - [Social Login](#social-login)
  - [OTP Login](#otp-login)
  - [Graphql Playground](#graphql-playground)
  - [Additional Information](#additional-information)

## Prerequisites

- Go 1.23.1 or later
- Docker (optional, for containerized deployment)
- SurrealDB (as the database)
- Minio (for object storage)

## Project Structure

The project follows a typical Golang project structure with gqlgen-specific directories:

```
.
├── .env
├── .gitignore
├── Dockerfile
├── Dockerfile.dev
├── README.md
├── cron.go
├── go.mod
├── go.sum
├── gqlgen.yml
├── graph/
│   ├── generated.go
│   ├── model/
│   ├── resolver.go
│   ├── schema.graphqls
│   ├── schema.resolvers.go
│   └── ...
├── server.go
└── tools.go
```

## Getting Started

1. Clone the repository:
   ```
   git clone <repository-url>
   cd <project-directory>
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Set up your environment variables in a `.env` file, using the `.env.example` file as a template:
```
PORT=8080
SMTP_SERVER=<your-smtp-server>
SMTP_USER=<your-smtp-user>
SMTP_PASS=<your-smtp-pass>
SURREALDB_URL=<your-surrealdb-url>
SURREALDB_USER=<your-surrealdb-user>
SURREALDB_PASSWORD=<your-surrealdb-password>
SURREALDB_NS=<your-surrealdb-ns>
SURREALDB_DB=<your-surrealdb-db>
MINIO_URL=<your-minio-url>
MINIO_ACCESS_KEY_ID=<your-minio-access-key-id>
MINIO_ACCESS_KEY_SECRET=<your-minio-access-key-secret>
MINIO_BUCKET_NAME=<your-minio-bucket-name>
```

4. Generate GraphQL code:
   ```
   go run github.com/99designs/gqlgen generate
   ```

5. Run the project:
   ```
   go run server.go
   ```

## How It Works

This project uses gqlgen to generate a GraphQL server based on the schema defined in `graph/schema.graphqls`. The main components are:

1. GraphQL Schema: Defines the API structure (queries and mutations).
2. Resolvers: Implement the logic for each query and mutation.
3. Models: Represent the data structures used in the API.
4. Middleware: Handles authentication and request processing.

The authentication system uses JSON Web Tokens (JWT) for secure communication between the client and server.

## Authentication Flow

1. User Registration:
   - Create a new user in the database.
   - Send verification email to the user (not implemented in this example).

2. Login:
   - User provides email and password.
   - Server verifies credentials against the database.
   - If valid, a new JWT token pair (access token and refresh token) is generated.

3. Token Usage:
   - Client includes the access token in the Authorization header for authenticated requests.
   - Server validates the token using the public key.

4. Token Refresh:
   - When the access token expires, the client can use the refresh token to obtain a new token pair.

5. Password Change and Reset:
   - Endpoints are provided for changing passwords and requesting password resets.

## Running the Project

1. Start your SurrealDB instance.

2. Run the server:
   ```
   go run server.go
   ```

3. The GraphQL playground will be available at `http://localhost:8080/graphql`.

## Development with Air

For hot-reloading during development, you can use Air:

1. Install Air:
   ```
   go install github.com/cosmtrek/air@latest
   ```

2. Run the project with Air:
   ```
   air
   ```

Or you can use the `docker-compose.yml` file to run the project with Air.
```
docker-compose up --build
```

## GraphQL Schema

The GraphQL schema defines the following main types:

- User
- Token
- PasswordChange

And the following operations:

- Query:
  - user: Fetch the current user's information

- Mutation:
  - loginWithEmailPassword
  - changePassword
  - requestChangePassword
  - refreshToken

For the complete schema, refer to `graph/schema.graphqls`


## Key Components

1. Middleware (graph/middlewares/middleware.go):
   - Handles JWT validation for authenticated requests.

2. Login Handler (graph/utils/controllers.go):
   - Manages user authentication and token generation.

3. Token Verification (graph/utils/appleToken.go):
   - Verifies JWT tokens using public keys.

4. Email Sending (graph/utils/gomail.go):
   - Handles sending verification and password reset emails.

5. Key Rotation (cron.go):
   - Periodically rotates JWT signing keys for enhanced security.
  
6. Minio Object Storage (graph/utils/minio.go):
   - Handles file uploads and downloads.
   - Used for storing user avatars, files, etc.

## Security Considerations

1. Password Hashing: The project uses Argon2 for password hashing (as seen in the login query).

2. JWT Key Rotation: Implemented in `cron.go` to periodically generate new signing keys.

3. Environment Variables: Sensitive information is stored in environment variables.

4. CORS: Configured in `server.go` to control cross-origin requests.

5. Rate Limiting: Not implemented in the provided code, often implemented in load balancer.

## Social Login

Social login with google, apple, facebook, etc works almost the same way as the email/password login.

The main difference is that you need to implement the social login logic in the resolvers, and they take in the token from the social login provider instead of the email and password. You have to refer to the official documentation of the social login provider to verify the token and get the user's information.

## OTP Login

OTP login is a login method that sends a temporary code to the user's phone number. This is useful for passwordless login.

To implement OTP login, you need to:

1. Send a verification code to the user's phone number.
2. Ask the user to enter the code in the app.
3. Verify the code and generate a JWT token.

So to implement this you will need to add 2 methods in the mutation resolver:

1. sendOTP: Sends a verification code to the user's phone number.
2. verifyOTP: Verifies the code and generates a JWT token.

You will also need to implement the logic to send and verify the OTP code. This is usually done with a third-party service like Twilio or Nexmo.

## Graphql Playground

You can use the Graphql playground to test the API. It is available at `http://localhost:8080/graphql` after running the project.


## Additional Information

For more information, refer to the [gqlgen documentation](https://gqlgen.com/getting-started/).

Golang documentation: https://pkg.go.dev/

SurrealDB documentation: https://surrealdb.com/docs/

Minio documentation: https://min.io/docs/

