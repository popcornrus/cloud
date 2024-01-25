<script>
    export let filesToDownload = [];
</script>

<div class="fixed bottom-3 w-[30rem] transition-[right] duration-300" class:right-3={filesToDownload.length > 0} class:-right-full={filesToDownload.length === 0}>
    <div class="flex flex-col bg-white border shadow-sm rounded-xl dark:bg-slate-800 dark:border-gray-700 pt-1">
        <div class="max-h-96 overflow-y-auto">
            {#if filesToDownload.length > 0}
                {#each filesToDownload as file}
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
                            {#if file.progress.progress < 1}
                                <div class="inline-flex items-center gap-x-2">
                                    <div class="animate-spin inline-block w-4 h-4 border border-current border-t-transparent text-blue-600 rounded-full dark:text-blue-500" role="status" aria-label="loading">
                                        <span class="sr-only">Loading...</span>
                                    </div>
                                </div>
                            {/if}
                        </div>
                        <div class="flex w-full h-2 bg-gray-200 rounded-full overflow-hidden dark:bg-gray-700" role="progressbar" aria-valuenow="100" aria-valuemin="0" aria-valuemax="100">
                            <div class="flex flex-col justify-center rounded-full overflow-hidden bg-teal-500 text-xs text-white text-center whitespace-nowrap transition duration-500"
                                 class:bg-teal-500={file.progress < 1}
                                 class:bg-green-500={file.progress === 1}
                                 style="width: {file.progress * 100}%"></div>
                        </div>
                    </div>
                {/each}
            {/if}
            <!-- End Progress Bar -->

            <div class="bg-gray-50 border-t border-gray-200 rounded-b-xl py-2 px-4 md:px-5 dark:bg-white/[.05] dark:border-gray-700">
                <div class="flex flex-wrap justify-between items-center gap-x-3">
                    <div>
                        <span class="text-sm font-semibold text-gray-800 dark:text-white"> { filesToDownload.length } files left </span>
                    </div>
                    <div class="-me-2.5">
                        <button type="button"
                                class="py-2 px-3 inline-flex items-center gap-x-1.5 text-sm font-medium rounded-lg border border-transparent text-gray-500 hover:bg-gray-200 hover:text-gray-800 disabled:opacity-50 disabled:pointer-events-none dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-white dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
                            <svg class="flex-shrink-0 w-4 h-4" xmlns="http://www.w3.org/2000/svg" width="24"
                                 height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                                 stroke-linecap="round" stroke-linejoin="round">
                                <rect width="4" height="16" x="6" y="4"/>
                                <rect width="4" height="16" x="14" y="4"/>
                            </svg>
                            Pause
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>