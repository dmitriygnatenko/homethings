import {getAuth} from './auth.js';

export function request(method, url) {
    const res = {
        data: {},
        status: undefined,
        error: undefined,
    };

    fetch(url, {
        method: method,
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': getAuth(),
        },

    }).then((response) => {
        res.status = response.status
        return response.json()
    }).then((data) => {
        res.data = data
    }).catch((e) => {
        let error = e.toString()
        error = error.replace(/["']/g,'')
        res.error = error
    })

    return res
}