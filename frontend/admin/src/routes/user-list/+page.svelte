<script>
    import { onMount, onDestroy, beforeUpdate, afterUpdate } from "svelte";
    import { invalidateAll } from "$app/navigation";
    import { page } from "$app/stores";

    import moment from "moment";
    // import "moment/dist/locale/ko";

    import "$lib/styles/table.css"

    export let data;

    const columns = data.columns;
    const grades = data.grades;

    let listCount = Number($page.url.searchParams.get("list-count")) || 20;
    $: users = data["userlist-data"]["user-list"];
    let previousPage = -1;
    $: currentPage = data["userlist-data"]["current-page"];
    $: totalPage = data["userlist-data"]["total-page"];
    $: jumpPrev = currentPage - 5 > 1 ? currentPage - 5 : 1;
    $: jumpNext = currentPage + 5 < totalPage ? currentPage + 5 : totalPage;

    let searchKeyword = $page.url.searchParams.get("search") || "";

    const adminIDX = 1;
    $: adminCount = 0;
    let selectAll = false;
    let selectedIndices = [];

    let editINDEX = -1;
    let showNewUser = false;
    let newUser = {};
    let editUser = {};

    function toggleSelectAll() {
        if (selectAll) {
            selectedIndices = [];
        } else {
            selectedIndices = users
                .filter((user) => parseInt(user["idx"]) > adminIDX)
                .map((user) => user["idx"]);
        }
    }

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
            (searchKeyword != "" ? `&search=${searchKeyword}` : "");
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

    async function deleteUser(index, mode = "resign") {
        const userIDX = users[index]["idx"];

        let uri = "/api/admin/user";
        uri += mode == "delete" ? "?mode=delete" : "";

        const r = await fetch(uri, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify([{ idx: userIDX }]),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        selectedIndices = [];
        invalidateAll();
    }

    async function deleteSelectedUsers(mode = "resign") {
        if (selectedIndices.length == 0) {
            alert("Selected nothing");
            return;
        }

        const userIndices = [];
        for (let i = 0; i < selectedIndices.length; i++) {
            userIndices.push({ idx: selectedIndices[i] });
        }

        let uri = "/api/admin/user";
        uri += mode == "delete" ? "?mode=delete" : "";

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
        // moment.locale("ko");
    });

    afterUpdate(() => {
        if (previousPage != currentPage) {
            selectedIndices = [];
            previousPage = currentPage;
        }

        adminCount = parseInt(users[0]["idx"]) > adminIDX ? 0 : 1;
        selectAll = selectedIndices.length == users.length - adminCount;
    });
</script>

<h1>Users</h1>

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

<button type="button" on:click={deleteSelectedUsers}>
    Resign selected users
</button>

<button type="button" on:click={() => deleteSelectedUsers("delete")}>
    Delete selected users
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

