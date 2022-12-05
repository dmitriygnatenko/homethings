const tokenKey = "x326gjdskfh8ndhskjfhs"

function getAuth() {
    //   token = localStorage.getItem(tokenKey)
    //   console.log(token)
    var user = "admin"
    var password = "12345"
    var token = user + ":" + password;
    const auth = "Basic " + btoa(token)
    return auth
}

function serverRequest(method, url) {
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
    }).catch((error) => {
        res.error = error
    })

    return res
}