<script>
    import Template from "$lib/components/explorer/files/File.svelte";
    import { authToken } from "$lib/stores/token.js";
    import { user } from "$lib/stores/user.js";

    import { File } from "$lib/classes/File.js";
    import ContextMenu from "$lib/components/explorer/files/ContextMenu.svelte";
    import {onMount} from "svelte";
    import Preview from "$lib/components/explorer/files/Preview.svelte";

    export let filesToDownload = [],
        search = null;

    let wsSocket = null;

    let fileInstance = null,
        fileUUID = null,
        filePreview = null,
        files = [],
        filesLoaded = new Promise((resolve) => {
            resolve([]);
        });

    const openContextMenu = (e) => {
        const contextMenu = document.querySelector("#file-context-menu");
        const fileElement = e.target.closest(".file");

        if (fileElement) {
            let fileRect = fileElement.getBoundingClientRect()

            contextMenu.style.display = "block";
            contextMenu.style.left = ((fileRect.left + (contextMenu.clientWidth / 2)) / 16) + "rem";
            contextMenu.style.top = ((fileRect.top + fileRect.height + 5) / 16) + "rem";

            fileUUID = fileElement.dataset.uuid;
        }
    }

    const loadFiles = async (query) => {
        files = await fileInstance.list(query);
    }

    onMount(async () => {
        fileInstance = new File($authToken)
        wsSocket = fileInstance.webSocket($user)
        filesLoaded = fileInstance.list();

        files = await filesLoaded ?? [];

        document.addEventListener("click", () => {
            const contextMenu = document.querySelector("#file-context-menu");
            if (contextMenu) {
                contextMenu.style.display = "none";
            }
        });

        wsSocket.addEventListener("message", async (e) => {
            const data = JSON.parse(e.data);

            if (data.event === "file:created") {
                await fileInstance.data(data.data.uuid);
                files = fileInstance.files;
                return
            }

            if (data.event === "file:deleted") {
                files = files.filter((file) => file.uuid !== data.data.uuid);
            }
        });
    })

    var delayTimer = null;

    $: if (search !== null) {
        clearTimeout(delayTimer);
        delayTimer = setTimeout(function() {
            if (search?.length === 0) {
                loadFiles()
            }

            if (search?.length >= 3) {
                loadFiles(search)
            }
        }, 500);
    }
</script>

{#await filesLoaded}
    <p>Loading files...</p>
{:then ok}
    <div class="grid 2xl:grid-cols-5 lg:grid-cols-3 sm:grid-cols-2 xs:grid-cols-2 grid-cols-1 gap-4">
        {#if ok === null || ok.length === 0}
            <p>No files found.</p>
        {:else}
            {#each files as file}
                <Template
                    oncontextmenu={openContextMenu}
                    file="{file}"
                    process="{fileInstance}"
                />
            {/each}
        {/if}
    </div>
{:catch error}
    <p>error: {error.message}</p>
{/await}

<ContextMenu fileInstance="{fileInstance}" fileUUID="{fileUUID}" bind:filesToDownload />