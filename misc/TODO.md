# Todo

* [ ] Logger - request, user login history - useridx, ip, regdate

* [ ] manager account
    * [ ] API
        * [ ] Edit board settings
        * [ ] Add/edit/remove user as member -> Add manager, member column to board setting table
    * [ ] DB
    * [ ] Admin page - add/edit manager account setting for each boards

* [ ] Resurrect account verification table usage
    * [ ] Email
    * [ ] Verifying

* [ ] Write/Update/Upload text/file size limit
    * [ ] Topic
    * [ ] Comment
* [ ] Resurrect session storage - etcd, ~~redis~~..
* [ ] Admin - page for temporary files flushing - `TOPIC_IDX` and `COMMENT_IDX` are `-1`
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
* [x] Grant,confirm routine for guest
    * [x] ip address - Add IP_ADDRESS to each db table creation functions
    * [x] Appear guest board in list selector
    * [x] Control topic
        * [x] Password check if user or author is guest
            * [x] Write
            * [x] Edit
            * [x] Delete
    * [x] Control comment
        * [x] Password check if author is guest
            * [x] DB
            * [x] bug - no authorname
            * [x] Write
            * [x] Edit
            * [x] Delete
    * [x] show edit/delete button of guest comment for normal user
    * [x] Admin allow comment control without password
        * [x] Edit
        * [x] Delete
    * [x] password encryption
        * Topic
            * [x] Write
            * [x] Edit
            * [x] Delete
        * Comment
            * [x] Write
            * [x] Edit
            * [x] Delete
* [x] Change less_eq, more_eq to le, ge
    * https://pkg.go.dev/text/template
    * [x] change less_eq, more_eq
    * [x] delete less_eq, more_eq
* [x] pgsql - correct bugs
* [x] Control edit/delete button for author only
* [x] Add timestamp at list
* [x] Change js prompt to input - password masking
    * [x] Topic edit
    * [x] Topic delete
    * [x] Comment delete
* [x] Show `empty topic` if list is empty
* [x] Attachment transaction when edit
    * [x] Topic write - cancel only, history.back not work
    * [x] Topic edit
    * [x] Comment
* [x] Append `TOPIC_IDX` to upload table
* [x] Write IDX to upload table
    * [x] Topic write
    * [x] Topic edit
    * [x] Comment write
    * [x] Comment edit
    * [x] Add `no index` when delete new upload
* [x] Rename all post, posting to topic
* [x] Add regdate to upload table
* [x] Add editor link target `_blank`

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
