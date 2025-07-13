# Data Ingestion
Overview

->The DataIngestion project is designed to fetch data from an external API, process it, and store it in an AWS S3 bucket
->It also provides an API to retrieve stored data from S3.

Features
->Fetch data from an external API.
->Store the fetched data in AWS S3 in JSON file format.
->Retrieve stored data from S3 via an HTTP API.
->.env-based configuration management.

Project Structure

->cmd/server/main.go: Entry point of the application.
->Configures the application and starts the HTTP server.
->pkg/dataParser/dataParser.go: Contains logic for fetching data from the external API and retrieving data from S3.
->pkg/storage/storage.go: Handles AWS S3 interactions, including reading and writing data.
->util/util.go: Utility functions for environment detection, configuration loading, and string sanitization.
->types: Contains shared types used across the project.



API Documentation

1. GET /gets3Data/{filename}

Description: Retrieves a file from S3 and returns its contents as JSON.
Path Parameters:
    filename (string): The name of the file to retrieve from S3.
Responses:
    200 OK: Returns the file contents as JSON.
    400 Bad Request: If the filename is missing.
    { "error": "Filename is required" }
    500 Internal Server Error: If the file  cannot be parsed.
    { "error": "Failed to parse file contents" }

2. Data Fetching Logic

Source: Fetches data from the external API (CYDRES_URL).
    Transformation:
    Converts the API response into a models.Posts structure.
    Adds metadata such as IngestedAt (current timestamp) and Source (e.g., "PlaceHolderAPI").

3. Data Ingestion to S3
Logic:
    Transforms the models.Posts structure into JSON.
    Writes the JSON data to an S3 bucket under a sanitized timestamp-based filename.



Running the Application

Run the application:
->go run cmd/server/main.go

Dependencies
->Gorilla Mux: For HTTP routing.
->AWS SDK for Go: For S3 interactions.
->Zap: For structured logging.


Q&A
1.What would you improve if you had more time?
    -To add elastic and kibana support to visualize the logs and monitor.
    -Jenkins support to continous developmemnt and deployment 
2.What were the hardest parts to implement and why?
    -Aws s3 storage part 
    -Prior knowledge of s3 api was missing 
3.What trade-offs did you consider?
    Hard coded s3 credentials , use of secrets and configMaps would have been a great option
