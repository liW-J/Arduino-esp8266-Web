import axios from "../config"
import type { tempLog, userInfo, fingerLog, operateLog } from '@/type/admin'
import { useToCamelCase } from '@/hooks/useToCamelCase';

export function login(name: string, password: string): Promise<string> {
    return new Promise((resolve, reject) => {
        axios.post('/admin/login', {
            name,
            password
        }).then(res => {
            if (res.data.code === 200) resolve(res.data.data.token)
            else reject(res.data.msg)
        })
    })
}

export function getTempLogList(): Promise<Array<tempLog>> {
    return new Promise((resolve, reject) => {
        axios.get('/admin/tempLog').then(res => {
            if (res.data.code === 200) {
                const data = useToCamelCase(res.data.data) as unknown as Array<tempLog>
                resolve(data)
            } else {
                reject(res.data.msg)
            }
        })
    })
}

export function getUserList(): Promise<Array<userInfo>> {
    return new Promise((resolve, reject) => {
        axios.get('/admin/userList').then(res => {
            if (res.data.code === 200) {
                const data = useToCamelCase(res.data.data) as unknown as Array<userInfo>
                resolve(data)
            } else {
                reject(res.data.msg)
            }
        })
    })
}

export function getFingerLogList(): Promise<Array<fingerLog>> {
    return new Promise((resolve, reject) => {
        axios.get('/admin/fingerLog').then(res => {
            if (res.data.code === 200) {
                const data = useToCamelCase(res.data.data) as unknown as Array<fingerLog>
                resolve(data)
            } else {
                reject(res.data.msg)
            }
        })
    })
}

export function getOperateLogList(): Promise<Array<operateLog>> {
    return new Promise((resolve, reject) => {
        axios.get('/admin/operateLog').then(res => {
            if (res.data.code === 200) {
                resolve(res.data.data)
            } else {
                reject(res.data.msg)
            }
        })
    })
}

export function userSignUp(user_name: string, work_num: string): Promise<string> {
    return new Promise((resolve, reject) => {
        axios.post('/admin/signUp', {
            user_name,
            work_num
        }).then(res => {
            if (res.data.code === 200) {
                resolve(res.data.msg)
            } else {
                reject(res.data.msg)
            }
        })
    })
}

export function deleteUser(work_num: string): Promise<string> {
    return new Promise((resolve, reject) => {
        axios.post('/admin/deleteUser', {
            work_num
        }).then(res => {
            if (res.data.code === 200) {
                resolve(res.data.msg)
            } else {
                reject(res.data.msg)
            }
        })
    })
}