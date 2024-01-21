<script>
    import { Uploader } from "$lib/classes/Uploader.js";
    import {onMount} from "svelte";
    import {authToken} from "$lib/stores/token.js";
    import DropZone from "$lib/components/explorer/files/DropZone.svelte";

    const inProgressStates = ['waiting', 'processing'],
        pausedStates = ['paused'],
        errorStates = ['error'],
        canceledStates = ['canceled'],
        doneStates = ['done'];

    let fileHandler = null;

    onMount(() => {
        fileHandler = new Uploader($authToken);
    })

    let filesRemain = 0,
        fileList = [];

    const fileUpdateProgressUploading = () => {
        fileList = fileHandler.list()
    }

    export let uploadfiles = [];

    $: if (uploadfiles.length > 0) {
        filesRemain = uploadfiles.filter(file => !doneStates.includes(file.state)).length;

        setTimeout(() => {
            HSStaticMethods.autoInit();
        }, 150)

        for (let i = 0; i < uploadfiles.length; i++) {
            const file = uploadfiles[i];
            fileHandler.add(file);
        }

        fileList = fileHandler.list();

        setInterval(fileUpdateProgressUploading, 50)
    }

    const actionWithFile = (action, uuid) => {
        fileHandler[action](uuid)
        fileList = fileHandler.list();
    }
</script>

<DropZone bind:uploadfiles={uploadfiles} />

