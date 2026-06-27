# Smart Parking & EV Charging Reservation System

A centralized platform for busy airports and malls to manage parking zones, specifically handling the high-demand reservation of limited EV charging spots. Built using **Go**, **Echo (v5)**, **GORM**, and **PostgreSQL**.

## 🛠️ Technology Stack

| Technology | Description |
| --- | --- |
| **Go (Golang)** | Version 1.22+ |
| **Echo** | `github.com/labstack/echo/v5` (High performance, minimalist web framework) |
| **GORM** | `gorm.io/gorm` (ORM for Go, using PostgreSQL driver) |
| **PostgreSQL** | Relational database (NeonDB or Supabase) |
| **Validator** | `github.com/go-playground/validator/v10` (Struct validation) |
| **JWT** | `github.com/golang-jwt/jwt/v5` (Standard token generation & verification) |
| **bcrypt** | `golang.org/x/crypto/bcrypt` (Password hashing, cost 10) |

---

## 🏛️ Project Structure (Clean Architecture)

This project strictly adheres to a clean architecture pattern separating concerns into distinct layers:

```text
gotickets/
├── cmd/
│   └── main.go                 # Application entry point & Dependency Injection
├── internal/
│   ├── config/                 # Environment variables and DB connection logic
│   ├── dto/                    # Data Transfer Objects (Request/Response structures)
│   ├── handler/                # HTTP layer (Echo endpoints, data binding, returns JSON)
│   ├── middlewares/            # JWT authentication and Role-based access control
│   ├── models/                 # Database schema models (GORM)
│   ├── repository/             # Database access operations (CRUD, Transactions)
│   ├── service/                # Core business logic (Validation, Calculations)
│   └── utils/                  # Shared utilities (Standardized JSON responses)
├── .env                        # Environment configuration
├── go.mod                      # Go module dependencies
└── README.md                   # Project documentation
```

### Layer Responsibilities

- **DTO**: Defines request payloads and response structures. Prevents exposing GORM models directly.
- **Handler**: Binds/validates incoming HTTP requests, extracts JWT claims, calls Services, and formats HTTP responses. Handlers **never** talk to the database directly.
- **Service**: Executes core business logic (e.g., hashing passwords, capacity checks, generating JWTs) and coordinates with Repositories.
- **Repository**: Handles all data access and GORM database operations (CRUD, row locks, transactions).

---

## 👥 User Roles & Permissions

| Role | Allowed Actions |
| --- | --- |
| **driver** | • Register and log in<br>• View all parking zones and availability<br>• Reserve a parking/EV spot<br>• View and cancel their own reservations |
| **admin** | • All driver permissions<br>• Create parking zones<br>• View all reservations in the system |

---

## 🔐 Concurrency & The "EV Spot Bottleneck"

A primary technical challenge in this platform is preventing overbooking for highly demanded EV spots. 
If `total_capacity` is 20, and two drivers attempt to book the very last spot at the exact same millisecond, standard SQL queries might read "19 active" for both requests, resulting in 21 cars booked.

**The Solution:**
This project utilizes a **GORM Database Transaction** combined with **Row-Level Locking (`FOR UPDATE`)** on the `parking_zones` record. The repository locks the specific zone row, safely counts active reservations, verifies capacity limits, creates the reservation, and then releases the lock atomically. 

---

## 🚀 Running the Project

### 1. Configure Environment
Create a `.env` file in the root directory and define the following variables:
```env
DSN="host=localhost user=postgres password=postgres dbname=gotickets port=5432 sslmode=disable TimeZone=Asia/Dhaka"
PORT=8080
JWT_SECRET="super_secret_jwt_key"
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Run the Server
```bash
go run cmd/main.go
```
*Alternatively, you can build and execute the binary:*
```bash
go build -o bin/gotickets ./cmd/main.go
./bin/gotickets
```

---

## 🌐 API Endpoints

### 🔹 Authentication
- `POST /api/v1/auth/register` (Public) - Register a new user (`driver` or `admin`).
- `POST /api/v1/auth/login` (Public) - Authenticate and receive a JWT.

### 🔹 Parking Zones
- `POST /api/v1/zones` (Admin Only) - Create a new parking zone.
- `GET /api/v1/zones` (Public) - View all zones and dynamically calculated `available_spots`.
- `GET /api/v1/zones/:id` (Public) - View details for a single zone.

### 🔹 Reservations
- `POST /api/v1/reservations` (Driver/Admin) - Book a spot (Executes row-level locked transaction).
- `GET /api/v1/reservations/my-reservations` (Driver/Admin) - View your own active bookings.
- `DELETE /api/v1/reservations/:id` (Driver/Admin) - Cancel your own booking.
- `GET /api/v1/reservations` (Admin Only) - View all reservations system-wide.
