# Todo

* [ ] User approval
    * [ ] Send random number via email
    * [ ] verify random number
* [ ] Change board read from ajax to template
* [ ] Download count
* [ ] Upload filename manage - Upload-xxxx -> hashed name
* [ ] Download original filename
* [ ] Delete uploaded file when edit content
* [ ] Delete all related `book`, `books`
* [ ] Comment enhancement - editor, file/image upload, etc.
* [ ] Remove board type column or query param
* [ ] Add avatar/profile image
* [ ] Add heroku env for rest of dbs except mysql
* [ ] Change datalist to select
* [ ] Add sheet type - custom board type
* [ ] beauty cuty html, css
* [ ] Sample menu - eg. board link
* [ ] Add more test
* [ ] Image limit width
* [ ] Captcha
* [ ] Log analyzer
* [ ] Add more info into session
* [ ] DB backup, migration
* [ ] Chatting???
* [ ] CORS???
* [ ] Embed css styled html

* [x] Email sender - for user approval
* [x] Add session store address, port for etcd, redis
* [x] Youtube link, embed
* [x] Comment remove
    - Remove working but, show 500 internal server error -> Missed exception when not logged in
    - Correct comment paging range<br />
    Set -1 (=last page) when add comment or first page reading
* [x] html
    * Remove `Back to admin` in all html - Use Admin link
    * Add `Frontpage menu` in all board html