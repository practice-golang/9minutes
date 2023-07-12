/* 
[
    {
        "email": "admin@please.modify",
        "grade": "admin",
        "idx": "1",
        "regdate": "20230710044230",
        "userid": "admin"
        ... and user defined ...
    },
    {
        "email": "edp1096@naver.com",
        "grade": "pending_user",
        "idx": "5",
        "regdate": "20230710044824",
        "userid": "bab2"
        ... and user defined ...
    }
]
 */

export const load = async ({ url, fetch }) => {
    const listCount = Number(url.searchParams.get("list-count")) || 10
    const page = Number(url.searchParams.get("page")) || 1
    const search = url.searchParams.get("search") || ""

    async function getColumns() {
        let columns = []

        const rc = await fetch("/api/admin/user-columns", {
            method: "GET",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
        })

        if (rc.ok) { columns = await rc.json() }

        return columns
    }

    async function getUsers(page, listCount, search) {
        let usersData = {}

        let uri = `/api/admin/user?page=${page}&list-count=${listCount}`
        if (search != "") { uri += `&search=${search}` }

        const rl = await fetch(uri, {
            method: "GET",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
        })

        if (rl.ok) { usersData = await rl.json() }

        return usersData
    }

    return { columns: getColumns(), "userlist-data": getUsers(page, listCount, search) }
}