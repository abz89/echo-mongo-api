# REST API with Golang and MongoDB using Echo Framework

This is a simple example of building RESTful API using Golang and MongoDB with the [Echo](https://echo.labstack.com/) web framework that using this [article](https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-echo-version-2gdg) as a reference.

## Prerequisites

- Golang installed on your machine. You can download it from [the official Golang website](https://go.dev/dl).
- MongoDB installed locally or MongoDB Atlas cluster. You can download MongoDB from [the official MongoDB website](https://www.mongodb.org/try/download/community).
- [Postman](https://www.postman.com/) or any other API testing tool.

## Getting Started

1. Clone this repository:

```bash
git clone https://github.com/abz89/echo-mongo-api
```

2. Navigate to the project directory

```bash
cd echo-mongo-api
```

3. Install the project dependencies

```bash
go mod tidy
```

4. Adjust `.env` file int the project root directory and add the following configuration

```bash
PORT=<APP PORT>
MONGO_URI=<MONGODB URI>
SECRET=<JWT SECRET>
ADMIN_USER=<ADMINISTRATOR USERNAME>
ADMIN_PASSWORD=<ADMINISTRATOR PASSWORD>
```

5. Run the application

```bash
go run main.go
```

## API Endpoints

- `POST /register`: Register a new user
- `POST /login`: Login a user for obtaining JWT token
- `POST /users`: Same as `POST /register`
- `GET /users`: Get all users (protected by basic auth)
- `GET /users/:userId`: Get a user by ID (protected by JWT auth & guarded by matching userId)
- `PUT /users/:userId`: Update a user by ID (protected by JWT auth & guarded by matching userId)
- `PATCH /users/:userId`: Update (partially) a user by ID (protected by JWT auth & guarded by matching userId)
- `DELETE /users/userId`: Delete a user by ID (protected by JWT auth & guarded by matching userId)

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvement, feel free to open an issue or create a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
