<script>
    import {env} from "$env/dynamic/public";
    import Image from "$lib/components/Image.svelte";
    import Plyr from 'plyr';
    import "plyr/dist/plyr.css";
    import axios from "axios";
    import {Toast} from "$lib";

    export let file = null,
        modalOpen = false;

    let progressBar = false,
        progress = 0,
        bytesProgress = 0;

    const download = async () => {
        const streamUrl = `${env.PUBLIC_BACKEND_URL}/api/v1/explorer/share/${file?.uuid}/download`;

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

                axios.get(streamUrl, {
                    responseType: "blob",
                    onDownloadProgress: (p) => {
                        progress = Math.round((p.loaded / p.total) * 100);
                        progressBar = true;

                        bytesProgress = p.loaded;
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
            onDownloadProgress: (p) => {
                progress = Math.round((p.loaded / p.total) * 100);
                progressBar = true;

                bytesProgress = p.loaded;
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

    $: if (file?.uuid?.length > 0) {
        modalOpen = true

        if (file?.type.IsVideo()) {
            setTimeout(() => {
                new Plyr('#player');
            }, 500)
        }
    }

    $: if (progress === 100) {
        progressBar = false;

        Toast({
            message: 'Download completed',
            type: 'success',
            duration: 3000
        })
    }
</script>

<div class="w-full h-full fixed top-0 start-0 overflow-x-hidden overflow-y-auto backdrop-blur">
    <div class="mt-7 opacity-100 duration-500 mt-0 ease-out transition-all sm:max-w-[45%] sm:w-full m-3 sm:mx-auto min-h-[calc(100%-3.5rem)] flex items-center">
        <div class="w-full flex flex-col bg-white border shadow-sm rounded-xl dark:bg-gray-800 dark:border-gray-700 dark:shadow-slate-700/[.7]">
            <div class="flex justify-center items-center py-3 px-4 border-b dark:border-gray-700">
                <h3 class="font-bold text-gray-800 dark:text-white">
                    {file?.name} | {file?.size.humanReadableSize()}
                </h3>
            </div>
            <div class="p-4 overflow-y-auto">
                {#if file?.type.IsImage()}
                    <Image src="{env.PUBLIC_BACKEND_URL}/api/v1/explorer/files/{file?.uuid}" alt="{file?.name}"/>
                {:else if file?.type.IsVideo()}
                    <video src="{env.PUBLIC_BACKEND_URL}/api/v1/explorer/files/{file?.uuid}" id="player">
                        <track kind="captions" src="{env.PUBLIC_BACKEND_URL}/api/v1/explorer/files/{file?.uuid}">
                    </video>
                {:else}
                    <div class="flex justify-center items-center w-full h-full p-24">
                        <div class="flex flex-col justify-center items-center gap-y-2 text-gray-500">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                                 stroke="currentColor" class="w-12 h-12">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                      d="M18.364 18.364A9 9 0 0 0 5.636 5.636m12.728 12.728A9 9 0 0 1 5.636 5.636m12.728 12.728L5.636 5.636"/>
                            </svg>
                            <span class="">No preview available</span>
                        </div>
                    </div>
                {/if}
            </div>

            <div class="flex justify-center items-center gap-x-2 py-3 px-4 border-t dark:border-gray-700">
                {#if !progressBar}
                    <button on:click={download} type="button"
                            class="py-3 px-4 inline-flex items-center gap-x-2 text-sm font-semibold rounded-lg border border-transparent bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:pointer-events-none dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
                        Download
                    </button>
                {:else}
                    <!-- File Uploading Progress Form -->
                    <div class="w-full">
                        <!-- Uploading File Content -->
                        <div class="mb-2 flex justify-between items-center">
                            <div class="flex items-center gap-x-3">
                                <div>
                                    <p class="text-sm font-medium text-gray-800 dark:text-white">Downloading...</p>
                                    <p class="text-xs text-gray-500 dark:text-gray-500">{bytesProgress.humanReadableSize()} / {file?.size.humanReadableSize()}</p>
                                </div>
                            </div>
                            {#if progress === 100}
                                <div class="inline-flex items-center gap-x-2">
                                    <svg class="flex-shrink-0 w-4 h-4 text-teal-500" xmlns="http://www.w3.org/2000/svg"
                                         width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                                        <path d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0zm-3.97-3.03a.75.75 0 0 0-1.08.022L7.477 9.417 5.384 7.323a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-.01-1.05z"/>
                                    </svg>
                                </div>
                            {:else}
                                <div class="animate-spin inline-block w-6 h-6 border-[3px] border-current border-t-transparent text-blue-600 rounded-full dark:text-blue-500" role="status" aria-label="loading">
                                    <span class="sr-only">Loading...</span>
                                </div>
                            {/if}
                        </div>
                        <!-- End Uploading File Content -->

                        <!-- Progress Bar -->
                        <div class="flex items-center gap-x-3 whitespace-nowrap">
                            <div class="flex w-full h-2 bg-gray-200 rounded-full overflow-hidden dark:bg-gray-700"
                                 role="progressbar" aria-valuenow="{progress}" aria-valuemin="0" aria-valuemax="100">
                                <div
                                    class:bg-blue-500={progress < 100}
                                    class:bg-green-500={progress === 100}
                                    class="flex flex-col justify-center rounded-full overflow-hidden text-xs text-white text-center whitespace-nowrap transition duration-500"
                                     style="width: {progress}%"></div>
                            </div>
                            <div class="w-6 text-end">
                                <span class="text-sm text-gray-800 dark:text-white">{progress}%</span>
                            </div>
                        </div>
                        <!-- End Progress Bar -->
                    </div>
                    <!-- End File Uploading Progress Form -->
                {/if}
            </div>
        </div>
    </div>
</div>

