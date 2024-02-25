# Upload API

## Description

REST API for uploading files, supporting two formats:

### 1. GridFS:

**Description:**
"GridFS" is a large-data storage system in MongoDB used for storing files that exceed the maximum size that a single file can be stored in regular MongoDB.

**Usage:**
When using "GridFS" for file uploads in the API, clients can directly send large files to MongoDB using the GridFS API.

**Benefits:**
- Efficient storage of large files in MongoDB.
- Can help reduce code redundancy in handling large files.

**More Info:** [GridFS Documentation](https://www.mongodb.com/docs/drivers/go/current/fundamentals/gridfs/)

### 2. FileServer:

**Description:**
"FileServer" is a method of storing files directly on your server, supporting both file uploads and access to files stored in a specified directory or path.

**Usage:**
When using "FileServer," you can receive files sent by clients and save them in a directory or a specified path on your server.

**Benefits:**
- Low friction in implementation and maintenance.
- Flexibility in specifying paths and managing files according to requirements.

## Dependencies

- **Programming Language:** Go
- **API Framework:** [Gin Framework](https://github.com/gin-gonic/gin)
- **Database:** MongoDB
  - Library: [mongo-driver](https://github.com/mongodb/mongo-go-driver)
- **Containerization:** Docker
- **Other Dependencies:**
  - github.com/coolbed/mgo-oid
  - github.com/go-ini/ini
  - github.com/joho/godotenv
  - github.com/spf13/viper
  - github.com/stretchr/testify

## Authors

[chanyapatshell@gmail.com](mailto:chanyapatshell@gmail.com)

## Version History

- 1.0.0
  - Initial Release

## License

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)