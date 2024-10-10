## Technologies Used
- **Golang** (1.23)
- **Gin** (Web framework)
- **GORM** (ORM for PostgreSQL)
- **PostgreSQL** (Database)
- **Ginkgo** and **Gomega** (Testing frameworks)
- **Testify** (Mocking framework)

## Prerequisites
Before you start, ensure that you have the following installed:
- **Golang** (1.23 or above)
- **PostgreSQL** (version 12 or higher)
- **Git**

## Setup Instructions

### Clone the Repository
First, clone the repository to your local machine:
```bash
git clone https://github.com/yourusername/inventory_management.git
cd inventory_management
```

### Run Migrations
To set up the database schema, run the SQL migration file:

```bash
soda migrate
```

This will create the necessary table in your PostgreSQL database.

## Running the Application

### Running in Development
To run the application in development mode, execute the following command:

```bash
make run
```

The server will start on `http://localhost:8080`.

## Running Tests

### End-to-End Tests
This project uses Ginkgo for end-to-end testing. To run the tests:

```bash
make test
```

This will run all tests across your project, providing verbose output.

### Running with Coverage
To run tests with coverage and generate a coverage report, use:

```bash
make coverage
```

You can open `coverage.html` in a browser to see detailed test coverage results.

### Running Sonar Scanner
To run tests with coverage and generate a coverage report, use:

```bash
make sonar
```

## Project Structure

```bash
├── api                 # Contains API handlers and DTOs (Data Transfer Objects).
│   ├── handler         # API handlers (controllers) and DTOs
│   ├── transformer     # Transforms entities to DTOs for responses
│   └── helper          # Utility functions for API handling
├── cmd                 # Responsible for bootstrapping and configuring the application.
│   └── api             # Main entry point of the application
├── internal            # Core business logic, use cases, and repositories.
│   ├── entity          # Data logic entities
│   ├── model           # Database models
│   ├── repository      # Data access layer (repositories)
│   └── usecase         # Application use cases (business logic)
├── migrations          # SQL migration files
├── pkg                 # Contains database utilities and configuration.
│   └── db              # Database connection setup
│   └── utility         # Utility Helper
└── tests               # Contain tests
    └── e2e             # End-to-end tests
```


## Contribution Guidelines

If you'd like to contribute, feel free to fork the repository and submit a pull request. Make sure to:
- Write tests for new features and bug fixes.
- Follow Go conventions and run `go fmt` to format your code.
- Include detailed commit messages.
