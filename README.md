# Go-Frame

Go-Frame is a digital photo frame application designed to display images in a continuous loop with configurable intervals. It features a lightweight Go backend and a dedicated Angular-based image display.

## Key Features

- **Automatic Image Cycling**: Rotates through images stored in the database.
- **Configurable Intervals**: Adjust the duration each image is displayed.
- **RESTful Management API**: Administer the system via API endpoints for:
  - Uploading new images.
  - Deleting existing images.
  - Reordering the display sequence.
  - Updating global configurations.
- **Embedded Web UI**: An Angular frontend is embedded within the Go binary for seamless deployment and image presentation.
- **Lightweight Persistence**: Uses BoltDB for fast and simple metadata storage.

## Technologies

- **Backend**: [Go](https://go.dev/) (v1.16+), [Gin Gonic](https://gin-gonic.com/) (HTTP Web Framework), [BoltDB](https://go.etcd.io/bbolt) (Embedded KV Store).
- **Frontend**: [Angular](https://angular.io/) (v13), TypeScript.

## Project Structure

- `cmd/go-frame-app/`: Main application source code.
- `web/`: Angular frontend source code.
- `scripts/`: Build and utility scripts.
- `images/`: Local storage for uploaded image files.

## Setup and Installation

### Prerequisites

- **Go**: v1.16 or higher.
- **Node.js & npm**: For building the frontend.

### Build and Run

1.  **Build the Frontend**:
    Navigate to the `web` directory and build the Angular application. This will output the compiled assets to the backend's web directory for embedding.
    ```bash
    cd web
    npm install
    npm run build
    ```
    *(Note: The build script is configured to output directly to `../cmd/go-frame-app/web`)*.

2.  **Build the Backend**:
    Navigate to the application root or `cmd/go-frame-app` and build the Go binary.
    ```bash
    go build ./cmd/go-frame-app
    ```

3.  **Run the Application**:
    Start the server. By default, it runs on port `8080`.
    ```bash
    ./go-frame-app
    ```
    Access the image display at `http://localhost:8080`.

## API Documentation

The management API is accessible under the `/admin/api` prefix. Key endpoints include:

- `GET /admin/api/image`: List all images.
- `POST /admin/api/image`: Upload a new image.
- `PUT /admin/api/image`: Update image display order.
- `DELETE /admin/api/image/:id`: Remove an image.
- `GET /admin/api/configuration`: Retrieve current config.
- `PUT /admin/api/configuration`: Update configuration.

The public image data is available at:
- `GET /api/image/current`: Get the currently active image metadata.

## Running Tests

To run the full test suite for the backend application, including unit and integration tests, run the following command from the `cmd/go-frame-app` directory:

```bash
cd cmd/go-frame-app
go test -coverpkg=./... -coverprofile=coverage.out ./...
```

To view the coverage report:

```bash
go tool cover -func=coverage.out
```

This ensures that all packages (`persistence`, `api`, `admin-api`, `static`) are tested and coverage is tracked across boundaries.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
