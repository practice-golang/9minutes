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

Current, Sqlite only


## Todo
- [x] Search
- [x] Edit
- [x] User add, edit, delete
- [x] Comment
- [x] Auth - token reissuing
- [x] File upload
- [x] Sign up
- [ ] Other DB except sqlite
    - [x] MySQL
    - [ ] Postgres
    - [ ] MS-SQL
- [ ] Block duplicate board table name
- [ ] API - list of board columns
- [ ] User page
- [ ] User defined template
- [ ] Add reg/last modified datetime, read count
- [ ] Code cleaning
    - [ ] Password hash
    - [ ] Diet dup codes
- [ ] ~~Email sending~~
- [ ] ~~Menu~~
- [ ] ~~Test~~


## License
[MIT License](http://www.opensource.org/licenses/MIT)