<div class="fixed bottom-3 right-3 w-[30rem]">
    <div class="flex flex-col bg-white border shadow-sm rounded-xl dark:bg-slate-800 dark:border-gray-700 pt-1" class:hidden={uploadfiles.length === 0}>
        <div class="max-h-96 overflow-y-auto">
            {#if uploadfiles.length > 0}
                {#each fileList as file}
                    <div class="p-4 md:p-5 space-y-7">
                        <div class="mb-2 flex justify-between items-center group">
                            <div class="flex items-center gap-x-3">
                                <span class="w-8 h-8 flex justify-center items-center border border-gray-200 text-gray-500 rounded-lg dark:border-neutral-700">
                                    <svg class="flex-shrink-0 w-5 h-5" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                                        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                                        <polyline points="17 8 12 3 7 8" />
                                        <line x1="12" x2="12" y1="3" y2="15" />
                                    </svg>
                                </span>
                                <div>
                                    <p class="text-sm font-medium text-gray-800 dark:text-white">{file.name}</p>
                                    <p class="text-xs text-gray-500 dark:text-gray-500">{ (file.size).humanReadableSize() }</p>
                                </div>
                            </div>
                            <div class="inline-flex items-center gap-x-2">
                                {#if inProgressStates.includes(file.state) }
                                    <div class="animate-spin inline-block w-4 h-4 border border-current border-t-transparent text-blue-600 rounded-full dark:text-blue-500" role="status" aria-label="loading">
                                        <span class="sr-only">Loading...</span>
                                    </div>
                                {:else if pausedStates.includes(file.state) }
                                    <a class="text-gray-500 hover:text-gray-800" href="#">
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4 text-gray-500 group-hover:text-gray-800">
                                            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 5.25v13.5m-7.5-13.5v13.5" />
                                        </svg>
                                    </a>
                                {:else if errorStates.includes(file.state)}
                                    <span class="inline-flex items-center gap-x-1.5 py-1.5 px-3 rounded-full text-xs font-medium bg-red-100 text-red-800 dark:bg-red-800/30 dark:text-red-500">
                                        Error
                                    </span>
                                {/if}
                                {#if !doneStates.includes(file.state)}
                                    {#if canceledStates.includes(file.state)}
                                        <span class="inline-flex items-center gap-x-1.5 py-1.5 px-3 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-800/30 dark:text-yellow-500">Cancelled</span>
                                    {/if}
                                    <div class="hs-dropdown relative inline-flex" class:hidden={canceledStates.includes(file.state)}>
                                        <button id="hs-dropdown-custom-icon-trigger" type="button" class="hs-dropdown-toggle flex justify-center items-center w-9 h-9 text-sm font-semibold rounded-lg border border-gray-200 bg-white text-gray-800 shadow-sm hover:bg-gray-50 disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:border-gray-700 dark:text-white dark:hover:bg-gray-800 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
                                            <svg class="flex-none w-4 h-4 text-gray-600" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="1"/><circle cx="12" cy="5" r="1"/><circle cx="12" cy="19" r="1"/></svg>
                                        </button>

                                        <div class="hs-dropdown-menu z-50 transition-[opacity,margin] duration hs-dropdown-open:opacity-100 opacity-0 hidden min-w-[15rem] bg-white shadow-md rounded-lg p-2 mt-2 dark:bg-gray-800 dark:border dark:border-gray-700" aria-labelledby="hs-dropdown-custom-icon-trigger">
                                            {#if inProgressStates.includes(file.state) }
                                                <a class="flex items-center gap-x-3.5 py-2 px-3 rounded-lg text-sm text-gray-800 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300 dark:focus:bg-gray-700" on:click|preventDefault={() => actionWithFile('pause', file.uuid)} href="#pause">
                                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4 text-gray-500 group-hover:text-gray-800">
                                                        <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 5.25v13.5m-7.5-13.5v13.5" />
                                                    </svg>
                                                    Pause uploading
                                                </a>
                                            {:else if pausedStates.includes(file.state) }
                                                <a class="flex items-center gap-x-3.5 py-2 px-3 rounded-lg text-sm text-gray-800 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300 dark:focus:bg-gray-700" on:click|preventDefault={() => actionWithFile('resume', file.uuid)} href="#resume">
                                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4 text-gray-500 group-hover:text-gray-800">
                                                        <path stroke-linecap="round" stroke-linejoin="round" d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.347a1.125 1.125 0 0 1 0 1.972l-11.54 6.347a1.125 1.125 0 0 1-1.667-.986V5.653Z" />
                                                    </svg>
                                                    Resume uploading
                                                </a>
                                            {/if}
                                            {#if errorStates.includes(file.state)}
                                                <a class="flex items-center gap-x-3.5 py-2 px-3 rounded-lg text-sm text-gray-800 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300 dark:focus:bg-gray-700" href="#">
                                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4 text-gray-500 group-hover:text-gray-800">
                                                        <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
                                                    </svg>
                                                    Retry uploading
                                                </a>
                                            {/if}
                                            <a class="flex items-center gap-x-3.5 py-2 px-3 rounded-lg text-sm text-gray-800 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300 dark:focus:bg-gray-700" on:click|preventDefault={() => actionWithFile('cancel', file.uuid)} href="#cancel">
                                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4 text-gray-500 group-hover:text-gray-800">
                                                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
                                                </svg>
                                                Cancel uploading
                                            </a>
                                        </div>
                                    </div>
                                {/if}
                            </div>
                        </div>
                        {#if !canceledStates.includes(file.state)}
                            <div class="flex w-full h-2 bg-gray-200 rounded-full overflow-hidden dark:bg-gray-700" role="progressbar" aria-valuenow="100" aria-valuemin="0" aria-valuemax="100">
                                <div class="flex flex-col justify-center rounded-full overflow-hidden bg-teal-500 text-xs text-white text-center whitespace-nowrap transition duration-500"
                                     class:bg-teal-500={file.state === 'processing'}
                                     class:bg-gray-500={file.state === 'waiting'}
                                     class:bg-red-500={file.state === 'error'}
                                     class:bg-yellow-500={file.state === 'paused'}
                                     class:bg-green-500={file.state === 'done'}
                                     style="width: {file.progress}%"></div>
                            </div>
                        {/if}
                    </div>
                {/each}
            {/if}
        </div>

        <div class="bg-gray-50 border-t border-gray-200 rounded-b-xl py-2 px-4 md:px-5 dark:bg-white/[.05] dark:border-gray-700">
            <div class="flex flex-wrap justify-between items-center gap-x-3">
                <div>
                    <span class="text-sm font-semibold text-gray-800 dark:text-white"> { filesRemain } files left </span>
                </div>
                <div class="-me-2.5">
                    <button type="button" class="py-2 px-3 inline-flex items-center gap-x-1.5 text-sm font-medium rounded-lg border border-transparent text-gray-500 hover:bg-gray-200 hover:text-gray-800 disabled:opacity-50 disabled:pointer-events-none dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-white dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
                        <svg class="flex-shrink-0 w-4 h-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <rect width="4" height="16" x="6" y="4" />
                            <rect width="4" height="16" x="14" y="4" />
                        </svg>
                        Pause
                    </button>
                </div>
            </div>
        </div>
    </div>
</div>
