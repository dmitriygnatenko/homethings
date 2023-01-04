"use strict"

const authTokenKey = "xuen4sFsjdfhKJHf"

export function setToken(login, password) {
    let token = login + ":" + password;
    localStorage.setItem(authTokenKey, btoa(token))
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