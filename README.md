# KUMPARAN SENIOR BACKEND TEST
API KUMPARAN CREATED BY IBAMAULANA MARET 2020

## SETUP GUIDE
** Copy `env.yaml.example` into `env.yaml` and `dbconfig.yml.example` into `dbconfig.yml`
** Please setup `sql connection` in `dbconfig.yml` and do migrate using ```sql-migrate up -env="dev"```
** Please setup `sql connection`, `elasticsearch connection`, and `redis port` in `env.yaml`
** Please run redis before execute main.go
** If got error from `github.com/adjust/rmq` please rewrite on `github.com/adjust/rmq/connection.go` on line 176 from `connection.redisClient.flushDb()` to `connection.redisClient.FlushDb()`