<div class="table-container">
    <table id="user-list-container">
        <!-- Column titles -->
        <thead>
            <tr>
                <td>
                    <input
                        type="checkbox"
                        bind:checked={selectAll}
                        on:change={toggleSelectAll}
                    />
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
                <td /><td />
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
                                {#each Object.entries(grades) as [key, grade]}
                                    <option value={key}>{grade.name}</option>
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
                    <button type="button" on:click={closeNewUser}>
                        Cancel
                    </button>
                    <button type="button" on:click={saveNewUser}>Save</button>
                </td>
            </tr>
        {/if}

        <tbody id="users-list-body">
            {#each users as user, index}
                <tr>
                    {#if editINDEX == index}
                        <!-- Edit user -->
                        <td />
                        {#each columns as col}
                            {#if col["column-code"] == "idx"}
                                <td class="colFixedMin">{editUser["idx"]}</td>
                            {:else if col["column-code"] == "grade"}
                                <td class="colFixedMid">
                                    <select
                                        bind:value={editUser[
                                            col["column-code"]
                                        ]}
                                    >
                                        {#each Object.entries(grades) as [key, grade]}
                                            <option value={key}>{grade.name}</option>
                                        {/each}
                                    </select>
                                </td>
                            {:else if col["column-code"] == "approval"}
                                <td class="colFixedMid">
                                    <select
                                        bind:value={editUser[
                                            col["column-code"]
                                        ]}
                                    >
                                        <option value="Y">Y</option>
                                        <option value="N" selected>N</option>
                                    </select>
                                </td>
                            {:else if col["column-code"] == "regdate"}
                                <td class="colFixedMid">
                                    {moment(
                                        editUser["regdate"],
                                        "YYYYMMDDhhmmss"
                                    ).format("YYYY-MM-DD")}
                                </td>
                            {:else}
                                <td class="colField">
                                    <input
                                        type="text"
                                        bind:value={editUser[
                                            col["column-code"]
                                        ]}
                                        placeholder={col["display-name"]}
                                    />
                                </td>
                            {/if}
                        {/each}
                        <td class="colFixedMid">
                            <button type="button" on:click={closeEditUser}>
                                Cancel
                            </button>
                            <button type="button" on:click={updateEditUser}>
                                Save
                            </button>
                        </td>
                    {:else}
                        <!-- Show user -->
                        <td class="colFixedMin">
                            {#if parseInt(user["idx"]) > 1}
                                <input
                                    type="checkbox"
                                    bind:group={selectedIndices}
                                    value={user["idx"]}
                                />
                            {/if}
                        </td>
                        {#each columns as col}
                            {#if col["column-code"] == "regdate"}
                                <td class="colFixedMid">
                                    {moment(
                                        user["regdate"],
                                        "YYYYMMDDhhmmss"
                                    ).format("YYYY-MM-DD")}
                                </td>
                            {:else}
                                <td class="colField">
                                    {user[col["column-code"]]}
                                </td>
                            {/if}
                        {/each}
                        <td class="colFixedMax">
                            <button
                                type="button"
                                on:click={() => {
                                    openEditUser(index);
                                }}
                            >
                                Edit
                            </button>

                            <span>|</span>

                            <button
                                type="button"
                                on:click={() => {
                                    deleteUser(index);
                                }}
                                disabled={parseInt(user["idx"]) == 1}
                            >
                                Resign
                            </button>

                            <button
                                type="button"
                                on:click={() => {
                                    deleteUser(index, "delete");
                                }}
                                disabled={parseInt(user["idx"]) == 1}
                            >
                                Delete
                            </button>
                        </td>
                    {/if}
                </tr>
            {/each}
        </tbody>
    </table>
</div>

<div id="page-container">
    <a
        href={`?page=1&list-count=${listCount}` +
            (searchKeyword != "" ? `&search=${searchKeyword}` : "")}
    >
        <span>&laquo;</span>
    </a>
    <a
        href={`?page=${jumpPrev}&list-count=${listCount}` +
            (searchKeyword != "" ? `&search=${searchKeyword}` : "")}
    >
        <span>&lt;</span>
    </a>

    <span>..</span>

    {#each [currentPage - 2, currentPage - 1] as page}
        {#if page >= 1}
            <a
                href={`?page=${page}&list-count=${listCount}` +
                    (searchKeyword != "" ? `&search=${searchKeyword}` : "")}
            >
                {page}
            </a>
        {/if}
    {/each}

    <b>{currentPage}</b>

    {#each [currentPage + 1, currentPage + 2] as page}
        {#if page <= totalPage}
            <a
                href={`?page=${page}&list-count=${listCount}` +
                    (searchKeyword != "" ? `&search=${searchKeyword}` : "")}
            >
                {page}
            </a>
        {/if}
    {/each}

    <span>..</span>

    <a
        href={`?page=${jumpNext}&list-count=${listCount}` +
            (searchKeyword != "" ? `&search=${searchKeyword}` : "")}
    >
        <span>&gt;</span>
    </a>
    <a
        href={`?page=${totalPage}&list-count=${listCount}` +
            (searchKeyword != "" ? `&search=${searchKeyword}` : "")}
    >
        <span>&raquo;</span>
    </a>
</div>
