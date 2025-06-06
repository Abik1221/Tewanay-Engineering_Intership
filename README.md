# Tewanay Engineering Internship - Restaurant Management Backend

A robust backend infrastructure for a Restaurant Management System built with **Golang** and the **Gin** web framework. This project is developed as part of the Tewanay Engineering Internship Program and is designed to provide a scalable, secure, and feature-rich API for managing restaurant operations.

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Environment Variables](#environment-variables)
- [Available Endpoints](#available-endpoints)
- [Authentication & Security](#authentication--security)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- **User Authentication**: Secure signup and login with JWT-based authentication.
- **Role Management**: Admin and user roles for access control.
- **Menu Management**: CRUD operations for restaurant menus.
- **Food Management**: Add, update, delete, and list food items.
- **Order Management**: Place, update, and track orders.
- **Invoice Management**: Generate and manage invoices for orders.
- **Table Management**: Manage restaurant tables and their statuses.
- **Ordered Items**: Track items ordered per order.
- **Swagger API Docs**: Interactive API documentation with Swagger UI.
- **MongoDB Integration**: Persistent storage using MongoDB.
- **Validation**: Request validation using go-playground/validator.
- **Middleware**: Custom authentication middleware for protected routes.

---

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin
- **Database**: MongoDB
- **Authentication**: JWT (github.com/golang-jwt/jwt)
- **Validation**: go-playground/validator
- **API Docs**: Swagger (swaggo/gin-swagger)
- **Other**: Docker (optional for deployment)

---

## Project Structure

```
.
├── controllers/         # Route handlers for each resource
├── database/            # MongoDB connection logic
├── docs/                # Swagger documentation files
├── helpers/             # Utility functions (e.g., JWT handling)
├── middlewares/         # Custom middleware (e.g., Auth)
├── models/              # Data models (MongoDB schemas)
├── routes/              # Route grouping and registration
├── services/            # (Planned) AI recommendation and analytics
├── main.go              # Application entry point
├── go.mod / go.sum      # Go module files
├── .env                 # Environment variables
└── README.md            # Project documentation
```

---

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.18+
- [MongoDB](https://www.mongodb.com/try/download/community)
- (Optional) [Docker](https://www.docker.com/)

### Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/abik1221/Tewanay-Engineering_Intership.git
   cd Tewanay-Engineering_Intership
   ```

2. **Install dependencies:**
   ```sh
   go mod tidy
   ```

3. **Configure environment variables:**

   Create a `.env` file in the root directory:

   ```
   PORT=8080
   SECRET_KEY=your_jwt_secret
   ```

   *(Replace `your_jwt_secret` with a secure random string.)*

4. **Start MongoDB** (if not already running):

   ```sh
   mongod
   ```

5. **Run the application:**
   ```sh
   go run main.go
   ```

6. **Access the API:**
   - API base URL: `http://localhost:8080`
   - Swagger docs: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## API Documentation

Interactive API documentation is available via Swagger UI:

- Visit: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

All endpoints, request/response schemas, and authentication requirements are documented.

---

## Environment Variables

| Variable     | Description                | Example         |
|--------------|----------------------------|-----------------|
| PORT         | Server port                | 8080            |
| SECRET_KEY   | JWT signing secret         | mysecretkey     |
| GEMINI_API_KEY | (Optional) AI API key    | ...             |

---

## Available Endpoints

### User & Auth

- `POST /users/signup` — Register a new user
- `POST /users/login` — Login and receive JWT tokens
- `GET /users` — List all users *(admin only)*
- `GET /users/:user_id` — Get user by ID *(admin/user)*

### Menu

- `GET /menus` — List all menus
- `GET /menus/:menu_id` — Get menu by ID
- `POST /menus` — Create menu *(admin)*
- `PATCH /menus/:menu_id` — Update menu *(admin)*
- `DELETE /menus/:menu_id` — Delete menu *(admin)*

### Food

- `GET /foods` — List all foods (paginated)
- `GET /foods/:food_id` — Get food by ID
- `POST /foods` — Create food *(admin)*
- `PATCH /foods/:food_id` — Update food *(admin)*
- `DELETE /foods/:food_id` — Delete food *(admin)*

### Orders

- `GET /orders` — List all orders
- `GET /orders/:order_id` — Get order by ID
- `POST /orders` — Create order
- `PATCH /orders/:order_id` — Update order
- `DELETE /orders/:order_id` — Delete order

### Invoices

- `GET /invoices` — List all invoices
- `GET /invoices/:invoice_id` — Get invoice by ID
- `POST /invoices` — Create invoice
- `PATCH /invoices/:invoice_id` — Update invoice
- `DELETE /invoices/:invoice_id` — Delete invoice

### Tables

- `GET /tables` — List all tables
- `GET /tables/:table_id` — Get table by ID
- `POST /tables` — Create table *(admin)*
- `PATCH /tables/:table_id` — Update table *(admin)*
- `DELETE /tables/:table_id` — Delete table *(admin)*

### Ordered Items

- `GET /order_items` — List all ordered items
- `GET /order_items/:order_item_id` — Get ordered item by ID
- `GET /orderItems-order/:order_id` — Get ordered items by order ID
- `POST /order_items` — Create ordered item
- `PATCH /order_items/:order_item_id` — Update ordered item
- `DELETE /order_items/:order_item_id` — Delete ordered item

---

## Authentication & Security

- **JWT Authentication**: Most endpoints require a valid JWT token in the `token` header.
- **Role-based Access**: Certain endpoints are restricted to admin users.
- **Password Hashing**: User passwords are securely hashed using bcrypt.
- **Input Validation**: All input data is validated for security and integrity.

---

## Contributing

Contributions are welcome! Please open issues or submit pull requests for improvements, bug fixes, or new features.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/YourFeature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/YourFeature`)
5. Open a pull request

---

## License

This project is licensed under the [MIT License](LICENSE).

---

## Acknowledgements

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [Swaggo](https://github.com/swaggo/swag)
- [Go Playground Validator](https://github.com/go-playground/validator)
- Tewanay Engineering Internship Program

---