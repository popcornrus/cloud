<script>
    import {env} from "$env/dynamic/public";
    import Image from "$lib/components/Image.svelte";
    import Plyr from 'plyr';
    import "plyr/dist/plyr.css";

    export let file = null,
        modalOpen = false;

    $: if (file?.uuid?.length > 0) {
        modalOpen = true

        if (file?.type.IsVideo()) {
            setTimeout(() => {
                new Plyr('#player');
            }, 500)
        }
    }
</script>

<div class:hidden={!modalOpen} class="w-full h-full fixed top-0 start-0 z-[80] overflow-x-hidden overflow-y-auto backdrop-blur bg-gray-600/60">
    <div class="mt-7 opacity-100 duration-500 mt-0 ease-out transition-all sm:max-w-[45%] sm:w-full m-3 sm:mx-auto min-h-[calc(100%-3.5rem)] flex items-center">
        <div class="w-full flex flex-col bg-white border shadow-sm rounded-xl dark:bg-gray-800 dark:border-gray-700 dark:shadow-slate-700/[.7]">
            <div class="flex justify-between items-center py-3 px-4 border-b dark:border-gray-700">
                <h3 class="font-bold text-gray-800 dark:text-white">
                    {file?.name}
                </h3>
                <button type="button" on:click={() => modalOpen = false} class="flex justify-center items-center w-7 h-7 text-sm font-semibold rounded-full border border-transparent text-gray-800 hover:bg-gray-100 disabled:opacity-50 disabled:pointer-events-none dark:text-white dark:hover:bg-gray-700 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600" data-hs-overlay="#hs-vertically-centered-modal">
                    <span class="sr-only">Close</span>
                    <svg class="flex-shrink-0 w-4 h-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
                </button>
            </div>
            <div class="p-4 overflow-y-auto">
                {#if file?.type.IsImage()}
                    <Image src="{env.PUBLIC_BACKEND_URL}/explorer/files/{file?.uuid}" alt="{file?.name}" />
                {:else if file?.type.IsVideo()}
                    <video src="{env.PUBLIC_BACKEND_URL}/explorer/files/{file?.uuid}" id="player"

                    >
                        <track kind="captions" src="{env.PUBLIC_BACKEND_URL}/explorer/files/{file?.uuid}">
                    </video>
                {:else}
                    <div class="flex justify-center items-center w-full h-full p-24">
                        <div class="flex flex-col justify-center items-center gap-y-2 text-gray-500">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-12 h-12">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 0 0 5.636 5.636m12.728 12.728A9 9 0 0 1 5.636 5.636m12.728 12.728L5.636 5.636" />
                            </svg>
                            <span class="">No preview available</span>
                        </div>
                    </div>
                {/if}
            </div>
            <div class="flex justify-end items-center gap-x-2 py-3 px-4 border-t dark:border-gray-700">
                <button type="button" on:click={() => modalOpen = false} class="py-2 px-3 inline-flex items-center gap-x-2 text-sm font-medium rounded-lg border border-gray-200 bg-white text-gray-800 shadow-sm hover:bg-gray-50 disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:border-gray-700 dark:text-white dark:hover:bg-gray-800 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600" data-hs-overlay="#hs-vertically-centered-modal">
                    Close
                </button>
            </div>
        </div>
    </div>
</div>