<script>
    import { onMount, onDestroy, beforeUpdate, afterUpdate } from "svelte";
    import { invalidateAll } from "$app/navigation";
    import { page } from "$app/stores";

    import moment from "moment";

    export let data;

    onMount(() => {
        console.log(data)
    });

    afterUpdate(() => {});
</script>

<h1>List</h1>
