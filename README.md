# Sneaker Collector App

The Sneaker Collector App is a Go application that scrapes sneaker information from various websites and stores the data in a PostgreSQL database. It also provides a simple API to retrieve the collected data and log entries. The brands included are:

- Nike
- Adidas
- New Balance

Data collected currently includes:

- Name
- Price
- Link
- Brand

## Table of Contents

- [Requirements](#requirements)
- [Getting Started](#gettingstarted)
- [Customization](#customization)
- [License](#license)

## Requirements

- Docker and Docker Compose
- Go 1.20+
- Postgres and PgAdmin

## Getting Started

1. Clone this repository to your local machine:

   ```sh
   git clone https://github.com/yourusername/sneaker-collector.git
   cd sneaker-collector
   ```

Create a content.json file at the root of the project directory containing your initial sneaker data (only applicable for Nike API).

2. Open the [docker-compose](docker-compose.yml) file and update the environment variables for the sneaker-db service to match your desired PostgreSQL settings.

3. Build and start the containers using Docker Compose:

   ```sh
   docker-compose up --build
   ```

This will start the PostgreSQL database, pgAdmin for managing the database.

4. Build the sneaker app:

   ```sh
   go build
   ```

This will build you Sneaker app then you open it.

5. The app will propmpt you to input the user, password and host.You have 2 tries before the app exits.

6. Access the Sneaker Collector App API by opening a web browser or using a tool like curl:

To get Shoes for a selected brand:

    http://localhost:port/protected?action=latest_shoes&brand=selected_brand

To get sneaker data:

    http://localhost:port/protected?action=sneaker_db_data

To trigger data refresh:

    http://localhost:port/protected?action=refresh_data

Access the pgAdmin web interface by navigating to http://localhost:8080 in your browser. Log in using the credentials defined in the [docker-compose](docker-compose.yml) file.

## Customization

To modify the scraping logic and data sources, update the appropriate functions in the scrapper package.

To customize the API behavior, update the handlers in the [main](main.go) Go file.

To change the database schema, update the SQL statements in the database module.

## License

This project is licensed under the [MIT License](LICENSE).
