import {env} from "$env/dynamic/public";

export class WebSocketClient {
    constructor(user) {
        this.user = user;
        this.events = null;

        return this.connect();
    }

    connect() {
        const socket = new WebSocket(env.PUBLIC_WEBSOCKET_URL);

        socket.onopen = (e) => {
            socket.send(JSON.stringify({
                channel: "files." + this.user.uuid,
                event: "action:subscribe",
                data: {}
            }))
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