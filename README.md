# ContentFlow CMS 

A powerful, high-performance Headless CMS built with Go and Fiber. Designed for flexibility, it supports rich content structures, multiple languages, and advanced content management features.

## Features

*   **Rich Content Blocks**: Support for structured, block-based content (similar to Notion/Editor.js) via JSON.
*   **Localization**: Built-in support for multi-language content with translation grouping.
*   **Taxonomies**: Organize content using robust **Categories** and **Tags**.
*   **Scheduled Publishing**: Schedule content to automatically go live at a specific date and time.
*   **Webhooks**: Real-time event triggers (`content.create`, `content.update`, `content.published`) to integrate with external systems (CI/CD, static site generators, etc.).
*   **Advanced Search**: Filter content by status, type, language, tags, and perform full-text searches.
*   **Authentication**: Secure, role-based access control using JWT (JSON Web Tokens).
*   **Media Management**: Simple and efficient file upload and association system.
*   **Performance**: Built on Fiber, one of the fastest Go web frameworks.
*   **API Documentation**: Auto-generated, interactive Swagger documentation.

## Getting Started

### Prerequisites

*   [Go](https://go.dev/dl/) 1.21 or higher

### Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/DenizBitmez/content-flow.git
    cd content-flow
    ```

2.  **Install dependencies**
    ```bash
    go mod download
    ```

### Running the Application

Start the development server:

```bash
go run cmd/server/main.go
```

The server will start at `http://localhost:3000`.

## API Documentation

Interactive API documentation is available via Swagger UI.

1.  Start the application.
2.  Navigate to **[http://localhost:3000/swagger](http://localhost:3000/swagger)**.
3.  **Authentication**:
    *   Register a user via `/api/auth/register`.
    *   Login via `/api/auth/login` to receive a `token`.
    *   Click the **Authorize** button at the top of the Swagger page.
    *   Enter your token in the format: `Bearer <YOUR_TOKEN>`.

## Architecture

*   **Language**: Go (Golang)
*   **Framework**: [Fiber](https://gofiber.io/)
*   **ORM**: [GORM](https://gorm.io/)
*   **Database**: SQLite (Default) / PostgreSQL / MySQL supported via GORM
*   **Authentication**: JWT (golang-jwt)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
