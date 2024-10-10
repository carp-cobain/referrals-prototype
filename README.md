# Referrals

A prototype web-service created to explore capturing signups for referral campaigns.

**Replication setup**

[Install](https://litestream.io/install/) litestream

Run a local minio server in docker

```sh
docker run -d --name minio-server -v "$(pwd)"/data:/data/ -p 9000:9000 -p 9001:9001 minio/minio server /data --console-address ":9001"
```

Go to http://localhost:9001/login

- user = `minioadmin`
- pass = `minioadmin`

Then, create a bucket: `referrals`

**Running with replication**

```sh
source scripts/envrc
make && make exec
```

See `Makefile` for `litestream` replication command details

### References:

- [Gin](https://gin-gonic.com)
- [Gorm](https://gorm.io)
- [SQLite](https://www.sqlite.org)
- [Litestream](https://litestream.io)
- [MinIO](https://min.io/)
