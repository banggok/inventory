
# Inventory Management System

An inventory management system for managing products, allowing users to create, retrieve, update, and delete product data. This system is built using Golang, Gin framework, GORM ORM, PostgreSQL, and includes E2E testing with Ginkgo.

## Table of Contents
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Prerequisites](#prerequisites)
- [Setup Instructions](#setup-instructions)
  - [Clone the Repository](#clone-the-repository)
  - [Set Up Environment Variables](#set-up-environment-variables)
  - [Run Migrations](#run-migrations)
- [Running the Application](#running-the-application)
  - [Running in Development](#running-in-development)
- [Running Tests](#running-tests)
  - [End-to-End Tests](#end-to-end-tests)
  - [Running with Coverage](#running-with-coverage)
- [API Endpoints](#api-endpoints)
  - [Create Product](#create-product)
  - [Get Product by ID](#get-product-by-id)
- [Project Structure](#project-structure)
- [Contribution Guidelines](#contribution-guidelines)

## Features
- Create new products
- Retrieve products by ID
- End-to-End tests for product management
- Graceful shutdown handling with SIGTERM/SIGINT
- Mock-based testing for use cases

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

### Set Up Environment Variables
Create a `.env` file in the root directory and add the following environment variables to configure your PostgreSQL database connection:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=yourusername
DB_PASSWORD=yourpassword
DB_NAME=inventory_db
```

Make sure PostgreSQL is running on your system and a database named `inventory_db` exists.

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
go run ./cmd/api/
```

The server will start on `http://localhost:8080`.

## Running Tests

### End-to-End Tests
This project uses Ginkgo for end-to-end testing. To run the tests:

```bash
ginkgo -r -v
```

This will run all tests across your project, providing verbose output.

### Running with Coverage
To run tests with coverage and generate a coverage report, use:

```bash
ginkgo -r -cover -coverpkg=./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

You can open `coverage.html` in a browser to see detailed test coverage results.

## API Endpoints

### Create Product

- **Endpoint**: `POST /api/v1/products`
- **Description**: Creates a new product.
- **Request Body**:
  ```json
  {
    "name": "Product Name"
  }
  ```
- **Response**:
  - **201 Created**:
    ```json
    {
      "id": 1,
      "name": "Product Name",
      "sku": "SKU-PRO-12345"
    }
    ```
  - **422 Unprocessable Entity** if the request body is invalid.

### Get Product by ID

- **Endpoint**: `GET /api/v1/products/:id`
- **Description**: Retrieves a product by its ID.
- **Response**:
  - **200 OK**:
    ```json
    {
      "id": 1,
      "name": "Product Name",
      "sku": "SKU-PRO-12345"
    }
    ```
  - **400 Bad Request** if the product ID is invalid.
  - **404 Not Found** if the product does not exist.

## Project Structure

- **cmd/api/**: The main entry point for the application.
- **api/**: Contains API handlers and DTOs.
- **internal/**: Core business logic, use cases, and repositories.
- **pkg/db/**: Database connection setup.
- **migrations/**: SQL files for database schema creation.
- **e2e/**: End-to-end tests.

## Contribution Guidelines

If you'd like to contribute, feel free to fork the repository and submit a pull request. Make sure to:
- Write tests for new features and bug fixes.
- Follow Go conventions and run `go fmt` to format your code.
- Include detailed commit messages.

