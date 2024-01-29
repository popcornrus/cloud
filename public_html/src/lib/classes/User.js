'use strict';

import {env} from "$env/dynamic/public";
import axios from "axios";
import {Toast} from "$lib/classes/Toast.js";

export class User {
    constructor(token) {
        this.axios = axios.create({
            baseURL: `${env.PUBLIC_BACKEND_URL_API}/users`,
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        })
    }

    async login(email, password) {
        return await this.axios.post(`/sign-in`, {
            email: email,
            password: password
        }).then(({data}) => {
            localStorage.setItem('token', data.data.token);

            Toast({
                type: 'success',
                message: 'Successfully logged in',
                duration: 3000,
            })

            return true;
        }).catch((error) => {
            Toast({
                type: 'error',
                message: error.response?.data.message ?? error.message,
                duration: 3000,
            })

            return false
        })
    }

    async signUp(username, email, password) {
        let toast = null;

        await this.axios.post(`/sign-up`, {
            username: username,
            email: email,
            password: password
        }).then(({ data }) => {
            localStorage.setItem('token', data.data.token);

            toast = {
                status: 'success',
                message: 'Successfully signed up',
            }
        }).catch((error) => {
            toast = {
                status: 'error',
                message: error.response.data.message,
            }
        })

        return toast;
    }

    async me() {
        setTimeout(() => {}, 1000)
        return await this.axios.get('/me');
    }
}