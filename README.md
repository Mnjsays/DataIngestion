# DataIngestion
Overview

The DataIngestion project is designed to fetch data from an external API, process it, and store it in an AWS S3 bucket(Different Storage provide implemenation in progress). 
It also provides an API to retrieve stored data from S3.

Features
Fetch data from an external API.
Store the fetched data in AWS S3 in JSON format.
Retrieve stored data from S3 via an HTTP API.
Environment-based configuration management.

Project Structure

cmd/server/main.go: Entry point of the application. Configures the application and starts the HTTP server.
pkg/dataParser/dataParser.go: Contains logic for fetching data from the external API and retrieving data from S3.
pkg/storage/storage.go: Handles AWS S3 interactions, including reading and writing data.
util/util.go: Utility functions for environment detection, configuration loading, and string sanitization.
types: Contains shared types and configurations used across the project.



API Endpoints
1. Fetch Data from S3
Endpoint: /gets3Data/{filename}
Method: GET
Description: Retrieves the specified file from the S3 bucket.
Path Parameter:
filename: Name of the file to retrieve.
Response:
200 OK: JSON data from the file.
400 Bad Request: If the filename is missing.
500 Internal Server Error: If there is an error fetching or parsing the file.


GET http://localhost:8094/gets3Data/samplefilename


How It Works
Data Fetching:  
The DataFetch function fetches data from the external API (CydresUrl).
The data is processed and stored in an S3 bucket using the AwsStorage function.
Data Retrieval:  
The DataRetriever function retrieves a file from the S3 bucket based on the filename provided in the URL.
AWS S3 Integration:  
The AwsRead function reads files from the S3 bucket.
The AwsStorage function uploads files to the S3 bucket.
Running the Application
Set the DATAENV environment variable to the desired environment (e.g., local).
Ensure the config.yaml file is properly configured.
Run the application:
go run cmd/server/main.go
Access the API at http://localhost:8094.
Error Handling
Logs errors using zap logger.
Returns appropriate HTTP status codes for API errors.
Dependencies
Gorilla Mux: For HTTP routing.
AWS SDK for Go: For S3 interactions.
Zap: For structured logging.
YAML.v3: For configuration parsing.
Future Enhancements
Add authentication for API endpoints.
Implement retries for S3 operations.
Add unit tests for critical functions.