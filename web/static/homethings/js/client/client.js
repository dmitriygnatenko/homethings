'use strict'

import {getToken} from '../auth/auth.js';

export const statusOK = 200

export const methodGet = "GET"
export const methodPost = "POST"
export const methodPut = "PUT"
export const methodDelete = "DELETE"

export const routeCheckAuth = "/api/v1/auth/check"
export const routeAddPlace = "/api/v1/places"
export const routeUpdatePlace = "/api/v1/places/{id}"
export const routeDeletePlace = "/api/v1/places/{id}"
export const routeGetPlace = "/api/v1/places/{id}"
export const routeGetPlaces = "/api/v1/places"
export const routeGetPlacesTree = "/api/v1/places/tree"
export const routeAddThing = "/api/v1/things"
export const routeGetPlaceThings = "/api/v1/places/{id}/things"

export function request(method, url, data) {
    let xhr = new XMLHttpRequest();

    xhr.open(method, url, false);

    xhr.setRequestHeader('Accept', 'application/json')
    xhr.setRequestHeader('Content-Type', 'application/json')
    xhr.setRequestHeader('Authorization', 'Basic ' + getToken())

    if (data !== undefined) {
        xhr.send(JSON.stringify(data));
    } else {
        xhr.send();
    }

    return {
        data: JSON.parse(xhr.responseText),
        status: xhr.status,
    };
}