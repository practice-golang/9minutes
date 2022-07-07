# 9minutes

Bulletin board

## Goal

* Create
    * HTML file server
    * Simple Board with user account managing
* Escape from vaporware :-p

## Getting started

* Build and run
* Then run again
* Open `http://localhost:5525` in browser
* Sign in with `admin` / `admin`


## Modify HTML

* When exist html files in storage, 9m load html files in real storage instead of embedded html
* Get embedded html files
```zsh
$ ./9m -get html
```
* Then, edit html files
* Files in admin directory are not exported hence, If you want to modify them, download this repository and modify and build


## API, Route

See `setup.go` and `router_*.go`

* Because of regex, used custom router


## HTML, Static, Upload

See `/embed`, `/html`, `/static`, `/upload`


## DB

* SQLite
* MySQL - MariaDB 10.5
* Postgres - PostgreSQL 12.3
* SQL Server - MS SQL Server express 2014 & 2019 (Not tested yet)


## Todo
* [] Try html/template at board list and content
* [] File attatchment
* [] Shared session
* [] User approval - email sending

## License

[3-Clause BSD](https://opensource.org/licenses/BSD-3-Clause)
