const authTokenKey = "xuen4sFsjdfhKJHf"

export function setAuth(login, password) {
    let token = login + ":" + password;
    localStorage.setItem(authTokenKey, btoa(token))
}

export function getAuth() {
    let token = localStorage.getItem(authTokenKey)
    if (token == null) {
        return ""
    }

    return "Basic " + token
}
