package handler

// func Login(c *fiber.Ctx) error {
// 	failBody := `<meta http-equiv="refresh" content="2; url=/"></meta>`

// 	destination := c.FormValue("destination")
// 	if destination == "" {
// 		destination = "/"
// 	}

// 	userid := c.FormValue("userid")
// 	password := c.FormValue("password")

// 	if userid == "" || password == "" {
// 		return c.Status(http.StatusBadRequest).SendString(failBody + `Missing parameter`)
// 	}

// 	var users []model.UserData
// 	table := db.GetFullTableName(db.Info.UserTable)
// 	dbtype := db.GetDatabaseTypeString()

// 	column := np.CreateString(model.UserData{}, dbtype, "", false)
// 	where := np.CreateString(map[string]interface{}{"USERNAME": nil}, dbtype, "", false)
// 	whereAND := np.CreateString(map[string]interface{}{"GRADE": nil}, dbtype, "", false)

// 	sql := `
// 	SELECT
// 		` + column.Names + `
// 	FROM ` + table + `
// 	WHERE ` + where.Names + `='` + userid + `'
// 		AND ` + whereAND.Names + `!='` + "resigned_user" + `'`

// 	rows, err := db.Con.Query(sql)
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).SendString(failBody + `DB error or User may not exists`)
// 	}
// 	defer rows.Close()

// 	err = scan.Rows(&users, rows)
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).SendString(failBody + `User may not exists`)
// 	}

// 	if len(users) == 0 {
// 		return c.Status(http.StatusBadRequest).SendString(failBody + `User not exists`)
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password.String), []byte(password))
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).SendString(failBody + `Check password`)
// 	}

// 	// ua := useragent.Parse(c.Get("User-Agent"))
// 	// deviceType := ""
// 	// switch true {
// 	// case ua.Desktop:
// 	// 	deviceType = "desktop"
// 	// case ua.Mobile:
// 	// 	deviceType = "mobile"
// 	// case ua.Tablet:
// 	// 	deviceType = "tablet"
// 	// case ua.Bot:
// 	// 	deviceType = "bot"
// 	// }

// 	// authinfo := model.AuthInfo{
// 	// 	Name:       null.NewString(userid, true),
// 	// 	IpAddr:     null.NewString(c.IP(), true),
// 	// 	Device:     null.NewString(ua.Device, true),
// 	// 	DeviceType: null.NewString(deviceType, true),
// 	// 	Os:         null.NewString(ua.OS, true),
// 	// 	Browser:    null.NewString(ua.Name, true),
// 	// 	Duration:   null.NewInt(60*60*24*7, true),
// 	// 	// Duration: null.NewInt(10, true), // 10 seconds test
// 	// }

// 	// auth.SetupCookieToken(c.ResponseWriter, authinfo)
// 	// auth.SetCookieSession(c, authinfo)

// 	return c.Status(http.StatusOK).SendString(`<meta http-equiv="refresh" content="0; url=` + destination + `"></meta>`)
// }

// // Logout - Expire cookie
// func Logout(c *router.Context) {
// 	auth.ExpireCookie(c.ResponseWriter)
// 	auth.DestroyCookieSession(c)

// 	// authinfo := model.AuthInfo{}
// 	// if c.AuthInfo != nil {
// 	// 	authinfo = c.AuthInfo.(model.AuthInfo)
// 	// }
// 	// c.Text(http.StatusOK, "Good bye "+authinfo.Name.String)
// 	c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=/"></meta>`))
// }

// // Signup - Create new user
// func Signup(c *router.Context) {
// 	var err error

// 	now := time.Now().Format("20060102150405")
// 	columnsCount, _ := crud.GetUserColumnsCount()

// 	userIDX := ""
// 	userid := ""
// 	useremail := ""

// 	rbody, err := io.ReadAll(c.Body)
// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	var captchaData map[string]string
// 	err = json.Unmarshal(rbody, &captchaData)
// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	captchaResult := captcha.VerifyString(captchaData["captcha-id"], captchaData["captcha-answer"])

// 	if !captchaResult {
// 		c.Text(http.StatusBadRequest, "Captcha is not correct")
// 		return
// 	}

// 	switch columnsCount {
// 	case model.UserDataFieldCount:
// 		var userData model.UserData

