# Todo

* [x] Control edit/delete button for author only
* [ ] Grant,confirm routine for guest
    * [ ] Control by ip address if guest & table column
        * [x] Add IP_ADDRESS to each db table creation functions
        * [ ] Edit password for guest
    * [x] Appear guest board in list selector
    * [ ] Edit/Delete topic
    * [ ] Edit/Delete comment
    * [ ] Write use password if guest
    * [ ] Edit use password if guest

* [ ] Logger
* [x] Rename all post, posting to topic
* [ ] Separate upload FILES and IMAGES
    * [ ] Create Image upload api - copy from upload api
    * [ ] Add image upload path for web editor
* [ ] Add user login history - useridx, ip, regdate
* [ ] Board managing page - Add manager account setting for each boards
    * [ ] API
        * [ ] Edit board settings
        * [ ] Add/edit/remove user as member -> Add manager, member column to board setting table
    * [ ] DB
    * [ ] Admin page
* [ ] Resurrect session storage - etcd, redis..
* [ ] Add like/dislike button
    * [ ] topic
    * [ ] comment

* [x] Change to go-fiber
* [x] Admin - Use Svelte component for Page, HTML
* [x] Board list generation for menu items
* [x] Comment edit, Delete
    * [x] Comment edit
    * [x] Comment delete
* [x] All board table names to lowercase
* [x] Editor
    * [x] Multiple upload - https://github.com/nhn/tui.editor/issues/1401#issuecomment-785557945
* [x] Board list cache
    * [x] When execute
    * [x] When create
    * [x] When change
* [x] User grade from numbers
    * [x] API - User grade list
    * Refs. - consts/contant.go

* [ ] Banned user time count - When page open or login, check the ban time is gone
* [ ] Clean up HandleHTML function
* [ ] escape/unscape comment writing/reading
* [ ] Move and clean up folders under static - admin, html, static..
* [ ] Show/Edit default and custom column in mypage
* [ ] Choose use or delete approval column
* [ ] User list
    * [ ] Add user custom column not showing in admin page
    * [ ] Edit user custom column not showing in admin page
* [ ] Dump
    * [ ] DB dump, migration - Need to try to taste [goose](https://github.com/pressly/goose), [sql-migrate](https://github.com/rubenv/sql-migrate), [migrate](https://github.com/golang-migrate/migrate), [dbmate](https://github.com/amacneil/dbmate)
    * [ ] files dump, migration - No idea yet
* [ ] Make user column definitions move to up or down
* [ ] Add content type in admin page - Page
* [ ] Rearrange route paths

* Waive
    * Board
        * ~~Comment enhancement - editor, file/image upload, etc.~~
        * ~~Add datagrid type, separate from board~~
    * User
        * ~~Add avatar as profile image~~
        * ~~Add user pic as profile image~~
    * Attachment
        * ~~Download count~~
    * Other
        * ~~Add more test~~
        * ~~Log analyzer~~
        * ~~Chatting~~
        * ~~Oracle - grip connection longtime~~
            * https://easyoradba.com/2021/03/10/shell-script-to-keep-oracle-always-free-autonomous-database-alive-with-sqlcl/
