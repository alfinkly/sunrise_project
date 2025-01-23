# Geolocation Service Project

## Overview

This project is a geolocation service that provides geolocation information based on IP addresses. It includes both a backend API built with Go and Gin, and a frontend built with HTML, CSS, and JavaScript.

## Prerequisites

* Docker
* Docker Compose

## Getting Started

### Environment Variables

Create a `.env` file in the root directory of the project based on the `.env_template` file. Update the values as needed.

### Running the Project

1. Build and start the services using Docker Compose:
   ```sh
   docker-compose up --build
   ```

2. The services will be available at the following URLs:
    * Frontend: http://localhost:80
    * API: http://localhost:8080

3. You can test it with Postman
   * Import to your work space [Sunrise Project.postman_collection.json](Sunrise%20Project.postman_collection.json)

## API Endpoints

* `GET /location/:ip` - Get geolocation information for a specific IP address.
* `GET /locations` - Get geolocation information for all stored IP addresses.
* `GET /` - Get a secret value.

## Project Structure

* `cmd/api/main.go` - Entry point for the backend API.
* `cmd/api/docs` - Swagger documentation for the API.
* `internal/dao` - Data access objects.
* `internal/handler` - HTTP handlers.
* `internal/platform` - Platform-specific code, such as database connections.
* `internal/repository` - Repositories for accessing data.
* `internal/service` - Business logic and services.
* `web` - Frontend code.

## Frontend

The frontend is a simple HTML page that interacts with the backend API to display geolocation information. It is served by an Nginx server.

## Docker

The project uses Docker for containerization. The `docker-compose.yml` file defines the services:
* `db` - PostgreSQL database.
* `api` - Backend API.
* `web` - Frontend web server.

## Additional Information

* The backend API uses Gin for routing and GORM for database interactions.
* The frontend uses JavaScript to fetch data from the backend API and display it on the page.