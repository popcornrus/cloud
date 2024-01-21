<script>
    import Image from "$lib/components/Image.svelte";
    import {env} from "$env/dynamic/public";

    export let file;

    const preview = `${env.PUBLIC_BACKEND_URL}/api/v1/explorer/files/${file.uuid}/preview?h=128&w=300&a=resize`
</script>

<div class="flex flex-col bg-white border shadow-sm rounded-xl dark:bg-slate-900 dark:border-gray-700 dark:shadow-slate-700/[.7] hover:scale-105 transition cursor-pointer file overflow-hidden"
     on:contextmenu|preventDefault={$$restProps?.oncontextmenu}
     on:click|preventDefault={$$restProps?.onclick}
     id="file-{file.uuid}"
     data-uuid="{file.uuid}"
     role="button"
     tabindex="-1"
>
    <div class="flex items-center justify-center w-full h-32 bg-gray-100 dark:bg-slate-800 group">
        {#if file.type.IsImage()}
            <Image src="{preview}" alt="{file.name}" class="object-contain w-full h-full" />
        {:else if file.type.IsVideo()}
            <Image src="{preview}" alt="{file.name}" class="object-contain w-full h-full" />
            <div class="absolute bg-white rounded-full shadow group-hover:text-gray-500 group-hover:scale-105 transition">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1" stroke="currentColor" class="w-12 h-12">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15.91 11.672a.375.375 0 0 1 0 .656l-5.603 3.113a.375.375 0 0 1-.557-.328V8.887c0-.286.307-.466.557-.327l5.603 3.112Z" />
                </svg>
            </div>
        {/if}
    </div>

    <div class="p-4 md:p-5">
        <div class="mb-2">
            <p class="text-md font-bold text-gray-800 dark:text-white truncate file-title" title="{file.name}">{file.name}</p>
            <input type="text" style="display:none" value="{file.name}" class="text-md focus:outline-0 font-bold file-input block w-full disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:border-gray-700 dark:text-gray-400 dark:focus:ring-gray-600" placeholder="{file.name}">
        </div>
        <div class="flex justify-between">
            <p class="text-sm text-gray-500 dark:text-gray-400" title="Size">{file.size.humanReadableSize()}</p>
            <p class="text-sm text-gray-500 dark:text-gray-400" title="Mime Type">{file.type}</p>
        </div>
    </div>
</div>