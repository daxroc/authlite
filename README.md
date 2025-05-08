# AuthLite

AuthLite is a simple lightweight forwardAuth service designed to simplify providing bearer token management for your applications. It provides basic features for managing tokens while maintaining a minimal footprint.

## Features

- Minimal dependencies and easy to integrate.
- Supports Live Reloading

## Getting Started

### Prerequisites

- [Go](https://golang.org/) 1.21 or later
- [Docker](https://www.docker.com/)
- [GitHub Actions](https://github.com/features/actions) for CI/CD

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/authlite.git
   cd authlite

2. Build the project:

```
go build -o authlite
```

3. Run the application:

```
./authlite
```

4. Running Tests
To run the unit tests, use the following command:

```
go test -v ./...
```

Docker

- Unstable Images: Built for pull requests and tagged as unstable.
- Release Images: Built for pushes to the main branch and tagged releases. Images are tagged with both latest and the version number.

