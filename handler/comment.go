package handler

import (
	"9minutes/config"
	"9minutes/internal/crud"
	"9minutes/model"
	"strconv"

	"gopkg.in/guregu/null.v4"
)

func GetCommentsList(boardCode string, postingIDX string, queries map[string]string) (model.CommentPageData, error) {
	var err error

	board := model.Board{}
	board.BoardCode = null.StringFrom(boardCode)
	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return model.CommentPageData{}, err
	}

	page := 1
	count := config.CommentCountPerPage
	if queries["page"] != "" {
		page, err = strconv.Atoi(queries["page"])
		if err != nil {
			return model.CommentPageData{}, nil
		}
	}
	if queries["count"] != "" {
		count, err = strconv.Atoi(queries["count"])
		if err != nil {
			return model.CommentPageData{}, nil
		}
	}

	commentSearch := queries["search"]

	commentOptions := model.CommentListingOptions{}
	commentOptions.Search = null.StringFrom(commentSearch)

	commentOptions.Page = null.IntFrom(int64(page - 1))
	commentOptions.ListCount = null.IntFrom(int64(count))

	comments, err := crud.GetComments(board, postingIDX, commentOptions)
	if err != nil {
		return model.CommentPageData{}, err
	}

	pageList := []int{}
	pageShowGap := 2
	for i := comments.CurrentPage - pageShowGap; i <= comments.CurrentPage+pageShowGap; i++ {
		if i > 0 && i <= comments.CurrentPage+pageShowGap && i <= comments.TotalPage {
			pageList = append(pageList, i)
		}
	}
	comments.PageList = pageList

	pageJumpGap := 5
	comments.JumpPrev = comments.CurrentPage - pageJumpGap
	comments.JumpNext = comments.CurrentPage + pageJumpGap
	if comments.JumpPrev < 1 {
		comments.JumpPrev = 1
	}
	if comments.JumpNext > comments.TotalPage {
		comments.JumpNext = comments.TotalPage
	}

	return comments, nil
}
