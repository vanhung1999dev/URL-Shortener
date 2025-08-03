# 🚀 URL Shortener

A simple and efficient URL shortening service built with Go, Redis, and Docker. This project allows users to convert long URLs into compact short links and provides basic analytics, including click tracking, rate limiting, and support for distributed deployment.

---

## 📑 Table of Contents

- [Features](#features)  
- [Requirements](#requirements)  
- [Installation](#installation)  
- [Usage](#usage)  
- [Configuration](#configuration)  
- [Troubleshooting & FAQ](#troubleshooting--faq)  
- [Contributing](#contributing)  
- [License](#license)  
- [Contact / Support](#contact--support)  
- [Acknowledgments](#acknowledgments)

---

## ✨ Features

- Shorten URLs using unique base36 Snowflake IDs
- Prevent duplicate shortening of the same URL
- Rate-limited per IP address (default: 1 request/sec)
- Basic analytics: global and per-shortID click counters
- Redis-backed data storage (2 DBs: main and analytics)
- Distributed system ready (with `MACHINE_ID` support)
- Docker and Docker Compose for easy setup

---

## ✅ Requirements

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- Make (optional, for command shortcuts)

---

## ⚙️ Installation

Clone the repository:

```bash
git clone https://github.com/your-username/url-shortener.git
cd url-shortener
```

Create a `.env` file in the root directory:

```env
# Environment Configuration
DB_ADDR=db:6379
DB_PASS=
APP_PORT=:3000
API_QUOTA=10
MACHINE_ID=1
DOMAIN=http://localhost:3000
```

> ⚠️ `MACHINE_ID` is used for distributed ID generation via Sonyflake. Each instance in a multi-node deployment should have a **unique MACHINE_ID** (between 0–127).

Start the project using Docker Compose:

```bash
make up
```

The service should now be running at `http://localhost:3000`.

---

## 🚀 Usage

### Shorten a URL

**Endpoint:** `POST /api/v1`

**Request body:**

```json
{
  "long_url": "https://example.com"
}
```

**Response:**

```json
{
  "short_url": "http://localhost:3000/abc123"
}
```

### Redirect using a short URL

**Endpoint:** `GET /:shortID`

Example: `GET http://localhost:3000/abc123` → redirects to the original URL.

---

## 🔧 Configuration

The following environment variables can be configured in your `.env` file:

| Variable     | Description                                           | Default                |
|--------------|-------------------------------------------------------|------------------------|
| `DB_ADDR`    | Redis server address                                  | `db:6379`              |
| `DB_PASS`    | Redis password (if any)                               | *(empty)*              |
| `APP_PORT`   | Port for the app to run on                            | `:3000`                |
| `API_QUOTA`  | Max requests per second per IP (rate limit)           | `10`                   |
| `MACHINE_ID` | Unique ID for this node (used in Sonyflake)           | `1`                    |
| `DOMAIN`     | Public domain used in short URL generation            | `http://localhost:3000`|

---

## 🧯 Troubleshooting & FAQ

**Q: I keep getting a rate limit error.**  
A: The app enforces a rate limit per IP (default: 10 requests/second). Wait a second and try again.

**Q: Why does my shortened URL expire?**  
A: URLs are stored in Redis with a TTL of 24 hours. Modify TTL in the code if you want permanent storage.

**Q: Redis not connecting?**  
A: Make sure Docker is running and port `6379` is available. Check logs using:

```bash
docker-compose logs
```

---

## 🤝 Contributing

Contributions are welcome! To contribute:

1. Fork the repo
2. Create a new branch (`git checkout -b feature/your-feature`)
3. Commit your changes
4. Push to your branch (`git push origin feature/your-feature`)
5. Open a pull request

### Code Style

- Format code using `gofmt`
- Keep commit messages meaningful

---

## 📝 License

This project is licensed under the [MIT License](LICENSE).

---

## 📬 Contact / Support

For issues or questions, open an [issue on GitHub](https://github.com/your-username/url-shortener/issues).

---

## 🙏 Acknowledgments

- [Go Fiber](https://gofiber.io/)
- [Redis](https://redis.io/)
- [Sonyflake ID Generator](https://github.com/sony/sonyflake)
- [govalidator](https://github.com/asaskevich/govalidator)
