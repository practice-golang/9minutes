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
- [ ] Comment - edit, delete
- [x] User add, edit, delete
- [ ] User auth
- [ ] Comment paging
- [ ] Code cleaning
- [ ] Add reg/last modified datetime, read count
- [ ] User defined template
- [ ] Other DB except sqlite
- [ ] ~~Email~~
- [ ] ~~Menu~~
- [ ] ~~Test~~


## License
[MIT License](http://www.opensource.org/licenses/MIT)
