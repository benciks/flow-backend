# flow-backend
A GraphQL Server for mobile application about time and task management called [Flow](https://github.com/benciks/flow-native). This server was built as part of bachelor thesis.

## Getting Started
To get a local copy up and running follow these simple steps.

### Prerequisites
- Docker
- Docker Compose

### Running the server
1. Clone the repo
2. Run the following command to start the server
```$ docker-compose up```. This will run the environment containing necessary tools.
3. The server will be running on `http://localhost:3000`

Note: It is encouraged to modify the Dockerfile environment variables to match your certificate details.
### Development
It is possible to run the server alone, however, following conditions must be met:
- Running and configured [TaskD server](https://github.com/GothenburgBitFactory/taskserver)
- Running and configured [timew-sync-server](https://github.com/timewarrior-synchronize/timew-sync-server)

Once the conditions are met, follow these steps:
1. Clone the repo
2. Copy the `.env.example` file to `.env` and fill in the necessary information
3. Run the migrations by running the following command
```$ go run cmd/db/main.go```
3. Run the following command to start the server
```$ go run cmd/server/main.go```. This will run the server in development mode.