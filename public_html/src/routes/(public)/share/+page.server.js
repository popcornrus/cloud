import {env} from "$env/dynamic/public";
import axios from "axios";

/** @type {import('./$types').Actions} */
export const actions = {
    pin: async (event) => {
        const data = await event.request.formData()

        let hashPin = await data.get('pinCode').Hash256()

        return await axios.get(`${env.PUBLIC_BACKEND_URL}/explorer/share/${data.get('shareId')}/${hashPin.slice(0, 32)}`)
            .then(res => {
                return JSON.stringify(res.data.data)
            })
    }
};