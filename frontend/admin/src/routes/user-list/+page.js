export const load = async () => {
    const r = await fetch("/api/admin/user", {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
    })

    if (r.ok) {
        const response = await r.json()
        console.log(response)
    }

    let result = {
        idx: 'idx',
        userid: "userid",
        email: "email",
        grade: "grade",
        approval: "approval",
        link: "link",
        page: "page",
    }

    return result
}