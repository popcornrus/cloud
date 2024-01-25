<script>
    import Uploader from "$lib/components/explorer/files/Uploader.svelte";
    import FileSection from "$lib/components/explorer/files/Section.svelte";
    import Download from "$lib/components/explorer/files/Download.svelte";

    let filesToDownload = [],
        searchValue = null;
</script>

<div class="flex justify-between mb-4 items-center">
    <h2 class="text-2xl font-bold">Explorer</h2>
    <div class="relative">
        <input type="text" bind:value={searchValue} id="hs-leading-icon" name="hs-leading-icon" class="border py-3 px-4 ps-11 block w-full border-gray-200 shadow-sm rounded-lg text-sm disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:border-gray-700 dark:text-gray-400 dark:focus:ring-gray-600" placeholder="Search...">
        <div class="absolute inset-y-0 start-0 flex items-center pointer-events-none z-20 ps-4">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="flex-shrink-0 h-4 w-4 text-gray-400 dark:text-gray-600">
                <path stroke-linecap="round" stroke-linejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" />
            </svg>
        </div>
    </div>
</div>

<FileSection bind:filesToDownload search="{searchValue}" />

<Download bind:filesToDownload />
<Uploader />
