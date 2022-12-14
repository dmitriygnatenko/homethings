import {getAuth} from './auth.js';

export function request(method, url) {
    let xhr = new XMLHttpRequest();
    xhr.open(method, url, false);
    xhr.setRequestHeader('Accept', 'application/json')
    xhr.setRequestHeader('Content-Type', 'application/json',)
    xhr.setRequestHeader('Authorization', getAuth())
    xhr.send();

    return {
        data: xhr.responseText,
        status: xhr.status,
    };
}