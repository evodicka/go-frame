# Post-Task Actions

Whenever a task is completed, several steps should be taken to ensure project integrity:

1.  **Format Code**: Run `go fmt ./...` for backend changes.
2.  **Run Tests**: Execute backend tests in `cmd/go-frame-app` to verify no regressions.
3.  **Verify Build**: Ensure the application still builds using `go build ./cmd/go-frame-app`.
4.  **Frontend Build**: If frontend changes were made, verify them by running `npm run build` in the `web` directory.
