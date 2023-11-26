export const load = async ({ url, fetch }) => {
    const listCount = Number(url.searchParams.get("list-count")) || 20
    const page = Number(url.searchParams.get("page")) || 1
    const search = url.searchParams.get("search") || ""

    async function getUserColumns() {
        let columns = []

        const rc = await fetch("/api/admin/user-columns", {
            method: "GET",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
        })

        if (rc.ok) { columns = await rc.json() }

        return columns
    }

    async function getUserGrades() {
        let grades = []

        const rg = await fetch("/api/admin/user-grades", {
            method: "GET",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
        })

        if (rg.ok) {
            let gradesArr = Object.entries(await rg.json()).sort((a, b) => { return a[1].point - b[1].point })
            for (let el of gradesArr) { grades.push(el[1]) }
        }

        return grades
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

    return {
        columns: getUserColumns(),
        grades: getUserGrades(),
        "userlist-data": getUsers(page, listCount, search)
    }
}