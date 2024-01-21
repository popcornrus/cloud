<script>
    import {onMount} from "svelte";

    export let uploadfiles = [];

    let dropzoneVisible = false;

    function allowDrag(e) {
        // eslint-disable-next-line no-constant-condition
        if (true) {
            e.dataTransfer.dropEffect = 'copy';
            e.preventDefault();
        }
    }

    function handleDrop(e) {
        e.preventDefault();
        dropzoneVisible = false;

        let result = [];

        if (e.dataTransfer.files.length) {
            let upload = e.dataTransfer.files;

            for (let i = 0; i < upload.length; i++) {
                result.push({
                    uuid: Math.random().toString(36).slice(2, 9),
                    name: upload[i].name,
                    size: upload[i].size,
                    progress: 0,
                    state: 'waiting',
                    binary: upload[i],
                });
            }
        }

        uploadfiles = result;
    }

    onMount(() => {
        var dropZone = document.getElementById('dropZone');

        window.addEventListener('dragenter', function (e) {
            dropzoneVisible = true;
        });

        dropZone.addEventListener('dragenter', allowDrag);
        dropZone.addEventListener('dragover', allowDrag);

        dropZone.addEventListener('dragleave', function (e) {
            dropzoneVisible = false;
        });

        dropZone.addEventListener('drop', handleDrop);
    })
</script>

<div id="dropZone" class="fixed top-0 left-0 w-full h-full z-[90] bg-gray-800/40 backdrop-blur flex justify-center items-center text-black/60 transition duration-300" class:left-full={!dropzoneVisible} class:opacity-0={!dropzoneVisible}>
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="3.5" stroke="currentColor" class="w-24 h-24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5m-13.5-9L12 3m0 0 4.5 4.5M12 3v13.5" />
    </svg>
</div>