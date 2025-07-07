# Gator RSS Aggregator

A command-line RSS feed aggregator built in Go. Gator allows you to subscribe to RSS feeds, automatically fetch new posts in the background, and browse them right from your terminal.

*Built as part of the [Boot.dev](https://boot.dev) "Learn Go" course with guidance from Boots, the sagely wizard bear teaching assistant.*

## Features

- **Multi-user support** - Multiple users can manage their own feeds
- **Background aggregation** - Automatically fetches new posts at configurable intervals
- **Feed management** - Add, follow, unfollow, and list RSS feeds
- **Post browsing** - View recent posts from your subscribed feeds
- **PostgreSQL storage** - Reliable database storage for feeds and posts

## What it does

Gator helps you stay up-to-date with your favorite RSS feeds without leaving the command line. Add feeds from blogs, news sites, or any RSS-enabled website, then let Gator continuously fetch new posts in the background while you work.

## Prerequisites

Before you can run Gator, you'll need:

### Go
- **Go 1.21+** - [Download and install Go](https://golang.org/dl/)
- Verify installation: `go version`

### PostgreSQL
- **PostgreSQL 12+** - [Download PostgreSQL](https://www.postgresql.org/download/)
- Create a database for Gator:
  ```sql
  CREATE DATABASE gator;
  ```
Note your connection details (host, port, username, password, database name)

### Database Migration Tool
- **goose** - For running database migrations
  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```

### Environment Setup
- Set your database connection string as an environment variable:
  ```bash
  export DATABASE_URL="postgres://username:password@localhost:5432/gator?sslmode=disable"
  ```
  Or on Windows:
  ```cmd
  set DATABASE_URL=postgres://username:password@localhost:5432/gator?sslmode=disable
  ```

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/gator.git
   cd gator
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run database migrations:**
   ```bash
   goose -dir sql/schema postgres $DATABASE_URL up
   ```

4. **Build the application:**
   ```bash
   go build -o gator
   ```

5. **Verify installation:**
   ```bash
   ./gator --help
   ```

## Usage

### Getting Started

1. **Register a user:**
   ```bash
   ./gator register john
   ```

2. **Add your first feed:**
   ```bash
   ./gator addfeed "TechCrunch" "https://techcrunch.com/feed/"
   ```

3. **Start the aggregator (in a separate terminal):**
   ```bash
   ./gator agg 30s
   ```
   This will fetch new posts every 30 seconds. Leave this running in the background.

4. **Browse your posts:**
   ```bash
   ./gator browse 5
   ```

### Available Commands

#### User Management
- `register <username>` - Register a new user and set as current user
- `login <username>` - Switch to an existing user
- `reset` - Reset the database (removes all users, feeds, and posts)

#### Feed Management
- `addfeed <name> <url>` - Add a new RSS feed and automatically follow it
- `feeds` - List all available feeds in the system
- `follow <url>` - Follow an existing feed
- `following` - List feeds you're currently following
- `unfollow <url>` - Unfollow a feed

#### Content Management
- `agg <duration>` - Start the aggregator to fetch posts (e.g., `30s`, `1m`, `5m`)
- `browse [limit]` - Browse recent posts from your followed feeds (default limit: 2)