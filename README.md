# Data Ingestion
Overview

->The DataIngestion project is designed to fetch data from an external API, process it, and store it in an AWS S3 bucket(Different Storage provide implemenation in progress). 
->It also provides an API to retrieve stored data from S3.

Features
->Fetch data from an external API.
->Store the fetched data in AWS S3 in JSON format.
->Retrieve stored data from S3 via an HTTP API.
->.env-based configuration management.

Project Structure

->cmd/server/main.go: Entry point of the application.
->Configures the application and starts the HTTP server.
->pkg/dataParser/dataParser.go: Contains logic for fetching data from the external API and retrieving data from S3.
->pkg/storage/storage.go: Handles AWS S3 interactions, including reading and writing data.
->util/util.go: Utility functions for environment detection, configuration loading, and string sanitization.
->types: Contains shared types used across the project.



API Endpoints
1.  Fetch Data from S3
    Endpoint: /gets3Data/{filename}
    Method: GET
    Description: Retrieves the specified file data from the S3 bucket.
    Path Parameter:
    filename: Name of the file to retrieve.

GET http://localhost:8094/gets3Data/samplefilename

Data Retrieval:  
->The DataRetriever function retrieves a file from the S3 bucket based on the filename provided in the URL.
->AWS S3 Integration:  

Running the Application

Run the application:
->go run cmd/server/main.go

Dependencies
->Gorilla Mux: For HTTP routing.
->AWS SDK for Go: For S3 interactions.
->Zap: For structured logging.

