"use strict"

const authTokenKey = "xuen4sFsjdfhKJHf"

export function setToken(token) {
    localStorage.setItem(authTokenKey, token)
}

export function getToken() {
    let token = localStorage.getItem(authTokenKey)
    if (token == null) {
        return ""
    }

    return token
}

export function clearToken() {
    localStorage.removeItem(authTokenKey)
}