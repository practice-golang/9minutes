<script>
    import { onMount, onDestroy } from "svelte";
    import { invalidateAll } from "$app/navigation";
    import { page } from "$app/stores";

    import "moment/dist/locale/ko";
    import moment from "moment";

    export let data;

    const columns = data.columns;

    let listCount = Number($page.url.searchParams.get("list-count")) || 10;
    $: users = data["userlist-data"]["user-list"];
    $: currentPage = data["userlist-data"]["current-page"];
    $: totalPage = data["userlist-data"]["total-page"];
    $: jumpPrev = currentPage - 5 > 1 ? currentPage - 5 : 1;
    $: jumpNext = currentPage + 5 < totalPage ? currentPage + 5 : totalPage;

    let searchKeyword = $page.url.searchParams.get("search") || "";

    let selectedIndices = [];

    let editINDEX = -1;
    let showNewUser = false;
    let newUser = {};
    let editUser = {};

    const userGRADES = {
        admin: "Admin",
        manager: "Manager",
        regular_user: "Regular User",
        pending_user: "Pending User",
        banned_user: "Banned User",
        resigned_user: "Resigned User",
    };

    function search() {
        if (searchKeyword != "") {
            let params = "?search=" + searchKeyword;
            if (listCount > 10) {
                params += "&list-count=" + listCount;
            }
            location.href = params;
        }
    }

    function pressEnter(e) {
        if (e.key == "Enter") {
            console.log(searchKeyword);
            search();
        }
    }

    function paramsChange() {
        location.href =
            `?list-count=${listCount}` +
            (searchKeyword != "" ? `&search-keyword=${searchKeyword}` : "");
    }

    function closeNewUser() {
        newUser = {};
        showNewUser = false;
    }

    async function saveNewUser() {
        const uri = "/api/admin/user";
        const r = await fetch(uri, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(newUser),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        closeNewUser();
        invalidateAll();
    }

    function openEditUser(index) {
        editINDEX = index;
        editUser = {};
        for (const k in users[index]) {
            editUser[k] = users[index][k];
        }
    }

    function closeEditUser() {
        editUser = {};
        editINDEX = -1;
    }

    async function updateEditUser() {
        const columnTypes = {};
        for (const c of columns) {
            columnTypes[c["column-code"]] = c["column-type"];
        }

        for (const k in editUser) {
            switch (true) {
                case columnTypes[k] === "number-integer":
                    editUser[k] = parseInt(editUser[k]);
                    if (isNaN(editUser[k])) {
                        delete editUser[k];
                    }
                    break;
                case columnTypes[k] === "number-real":
                    editUser[k] = parseFloat(editUser[k]);
                    if (isNaN(editUser[k])) {
                        delete editUser[k];
                    }
                    break;
            }
        }

        const uri = "/api/admin/user";
        const r = await fetch(uri, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify([editUser]),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        closeEditUser();
        invalidateAll();
    }

    async function deleteUser(index) {
        const userIDX = users[index]["idx"];

        const uri = "/api/admin/user";
        const r = await fetch(uri, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify([{ idx: parseInt(userIDX) }]),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        selectedIndices = [];
        invalidateAll();
    }

    async function deleteSelectedUsers() {
        if (selectedIndices.length == 0) {
            alert("No columns selected");
            return;
        }

        const userIndices = [];
        for (let i = 0; i < selectedIndices.length; i++) {
            userIndices.push({ idx: parseInt(selectedIndices[i]) });
        }

        const uri = "/api/admin/user";
        const r = await fetch(uri, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(userIndices),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        selectedIndices = [];
        invalidateAll();
    }

    onMount(() => {
        moment.locale("ko");
    });
</script>

<h1>Admin / User list</h1>

<h1>Users list</h1>

<button
    type="button"
    on:click={() => {
        newUser["grade"] = "pending_user";
        showNewUser = true;
    }}
>
    Add user
</button>

<span>|</span>

<label for="search">Search:</label>
<input
    type="text"
    id="search"
    on:keyup={pressEnter}
    bind:value={searchKeyword}
    placeholder="Search for..."
/>
<button type="button" on:click={search}>Search</button>

<span>|</span>

<label for="set-list-count">List:</label>
<select id="set-list-count" bind:value={listCount} on:change={paramsChange}>
    {#each [5, 10, 20, 30, 50, 80] as listCountNum}
        <option value={listCountNum}>{listCountNum}</option>
    {/each}
</select>

<table id="users-list-container">
    <!-- Column titles -->
    <thead>
        <tr>
            <td>
                <input type="checkbox" />
            </td>
            {#each columns as col}
                <th>{col["display-name"]}</th>
            {/each}
            <th>Control</th>
        </tr>
    </thead>

    {#if showNewUser}
        <!-- Add user -->
        <tr id="add-user">
            <td />
            <td />
            {#each columns as col}
                {#if col["column-code"] == "idx"}
                    {""}
                {:else if col["column-code"] == "password"}
                    <td>
                        <input
                            type="password"
                            bind:value={newUser["password"]}
                            placeholder={col["display-name"]}
                        />
                    </td>
                {:else if col["column-code"] == "grade"}
                    <td>
                        <select bind:value={newUser["grade"]}>
                            {#each Object.entries(userGRADES) as [key, name]}
                                <option value={key}>{name}</option>
                            {/each}
                        </select>
                    </td>
                {:else if col["column-code"] == "approval"}
                    <td>
                        <select bind:value={newUser["approval"]}>
                            <option value="y">Y</option>
                            <option value="n" selected>N</option>
                        </select>
                    </td>
                {:else if col["column-code"] == "regdate"}
                    <td />
                {:else}
                    <td>
                        <input
                            type="text"
                            bind:value={newUser[col["column-code"]]}
                            placeholder={col["display-name"]}
                        />
                    </td>
                {/if}
            {/each}
            <td>
                <button type="button" on:click={closeNewUser}>Cancel</button>
                <button type="button" on:click={saveNewUser}>Save</button>
            </td>
        </tr>
    {/if}

    <tbody id="users-list-body">
        {#each users as user, index}
            {#if editINDEX == index}
                <!-- Edit user -->
                <tr>
                    <td />
                    {#each columns as col}
                        {#if col["column-code"] == "idx"}
                            <td>{editUser["idx"]}</td>
                        {:else if col["column-code"] == "grade" || col["column-code"] == "approval"}
                            <td>
                                <input
                                    type="text"
                                    list="{col['column-code']}-list"
                                    bind:value={editUser[col["column-code"]]}
                                    placeholder={col["display-name"]}
                                />
                            </td>
                        {:else if col["column-code"] == "regdate"}
                            <td>
                                {moment(
                                    editUser["regdate"],
                                    "YYYYMMDDhhmmss"
                                ).format("YYYY-MM-DD")}
                            </td>
                        {:else}
                            <td>
                                <input
                                    type="text"
                                    bind:value={editUser[col["column-code"]]}
                                    placeholder={col["display-name"]}
                                />
                            </td>
                        {/if}
                    {/each}
                    <td>
                        <button type="button" on:click={closeEditUser}>
                            Cancel
                        </button>
                        <button type="button" on:click={updateEditUser}>
                            Save
                        </button>
                    </td>
                </tr>
            {:else}
                <!-- Show user -->
                <tr>
                    <td>
                        <input type="checkbox" />
                    </td>
                    {#each columns as col}
                        {#if col["column-code"] == "regdate"}
                            <td>
                                {moment(
                                    user["regdate"],
                                    "YYYYMMDDhhmmss"
                                ).format("YYYY-MM-DD")}
                            </td>
                        {:else}
                            <td>{user[col["column-code"]]}</td>
                        {/if}
                    {/each}
                    <td>
                        <button
                            type="button"
                            on:click={() => {
                                openEditUser(index);
                            }}
                        >
                            Edit
                        </button>
                        <button
                            type="button"
                            on:click={() => {
                                deleteUser(index);
                            }}
                        >
                            Delete
                        </button>
                    </td>
                </tr>
            {/if}
        {/each}
    </tbody>
</table>

<datalist id="grade-list">
    <option value="admin">Admin</option>
    <option value="manager">Manager</option>
    <option value="regular_user">Regular user</option>
    <option value="pending_user">Pending user</option>
    <option value="banned_user">Banned user</option>
    <option value="resigned_user">Resigned user</option>
</datalist>

<datalist id="approval-list">
    <option value="Y">Y</option>
    <option value="N">N</option>
</datalist>

<div id="pages-container">
    <div lr-loop="pages">
        <a href="?page=1&list-count={listCount}">
            <span>&laquo;</span>
        </a>
        <a href="?page={jumpPrev}&list-count={listCount}">
            <span>&lt;</span>
        </a>

        <span>..</span>

        {#each [currentPage - 2, currentPage - 1] as page}
            {#if page >= 1}
                <a href="?page={page}&list-count={listCount}">{page}</a>
            {/if}
        {/each}

        <b lr-if="page == usersData['current-page']">
            <a href={data.link} rel="external">
                {currentPage}
            </a>
        </b>

        {#each [currentPage + 1, currentPage + 2] as page}
            {#if page <= totalPage}
                <a href="?page={page}&list-count={listCount}">{page}</a>
            {/if}
        {/each}

        <span>..</span>

        <a href="?page={jumpNext}&list-count={listCount}">
            <span>&gt;</span>
        </a>
        <a href="?page={totalPage}&list-count={listCount}">
            <span>&raquo;</span>
        </a>
    </div>
</div>
