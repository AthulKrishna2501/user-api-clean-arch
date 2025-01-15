# 🏗️ User API with Clean Architecture  

This project is a **User API** built with **Clean Architecture** principles to ensure modularity, scalability, and ease of maintenance. It includes **unit testing** for key functionalities and integrates **CI/CD pipelines** for streamlined development and deployment.  

---

## 🚀 Features  

- **Clean Architecture**: Layered design for separation of concerns and flexibility.  
- **User API**: CRUD operations for user management with JWT-based authentication.  
- **Unit Testing**: Comprehensive tests for controllers, repositories, and services.  
- **CI/CD Integration**: Automated workflows using GitHub Actions for testing and deployment.  
- **Logging**: Centralized and efficient logging for better debugging and monitoring.  

---

## 📂 Directory Structure  

```plaintext
athulkrishna2501-user-api-clean-arch/
├── README.md                  # Project documentation
├── go.mod                     # Dependency management
├── go.sum                     # Dependency versions
├── cmd/
│   └── app/
│       └── main.go            # Application entry point
├── internal/
│   ├── app/
│   │   ├── config/            # Configuration handling
│   │   ├── controllers/       # API controllers and tests
│   │   └── utils/             # Utility functions (e.g., JWT handling)
│   ├── core/
│   │   ├── database/          # Database connection logic
│   │   ├── models/            # Data models and validation
│   │   ├── repository/        # Repository layer for database interaction
│   │   └── services/          # Business logic layer
│   ├── logger/                # Logging implementation
│   └── mocks/                 # Mock implementations for testing
└── .github/
    └── workflows/
        └── main.yml           # CI/CD pipeline configuration
