# Price Tracker

A Go application for tracking prices with EMMSA API integration.

## Features

- EMMSA API client implementation
- Price tracking functionality
- RESTful API endpoints

## Prerequisites

- Go 1.16 or higher
- Git

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/price-tracker.git
   cd price-tracker
   ```

2. Build the application:
   ```bash
   go build -o price-tracker ./cmd/price-tracker
   ```

## Configuration

Create a `.env` file in the root directory with the following variables:

```
API_KEY=your_emmsa_api_key
API_BASE_URL=https://api.emmsa.com/v1
```

## Usage

Run the application:

```bash
./price-tracker
```

The application will start a server on the default port (8080).

## API Endpoints

- `GET /prices` - Get all tracked prices
- `GET /prices/{id}` - Get a specific price by ID
- `POST /prices` - Add a new price to track
- `PUT /prices/{id}` - Update a tracked price
- `DELETE /prices/{id}` - Stop tracking a price

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -ldflags "-s -w" -o price-tracker ./cmd/price-tracker
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
