<script>
    import {onMount} from "svelte";

    export let fileUUID, fileInstance, filesToDownload = [];

    onMount(() => {
        const contextMenu = document.getElementById('file-context-menu');

        window.addEventListener('scroll', () => {
            contextMenu.style.display = 'none';
        });

        window.addEventListener('resize', () => {
            contextMenu.style.display = 'none';
        });

        window.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                contextMenu.style.display = 'none';
            }
        });
    });

    const deleteModalParams = {
        show: false,
        file: null,
        confirmed: false
    }

    const downloadInProgress = (file, { progress}) => {
        file.progress = progress;

        if (filesToDownload.find(f => f.uuid === file.uuid)) {
            if (progress === 1) {
                setTimeout(() => {
                    filesToDownload = filesToDownload.filter(f => f.uuid !== file.uuid);
                }, 3000)
            }

            filesToDownload = filesToDownload.map(f => {
                if (f.uuid === file.uuid) {
                    f.progress = progress;
                }
                return f;
            });
        } else {
            filesToDownload = [...filesToDownload, file]
        }
    }

    const deleteModal = (uuid) => {
        deleteModalParams.show = true;
        deleteModalParams.file =  fileInstance.findByUuid(uuid);

        const confirmDelete = document.getElementById('confirm-delete');
        confirmDelete.addEventListener('click', () => {
            fileInstance.delete(uuid);

            const el = document.getElementById(`file-${uuid}`);
            el.classList.add('animate-pulse');
            setTimeout(() => {
                el.remove();
            }, 1000);
        });
    }
</script>

<div class="absolute t-full bg-white shadow border rounded" id="file-context-menu" style="display:none">
    <div class="flex align-middle gap-x-0.5 p-2">
        <button title="Download" on:click|preventDefault={() => fileInstance.download(fileUUID, downloadInProgress)} class="ql-bold w-8 h-8 inline-flex justify-center items-center gap-x-2 text-sm font-semibold rounded-full border border-transparent text-gray-800 hover:bg-gray-100 disabled:opacity-50 disabled:pointer-events-none dark:text-white dark:hover:bg-gray-700 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
            <svg title="Download" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3" />
            </svg>
        </button>

        <button title="Rename" on:click|preventDefault={() => fileInstance.rename(fileUUID)} class="ql-italic w-8 h-8 inline-flex justify-center items-center gap-x-2 text-sm font-semibold rounded-full border border-transparent text-gray-800 hover:bg-gray-100 disabled:opacity-50 disabled:pointer-events-none dark:text-white dark:hover:bg-gray-700 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600" type="button">
            <svg title="Rename" class="w-6 h-6" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M11.75 2C11.3358 2 11 2.33579 11 2.75C11 3.16421 11.3358 3.5 11.75 3.5H13.25V24.5H11.75C11.3358 24.5 11 24.8358 11 25.25C11 25.6642 11.3358 26 11.75 26H16.25C16.6642 26 17 25.6642 17 25.25C17 24.8358 16.6642 24.5 16.25 24.5H14.75V3.5H16.25C16.6642 3.5 17 3.16421 17 2.75C17 2.33579 16.6642 2 16.25 2H11.75Z" fill="#212121"/>
                <path d="M6.25 6.01958H12.25V7.51958H6.25C5.2835 7.51958 4.5 8.30308 4.5 9.26958V18.7696C4.5 19.7361 5.2835 20.5196 6.25 20.5196H12.25V22.0196H6.25C4.45507 22.0196 3 20.5645 3 18.7696V9.26958C3 7.47465 4.45507 6.01958 6.25 6.01958Z" fill="#212121"/>
                <path d="M21.75 20.5196H15.75V22.0196H21.75C23.5449 22.0196 25 20.5645 25 18.7696V9.26958C25 7.47465 23.5449 6.01958 21.75 6.01958H15.75V7.51958H21.75C22.7165 7.51958 23.5 8.30308 23.5 9.26958V18.7696C23.5 19.7361 22.7165 20.5196 21.75 20.5196Z" fill="#212121"/>
            </svg>
        </button>

        <button title="Move" on:click|preventDefault={() => fileActions.move(fileUUID)} class="ql-italic w-8 h-8 inline-flex justify-center items-center gap-x-2 text-sm font-semibold rounded-full border border-transparent text-gray-800 hover:bg-gray-100 disabled:opacity-50 disabled:pointer-events-none dark:text-white dark:hover:bg-gray-700 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600" type="button">
            <svg title="Move" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 9V5.25A2.25 2.25 0 0 1 10.5 3h6a2.25 2.25 0 0 1 2.25 2.25v13.5A2.25 2.25 0 0 1 16.5 21h-6a2.25 2.25 0 0 1-2.25-2.25V15M12 9l3 3m0 0-3 3m3-3H2.25" />
            </svg>
        </button>

        <button title="Delete" on:click|preventDefault={() => deleteModal(fileUUID)} class="ql-italic w-8 h-8 inline-flex justify-center items-center gap-x-2 text-sm font-semibold rounded-full border border-transparent text-gray-800 hover:bg-gray-100 disabled:opacity-50 disabled:pointer-events-none dark:text-white dark:hover:bg-gray-700 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600" type="button">
            <svg title="Delete"  xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
            </svg>
        </button>
    </div>
