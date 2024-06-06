# Lazymigrate

![logo](readme/logo.png)

Terminal ui for managing migrations based on [golang-migrate](https://github.com/golang-migrate/migrate)

![demo](readme/demo.gif)

### Configuration

Now lazymigrate supports PostgreSQL, MySQL and sqlite

lazymigrate can guess connection string if you have specific fields in __.env__ file

[examples](#Examples)

Also it parses env variables from __.lazymigrate__ file

You can set up config directly with variables:

```
LAZYMIGRATE_URL=postgres://root:root@localhost:5432/root?sslmode=disable
LAZYMIGRATE_SOURCE=migrations
```

### Examples

Naming based on docker images envs naming

Postgres

```
POSTGRES_USER=root
POSTGRES_PASSWORD=root
POSTGRES_HOST=localhost
POSTGRES_DB=root
POSTGRES_PORT=5432
```

MySQL

```
MYSQL_USER=root
MYSQL_PASSWORD=root
MYSQL_HOST=localhost
MYSQL_DATABASE=root
MYSQL_PORT=3306
```

Sqlite

```
SQLITE_DB=app.db
```