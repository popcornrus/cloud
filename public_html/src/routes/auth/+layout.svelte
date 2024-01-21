<script>
    import {onMount} from "svelte";
    import {goto} from "$app/navigation";

    onMount(() => {
        const localToken = localStorage.getItem('token')

        if (localToken !== null) {
            goto('/explorer')
        }
    });
</script>

<div class="flex items-center justify-center min-h-[100vh]">
    <slot></slot>
</div>
