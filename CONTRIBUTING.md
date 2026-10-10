# Contributing to DevPulse

Thank you for your interest in contributing to **DevPulse**! We welcome contributions of all kinds: bug fixes, new features, UI enhancements, documentation improvements, and feedback.

## Development Setup

1. **Prerequisites:**
   - Go 1.22 or higher installed.
   - Git installed.

2. **Clone and run:**
   ```bash
   git clone https://github.com/MisterPugs/devpulse.git
   cd devpulse
   go mod tidy
   go run .
   ```

3. **Running Tests & Building:**
   ```bash
   go test ./...
   go build -o devpulse .
   ```

## Commit Guidelines

We recommend using conventional commits:
- `feat: add support for UDP sockets`
- `fix: resolve PID resolution on macOS`
- `docs: update installation instructions`
- `style: refine tab border styles`

## Opening a Pull Request

1. Create a branch: `git checkout -b feature/my-feature`
2. Commit your changes: `git commit -m 'feat: my new feature'`
3. Push to your fork: `git push origin feature/my-feature`
4. Open a Pull Request on GitHub.
