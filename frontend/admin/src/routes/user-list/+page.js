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

export const load = async ({ fetch }) => {
    let columns = []
    let users = []

    const rc = await fetch("/api/admin/user-columns", {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
    })

    if (rc.ok) { columns = await rc.json() }

    const rl = await fetch("/api/admin/user", {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
    })

    if (rl.ok) { users = await rl.json() }

    return { columns: columns, users: users }
}