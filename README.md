
---

## 🏗 Backend Architecture (Clean Architecture)

The backend is engineered using **Clean Architecture** (Uncle Bob) principles. This ensures the IMCI logic remains decoupled from the web framework, database, and external dependencies.

### Architectural Layers
* **Domain (Entities):** The core business models (e.g., `Child`, `Assessment`, `Classification`) and the logic defined by the WHO IMCI booklet.
* **Usecases:** Orchestrates the flow of data. Contains the specific business rules for classifying a child based on symptoms (e.g., `ClassifySickChild`, `CalculateNutritionStatus`).
* **Repositories (Interfaces):** Defines the contracts for data persistence. The domain layer doesn't care if you use PostgreSQL, MongoDB, or an In-memory store.
* **Controllers (Delivery):** Handles the HTTP transport layer, input validation, and mapping JSON requests to internal Use Case structures.

### Directory Structure
```bash
backend/
├── cmd/                # Entry point (main.go) and dependency injection
├── internal/
│   ├── domain/         # Entities and Repository Interfaces
│   ├── usecase/        # Business Logic / Interactors
│   ├── repository/     # Implementation of data stores (e.g., Postgres/Gorm)
│   └── controller/     # HTTP Handlers & Middlewares
├── pkg/                # Shared utilities (validation, logger)
└── api/                # Swagger/OpenAPI documentation
```

---

## 🛠 Tech Stack

* **Language:** [Go (Golang)](https://go.dev/) 1.21+
* **API Framework:** [Gin Gonic](https://github.com/gin-gonic/gin) (or Fiber)
* **Database:** [PostgreSQL](https://www.postgresql.org/) 
* **Testing:** Go sub-tests for table-driven logic verification.

---

## 🧪 Implementation Detail: IMCI Logic Engine

The core strength of this backend is the **Deterministic Logic Engine**. Unlike standard CRUD apps, every input is validated against the **IMCI Color-Coded Matrix**.

> **Pro Tip:** We use the **Repository Pattern** to ensure that patient history can be stored for longitudinal health tracking, while the **Usecase layer** ensures that a "Red" classification always triggers a high-priority flag in the system.

### Example Usecase: Cough Assessment
```go
// Simplified logic within the usecase layer
func (u *IMCIUsecase) Execute(input AssessmentInput) (Classification, error) {
    if input.HasGeneralDangerSigns() {
        return domain.SevereDisease, nil // Pink/Red Row
    }
    // ... further logic based on booklet thresholds
}
```

---

## 🚀 Getting Started

### 1. Environment Setup
Create a `.env` file in the root directory:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=secret
DB_NAME=imci_db
PORT=8080
```

### 2. Installation & Migrations
```bash
# Install dependencies
go mod download

# Run database migrations (if applicable)
go run cmd/migrate/main.go

# Start the development server
go run cmd/api/main.go
```

---

## 📉 Testing Strategy

Given the medical nature of this project, we maintain high test coverage:
* **Unit Tests:** Testing the classification logic with hundreds of symptom combinations.
* **Integration Tests:** Testing the Controller -> Usecase -> Repository flow.
* **Mocking:** Using `gomock` or `testify` to mock database calls during logic testing.

```bash
go test -v ./internal/usecase/...
```

---

## 🤝 Contributing

We follow strict **Uncle Bob's Clean Architecture** guidelines. Please ensure that:
1.  Entities do not import Usecases.
2.  Usecases do not import Controllers.
3.  All business logic changes are accompanied by unit tests.

---
