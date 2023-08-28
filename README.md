# 9minutes

Bulletin board


## Goal

* Create
    * HTML file server
    * Simple Board with user account managing
* Escape from vaporware :-p


## Getting started

* Download - See `Sample bins` at `Release`

* Run
```sh
$ ./9minutes
```
* Modify `config.ini` then run again `9minutes` binary
* Open `http://localhost:5525` in browser
* Sign in with `admin` / `admin`


## HTML modification

* When html files are placed in storage, 9minutes load html files in real storage instead of embedded html
* Get embedded html files
```sh
$ ./9minutes -get html
```
* Then, edit html files
* `Files in admin directory` are not exported hence, if you want to modify them, clone this repository then modify and build


## API, Route

See `setup.go` and `router_*.go`

* Because of regex, used custom router


## DB

* SQLite3 - Tested
* MySQL - Tested with MariaDB 10.5
* Postgres - Tested with PostgreSQL 12.3
* SQL Server - Tested with MS SQL Server express 2014 & 2019
* Oracle - Tested with 12c as local, 19c as autonomous database of oracle cloud. 11g or before not support


## File access priority of router

* When router request triggered, router find and load the html file as following table

|        | embed | Find fs first |
|--------|-------|---------------|
| embed  | yes   | no            |
| html   | yes   | yes           |
| static | no    | yes           |


## Email sending

* For the purpose of user verification and password reset
* See `config.ini`
* You can choose `smtp` or `direct sending`
* When use `direct sending`, you should learn about domain, `DKIM`, `spf record`, `PTR record`
    * Also most of cloud service blocked port 25 so you probably can not use it
* Get DKIM key files - If not yet generate dkim keys
```sh
$ ./9minutes -get dkim
```


## Build

* Set `GOBIN` to `./bin`
* Windows
```
$ mingw32-make.exe
```
* Linux, Mac
```
$ make
```


## Etc.

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

### Heroku

* Append following variables to environment setting of dyno
* Oracle wallet is not supported
```
PORT
DATABASE_TYPE
DATABASE_ADDRESS
DATABASE_PORT
DATABASE_PROTOCOL -> tcp
DATABASE_NAME
DATABASE_ID
DATABASE_PASSWORD
```


## License

[3-Clause BSD](https://opensource.org/licenses/BSD-3-Clause)
