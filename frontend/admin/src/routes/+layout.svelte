<script>
    import { base } from "$app/paths";
</script>

<svelte:head>
    <meta name="google" content="notranslate" />
    <title>Admin - 9m</title>
</svelte:head>

<a href="/" rel="external">Home</a>
<a href={base}>Admin</a>
|
<a href="{base}/board-list">Boards</a>
<a href="{base}/user-columns">User columns</a>
<a href="{base}/user-list">User list</a>

<hr />

<slot />

<hr />

<a href="{base}/board-list">Boards</a>
<a href="{base}/user-columns">User columns</a>
<a href="{base}/user-list">User list</a>
