# gator
BootDev guilded project for a RSS feed aggregator

had an error with Goose not finding the migration files. A certain bear was leading me on a wild goose chase until we discovered that I had to use -dir
~/go/bin/goose -v -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" status