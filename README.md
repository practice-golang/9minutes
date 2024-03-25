# 9minutes

Bulletin board - Small forum application with user account managing

## Getting started
* Download binary or build from source
* Run
```sh
./9minutes
```
* Modify `config.ini` then run again `9minutes` binary
* Open `http://localhost:5525` in browser
* Sign in with initial account - `admin` / `admin`


## HTML modification
* When html files are placed in storage, 9minutes load html files in real storage instead of embedded html
* Get embedded html files by follow command then edit them
```sh
./9minutes -get html
```
* Files in `admin` directory are not exported hence, you should edit source


## API, Route
See `setup.go` and `router.go`


## Database
* Requirement - Account which possible to create/drop database, schema, table
* Support
    * SQLite3
    * MySQL, MariaDB >= 10.3
    * PostgresSQL >= 12.3
    * SQL Server >= 2014
    * Oracle >= 12c


## Email sending

* For the purpose of user verification and password reset
* See `config.ini`
* You can choose `smtp` or `direct sending`
* When use `direct sending`, you should learn about domain, `DKIM`, `spf record`, `PTR record`
    * Also most of cloud service blocked port 25 so you probably can not use it
* Get DKIM files - If not yet have dkim files
```sh
$ ./9minutes -get dkim
```


## Build

```sh
make
```


### Reverse proxy examples

* NginX
```nginx
# https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/
server {
    listen 80;
    listen [::]:80;

    server_name domain.name;

    location / {
        proxy_set_header Host $host;
        proxy_set_header Accept-Encoding "";
        proxy_set_header Real-Ip $remote_addr;
        proxy_pass http://localhost:5525;
    }

    error_page 404 /404.html;
    error_page 500 502 503 504 /50x.html;
}
```

* Caddy
```powershell
# https://caddyserver.com/docs/quick-starts/reverse-proxy
caddy reverse-proxy --from :80 --to 127.0.0.1:5525
```


## License

[3-Clause BSD](https://opensource.org/licenses/BSD-3-Clause)
