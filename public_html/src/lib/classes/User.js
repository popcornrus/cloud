'use strict';

import {env} from "$env/dynamic/public";
import axios from "axios";

export class User {
    constructor(token) {
        this.axios = axios.create({
            baseURL: `${env.PUBLIC_BACKEND_URL}/api/v1/users`,
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        })
    }

    async login(email, password) {
        let toast = null;

        await this.axios.post(`/sign-in`, {
            email: email,
            password: password
        }).then(({data}) => {
            localStorage.setItem('token', data.data.token);

            toast = {
                status: 'success',
                message: 'Successfully logged in',
            }
        }).catch((error) => {
            toast = {
                status: 'error',
                message: error.response?.data.message ?? error.error(),
            }
        })

        return toast;
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