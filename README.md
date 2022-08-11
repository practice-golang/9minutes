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
* Then run again after modify `9minutes.ini`
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
* SQL Server - MS SQL Server express 2014 & 2019


## Paas - Heroku

* Append following variables to environment setting in dyno
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


## Email sending

* For the purpose of user verification and password reset
* See `9minutes.ini`
* You can choose `smtp` or `direct sending`
* When use `direct sending`, you should learn about domain, `DKIM`, `spf record`, `PTR record`
    * Also most of cloud service blocked port 25 so probably you can not use it


## Build

* Set `GOBIN` to `./bin`
* Windows
    * build
    ```
    $ mingw32-make.exe
    ```
    * test
    ```
    $ mingw32-make.exe test
    ```
    * build for all platform
    ```
    $ mingw32-make.exe dist
    ```
* Linux, Mac
    * build
    ```
    $ make
    ```
    * test
    ```
    $ make test
    ```
    * build for all platform
    ```
    $ make dist
    ```


## License

[3-Clause BSD](https://opensource.org/licenses/BSD-3-Clause)
