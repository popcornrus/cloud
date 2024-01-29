<script>
    import {onMount} from "svelte";
    import axios from "axios";
    import {env} from "$env/dynamic/public";
    import Preview from "$lib/components/explorer/share/Preview.svelte";
    export let data;

    let pinCode = [];
    let file = null;

    onMount(() => {
        HSStaticMethods.autoInit();

        if (data.share?.file) {
            file = data.share.file;
        }
    })

    const checkPinCode = async () => {
        const fd = new FormData();

        fd.append('pinCode', pinCode.join(''));
        fd.append('shareId', data.share.uuid);

        const response = await axios.post(`${env.PUBLIC_FRONTEND_URL}/share?/pin`, fd, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        })

        file = JSON.parse(JSON.parse(response.data.data)[0])?.file
    }

    $: if (pinCode.length === 4) {
        checkPinCode();
    }
</script>

{#if data.share.pin_code && !file}
<div class="absolute inset-0 m-auto w-max px-4 bg-gray-500/60 rounded-xl h-max py-2">
    <label class="text-center block font-bold text-white/80 text-2xl uppercase" for="">Pin Code</label>
    <div class="flex justify-center items-center h-20">
        <div class="flex space-x-3 h-12" data-hs-pin-input>
            <input type="text" bind:value={pinCode[0]} class="focus:outline-none block w-[38px] text-center bg-2ray-100 border-transparent rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-gray-700 dark:border-transparent dark:text-gray-400 dark:focus:ring-gray-600" placeholder="⚬" data-hs-pin-input-item>
            <input type="text" bind:value={pinCode[1]} class="focus:outline-none block w-[38px] text-center bg-2ray-100 border-transparent rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-gray-700 dark:border-transparent dark:text-gray-400 dark:focus:ring-gray-600" placeholder="⚬" data-hs-pin-input-item>
            <input type="text" bind:value={pinCode[2]} class="focus:outline-none block w-[38px] text-center bg-2ray-100 border-transparent rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-gray-700 dark:border-transparent dark:text-gray-400 dark:focus:ring-gray-600" placeholder="⚬" data-hs-pin-input-item>
            <input type="text" bind:value={pinCode[3]} class="focus:outline-none block w-[38px] text-center bg-2ray-100 border-transparent rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-gray-700 dark:border-transparent dark:text-gray-400 dark:focus:ring-gray-600" placeholder="⚬" data-hs-pin-input-item>
        </div>
    </div>
</div>
{:else if file}
<div class="absolute inset-0 m-auto w-max bg-gray-500/60 px-4 bg-white rounded-xl h-max py-2">
    <div class="flex space-x-3 h-12">
        <Preview file={file} />
    </div>
</div>
{/if}
