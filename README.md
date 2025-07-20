
# RSS Feed Parser
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Description

A simple Go REST API service to import, parse, and manage RSS feeds.

## Installation

### Prerequisites

[Docker](https://www.docker.com/products/docker-desktop/) - required to run the application and a PostgreSQL database container

### Clone the Repository

```bash
git clone https://github.com/marchuknikolay/rss-parser.git
cd rss-parser
```

### Run the Application

```bash
docker-compose up
```

### Access the Application

Once the app is running, the server will be accessible at:

```bash
http://localhost:8080/
```
## API Reference

### Channels

#### Get All Channels

```http
GET /channels/
```

Returns a list of all RSS channels.

---

#### Import RSS Feed

```http
POST /channels/
```

Imports a new RSS feed.

**Request Body:**

| Parameter   | Type     | Description                |
|-------------|----------|----------------------------|
| `url`       | `string` | **Required**. RSS feed URL |

---

#### Get Items by Channel ID

```http
GET /channels/${id}/
```

Returns items belonging to the specified channel.

| Parameter | Type   | Description                  |
|-----------|--------|------------------------------|
| id        | int    | **Required**. Channel ID     |

---

#### Update Channel

```http
PUT /channels/${id}/
```

Updates the specified channel data.

| Parameter | Type   | Description                  |
|-----------|--------|------------------------------|
| id        | int    | **Required**. Channel ID     |

**Request Body (JSON):**

```json
{
  "title": "New title",
  "language": "New language",
  "description": "New description"
}
```

| Parameter       | Type   | Description                   |
|-----------------|--------|-------------------------------|
| title           | string | New title for the channel     |
| language        | string | Language code (e.g., "en")    |
| description.    | string | Description of the channel    |

---

#### Delete Channel

```http
DELETE /channels/${id}/
```

Deletes the specified channel.

| Parameter | Type | Description              |
|-----------|------|--------------------------|
| id        | int  | **Required**. Channel ID |

---

### Items

#### Get All Items

```http
GET /items/
```

Returns a list of all RSS items.

---

#### Get Item by ID

```http
GET /items/${id}/
```

Returns a single item by ID.

| Parameter | Type | Description              |
|-----------|------|--------------------------|
| id        | int  | **Required**. Item ID    |

---

#### Update Item

```http
PUT /items/${id}/
```

Updates the specified item data.

| Parameter | Type | Description              |
|-----------|------|--------------------------|
| id        | int  | **Required**. Item ID    |

**Request Body (JSON):**

```json
{
  "title": "New title",
  "description": "New description",
  "pub_date": "New publication date"
}
```

| Parameter       | Type   | Description                            |
|-----------------|--------|----------------------------------------|
| title           | string | New title of the item                  |
| description     | string | New description                        |
| pub_date        | string | Publication date in ISO format (YYYY-MM-DDTHH:MM) |

---

#### Delete Item

```http
DELETE /items/${id}/
```

Deletes the specified channel.

| Parameter | Type | Description              |
|-----------|------|--------------------------|
| id        | int  | **Required**. Item ID    |

## Packages
- [Echo](https://pkg.go.dev/github.com/labstack/echo/v4@v4.13.4) - Go web framework
- [pgx](https://pkg.go.dev/github.com/jackc/pgx/v5@v5.7.5) - PostgreSQL Driver and Toolkit
- [goose](https://pkg.go.dev/github.com/pressly/goose@v2.7.0+incompatible) - database migration tool
- [envdecode](https://pkg.go.dev/github.com/joeshaw/envdecode) - populating structs from environment variables

## License

[MIT](https://opensource.org/license/mit)