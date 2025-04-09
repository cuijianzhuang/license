# License Management Service

A Go-based service for managing software license validation and authentication.

![GitHub commit activity](https://img.shields.io/github/commit-activity/t/nannanStrawberry314/license?color=blue)
![GitHub forks](https://img.shields.io/github/forks/nannanStrawberry314/license?style=flat&color=brightgreen)
![GitHub stars](https://img.shields.io/github/stars/nannanStrawberry314/license?color=orange)
![GitHub pull requests](https://img.shields.io/github/issues-pr/nannanStrawberry314/license?color=red)
![Docker Pulls](https://img.shields.io/docker/pulls/raspberrycheese/license?color=blueviolet)

[English](README.md) | [中文](README_CN.md)

## Features

- License generation and validation for various software products
- Support for JetBrains products, GitLab, FinalShell, MobaXterm, and JRebel
- RESTful API interface built with Gin framework
- Scheduled tasks with cron
- Database storage with GORM (MySQL/SQLite support)
- Secure encryption using RSA

## Requirements

- Go 1.21 or higher
- MySQL database (or SQLite for development)
- Docker (optional, for containerized deployment)

## Installation

### Option 1: Direct Installation

1. Clone the repository
   ```
   git clone https://github.com/nannanStrawberry314/license.git
   cd license
   ```

2. Install dependencies
   ```
   go mod download
   ```

3. Configure environment variables (copy .env.example to .env and modify as needed)

4. Build and run
   ```
   go build -o license-server
   ./license-server
   ```

### Option 2: Docker Deployment

1. Build the Docker image
   ```
   docker build -t license-server .
   ```

2. Run using docker-compose
   ```
   docker-compose up -d
   ```

## Configuration

Configuration is handled through environment variables and the `.env` file:

- `HTTP_HOST`: The host address to bind the server
- `HTTP_PORT`: The port to listen on
- `DB_TYPE`: Database type (mysql or sqlite)
- `DB_DSN`: Database connection string

## API Endpoints

The API provides various endpoints for license management:

- `POST /v1/generate`: Generate a new license
- `POST /v1/validate`: Validate an existing license
- `GET /v1/status`: Check the service status

Refer to the API documentation for detailed usage instructions.

## Development

To contribute to this project:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

This project is proprietary software. All rights reserved. 

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=nannanStrawberry314/license&type=Date)](https://www.star-history.com/#nannanStrawberry314/license&Date) 