</div>

<div class:hidden={!deleteModalParams.show} class="w-full h-full fixed top-0 start-0 z-[80] overflow-x-hidden overflow-y-auto backdrop-blur bg-gray-600/60">
    <div class="mt-7 opacity-100 duration-500 mt-0 ease-out transition-all sm:max-w-[35%] sm:w-full m-3 sm:mx-auto min-h-[calc(100%-3.5rem)] flex items-center">
        <div class="w-full flex flex-col bg-white border shadow-sm rounded-xl dark:bg-gray-800 dark:border-gray-700 dark:shadow-slate-700/[.7]">
            <div class="flex justify-between items-center py-3 px-4 border-b dark:border-gray-700">
                <h3 class="font-bold text-gray-800 dark:text-white">
                    Confirm File Deletion
                </h3>
                <button type="button" on:click={() => deleteModalParams.show = false} class="flex justify-center items-center w-7 h-7 text-sm font-semibold rounded-full border border-transparent text-gray-800 hover:bg-gray-100 disabled:opacity-50 disabled:pointer-events-none dark:text-white dark:hover:bg-gray-700 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600" data-hs-overlay="#hs-vertically-centered-modal">
                    <span class="sr-only">Close</span>
                    <svg class="flex-shrink-0 w-4 h-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
                </button>
            </div>
            <div class="p-4 overflow-y-auto">
                <div class=" mb-4 text-center font-bold bg-yellow-100 border border-yellow-200 text-sm text-yellow-800 rounded-lg p-4 dark:bg-yellow-800/10 dark:border-yellow-900 dark:text-yellow-500" role="alert">
                    This action is irreversible.<br />All data associated with this file will be permanently deleted and cannot be recovered.
                </div>

                <p class="text-center">You are about to delete the following file</p>
                <p class="text-center"><strong>{deleteModalParams.file?.name}</strong></p>
            </div>
            <div class="flex justify-end items-center gap-x-2 py-3 px-4 border-t dark:border-gray-700">
                <button type="button"  on:click={() => deleteModalParams.show = false} class="py-3 px-4 inline-flex items-center gap-x-2 text-sm font-semibold rounded-lg border border-transparent bg-gray-100 text-gray-800 hover:bg-gray-200 disabled:opacity-50 disabled:pointer-events-none dark:bg-white/10 dark:hover:bg-white/20 dark:text-white dark:hover:text-white dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
                    No, Keep It.
                </button>
                <button type="button" id="confirm-delete" on:click={() => deleteModalParams.show = false} class="py-3 px-4 inline-flex items-center gap-x-2 text-sm font-semibold rounded-lg border border-transparent bg-red-100 text-red-800 hover:bg-red-200 disabled:opacity-50 disabled:pointer-events-none dark:hover:bg-red-900 dark:text-red-500 dark:hover:text-red-400 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
                    Yes, Delete!
                </button>
            </div>
        </div>
    </div>
</div>