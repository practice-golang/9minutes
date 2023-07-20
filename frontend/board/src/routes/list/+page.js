export const load = async ({ url, fetch }) => {
    const boardCode = url.searchParams.get("board_code") || ""
    const listCount = Number(url.searchParams.get("list-count")) || 10
    const page = Number(url.searchParams.get("page")) || 1
    const search = url.searchParams.get("search") || ""

    async function getContentList(boarCode, page, listCount, search) {
        let usersData = {}
        if (page == "") { page = 1 }

        let uri = `/api/board/${boarCode}?page=${page}&list-count=${listCount}`
        if (search != "") { uri += `&search=${search}` }

        const r = await fetch(uri, {
            method: "GET",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
        })

        if (r.ok) { usersData = await r.json() }

        console.log(usersData)

        return usersData
    }

    console.log("WTF???")

    return {
        "content-list": getContentList(boardCode, page, listCount, search)
    }
}