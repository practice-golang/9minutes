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

	board := model.Board{}
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

	return list, err
}

func GetContentData(boardCode string, idx int, queries map[string]string) (map[string]interface{}, error) {
	var err error

	result := make(map[string]interface{})

	board := model.Board{}
	content := model.Content{}
	comments := model.CommentPageData{}

	board.BoardCode = null.StringFrom(boardCode)
	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return nil, err
	}

	content, err = crud.GetContent(board, strconv.Itoa(idx))
	if err != nil {
		return nil, err
	}

	content.Views = null.IntFrom(content.Views.Int64 + 1)
	err = crud.UpdateContent(board, content, "viewcount")
	if err != nil {
		return nil, err
	}

	commentSearch := queries["search"]

	commentOptions := model.CommentListingOptions{}
	commentOptions.Search = null.StringFrom(commentSearch)

	commentOptions.Page = null.IntFrom(-1)
	commentOptions.ListCount = null.IntFrom(int64(config.CommentCountPerPage))

	comments, err = crud.GetComments(board, content, commentOptions)
	if err != nil {
		return nil, err
	}

	result["content"] = content
	result["comments"] = comments

	return result, err
}
