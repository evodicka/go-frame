# Suggested Commands

## Building
- **Build UI**: `cd web && npm install && npm run build`
- **Build Backend**: `go build ./cmd/go-frame-app`

## Running
- **Start Application**: `./go-frame-app` (runs on port 8080 by default)

## Testing
- **Run Backend Tests**: `cd cmd/go-frame-app && go test -coverpkg=./... -coverprofile=coverage.out ./...`
- **Check Coverage**: `go tool cover -func=coverage.out` (after running tests)

## Utils
- **List Files**: `ls -R`
- **Search**: `grep -r "pattern" .`
