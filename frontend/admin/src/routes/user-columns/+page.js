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
        displayName: 'display Name',
        columnCode: "",
        columnType: "",
        columnName: "",
        sortOrder: "",
    }

    return result
}
