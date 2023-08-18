package handler

import (
	"9minutes/config"
	"9minutes/internal/crud"
	"9minutes/model"
	"strconv"

	"gopkg.in/guregu/null.v4"
)

func GetPostingList(boardCODE string, queries map[string]string) (model.PostingPageData, error) {
	var err error

	board := model.Board{BoardCode: null.StringFrom(boardCODE)}
	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return model.PostingPageData{}, err
	}

	page := 1
	count := config.PostingListCountPerPage
	if queries["page"] != "" {
		page, err = strconv.Atoi(queries["page"])
		if err != nil {
			return model.PostingPageData{}, nil
		}
	}
	if queries["count"] != "" {
		count, err = strconv.Atoi(queries["count"])
		if err != nil {
			return model.PostingPageData{}, nil
		}
	}

	listingOptions := model.PostingListingOptions{}
	listingOptions.Search = null.StringFrom(queries["search"])

	listingOptions.Page = null.IntFrom(int64(page - 1))
	listingOptions.ListCount = null.IntFrom(int64(count))

	list, err := crud.GetPostingList(board, listingOptions)

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

	return list, err
}

func GetPostingData(boardCode, idx string) (model.Posting, error) {
	var err error

	board := model.Board{BoardCode: null.StringFrom(boardCode)}
	posting := model.Posting{}

	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return model.Posting{}, err
	}

	posting, err = crud.GetPosting(board, idx)
	if err != nil {
		return model.Posting{}, err
	}

	posting.Views = null.IntFrom(posting.Views.Int64 + 1)
	err = crud.UpdatePosting(board, posting, "viewcount")
	if err != nil {
		return model.Posting{}, err
	}

	return posting, nil
}

func GetCommentsList(boardCode string, postingIDX string, queries map[string]string) (model.CommentPageData, error) {
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

	comments, err := crud.GetComments(board, postingIDX, commentOptions)
	if err != nil {
		return model.CommentPageData{}, err
	}

	return comments, nil
}
