import axios from "../config"

export function userLogin(): Promise<string> {
    return new Promise((resolve, reject) => {
        axios.get('/user/login').then(res => {
            if (res.data.code === 200) {
                resolve(res.data.data)
            } else {
                reject(res.data.msg)
            }
        })
    })
}

export function recordTemp(): Promise<string> {
    return new Promise((resolve, reject) => {
        axios.post('/user/recordTemp').then(res => {
            if (res.data.code === 200) {
                resolve(res.data.data.temperature)
            } else {
                reject(res.data.msg)
            }
        })
    })
}

export function addFinger(): Promise<string> {
    return new Promise((resolve, reject) => {
        axios.post('/user/addFinger').then(res => {
            if (res.data.code === 200) {
                resolve(res.data.msg)
            } else {
                reject(res.data.msg)
            }
        })
    })
}

export function updateFinger(update_num: number): Promise<string> {
    return new Promise((resolve, reject) => {
        axios.post('/user/updateFinger', {
            update_num
        }).then(res => {
            if (res.data.code === 200) {
                resolve(res.data.msg)
            } else {
                reject(res.data.msg)
            }
        })
    })
}


