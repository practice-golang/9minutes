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


## Email sending for user verification

* See `9minutes.ini`
* Choose smtp or direct sending
* When use direct sending,
    * you should have own domain
    * IP address of `9minutes` is same as domain's
    * you should learn about following and set
        * `DKIM`
        * `spf record`
        * Also `PTR record` if possible to ask your internet service provider


## Build

* `GOBIN` must be set to `./bin`
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
* Linux, MAC
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
