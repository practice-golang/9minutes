/* 
[
    {
        "idx": 1,
        "display-name": "Idx",
        "column-code": "idx",
        "column-type": "integer",
        "column-name": "IDX",
        "sort-order": 1
    },
    {
        "idx": 2,
        "display-name": "UserID",
        "column-code": "userid",
        "column-type": "text",
        "column-name": "USERID",
        "sort-order": 2
    }
    ...
]
 */

export const load = async ({ fetch }) => {
    let columns = []

    const r = await fetch("/api/admin/user-columns", {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
    })

    if (r.ok) {
        const response = await r.json()
        columns = response
    }

    return { columns: columns }
}
