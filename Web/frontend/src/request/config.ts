import _axios from 'axios'

const axios = _axios.create({
    baseURL: 'http://ip:9090/'
})

function initState() {
    return {
        token: localStorage.getItem('token'),
        adminToken: localStorage.getItem('adminToken')
    }
}

const state = initState()

axios.interceptors.request.use(
    function (config) {
        if (document.location.hash.includes('admin')) {
            if (!state.adminToken) {
                state.adminToken = localStorage.getItem('adminToken')
            }
            config.headers!.token = state.adminToken
        } else {
            if (!state.token) {
                state.token = localStorage.getItem('token')
            }
            config.headers!.token = state.token
        }
        return config
    },
    async (err) => {
        console.log(err)
    }
)

axios.interceptors.response.use(
    function (response) {
        if (response?.data?.code === 401) {
            if (document.location.hash.includes('admin')) {
                window.location.hash = '/admin/login'
            } else {
                window.location.hash = '/'
            }
        }
        return response
    },
    async (err) => {
        console.log(err)
    }
)

export default axios