import axios from "axios";
import {env} from "$env/dynamic/public";
import moment from "moment";

export let ShareModal = async (file) => {
    const {default: Share} = await import('$lib/components/explorer/files/Share.svelte');

    new Share({
        target: document.body,
        props: {
            file: file,
            modalOpen: true,
        }
    })
}

export class Share {
    constructor(token, fileUUID) {
        this.axios = axios.create({
            baseURL: `${env.PUBLIC_BACKEND_URL_API}/explorer/share`,
            timeout: 1000,
            headers: {'Authorization': `Bearer ${token}`}
        })

        this.share = {
            fileUUID: fileUUID,
            uuid: null,
            state: 'prepare',
            type: null,
            url: null,
            pinCode: null,
            expiresAt: null,
            downloadLimit: null,
            downloadCount: null,
        }
    }

    async get() {
        await this.axios.get(`/${this.share.fileUUID}/data`).then(({data}) => {
            this.share.state = 'edit'
            this.share.uuid = data.data.uuid;
            this.share.type = data.data.type;
            this.share.url = `${env.PUBLIC_FRONTEND_URL}/share/${data.data.uuid}`;
            this.share.pinCode = data.data.pin_code ?? null;
            this.share.expiresAt = data.data.expires_at ?? null;
            this.share.downloadLimit = data.data.download_limit ?? null;
            this.share.downloadCount = data.data.download_count ?? null;
        })

        return this.share;
    }

    async update() {
        await this.axios.put(`/${this.share.uuid}`, {
            type: this.share.type,
            download_limit: this.share.downloadLimit,
            expires_at: this.share.expiresAt,
            pin_code: this.share.pinCode,
        }).then(({data}) => {
            this.share.state = 'edit'
            this.share.uuid = data.data.uuid;
            this.share.url = `${env.PUBLIC_FRONTEND_URL}/share/${data.data.uuid}`;
        })

        return this.share;
    }

    async create() {
        await this.axios.post('/create', {
            type: this.share.type,
            file_uuid: this.share.fileUUID,
            download_limit: this.share.downloadLimit,
            expires_at: moment(this.share.expiresAt).toISOString(),
            pin_code: this.share.pinCode,
        }).then(({data}) => {
            this.share.state = 'edit'
            this.share.uuid = data.data.uuid;
            this.share.url = `${env.PUBLIC_FRONTEND_URL}/share/${data.data.uuid}`;
        })

        return this.share;
    }
}