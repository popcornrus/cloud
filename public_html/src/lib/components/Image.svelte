<script>
    export let src = "";

    let loaded = false;
    let failed = false;
    let loading = false;

    $: if (src !== null) {
        const img = new Image();
        img.src = src;

        loading = true;

        img.onload = () => {
            loading = false;
            loaded = true;
        };
        img.onerror = () => {
            loading = false;
            failed = true;

            console.log("Failed to load image: " + src)
        };
    }
</script>

{#if loaded}
    <img class="{$$restProps?.size} {$$restProps?.imgClass} object-cover" alt="{$$restProps?.alt}" src="{src}"/>
{:else if failed}
    <svg fill="#444" viewBox="0 0 32 32" id="icon" xmlns="http://www.w3.org/2000/svg" class="w-24 h-24">
        <path d="M30,3.4141,28.5859,2,2,28.5859,3.4141,30l2-2H26a2.0027,2.0027,0,0,0,2-2V5.4141ZM26,26H7.4141l7.7929-7.793,2.3788,2.3787a2,2,0,0,0,2.8284,0L22,19l4,3.9973Zm0-5.8318-2.5858-2.5859a2,2,0,0,0-2.8284,0L19,19.1682l-2.377-2.3771L26,7.4141Z"/>
        <path d="M6,22V19l5-4.9966,1.3733,1.3733,1.4159-1.416-1.375-1.375a2,2,0,0,0-2.8284,0L6,16.1716V6H22V4H6A2.002,2.002,0,0,0,4,6V22Z"/>
        <rect fill="none" width="32" height="32"/>
    </svg>
{:else if loading}
    <div class="flex animate-pulse">
        <div class="flex-shrink-0">
            <span class="{$$restProps?.size}  block bg-gray-200 rounded-full dark:bg-gray-700"></span>
        </div>
    </div>
{/if}