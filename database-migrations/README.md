# Database Migration (Migrate + SQL)

## 1. Installation :rocket:

For Windows OS can use [scoop](https://scoop.sh/)

```bash
$ scoop install migrate
```

![scoop install migrate](https://user-images.githubusercontent.com/112603532/221208992-58f8a84c-463f-405f-b256-d01c5d8fd70f.png)

For another installation methods, check the following link : [golang-migrate](https://github.com/golang-migrate/migrate)

## 2. Set Environment

```bash
$ export ENV_VAR_NAME='postgres://(username):(password)@(hostname):(port_number)/(database_name)?sslmode=disable'
```

Example :
```bash
$ export POSTGRESQL_URL='postgres://postgres:@Abc428660@localhost:5432/postgres?sslmode=disable'
```

## 3. Run Migration

### Creating New Table

```bash
$ migrate create -ext sql -dir migration_files_path create_your_table_name_table -format
```

Example :

```bash
$ migrate create -ext sql -dir examples create_user_table -format
```

### Migrating Up

```bash
$ migrate -database ${POSTGRESQL_URL} -path migration_files_path up
```

Example : 

```bash
$ migrate -database ${POSTGRESQL_URL} -path examples up
```

### Migrating Down

```bash
$ migrate -database ${POSTGRESQL_URL} -path migrations down
```

### Reference 

- [golang-migrate/migrate](https://github.com/golang-migrate/migrate) 
- [golang-migrate/migrate/postgresql](https://github.com/golang-migrate/migrate/tree/master/database/postgres)
