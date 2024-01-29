<script>
    import {onMount} from "svelte";
    import { Share } from "$lib/classes/Share";
    import { authToken } from "$lib/stores/token";
    import moment from "moment";

    export let modalOpen = false,
        file = null;

    let shareAction = new Share($authToken, file.uuid),
        share = shareAction.share;

    let usingPinCode = 0,
        pinCode = [];

    let date = null,
        time = null;

    let sharingType = '3';

    onMount(async () => {
        HSStaticMethods.autoInit();

        share = await shareAction.get();
        if (share.pinCode) {
            usingPinCode = true;
            pinCode = share.pinCode.split('');
        }

        if (share.expiresAt) {
            date = moment(share.expiresAt).format('YYYY-MM-DD');
            time = moment(share.expiresAt).format('HH:mm');
        }

        sharingType = share.type.toString();
    });

    $: {
        shareAction.share.type = parseInt(sharingType);

        switch (shareAction.share.type) {
            case 1:
                shareAction.share.expiresAt = null;
                shareAction.share.downloadLimit = null;
                break;
            case 2:
                if (!share.downloadLimit) {
                    share.downloadLimit = 1;
                }

                if (share.downloadLimit?.length === 0) {
                    share.downloadLimit = 1;
                } else {
                    share.downloadLimit = parseInt(share.downloadLimit?.toString().replace(/\D/g, ""));
                }

                shareAction.share.downloadLimit = share.downloadLimit;

                if (date && time) {
                    shareAction.share.expiresAt = moment(`${date} ${time}`).format('YYYY-MM-DD HH:mm:ss');
                } else {
                    shareAction.share.expiresAt = null;
                }
                break;
            case 3:
                shareAction.share.expiresAt = null;
                shareAction.share.downloadLimit = null;
                break;
        }

        if (usingPinCode && pinCode.length === 4) {
            shareAction.share.pinCode = pinCode.join('');
        } else {
            shareAction.share.pinCode = null;
        }
    }

    $: share = shareAction.share;
</script>

