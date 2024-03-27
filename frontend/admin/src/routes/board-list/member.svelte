<script>
    export let BoardMember;
    export let title;
    export let message;

    function closeDialog() {
        BoardMember.close();
    }
</script>

<dialog bind:this={BoardMember} on:close>
    <h3>{title}</h3>
    <p>{message}</p>
    <button on:click={closeDialog}>Cancel</button>
</dialog>

<style>
    dialog {
        border: 1px solid gray;
        width: 80vw;
        height: 80vh;
        position: fixed;
        padding: 20px;
    }
</style>
