# Go Todo App

A simple, responsive Todo application built with Go, Chi router, SQLite, and Templ templating engine.

## Features

- Create, read, update, and delete todo items
- Toggle completion status of todos
- Responsive design that works on mobile and desktop
- Clean, modern UI with gradient styling

## Technologies Used

- Go 1.24+
- Chi router for HTTP routing
- SQLite for data storage
- Templ for HTML templating
- Pure CSS for styling (no frameworks)

## Getting Started

### Prerequisites

- Go 1.24 or higher
- SQLite

### Installation

1. Clone the repository
   ```
   git clone https://github.com/yourusername/go-crud-app.git
   cd go-crud-app
   ```

2. Install dependencies
   ```
   go mod download
   ```

3. Generate template files
   ```
   go generate ./...
   ```

4. Run the application
   ```
   go run main.go
   ```

5. Open your browser and navigate to `http://localhost:3000`

## Project Structure

- `main.go` - Application entry point and router setup
- `handlers/` - HTTP request handlers
- `models/` - Data models
- `views/` - Templ templates
- `database/` - Database setup and initialization

## License

MIT