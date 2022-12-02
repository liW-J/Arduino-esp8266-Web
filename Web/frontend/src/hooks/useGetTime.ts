export const getFormatTime = (): string => {
    const date = new Date()
    const dateObj: Record<string, number | string> = {
        Y: date.getFullYear(),
        M: date.getMonth() + 1,
        D: date.getDate(),
        h: date.getHours(),
        m: date.getMinutes(),
    }

    for (let k in dateObj) {
        if (dateObj[k] < 10) {
            dateObj[k] = '0' + dateObj[k]
        }
    }

    return `${dateObj.Y}-${dateObj.M}-${dateObj.D} ${dateObj.h}:${dateObj.m}:`
}

export const getSecond = (): string => {
    let second = new Date().getSeconds() + ''
    second = second.length > 1 ? second : '0' + second
    return second
}