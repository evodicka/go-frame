# Project Overview

## Purpose
Go-Frame is a digital photo frame application designed to display images in a continuous loop with configurable intervals. It features a lightweight Go backend and an Angular-based image display.

## Tech Stack
- **Backend**: Go (v1.25), Gin Gonic (HTTP Web Framework), BoltDB (Embedded KV Store).
- **Frontend**: Angular (v13), TypeScript.

## Base Structure
- `cmd/go-frame-app/`: Main application source code (backend).
- `web/`: Angular frontend source code.
- `scripts/`: Build and utility scripts.
- `images/`: Local storage for uploaded image files.
