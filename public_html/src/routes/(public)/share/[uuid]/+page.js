import axios from "axios";
import {env} from "$env/dynamic/public";
import {error} from "@sveltejs/kit";

/** @type {import('./$types').PageLoad} */
export async function load(event) {
    let result = null
    await axios.get(`${env.PUBLIC_BACKEND_URL}/explorer/share/${event.params.uuid}`)
        .then(({data}) => {
            result = data
        })
        .catch(() => error(404, 'Not found'))

    return {
        share: result.data
    }
}