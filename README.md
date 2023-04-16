# Go Mixtape Trading

Go Mixtape Trading is a web application that allows users to trade mixtapes with each other.

## Installation

To install and run the application, you will need to have Go 1.16 or later installed on your system.

1. Clone the repository: `git clone https://github.com/perezdid/go-mixtape-trading.git`
2. Change into the project directory: `cd go-mixtape-trading`
3. Copy the `.env.example` file to `.env` and set the appropriate values for your environment.
4. Build and run the server: `go run cmd/server/main.go`

## Usage

Once the server is running, you can use a web browser or a command-line tool like `curl` to interact with the API. The following endpoints are available:

- `GET /status`: Returns the status of the server.
- `POST /login`: Authenticates a user and returns an access token.
- `GET /user/{id}`: Returns information about a user.
- `GET /playlist/{id}`: Returns information about a playlist.
- `POST /playlist`: Creates a new playlist.
- `PUT /playlist/{id}`: Updates an existing playlist.
- `DELETE /playlist/{id}`: Deletes a playlist.

For more information about each endpoint, refer to the documentation in the `internal/api/handlers` directory.

## Folder Structure

The project follows a standard Go project layout. Here is an overview of the folder structure:

- `cmd/server`: The main entry point for the server.
- `internal/api/handlers`: The API endpoints and request handlers.
- `internal/api/models`: The data models and database access functions.
- `config`: Configuration settings for the application.
- `middleware`: Reusable middleware functions.
- `utils`: Utility functions.

## TLS Configuration

This project uses HTTPS to secure communication between clients and the server. To set up HTTPS, you need to generate a TLS certificate and key file.

### Generating a Self-Signed Certificate and Key

For testing purposes, you can generate a self-signed certificate and key using the `openssl` command-line tool. To generate a new certificate and key, run the following command in a terminal:

```sh
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

This command will create a new RSA private key and a self-signed X.509 certificate that is valid for 365 days.

### Configuring the Server

Once you have generated the TLS certificate and key files, you need to configure the server to use them. In the main.go file, update the tlsConfig variable to load the certificate and key files:

```go
// Load TLS certificate and key
certFile := "cert.pem"
keyFile := "key.pem"
cert, err := tls.LoadX509KeyPair(certFile, keyFile)
if err != nil {
    log.Fatalf("error loading TLS certificate and key: %v", err)
}

tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{cert},
}
```

Make sure to replace cert.pem and key.pem with the paths to your own certificate and key files.

### Production Considerations

For production use, you should obtain a trusted TLS certificate from a certificate authority (CA) instead of using a self-signed certificate. A trusted TLS certificate ensures that client browsers and other software will trust your server's identity, reducing the risk of man-in-the-middle attacks and other security vulnerabilities.

You can obtain a trusted TLS certificate from a CA
