# 9minutes board

```
An expreimental 9minutes or GUBUN board for own personal use

Using struct or map, mainly map

WIP
```

## Goal

Escape from a vaporware


## Getting started

* Build and run
* Then stop and edit `9minutes.ini`
* Then run again
* Open `http://localhost:2510` in browser
* Sign in with `admin` / `admin`


## REST API

See `requests-admin.http` and `requests-contents.http` on `vscode` with `rest-client` extension

`GET` and `POST` make me confused so, APIs are not `RESTful`

* Read, List - `GET`, `POST`
* Write - `PUT`
* Edit - `PATCH`
* DELETE - `DELETE`


## HTML

See `html` files in `/static` as pages and `/templates` as templates


## DB

Sqlite, MySQL, Postgres, SQL Server
* SQLite - Using non-cgo package https://gitlab.com/cznic/sqlite
* MySQL - Tested MariaDB 10.5
* Postgres - Tested PostgreSQL 12.3
* SQL Server - Tested MS SQL Server express 2014 & 2019


## Todo
- [x] Search
- [x] Edit
- [x] User add, edit, delete
- [x] Comment
- [x] Auth - token reissuing
- [x] File upload
- [x] Sign up
- [x] Other DB except sqlite
    - [x] MySQL
    - [x] Postgres
    - [x] MS-SQL
- [ ] Bug
     - [x] Correct comment table renaming
     - [ ] Remove writer from search target
- [ ] API - Add get list of board columns
- [ ] Code cleaning
    - [ ] Password hash
    - [ ] Diet dup codes
- [ ] User defined template
- [ ] Block duplicate board table name
- [ ] Add reg/last modified datetime, read count
- [ ] ~~Test~~
- [ ] ~~File downloader~~
- [ ] ~~Delete elete post with attached files~~
- [ ] ~~Email sending~~
- [ ] ~~Menu~~
- [ ] ~~User page~~


## License
[MIT License](http://www.opensource.org/licenses/MIT)
