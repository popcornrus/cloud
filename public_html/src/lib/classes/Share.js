import axios from "axios";
import {env} from "$env/dynamic/public";

export class Share {
    constructor(token) {
        this.axios = axios.create({
            baseURL: `${env.PUBLIC_BACKEND_URL}/api/v1/explorer/share`,
            timeout: 1000,
            headers: {'Authorization': `Bearer ${token}`}
        })

        this.share = {
            uuid: null,
            state: 'prepare',
            type: null,
            url: null,
            pinCode: null,
            expiresAt: null,
        }
    }

    async create() {
        return this.axios.post('/create', this.share)
    }
}