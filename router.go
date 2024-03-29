package server

import (
	"9minutes/handler"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func checkAdmin(c *fiber.Ctx) error {
	usergrade, err := handler.GetSessionUserGrade(c)
	if err != nil {
		return err
	}

	if usergrade != "admin" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"status": 403, "message": "forbidden"})
	}

	return c.Next()
}

func setApiAdmin(a *fiber.App) {
	/* API Admin */
	gadmin := a.Group("/api/admin")
	gadmin.Use(checkAdmin)
	gadmin.Get("/health", handler.HealthCheckAPI)

	/* API Admin - User fileds */
	gauserfield := gadmin.Group("/user-columns") // required add auth middleware
	gauserfield.Get("/", handler.GetUserColumnsAPI)
	gauserfield.Post("/", handler.AddUserColumnAPI)
	gauserfield.Put("/", handler.UpdateUserColumnsAPI)
	gauserfield.Delete("/", handler.DeleteUserColumnsAPI)

	/* API Admin - User/Board grades */
	gauserggrades := gadmin.Group("/user-grades") // required add auth middleware
	gauserggrades.Get("/", handler.GetUserGrades)
	gaboardgrades := gadmin.Group("/board-grades") // required add auth middleware
	gaboardgrades.Get("/", handler.GetBoardGrades)

	/* API Admin - Users */
	gauser := gadmin.Group("/user") // required add auth middleware
	gauser.Get("/", handler.GetUserListAPI)
	gauser.Post("/", handler.AddUserAPI)
	gauser.Put("/", handler.UpdateUserAPI)
	gauser.Delete("/", handler.DeleteUserAPI)

	/* API Admin - Boards */
	gaboard := gadmin.Group("/board") // required add auth middleware
	gaboard.Get("/", handler.GetBoardsAPI)
	gaboard.Post("/", handler.AddBoardAPI)
	gaboard.Put("/", handler.UpdateBoardAPI)
	gaboard.Delete("/", handler.DeleteBoardAPI)

	/* API Board Member */
	gamember := gadmin.Group("/member") // required add auth middleware
	gamember.Get("/", handler.GetMemberListAPI)
	gamember.Post("/", handler.AddMemberAPI)
	gamember.Put("/", handler.UpdateMemberAPI)
	gamember.Delete("/", handler.DeleteMemberAPI)
}

func setApiBoard(a *fiber.App) {
	/* API Board */
	gbrd := a.Group("/api/board") // Require add session middleware

	/* API Board list */
	gbrd.Get("/list", handler.BoardListAPI)

	/* API Topic */
	gbrd.Get("/:board_code", handler.ListTopicAPI)
	gbrd.Get("/:board_code/topic/:idx", handler.ReadTopicAPI)
	gbrd.Post("/:board_code/topic", handler.WriteTopicAPI)
	gbrd.Put("/:board_code/topic/:idx", handler.UpdateTopicAPI)
	gbrd.Delete("/:board_code/topic/:idx", handler.DeleteTopicAPI)

	/* API Comment */
	gbrd.Get("/:board_code/:topic_idx/comment", handler.GetCommentsAPI)
	gbrd.Post("/:board_code/:topic_idx/comment", handler.WriteCommentAPI)
	gbrd.Put("/:board_code/:topic_idx/comment/:comment_idx", handler.UpdateCommentAPI)
	gbrd.Delete("/:board_code/:topic_idx/comment/:comment_idx", handler.DeleteCommentAPI)
}

func setApiUploader(r *fiber.App) {
	/* API Uploader */
	gupload := r.Group("/api/uploader") // Require add session middleware
	gupload.Post("/", handler.UploadFile)
	gupload.Post("/files-info", handler.FilesInfo)
	gupload.Delete("/", handler.DeleteFiles)

	// gu.POST(`/image$`, handler.UploadImage)

	// Delete all of files, images, title-image which is(are) uploaded during writing or editing on board, when cancel
}

func setAPIs(a *fiber.App) {
	/* User login API */
	gapi := a.Group("/api")
	gapi.Get("/health", handler.HealthCheckAPI)
	gapi.Post("/login", handler.LoginAPI)
	gapi.Get("/logout", handler.LogoutAPI)
	gapi.Post("/signup", handler.SignupAPI)
	gapi.Post("/password-reset", handler.ResetPasswordAPI)
	gapi.Get("/2fa/qr", handler.Get2FaQR)
	gapi.Get("/2fa/verify", handler.Verify2FA)

	/* API myinfo */
	gmyinfo := gapi.Group("/myinfo") // Require add session middleware
	gmyinfo.Get("/", handler.GetMyInfo)
	gmyinfo.Put("/", handler.UpdateMyInfo)
	gmyinfo.Delete("/", handler.QuitUser)
}

// setStaticFiles - Set static files
func setStaticFiles(a *fiber.App) {
	if IsStaticEmbed {
		configFiles := filesystem.Config{Root: http.FS(EmbedStatic), PathPrefix: "/static/files"}
		configAssets := filesystem.Config{Root: http.FS(EmbedHTML), PathPrefix: "/static/html/assets"}
		configAdminApp := filesystem.Config{Root: http.FS(EmbedHTML), PathPrefix: "/static/html/admin/_app"}

		a.Use("/favicon.png", filesystem.New(configAssets))
		a.Use("/files", filesystem.New(configFiles))
		a.Use("/assets", filesystem.New(configAssets))
		a.Use("/admin/_app", filesystem.New(configAdminApp))
	} else {
		a.Static("/favicon.png", HtmlPath+"/favicon.png")
		a.Static("/files", FilesPath)
		a.Static("/assets/", HtmlPath+"/assets")
		a.Static("/admin/_app/", HtmlPath+"/admin/_app")
	}

	a.Static("/upload", UploadPath)
}

func setPage(a *fiber.App) {
	a.Get("/*", handler.HandleHTML)
}