// 		// err = json.NewDecoder(c.Body).Decode(&userData)
// 		err = json.Unmarshal(rbody, &userData)
// 		if err != nil {
// 			c.Text(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		password, err := bcrypt.GenerateFromPassword([]byte(userData.Password.String), consts.BcryptCost)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		userData.Password = null.StringFrom(string(password))
// 		userData.RegDTTM = null.StringFrom(now)
// 		userData.Grade = null.StringFrom("pending_user")
// 		userData.Approval = null.StringFrom("N")

// 		if userData.UserName.String == "" {
// 			c.Text(http.StatusBadRequest, "UserId is empty")
// 			return
// 		}
// 		if userData.Email.String == "" {
// 			c.Text(http.StatusBadRequest, "Email is empty")
// 			return
// 		}
// 		if userData.Password.String == "" {
// 			c.Text(http.StatusBadRequest, "Password is empty")
// 			return
// 		}

// 		err = crud.AddUser(userData)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		userid = userData.UserName.String
// 		useremail = userData.Email.String

// 		userInsertResult, err := crud.GetUserByNameAndEmail(userid, useremail)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		userIDX = fmt.Sprint(userInsertResult.Idx.Int64)
// 	default:
// 		userData := make(map[string]interface{})

// 		// err = json.NewDecoder(c.Body).Decode(&userData)
// 		err = json.Unmarshal(rbody, &userData)
// 		if err != nil {
// 			c.Text(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		password, err := bcrypt.GenerateFromPassword([]byte(userData["password"].(string)), consts.BcryptCost)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		userData["password"] = string(password)
// 		userData["regdate"] = now
// 		userData["grade"] = "pending_user"
// 		userData["approval"] = "N"

// 		if userData["userid"].(string) == "" {
// 			c.Text(http.StatusBadRequest, "UserId is empty")
// 			return
// 		}
// 		if userData["email"].(string) == "" {
// 			c.Text(http.StatusBadRequest, "Email is empty")
// 			return
// 		}
// 		if userData["password"].(string) == "" {
// 			c.Text(http.StatusBadRequest, "Password is empty")
// 			return
// 		}

// 		err = crud.AddUserMap(userData)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		userid = userData["userid"].(string)
// 		useremail = userData["email"].(string)

// 		userInsertResult, err := crud.GetUserByNameAndEmailMap(userid, useremail)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		userIDX = userInsertResult.(map[string]interface{})["IDX"].(string)
// 	}

// 	verificationKEY := GetRandomString(32)
// 	verificationData := map[string]string{
// 		"USER_IDX": userIDX,
// 		"TOKEN":    verificationKEY,
// 	}

// 	err = crud.AddUserVerification(verificationData)
// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	// Send verification email
// 	domain := email.Info.Domain
// 	message := email.Message{
// 		Service:          email.Info.Service,
// 		AppendFromToName: false,
// 		From:             email.From{Email: email.Info.SenderInfo.Email, Name: email.Info.SenderInfo.Name},
// 		To:               email.To{Email: useremail, Name: userid},
// 		Subject:          "EnjoyTools - Email Verification",
// 		Body: `
// 		Please click the link below to verify your email address.
// 		<br />
// 		<a href='` + domain + `/verify?userid=` + userid + `&email=` + useremail + `&token=` + verificationKEY + `'>Click here</a>`,
// 		BodyType: email.HTML,
// 	}

// 	if email.Info.UseEmail {
// 		err = email.SendVerificationEmail(message)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 	}

// 	result := map[string]string{
// 		"result": "ok",
// 	}

// 	c.Json(http.StatusOK, result)
// }

// func UserVerification(c *router.Context) {
// 	var err error

// 	queries := c.URL.Query()
// 	userid := queries.Get("userid")
// 	useremail := queries.Get("email")
// 	token := queries.Get("token")

// 	if userid == "" || useremail == "" || token == "" {
// 		c.Text(http.StatusBadRequest, "Not enough parameters")
// 		return
// 	}

// 	result, err := crud.VerifyAndUpdateUser(userid, useremail, token)
// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	if !result {
// 		c.Text(http.StatusBadRequest, "Invalid request")
// 		return
// 	}

// 	c.Html(http.StatusOK, []byte(`
// 	<script>
// 		alert("Your email has been verified. You can now login.");
// 		location.href = "/";
// 	</script>
// 	`))
// }

// func ResetPassword(c *router.Context) {
// 	var err error

// 	columnsCount, _ := crud.GetUserColumnsCount()

