"use strict"

const authTokenKey = "46sgdfjhFRTFhagfhkdd3"

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