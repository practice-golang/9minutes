# Todo

* [ ] Add post index and board index to upload table

* [ ] Correct and finish approval
* [ ] Rearrange route paths
* [ ] Add Oracle DB support
* [ ] Add more test
* [ ] Log analyzer
* [ ] CORS???
* [ ] DB backup, migration

* Board
    * [ ] Add sheet type - custom board type
    * [ ] ~~Comment enhancement - editor, file/image upload, etc.~~ waive
    * [x] Delete uploaded file(s)
        * [x] Write
            * [x] Write then delete post
            * [x] Cancel write post
        * [x] Update
            * [x] Update then delete post
            * [x] Cancel write post
        * [x] Delete
    * [x] Remove query param - gallery type
    * [x] Image width limit

* [x] Add option to generate dkim.key
* ~~Chatting???~~ - waive
* User
    * ~~Add avatar as profile image~~ - waive
    * ~~Add user pic as profile image~~ - waive
    * [x] Add password reset
    * [x] Add resign
    * [x] Captcha
    * [x] Add more info in session
* ~~Page, HTML~~ - waive
    * Embed css styled html - beauty cuty html, css
    * Sample menu - eg. board link
* Attachment
    * [x] Upload filename change - Upload-xxxx -> hashed name
    * [x] Add both filename and storagename to db
    * [x] Download original filename
    * ~~Download count~~ - waive
* [x] Delete all related `book`, `books`
* [x] Change board list, read content from ajax to html/template except comments
* [x] ~~Change datalist to select - admin~~ / Gave up :p, Add restrictDatalist function rather than to change them
* [x] Add heroku env for rest of dbs except mysql
* [x] Email sender - for user approval
* [x] User approval
    * [x] Send random number via email
    * [x] verify random number
* [x] Add session store address, port for etcd, redis
* [x] Youtube link, embed
* [x] Comment remove
    - Remove working but, show 500 internal server error -> Missed exception when not logged in
    - Correct comment paging range<br />
    Set -1 (=last page) when add comment or first page reading
* [x] html
    * Remove `Back to admin` in all html - Use Admin link
    * Add `Frontpage menu` in all board html
* [x] html/template for content list
* [x] File attatchment
* [x] Shared session
* [x] User approval - email sending