import {getToken} from '../auth/auth.js';

export const statusOK = 200

export const methodGet = "GET"

export const routeCheckAuth = "/api/v1/auth/check"
export const routeGetPlacesTree = "/api/v1/places/tree"

export function request(method, url) {
    let xhr = new XMLHttpRequest();

    xhr.open(method, url, false);

    xhr.setRequestHeader('Accept', 'application/json')
    xhr.setRequestHeader('Content-Type', 'application/json')
    xhr.setRequestHeader('Authorization', 'Basic ' + getToken())

    xhr.send();

    return {
        data: xhr.responseText,
        status: xhr.status,
    };
}