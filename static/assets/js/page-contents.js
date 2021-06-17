let setInnerHTML = async function (elm, html) {
    elm.innerHTML = html
    Array.from(elm.querySelectorAll("script")).forEach(oldScript => {
        const newScript = document.createElement("script")
        Array.from(oldScript.attributes)
            .forEach(attr => newScript.setAttribute(attr.name, attr.value))
        newScript.appendChild(document.createTextNode(oldScript.innerHTML))
        oldScript.parentNode.replaceChild(newScript, oldScript)
    })

    return true
}
const restURIdomain = window.location.protocol + "//" + window.location.host

async function checkLogedIn() {
    let result = false
    let reissue = false
    let response = ""
    let r = ""

    let token = localStorage.getItem('token')

    if (token != null) {
        response = await fetch(restURIdomain + "/api/user/token/verify", {
            method: "POST",
            headers: { "Authorization": "Bearer " + token }
        })
        if (response.ok) {
            r = await response.json()
            if (r.msg == "OK") { result = true }
        } else {
            if (response.status != 200) {
                r = await response.json()
                if (r.msg == "Token is expired") { reissue = true }
            }
        }

        if (!result && reissue) {
            response = await fetch(restURIdomain + "/api/user/token/verify", {
                method: "GET",
                headers: { "Authorization": "Bearer " + token }
            })

            if (response.ok) {
                r = await response.json()
                localStorage.setItem("token", r.token)
                token = r.token
                result = true
            } else {
                localStorage.removeItem("token")
                result = false
            }
        }
    }

    const results = {
        result: result,
        token: token
    }

    return results
}

async function getPageContents() {
    const grant = await checkLogedIn()

    if (!grant.result) {
        location.href = "/users/login"
    }

    let response = await fetch(restURIdomain + "/contents-body/" + routeTarget, {
        method: "GET",
        headers: { "Authorization": "Bearer " + grant.token }
    })

    if (response.ok) {
        const r = await response.text()

        let result = await setInnerHTML(document.querySelector('#body-cavity'), r)
        if (result) {
            setInnerHTML = undefined
        }
    }
}

window.onload = () => {
    getPageContents()
}