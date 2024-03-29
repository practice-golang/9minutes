/* 
{
  "board-list": [
    {
      "idx": 1,
      "board-name": "misc1",
      "board-code": "misc1",
      "board-type": "board",
      "board-table": "board_misc1",
      "comment-table": "comment_misc1",
      "grant-read": "guest",
      "grant-write": "user_active",
      "grant-comment": "user_active",
      "grant-upload": "user_active",
      "fields": null
    }
  ],
  "total-page": 1,
  "current-page": 1
}
 */

const defaultCount = 10

async function getBoardGrades() {
    let grades = []

    const rg = await fetch("/api/admin/board-grades", {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
    })

    if (rg.ok) {
        let gradesArr = Object.entries(await rg.json()).sort((a, b) => { return a[1].rank - b[1].rank })
        for (let el of gradesArr) { grades.push(el[1]) }
    }

    return grades
}

async function getBoards(page, listCount, search) {
    let boardsData = {}

    let uri = `/api/admin/board?page=${page}&list-count=${listCount}`
    if (search != "") { uri += `&search=${search}` }

    const rl = await fetch(uri, {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
    })

    if (rl.ok) { boardsData = await rl.json() }
    if (boardsData["board-list"] == null) { boardsData["board-list"] = [] }

    return boardsData
}

export const load = async ({ url, fetch }) => {
    const listCount = Number(url.searchParams.get("list-count")) || defaultCount
    const page = Number(url.searchParams.get("page")) || 1
    const search = url.searchParams.get("search") || ""

    const columns = [
        { "display-name": "Index", "column-code": "idx", "column-name": "IDX" },
        { "display-name": "Name", "column-code": "board-name", "column-name": "BOARD_NAME" },
        { "display-name": "Code", "column-code": "board-code", "column-name": "BOARD_CODE" },
        { "display-name": "Type", "column-code": "board-type", "column-name": "BOARD_TYPE" },
        { "display-name": "Board table", "column-code": "board-table", "column-name": "BOARD_TABLE" },
        { "display-name": "Comment table", "column-code": "comment-table", "column-name": "COMMENT_TABLE" },
        { "display-name": "Grant read", "column-code": "grant-read", "column-name": "GRANT_READ" },
        { "display-name": "Grant write", "column-code": "grant-write", "column-name": "GRANT_WRITE" },
        { "display-name": "Grant comment", "column-code": "grant-comment", "column-name": "GRANT_COMMENT" },
        { "display-name": "Grant upload", "column-code": "grant-upload", "column-name": "GRANT_UPLOAD" },
    ]

    return {
        "default-count": defaultCount,
        columns: columns,
        grades: getBoardGrades(),
        "boardlist-data": getBoards(page, listCount, search)
    }
}