// 	userid := c.FormValue("userid")
// 	useremail := c.FormValue("email")

// 	if userid == "" {
// 		c.Text(http.StatusBadRequest, "UserId is empty")
// 		return
// 	}
// 	if useremail == "" {
// 		c.Text(http.StatusBadRequest, "Email is empty")
// 		return
// 	}

// 	password := GetRandomString(16)
// 	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	switch columnsCount {
// 	case model.UserDataFieldCount:
// 		user, err := crud.GetUserByNameAndEmail(userid, useremail)
// 		if err != nil {
// 			// c.Text(http.StatusInternalServerError, err.Error())
// 			c.Html(http.StatusOK, []byte(consts.MsgPasswordResetUserNotFound))
// 			return
// 		}

// 		user.Password = null.StringFrom(string(passwordHash))
// 		crud.UpdateUser(user)

// 	default:
// 		user, err := crud.GetUserByNameAndEmailMap(userid, useremail)
// 		if err != nil {
// 			// c.Text(http.StatusInternalServerError, err.Error())
// 			c.Html(http.StatusOK, []byte(consts.MsgPasswordResetUserNotFound))
// 			return
// 		}

// 		user.(map[string]interface{})["password"] = string(passwordHash)
// 		crud.UpdateUserMap(user.(map[string]interface{}))
// 	}

// 	// Send password reset email
// 	message := email.Message{
// 		Service:          email.Info.Service,
// 		AppendFromToName: false,
// 		From:             email.From{Email: email.Info.SenderInfo.Email, Name: email.Info.SenderInfo.Name},
// 		To:               email.To{Email: useremail, Name: userid},
// 		Subject:          "EnjoyTools - Password changed",
// 		Body: `
// 		The password for your account was changed on ` + time.Now().UTC().Format("2006-01-02 15:04:05 UTC") + `
// 		<br /><br />
// 		` + password,
// 		BodyType: email.HTML,
// 	}

// 	if email.Info.UseEmail {
// 		err = email.SendVerificationEmail(message)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 	}

// 	c.Html(http.StatusOK, []byte(consts.MsgPasswordResetEmail))
// }

// func HandleUserList(c *router.Context) {
// 	// Use struct with default columns or map with default and user defined columns
// 	columnsCount, _ := crud.GetUserColumnsCount()
// 	// columnsCount, _ := db.Obj.GetColumnCount(db.Info.UserTable)

// 	// queries := c.URL.Query()
// 	// search := queries.Get("search")

// 	var err error
// 	queries := c.URL.Query()

// 	listingOptions := model.UserListingOptions{}
// 	listingOptions.Search = null.StringFrom(queries.Get("search"))

// 	listingOptions.Page = null.IntFrom(1)
// 	listingOptions.ListCount = null.IntFrom(100)

// 	if queries.Get("count") != "" {
// 		countPerPage, err := strconv.Atoi(queries.Get("count"))
// 		if err != nil {
// 			c.Text(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		listingOptions.ListCount = null.IntFrom(int64(countPerPage))
// 	}

// 	if queries.Get("page") != "" {
// 		page := queries.Get("page")
// 		pageNum, err := strconv.Atoi(page)
// 		if err != nil {
// 			c.Text(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		listingOptions.Page = null.IntFrom(int64(pageNum))
// 	}

// 	listingOptions.Page.Int64--

// 	h, err := LoadHTML(c)
// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	var listJSON []byte

// 	switch columnsCount {
// 	case model.UserDataFieldCount:
// 		// result, err := crud.GetUsersList(search)
// 		// if err != nil {
// 		// 	c.Text(http.StatusInternalServerError, err.Error())
// 		// 	return
// 		// }

// 		// c.Json(http.StatusOK, result)

// 		list, err := crud.GetUsers(listingOptions)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 		}

// 		listJSON, _ = json.Marshal(list)
// 	default:
// 		// result, err := crud.GetUsersListMap(search)
// 		// if err != nil {
// 		// 	c.Text(http.StatusInternalServerError, err.Error())
// 		// 	return
// 		// }

// 		// c.Json(http.StatusOK, result)

// 		list, err := crud.GetUsersMap(listingOptions)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 		}

// 		listJSON, _ = json.Marshal(list)
// 	}

// 	h = bytes.ReplaceAll(h, []byte("$USER_LIST$"), listJSON)

// 	c.Html(http.StatusOK, h)
// }
