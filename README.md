# 9minutes

Bulletin board

## Goal

* Create
    * HTML file server
    * Simple Board with user account managing
* Escape from vaporware :-p

## Getting started

* Build and run
```sh
$ go build .
$ ./9minutes
```
* Then run again after modify 9minutes.ini
```sh
$ ./9minutes
```
* Open `http://localhost:5525` in browser
* Sign in with `admin` / `admin`


## HTML modification

* When html files are placed in storage, 9minutes load html files in real storage instead of embedded html
* Get embedded html files
```sh
$ ./9minutes -get html
```
* Then, edit html files
* Files in admin directory are not exported hence, If you want to modify them, download this repository and modify and build


## API, Route

See `setup.go` and `router_*.go`

* Because of regex, used custom router


## DB

* SQLite
* MySQL - MariaDB 10.5
* Postgres - PostgreSQL 12.3
* SQL Server - MS SQL Server express 2014 & 2019 (Not tested yet)


## Todo
* [ ] Try html/template at board list
* [x] File attatchment
* [x] Shared session
* [ ] User approval - email sending

## License

[3-Clause BSD](https://opensource.org/licenses/BSD-3-Clause)
