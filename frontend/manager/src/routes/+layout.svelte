<script>
    import { base } from "$app/paths";
</script>

<svelte:head>
    <meta name="google" content="notranslate" />
    <title>Manager - 9m</title>
</svelte:head>

<a href="/" rel="external">Home</a>
<a href={base}>Manager</a>
|
<a href="{base}/user-list">User list</a>

<hr />

<slot />

<hr />

<a href="{base}/user-list">User list</a>
