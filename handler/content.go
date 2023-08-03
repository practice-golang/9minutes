package handler

import (
	"9minutes/config"
	"9minutes/internal/crud"
	"9minutes/model"
	"strconv"

	"gopkg.in/guregu/null.v4"
)

func GetContentsList(boardCODE string, queries map[string]string) (model.ContentPageData, error) {
	var err error

	board := model.Board{BoardCode: null.StringFrom(boardCODE)}
	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return model.ContentPageData{}, err
	}

	page := 1
	count := config.ContentListCountPerPage
	if queries["page"] != "" {
		page, err = strconv.Atoi(queries["page"])
		if err != nil {
			return model.ContentPageData{}, nil
		}
	}
	if queries["count"] != "" {
		count, err = strconv.Atoi(queries["count"])
		if err != nil {
			return model.ContentPageData{}, nil
		}
	}

	listingOptions := model.ContentListingOptions{}
	listingOptions.Search = null.StringFrom(queries["search"])

	listingOptions.Page = null.IntFrom(int64(page - 1))
	listingOptions.ListCount = null.IntFrom(int64(count))

	list, err := crud.GetContentList(board, listingOptions)

	pageList := []int{}
	pageShowGap := 2
	for i := list.CurrentPage - pageShowGap; i <= list.CurrentPage+pageShowGap; i++ {
		if i > 0 && i <= list.CurrentPage+pageShowGap && i <= list.TotalPage {
			pageList = append(pageList, i)
		}
	}
	list.PageList = pageList

	pageJumpGap := 5
	list.JumpPrev = list.CurrentPage - pageJumpGap
	list.JumpNext = list.CurrentPage + pageJumpGap
	if list.JumpPrev < 1 {
		list.JumpPrev = 1
	}
	if list.JumpNext > list.TotalPage {
		list.JumpNext = list.TotalPage
	}

	list.BoardCode = board.BoardCode.String

	return list, err
}

func GetContentData(boardCode, idx string) (model.Content, error) {
	var err error

	board := model.Board{BoardCode: null.StringFrom(boardCode)}
	content := model.Content{}

	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return model.Content{}, err
	}

	content, err = crud.GetContent(board, idx)
	if err != nil {
		return model.Content{}, err
	}

	content.Views = null.IntFrom(content.Views.Int64 + 1)
	err = crud.UpdateContent(board, content, "viewcount")
	if err != nil {
		return model.Content{}, err
	}

	return content, nil
}

func GetCommentsList(boardCode string, contentIDX string, queries map[string]string) (model.CommentPageData, error) {
	var err error

	board := model.Board{}
	board.BoardCode = null.StringFrom(boardCode)
	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return model.CommentPageData{}, err
	}

	commentSearch := queries["search"]

	commentOptions := model.CommentListingOptions{}
	commentOptions.Search = null.StringFrom(commentSearch)

	commentOptions.Page = null.IntFrom(-1)
	commentOptions.ListCount = null.IntFrom(int64(config.CommentCountPerPage))

	comments, err := crud.GetComments(board, contentIDX, commentOptions)
	if err != nil {
		return model.CommentPageData{}, err
	}

	return comments, nil
}
