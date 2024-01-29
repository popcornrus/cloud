<script>
    import {afterNavigate} from "$app/navigation";

    afterNavigate(() => {
        HSStaticMethods.autoInit();
    });
</script>

<slot></slot>