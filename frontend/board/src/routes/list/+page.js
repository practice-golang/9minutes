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

    async function getContentList(page, listCount, search) {
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

    return {
        columns: getColumns(),
        "userlist-data": getContentList(page, listCount, search)
    }
}