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

        this.WebSocket = null;
    }

    async list() {
        await this.axios.get(``).then(({data}) => {
            this.files = data.data ?? [];
        })

        return this.files;
    }

    async data(uuid) {
        const _self = this
        let file = new Promise(() => {});

        await this.axios.get(uuid + `/data`).then(({data}) => {
            file = data.data

            if (!_self.files.find(file => file.uuid === data.data.uuid)) {
                _self.files.push(data.data)
            }
        })

        return await file
    }

    findByUuid(uuid) {
        return this.files.find(file => file.uuid === uuid);
    }

    async download(uuid) {
        const file = await this.findByUuid(uuid)

        const streamUrl = `${uuid}/download`;

        const supportsFileSystemAccess =
            "showSaveFilePicker" in window &&
            (() => {
                try {
                    return window.self === window.top;
                } catch {
                    return false;
                }
            })();

        if (supportsFileSystemAccess) {
            try {
                const handle = await window.showSaveFilePicker({
                    suggestedName: file.name
                });

                this.axios.get(streamUrl, {
                    responseType: "blob",
                    onDownloadProgress: (progressEvent) => {
                    }
                }).then(async response => {
                    const writable = await handle.createWritable();
                    await writable.write(response.data);
                    await writable.close();
                }).catch(error => {
                    console.error("Error fetching and streaming:", error);
                });

                return;
            } catch (err) {
                console.error("Error saving file:", err);
                return
            }
        }

        let blob = await this.axios.get(streamUrl, {
            responseType: "blob",
            onDownloadProgress: (progressEvent) => {
            }
        }).then(async response => {
            return response.data;
        }).catch(error => {
            console.error("Error fetching and streaming:", error);
        });

        // Fallback if the File System Access API is not supported…
        // Create the blob URL.
        const blobURL = URL.createObjectURL(blob);
        // Create the `<a download>` element and append it invisibly.
        const a = document.createElement("a");
        a.href = blobURL;
        a.download = file.name;
        a.style.display = "none";
        document.body.append(a);
        // Programmatically click the element.
        a.click();

        // Revoke the blob URL and remove the element.
        setTimeout(() => {
            URL.revokeObjectURL(blobURL);
            a.remove();
        }, 1000);
    }

    async rename(uuid) {
        const _self = this
        const appendExtension = (fileName, extension) => {
            if (fileName.endsWith(extension)) {
                return fileName;
            }

            return `${fileName}.${extension}`;
        }

        const file = document.getElementById(`file-${uuid}`);
        const fileName = file.querySelector('.file-title');
        const fileInput = file.querySelector('.file-input');

        fileName.style.display = 'none';
        fileInput.style.display = 'block';
        fileInput.focus();

        fileInput.value = fileName.innerText;

        fileInput.addEventListener('blur', () => {
            fileName.style.display = 'block';
            fileInput.style.display = 'none';
        });

        fileInput.addEventListener('keydown', async (e) => {
            if (e.key === 'Enter') {
                fileName.innerText = appendExtension(fileInput.value, fileName.innerText.Extension());

                await _self.axios.patch(`/${uuid}`, {
                    name: name
                }).then(() => {

                });

                fileName.style.display = 'block';
                fileInput.style.display = 'none';
            }
        });
    }

    async move(uuid) {
        console.log('move')
    }

    async delete(uuid) {
        const _self = this

        await this.axios.delete(`/${uuid}`).then(() => {
            _self.files = _self.files.filter(file => file.uuid !== uuid);
        })
    }

    webSocket(user) {
        const socket = new WebSocket(`${env.PUBLIC_BACKEND_URL.replace('http', 'ws')}/ws/echo`);

        this.WebSocket = socket

        socket.onopen = function(e) {
            console.log("[open] Connection established");

            socket.send(JSON.stringify({
                channel: "files." + user.uuid,
                event: "action:subscribe",
                data: {}
            }))
        };

        socket.onmessage = function(event) {
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

        return socket
    }
}