<div class:hidden={!modalOpen} class="w-full h-full fixed top-0 start-0 z-[90] overflow-x-hidden overflow-y-auto backdrop-blur bg-gray-600/60" id="share-modal" data-uuid="{share.uuid}">
    <div class="mt-7 opacity-100 duration-500 mt-0 ease-out transition-all sm:max-w-[50%] xl:max-w-[30%] sm:w-full m-3 sm:mx-auto min-h-[calc(100%-3.5rem)] flex justify-center items-center">
        <div class="flex flex-col bg-white border shadow-sm rounded-xl pointer-events-auto dark:bg-gray-800 dark:border-gray-700 dark:shadow-slate-700/[.7] w-full">
            <div class="flex justify-between items-center py-3 px-4 border-b dark:border-gray-700">
                <h3 class="font-bold text-gray-800 dark:text-white">
                    Share file "{file.name}"
                </h3>
                <button type="button" on:click={() => modalOpen = false} class="flex justify-center items-center w-7 h-7 text-sm font-semibold rounded-full border border-transparent text-gray-800 hover:bg-gray-100 disabled:opacity-50 disabled:pointer-events-none dark:text-white dark:hover:bg-gray-700 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600" data-hs-overlay="#hs-small-modal">
                    <span class="sr-only">Close</span>
                    <svg class="flex-shrink-0 w-4 h-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
                </button>
            </div>
            <div class="p-4 overflow-y-auto">
                <div class="flex justify-between items-center mb-4">
                    <div class="flex items-center h-12">
                        <label class="relative inline-flex items-center cursor-pointer h-max">
                            <input type="checkbox" class="sr-only peer" bind:checked={usingPinCode}>
                            <span class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></span>
                            <span class="ms-3 text-sm text-gray-900 dark:text-gray-300">Use access pin code</span>
                        </label>
                    </div>

                    <div class="flex space-x-1 h-12 font-bold" data-hs-pin-input class:hidden={!usingPinCode}>
                        <input bind:value={pinCode[0]} type="text" class="block w-[38px] text-center bg-transparent border-t-transparent border-b-2 border-x-transparent border-b-gray-200 text-lg focus:border-t-transparent focus:border-x-transparent focus:border-b-blue-500 focus:ring-0 focus:outline-none disabled:opacity-50 disabled:pointer-events-none dark:border-b-gray-700 dark:text-gray-400 dark:focus:ring-gray-600 dark:focus:border-b-gray-600" placeholder="⚬" data-hs-pin-input-item>
                        <input bind:value={pinCode[1]} type="text" class="block w-[38px] text-center bg-transparent border-t-transparent border-b-2 border-x-transparent border-b-gray-200 text-lg focus:border-t-transparent focus:border-x-transparent focus:border-b-blue-500 focus:ring-0 focus:outline-none disabled:opacity-50 disabled:pointer-events-none dark:border-b-gray-700 dark:text-gray-400 dark:focus:ring-gray-600 dark:focus:border-b-gray-600" placeholder="⚬" data-hs-pin-input-item>
                        <input bind:value={pinCode[2]} type="text" class="block w-[38px] text-center bg-transparent border-t-transparent border-b-2 border-x-transparent border-b-gray-200 text-lg focus:border-t-transparent focus:border-x-transparent focus:border-b-blue-500 focus:ring-0 focus:outline-none disabled:opacity-50 disabled:pointer-events-none dark:border-b-gray-700 dark:text-gray-400 dark:focus:ring-gray-600 dark:focus:border-b-gray-600" placeholder="⚬" data-hs-pin-input-item>
                        <input bind:value={pinCode[3]} type="text" class="block w-[38px] text-center bg-transparent border-t-transparent border-b-2 border-x-transparent border-b-gray-200 text-lg focus:border-t-transparent focus:border-x-transparent focus:border-b-blue-500 focus:ring-0 focus:outline-none disabled:opacity-50 disabled:pointer-events-none dark:border-b-gray-700 dark:text-gray-400 dark:focus:ring-gray-600 dark:focus:border-b-gray-600" placeholder="⚬" data-hs-pin-input-item>
                    </div>
                </div>

                <div class="mb-4">
                    <div class="mb-2">
                        <label for="sharing-type" class="block text-sm font-medium mb-2 dark:text-white">Type of sharing</label>
                        <select bind:value={sharingType} id="sharing-type" class="border py-3 px-4 pe-9 block w-full border-gray-200 rounded-lg focus:outline-none text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:border-gray-700 dark:text-gray-400 dark:focus:ring-gray-600">
                            <option value="1">Delete after open</option>
                            <option value="2">Set number of openings</option>
                            <option value="3">Don't delete share link</option>
                        </select>
                    </div>

                    {#if parseInt(sharingType) === 2}
                        <div class="bg-white border border-gray-200 rounded-lg dark:bg-slate-700 dark:border-gray-700 mb-2">
                            <div class="w-full flex justify-between items-center gap-x-1">
                                <div class="grow py-2 px-3">
                                    <input class="focus:outline-none w-full p-0 bg-transparent border-0 text-gray-800 focus:ring-0 dark:text-white" min="1" type="text" bind:value="{share.downloadLimit}">
                                </div>
                                <div class="flex items-center -gap-y-px divide-x divide-gray-200 border-s border-gray-200 dark:divide-gray-700 dark:border-gray-700">
                                    <button on:click={() => share.downloadLimit > 1 ? share.downloadLimit-- : 0 } type="button" class="w-10 h-10 inline-flex justify-center items-center gap-x-2 text-sm font-medium last:rounded-e-lg bg-white text-gray-800 hover:bg-gray-50 disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:text-white dark:hover:bg-gray-800 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
                                        <svg class="flex-shrink-0 w-3.5 h-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"/></svg>
                                    </button>
                                    <button on:click={() => share.downloadLimit++} type="button" class="w-10 h-10 inline-flex justify-center items-center gap-x-2 text-sm font-medium last:rounded-e-lg bg-white text-gray-800 hover:bg-gray-50 disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:text-white dark:hover:bg-gray-800 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
                                        <svg class="flex-shrink-0 w-3.5 h-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"/><path d="M12 5v14"/></svg>
                                    </button>
                                </div>
                            </div>
                        </div>
                    {/if}

                    {#if ![1,3].includes(parseInt(sharingType))}
                        <div class="sm:flex rounded-lg shadow-sm">
                            <input type="date" min="{moment().format('YYYY-MM-DD')}" bind:value={date} class="focus:outline-none border py-3 px-4 block w-full border-gray-200 shadow-sm -mt-px -ms-px first:rounded-t-lg last:rounded-b-lg sm:first:rounded-s-lg sm:mt-0 sm:first:ms-0 sm:first:rounded-se-none sm:last:rounded-es-none sm:last:rounded-e-lg text-sm relative focus:z-10 focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:border-gray-700 dark:text-gray-400 dark:focus:ring-gray-600">
                            <input type="time" bind:value={time} class="focus:outline-none border py-3 px-4 block w-full border-gray-200 shadow-sm -mt-px -ms-px first:rounded-t-lg last:rounded-b-lg sm:first:rounded-s-lg sm:mt-0 sm:first:ms-0 sm:first:rounded-se-none sm:last:rounded-es-none sm:last:rounded-e-lg text-sm relative focus:z-10 focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:border-gray-700 dark:text-gray-400 dark:focus:ring-gray-600">
                        </div>
                    {/if}
                </div>

                <div>
                    <label for="hs-leading-icon" class="block text-sm font-medium mb-2 dark:text-white">Share URL</label>
                    <div class="relative text-gray-600">
                        <input type="text" id="hs-leading-icon" bind:value={share.url} name="hs-leading-icon" readonly class="py-3 px-4 focus:outline-none ps-11 block text-inherit w-full border-gray-200 border shadow-sm rounded-lg text-sm disabled:pointer-events-none dark:bg-slate-900 dark:border-gray-700 dark:text-gray-400 dark:focus:ring-gray-600" placeholder="There will be a link for sharing">
                        <button class="absolute inset-y-0 start-0 flex items-center pointer-events-none z-20 ps-4">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="flex-shrink-0 h-4 w-4 dark:text-gray-600">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M15.666 3.888A2.25 2.25 0 0 0 13.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 0 1-.75.75H9a.75.75 0 0 1-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 0 1-2.25 2.25H6.75A2.25 2.25 0 0 1 4.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 0 1 1.927-.184" />
                            </svg>
                        </button>
                    </div>
                </div>
            </div>
            <div class="flex justify-end items-center gap-x-2 py-3 px-4 border-t dark:border-gray-700">
                <button type="button" on:click={() => modalOpen = false} class="py-2 px-3 inline-flex items-center gap-x-2 text-sm font-medium rounded-lg border border-gray-200 bg-white text-gray-800 shadow-sm hover:bg-gray-50 disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:border-gray-700 dark:text-white dark:hover:bg-gray-800 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600" data-hs-overlay="#hs-small-modal">
                    Close
                </button>
                {#if share.state === 'prepare'}
                    <button on:click={async () => share = await shareAction.create() } type="button" class="py-2 px-3 inline-flex items-center gap-x-2 text-sm font-semibold rounded-lg border border-transparent bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:pointer-events-none dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
                        Generate link
                    </button>
                {:else}
                    <button on:click={async () => share = await shareAction.update() } type="button" class="py-2 px-3 inline-flex items-center gap-x-2 text-sm font-semibold rounded-lg border border-transparent bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:pointer-events-none dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
                        Update link
                    </button>
                {/if}
            </div>
        </div>
    </div>
</div>