import axios from "axios";
import {env} from "$env/dynamic/public";

export class Uploader {
    constructor(token) {
        this.axios = axios.create({
            baseURL: `${env.PUBLIC_BACKEND_URL}/api/v1/explorer/files`,
            timeout: 1000,
            headers: {'Authorization': `Bearer ${token}`}
        })

        this.files = [];
    }

    list() {
        return this.files;
    }

    async add(file) {
        this.files.push(file);

        this.prepare(file.uuid);

        this.__storeInLocalStorage();
    }

    async pause(uuid) {
        const fileIndex = this.files.findIndex(file => file.uuid === uuid);
        this.files[fileIndex].state = 'paused';

        this.__storeInLocalStorage();
    }

    async resume(uuid) {
        const fileIndex = this.files.findIndex(file => file.uuid === uuid);
        this.files[fileIndex].state = 'processing';

        this.__storeInLocalStorage();
    }

    async cancel(uuid) {
        const fileIndex = this.files.findIndex(file => file.uuid === uuid);
        this.files[fileIndex].state = 'canceled';

        this.__storeInLocalStorage();
    }

    async prepare(uuid) {
        const self = this

        const fileIndex = this.files.findIndex(file => file.uuid === uuid);
        this.files[fileIndex].state = 'processing';

        this.axios.post('/prepare', {
            name: this.files[fileIndex].name,
            size: this.files[fileIndex].size,
            type: this.files[fileIndex].binary.type,
        }).then(({data}) => {
            this.files[fileIndex].chunkSize = data.data?.chunk_size;
            this.files[fileIndex].uploadUrl = data.data?.url;
            this.files[fileIndex].state = 'processing';

            self.upload(uuid)
        })


        this.__storeInLocalStorage();
    }

    async upload(uuid) {
        const fileIndex = this.files.findIndex(file => file.uuid === uuid);

        const file = this.files[fileIndex];

        const chunkSize = file.chunkSize;
        const chunks = Math.ceil(file.size / chunkSize);

        let percent = 0;

        this.files[fileIndex].state = 'processing';

        let counterOfRequests = 0;
        let chunkProgress = 0;


        for (let i = 0; i < chunks; i++) {
            if (counterOfRequests >= 75) {
                await new Promise(resolve => setTimeout(resolve, 1000));
            }

            const start = i * chunkSize;
            const end = Math.min(file.size, start + chunkSize);

            const chunk = file.binary.slice(start, end);

            this.axios.post(file.uploadUrl, {
                chunk: chunk,
            }, {
                timeout: 3000,
                maxContentLength: 2 * 1024 * 1024,
                headers: {
                    'Content-Type': 'multipart/form-data',
                    'Content-Range': `bytes ${start}-${end - 1}/${file.size}`
                }
            }).then(() => {
                chunkProgress++;
                percent = Math.round(chunkProgress / chunks * 100);
                this.files[fileIndex].progress = percent;

                counterOfRequests--
            })

            counterOfRequests++;
        }

        this.files[fileIndex].state = 'done';
        this.__storeInLocalStorage();
    }

    __storeInLocalStorage() {
        localStorage.setItem('files', JSON.stringify(this.files));
    }
}