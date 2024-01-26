<script>
    import Image from "$lib/components/Image.svelte";
    import {env} from "$env/dynamic/public";

    import CollectingSVG from "$lib/assets/collecting.svg";

    export let file = null,
        process = null,
        checked = false;

    const share = async () => {
        const {default: Share} = await import('$lib/components/explorer/files/Share.svelte');

        new Share({
            target: document.body,
            props: {
                file: file,
                modalOpen: true,
            }
        })
    }

    const preview = async () => {
        const {default: Preview} = await import('$lib/components/explorer/files/Preview.svelte');

        new Preview({
            target: document.body,
            props: {
                file: file,
                modalOpen: true,
            }
        })
    }

    $: if (file && process !== null) {
        process.WebSocket.addEventListener('message', (e) => {
            if (e.data.indexOf(file.uuid) === -1) {
                return;
            }

            const data = JSON.parse(e.data);


            if (data.event === "file:update") {
                file.state = data.data.state;

                if (file.state === 'done') {
                    checked = true;

                    setTimeout(() => {
                        checked = false
                    }, 3000)
                }

                return
            }

            if (data.event === 'file:preview') {
                file.preview = `${env.PUBLIC_BACKEND_URL}/api/v1/explorer/files/${file.uuid}/preview?h=128&w=300&a=resize`
            }
        })
    }

    $: if (file !== null) {
        file.preview = file.state === 'pending' ? '' : `${env.PUBLIC_BACKEND_URL}/api/v1/explorer/files/${file.uuid}/preview?h=128&w=300&a=resize`
    }
</script>

<div class="flex flex-col bg-white border shadow-sm rounded-xl dark:bg-slate-900 dark:border-gray-700 dark:shadow-slate-700/[.7] hover:scale-105 transition cursor-pointer file overflow-hidden"
     on:contextmenu|preventDefault={$$restProps?.oncontextmenu}
     id="file-{file.uuid}"
     data-uuid="{file.uuid}"
     data-state="{file.state}"
     role="button"
     tabindex="-1"
>
    <div class="flex items-center justify-center w-full h-32 bg-gray-100 dark:bg-slate-800 group"
         on:click|preventDefault={preview}>
        {#if file.state === 'collecting'}
            <object class="w-12" data="{CollectingSVG}" type="image/svg+xml" />
        {:else}
            {#if file.type.IsImage()}
                <Image src="{file.preview ?? ''}" alt="{file.name}" class="object-contain w-full h-full" />
            {:else if file.type.IsVideo()}
                <Image src="{file.preview ?? ''}" alt="{file.name}" class="object-contain w-full h-full" />
                <div class="absolute bg-white rounded-full shadow group-hover:text-gray-500 group-hover:scale-105 transition">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1" stroke="currentColor" class="w-12 h-12">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                        <path stroke-linecap="round" stroke-linejoin="round" d="M15.91 11.672a.375.375 0 0 1 0 .656l-5.603 3.113a.375.375 0 0 1-.557-.328V8.887c0-.286.307-.466.557-.327l5.603 3.112Z" />
                    </svg>
                </div>
            {:else}
                <div class="absolute group-hover:text-gray-500 group-hover:scale-105 transition text-center text-black/70">
                <span class="shadow bg-white/40 block rounded mb-1 py-1 flex justify-center">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-12 h-12">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z" />
                    </svg>
                </span>
                    <span class="uppercase font-bold shadow bg-white/40 px-2 rounded py-0.5">.{file.name.Extension()}</span>
                </div>
            {/if}
        {/if}
    </div>

    <div class="p-2 md:p-3">
        <div class="mb-2">
            <p class="text-md font-bold text-gray-800 dark:text-white truncate file-title" title="{file.name}">{file.name}</p>
            <input type="text" style="display:none" value="{file.name}" class="text-md focus:outline-0 font-bold file-input block w-full disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:border-gray-700 dark:text-gray-400 dark:focus:ring-gray-600" placeholder="{file.name}">
        </div>
        <div class="grid grid-cols-2 items-center mb-2">
            <p class="text-sm text-gray-500 dark:text-gray-400" title="Size">{file.size.humanReadableSize()}</p>
            <p class="text-sm text-gray-500 dark:text-gray-400 uppercase text-right" title="Mime Type">{file.name.Extension()}</p>
        </div>

        <div class="flex items-center justify-center min-h-6">
            {#if checked}
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6 text-green-700">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
                </svg>
            {/if}
            {#if file.state === 'converting'}
                <div class="animate-spin inline-block w-6 h-6 border-2 border-current border-t-transparent text-yellow-600 rounded-full"
                     role="status"
                     aria-label="converting"
                     title="Converting in WebM"
                >
                    <span class="sr-only">Converting...</span>
                </div>
            {/if}
            {#if file.state === 'done'}
                <button role="button" on:click={share}>
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6 text-gray-500 hover:text-gray-400 transition duration-300">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M7.217 10.907a2.25 2.25 0 1 0 0 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186 9.566-5.314m-9.566 7.5 9.566 5.314m0 0a2.25 2.25 0 1 0 3.935 2.186 2.25 2.25 0 0 0-3.935-2.186Zm0-12.814a2.25 2.25 0 1 0 3.933-2.185 2.25 2.25 0 0 0-3.933 2.185Z" />
                    </svg>
                </button>
            {/if}
        </div>
    </div>
</div>