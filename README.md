# 📊 Data Ingestion

## 🔍 Overview

The **DataIngestion** project is designed to:

* Fetch data from an external API
* Transform the data with metadata
* Store it in an **AWS S3** bucket in JSON format
* Provide an **HTTP API** to retrieve stored files from S3

---

## 🚀 Features

* ✅ Fetches data from an external API (`CYDRES_URL`)
* ✅ Stores the transformed data as JSON files in AWS S3
* ✅ Retrieves stored data via a RESTful HTTP API
* ✅ Uses `.env`-based configuration for flexibility

---

## 📁 Project Structure

| Path                           | Description                                            |
| ------------------------------ | ------------------------------------------------------ |
| `cmd/server/main.go`           | Entry point; configures and starts the HTTP server     |
| `pkg/dataParser/dataParser.go` | Fetches data from external API; retrieves data from S3 |
| `pkg/storage/storage.go`       | Handles AWS S3 interactions (read/write)               |
| `util/util.go`                 | Utility functions (env detection, config, sanitizing)  |
| `types/`                       | Shared type definitions                                |

---

## 🧪 API Documentation

### GET `/gets3Data/{filename}`

**Description:**
Retrieves a file from S3 and returns its contents as JSON.

**Path Parameters:**

* `filename` (string) — Name of the file in the S3 bucket

**Responses:**

| Status                      | Description                   | Example JSON                                   |
| --------------------------- | ----------------------------- | ---------------------------------------------- |
| `200 OK`                    | Success                       | *(Returns file contents)*                      |
| `400 Bad Request`           | Missing or invalid filename   | `{ "error": "Filename is required" }`          |
| `500 Internal Server Error` | File parse or S3 read failure | `{ "error": "Failed to parse file contents" }` |

---

## 🔄 Data Flow

### 1. 🔃 **Data Fetching Logic**

* **Source:** External API defined in `CYDRES_URL`
* **Transformation:**

  * Wraps raw data in a `models.Posts` structure
  * Adds:

    * `IngestedAt` → current timestamp
    * `Source` → e.g. `"PlaceHolderAPI"`

### 2. ☁️ **Data Ingestion to AWS S3**

* Transforms data to JSON
* Saves the JSON to S3 with a **timestamp-based sanitized filename**

---

## ⚙️ Running the Application

```bash
go run cmd/server/main.go
```
Once the application starts:

It fetches data from the external API defined in CYDRES_URL

Transforms and uploads the data as a JSON file to your configured S3 bucket

Logs the uploaded filename (without .json extension) with timestamp to dataingestion.log. Use this filename in the GET /gets3Data/{filename} API call to retrieve the stored data

---

## 📆 Dependencies

* [Gorilla Mux](https://github.com/gorilla/mux) — HTTP routing
* [AWS SDK for Go](https://aws.github.io/aws-sdk-go-v2/) — S3 operations
* [Zap](https://github.com/uber-go/zap) — Structured logging

---

## ❓ Q\&A

### 1. ❓ What would you improve if you had more time?

* Add **Elasticsearch + Kibana** for real-time log search and visualization
* Integrate **Jenkins CI/CD** for automated build and deployment pipelines
* Move credentials to **AWS Secrets Manager** or Kubernetes **Secrets + ConfigMaps**

---

### 2. ❓ What were the hardest parts to implement and why?

* **AWS S3 integration**: Initially challenging due to unfamiliarity with AWS SDK and credential management

---

### 3. ❓ What trade-offs did you consider?

| Aspect      | Trade-off Made                     | Ideal Alternative                          |
| ----------- | ---------------------------------- | ------------------------------------------ |
| Logging     | File + stdout logging only         | Ship logs to centralized tools like ELK    |
| Storage     | S3 for simplicity and durability   | DB or Elastic if indexing/query is needed  |
