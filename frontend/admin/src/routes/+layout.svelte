<script>
    import { base } from "$app/paths";
</script>

<a href={base}>Admin</a>
<a href="{base}/board-list">Boards</a>
<a href="{base}/user-columns">User columns</a>
<a href="{base}/user-list">User list</a>

<hr />
<slot />
<hr />

<a href={base}>Admin</a>
<a href="{base}/board-list">Boards</a>
<a href="{base}/user-columns">User columns</a>
<a href="{base}/user-list">User list</a>
