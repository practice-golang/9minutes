/* 
[
    {
        "email": "admin@please.modify",
        "grade": "admin",
        "idx": "1",
        "reg-dttm": "20230710044230",
        "userid": "admin"
    },
    {
        "email": "edp1096@naver.com",
        "grade": "pending_user",
        "idx": "5",
        "reg-dttm": "20230710044824",
        "userid": "bab2"
    }
]
 */

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