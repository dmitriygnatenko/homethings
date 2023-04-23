"use strict"

import {getToken} from "../auth/auth.js";

export const statusOK = 200

export const methodGet = "GET"
export const methodPost = "POST"
export const methodPut = "PUT"
export const methodDelete = "DELETE"

// Auth
export const routeLogin = "/api/v1/auth/login"
export const routeCheckAuth = "/api/v1/auth/check"
// Users
export const routeAddUser = "/api/v1/users"
export const routeUpdateUser = "/api/v1/users"
// Places
export const routeGetPlaces = "/api/v1/places"
export const routeGetPlacesTree = "/api/v1/places/tree"
export const routeGetPlace = "/api/v1/places/{placeId}"
export const routeAddPlace = "/api/v1/places"
export const routeUpdatePlace = "/api/v1/places/{placeId}"
export const routeDeletePlace = "/api/v1/places/{placeId}"
export const routeGetNestedPlaces = "/api/v1/places/{parentPlaceId}/nested"
// Things
export const routeGetThing = "/api/v1/things/{thingId}"
export const routeGetPlaceThings = "/api/v1/things/place/{placeId}"
export const routeAddThing = "/api/v1/things"
export const routeUpdateThing = "/api/v1/things/{thingId}"
export const routeDeleteThing = "/api/v1/things/{thingId}"
export const routeSearchThings = "/api/v1/things/search/{search}"
// Images
export const routeGetPlaceImages = "/api/v1/images/place/{placeId}"
export const routeGetThingImages = "/api/v1/images/thing/{thingId}"
export const routeDeletePlaceImages = "/api/v1/images/place/{imageId}"
export const routeDeleteThingImages = "/api/v1/images/thing/{imageId}"
export const routeAddImage = "/api/v1/images"
// Tags
export const routeGetTags = "/api/v1/tags"
export const routeGetTag = "/api/v1/tags/{tagId}"
export const routeGetThingTags = "/api/v1/tags/thing/{thingId}"
export const routeAddTag = "/api/v1/tags"
export const routeUpdateTag = "/api/v1/tags/{tagId}"
export const routeDeleteTag = "/api/v1/tags/{tagId}"
export const routeAddThingTag = "/api/v1/tags/{tagId}/thing/{thingId}"
export const routeDeleteThingTag = "/api/v1/tags/{tagId}/thing/{thingId}"

export function jsonRequest(method, url, data) {
    let xhr = new XMLHttpRequest();
    const token = getToken()
    const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;

    xhr.open(method, url, false);

    xhr.setRequestHeader("Accept", "application/json")
    xhr.setRequestHeader("Content-Type", "application/json")
    xhr.setRequestHeader("Timezone", timezone)

    if (token !== "") {
        xhr.setRequestHeader("Authorization", "Bearer " + token)
    }

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

export function formDataRequest(method, url, data) {
    let xhr = new XMLHttpRequest();
    const token = getToken()
    const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;

    xhr.open(method, url, false);

    xhr.setRequestHeader("Accept", "application/json")
    xhr.setRequestHeader("Timezone", timezone)

    if (token !== "") {
        xhr.setRequestHeader("Authorization", "Bearer " + token)
    }

    xhr.send(data);

    return {
        data: JSON.parse(xhr.responseText),
        status: xhr.status,
    };
}