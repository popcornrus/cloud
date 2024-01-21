import axios from "axios";
import {env} from "$env/dynamic/public";

export class File {
    constructor(token) {
        this.axios = axios.create({
            baseURL: `${env.PUBLIC_BACKEND_URL}/api/v1/explorer/files`,
            timeout: 1000,
            headers: {'Authorization': `Bearer ${token}`}
        })

        this.token = token;
        this.files = [];
    }

    async list() {
        await this.axios.get(``).then(({data}) => {
            this.files = data.data;
        })

        return this.files;
    }

    findByUuid(uuid) {
        return this.files.find(file => file.uuid === uuid);
    }

    async download(uuid) {
        console.log('download')
    }

    async rename(uuid, name) {
        await this.axios.patch(`/${uuid}`, {
            name: name
        }).then(() => {

        });
    }

    async move(uuid) {
        console.log('move')
    }

    async delete(uuid) {
        console.log('delete')
    }

    webSocket(user) {
        const socket = new WebSocket(`${env.PUBLIC_BACKEND_URL.replace('http', 'ws')}/ws/echo`);

        socket.onopen = function(e) {
            console.log("[open] Connection established");

            socket.send(JSON.stringify({
                channel: "files." + user.uuid,
                event: "action:subscribe",
                data: {}
            }))
        };

        socket.onmessage = function(event) {
            console.log(`[message] Data received from server: ${event.data}`);
        };

        socket.onclose = function(event) {
            if (event.wasClean) {
                console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
            } else {
                console.log('[close] Connection died');
            }
        };

        socket.onerror = function(error) {
            console.log(`[error] ${error.message}`);
        };
    }
}