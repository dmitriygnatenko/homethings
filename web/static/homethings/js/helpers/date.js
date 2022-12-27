'use strict'

export function formatDate(str) {
    let date = new Date(str);
    let res = date.getDate() + " "

    switch (date.getMonth()) {
        case 0:
            res += "января"
            break
        case 1:
            res += "февраля"
            break
        case 2:
            res += "марта"
            break
        case 3:
            res += "апреля"
            break
        case 4:
            res += "мая"
            break
        case 5:
            res += "июня"
            break
        case 6:
            res += "июля"
            break
        case 7:
            res += "августа"
            break
        case 8:
            res += "сентября"
            break
        case 9:
            res += "октября"
            break
        case 10:
            res += "ноября"
            break
        case 11:
            res += "декабря"
    }

    res += " " + date.getFullYear()
    return res
}