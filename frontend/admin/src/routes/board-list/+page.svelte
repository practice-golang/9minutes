<script>
    import { onMount, onDestroy, beforeUpdate, afterUpdate } from "svelte";
    import { invalidateAll } from "$app/navigation";
    import { page } from "$app/stores";

    import "$lib/styles/table.css";

    export let data;

    const columns = data.columns;
    const grades = data.grades;

    let listCount = Number($page.url.searchParams.get("list-count")) || 10;
    $: boards = data["boardlist-data"]["board-list"];

    let previousPage = -1;
    $: currentPage = data["boardlist-data"]["current-page"];
    $: totalPage = data["boardlist-data"]["total-page"];
    $: jumpPrev = currentPage - 5 > 1 ? currentPage - 5 : 1;
    $: jumpNext = currentPage + 5 < totalPage ? currentPage + 5 : totalPage;

    let searchKeyword = $page.url.searchParams.get("search") || "";

    let selectAll = false;
    let selectedIndices = [];

    let editINDEX = -1;
    let showNewBoard = false;
    let newBoard = {};
    let editBoard = {};

    const grantSelections = [
        "grant-read",
        "grant-write",
        "grant-comment",
        "grant-upload",
    ];

    $: {
        if (
            newBoard["board-code"] != undefined &&
            newBoard["board-code"].length > 0
        ) {
            newBoard["board-table"] =
                "board_" + newBoard["board-code"].toLowerCase();
            newBoard["comment-table"] =
                "comment_" + newBoard["board-code"].toLowerCase();
        }

        if (
            editBoard["board-code"] != undefined &&
            editBoard["board-code"].length > 0
        ) {
            editBoard["board-table"] =
                "board_" + editBoard["board-code"].toLowerCase();
            editBoard["comment-table"] =
                "comment_" + editBoard["board-code"].toLowerCase();
        }
    }

    function toggleSelectAll() {
        if (selectAll) {
            selectedIndices = [];
        } else {
            selectedIndices = boards.map((board) => board["idx"]);
        }
    }

    function closeNewBoard() {
        newBoard = {};
        showNewBoard = false;
    }

    async function saveNewBoard() {
        const uri = "/api/admin/board";
        const r = await fetch(uri, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(newBoard),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        closeNewBoard();
        invalidateAll();
    }

    function moveToListView(index) {
        window.open(
            "/board/list?board_code=" + boards[index]["board-code"],
            "_blank",
        );
    }

    function openEditBoard(index) {
        editINDEX = index;
        editBoard = {};
        for (const k in boards[index]) {
            editBoard[k] = boards[index][k];
        }
        console.log(editBoard);
    }

    function closeEditBoard() {
        editBoard = {};
        editINDEX = -1;
    }

    async function updateEditBoard() {
        for (const k in editBoard) {
            if (typeof editBoard[k] == undefined || editBoard[k] == null) {
                continue;
            }
            if (typeof editBoard[k] != "string") {
                editBoard[k] = editBoard[k].toString();
            }
        }

        const uri = "/api/admin/board";
        const r = await fetch(uri, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify([editBoard]),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        closeEditBoard();
        invalidateAll();
    }

    async function deleteBoard(index) {
        const boardIDX =
            typeof boards[index]["idx"] == "number"
                ? boards[index]["idx"].toString()
                : boards[index]["idx"];

        let uri = "/api/admin/board";

        const r = await fetch(uri, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify([{ idx: boardIDX }]),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        selectedIndices = [];
        invalidateAll();
    }

    async function deleteSelectedBoards() {
        if (selectedIndices.length == 0) {
            alert("Selected nothing");
            return;
        }

        const boardIndices = [];
        for (let i = 0; i < selectedIndices.length; i++) {
            const targetIDX =
                typeof selectedIndices[i] == "number"
                    ? selectedIndices[i].toString()
                    : selectedIndices[i];
            boardIndices.push({ idx: targetIDX });
        }

        let uri = "/api/admin/board";

        const r = await fetch(uri, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(boardIndices),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        selectedIndices = [];
        invalidateAll();
    }

    onMount(() => {});

    afterUpdate(() => {
        if (previousPage != currentPage) {
            selectedIndices = [];
            previousPage = currentPage;
        }

        selectAll = selectedIndices.length == boards.length;
    });
</script>

<h1>Boards</h1>

<div>
    <button
        type="button"
        on:click={() => {
            newBoard = { "board-type": "board" };
            showNewBoard = true;
        }}
    >
        Add board
    </button>

    <span>|</span>

    <button type="button" on:click={deleteSelectedBoards}>
        Delete selected boards
    </button>

    <span>|</span>

    <label for="search">Search:</label>
    <input
        type="text"
        id="search"
        onkeyup="pressEnter()"
        placeholder="Search for..."
    />
    <button type="button" onclick="search()">Search</button>
</div>

<table id="boards-list-container">
    <thead>
        <tr>
            <th>
                <input
                    type="checkbox"
                    bind:checked={selectAll}
                    on:change={toggleSelectAll}
                />
            </th>

            {#each columns as col}
                <th>{col["display-name"]}</th>
            {/each}

            <th>Control</th>
        </tr>
    </thead>

    {#if showNewBoard}
        <!-- Add board -->
        <tr id="add-board">
            <td />
            <td />

            {#each columns as col}
                {#if col["column-code"] == "idx"}
                    {""}
                {:else if col["column-code"] == "board-type"}
                    <td>
                        <select bind:value={newBoard[col["column-code"]]}>
                            <option value="board">Board</option>
                            <option value="gallery">Gallery</option>
                        </select>
                    </td>
                {:else if col["column-code"] == "board-table" || col["column-code"] == "comment-table"}
                    <td>
                        <input
                            bind:value={newBoard[col["column-code"]]}
                            disabled
                        />
                    </td>
                {:else if grantSelections.includes(col["column-code"])}
                    <td>
                        <select bind:value={newBoard[col["column-code"]]}>
                            {#each Object.entries(grades) as [key, grade]}
                                <option value={grade.code}>{grade.name}</option>
                            {/each}
                        </select>
                    </td>
                {:else if col["column-code"] == "regdate"}
                    <td />
                {:else}
                    <td>
                        <input
                            type="text"
                            bind:value={newBoard[col["column-code"]]}
                            placeholder={col["display-name"]}
                        />
                    </td>
                {/if}
            {/each}

            <td>
                <button type="button" on:click={closeNewBoard}>Cancel</button>
                <button type="button" on:click={saveNewBoard}>Save</button>
            </td>
        </tr>
    {/if}

    <tbody id="boards-list-body">
        {#each boards as board, index}
            <tr>
                {#if editINDEX == index}
                    <!-- Edit board -->
                    <td />

                    {#each columns as col}
                        {#if col["column-code"] == "idx"}
                            <td>{editBoard["idx"]}</td>
                        {:else if col["column-code"] == "board-type"}
                            <td>
                                <select
                                    bind:value={editBoard[col["column-code"]]}
                                >
                                    <option value="board">Board</option>
                                    <option value="gallery">Gallery</option>
                                </select>
                            </td>
                        {:else if col["column-code"] == "board-table" || col["column-code"] == "comment-table"}
                            <td>
                                <input
                                    bind:value={editBoard[col["column-code"]]}
                                    disabled
                                />
                            </td>
                        {:else if grantSelections.includes(col["column-code"])}
                            <td>
                                <select
                                    bind:value={editBoard[col["column-code"]]}
                                >
                                    {#each Object.entries(grades) as [key, grade]}
                                        <option value={grade.code}>
                                            {grade.name}
                                        </option>
                                    {/each}
                                </select>
                            </td>
                        {:else if col["column-code"] == "regdate"}
                            <td />
                        {:else}
                            <td>
                                <input
                                    type="text"
                                    bind:value={editBoard[col["column-code"]]}
                                    placeholder={col["display-name"]}
                                />
                            </td>
                        {/if}
                    {/each}

                    <td>
                        <button type="button" on:click={closeEditBoard}>
                            Cancel
                        </button>
                        <button type="button" on:click={updateEditBoard}>
                            Save
                        </button>
                    </td>
                {:else}
                    <!-- Show board -->
                    <td>
                        <input
                            type="checkbox"
                            bind:group={selectedIndices}
                            value={board["idx"]}
                        />
                    </td>

                    {#each columns as col}
                        <td>
                            {#if grantSelections.includes(col["column-code"])}
                                {#each Object.entries(grades) as [key, grade]}
                                    {#if grade.code == board[col["column-code"]]}
                                        {grade.name}
                                    {/if}
                                {/each}
                            {:else}
                                {board[col["column-code"]]}
                            {/if}
                        </td>
                    {/each}

                    <td>
                        <button
                            type="button"
                            on:click={() => {
                                moveToListView(index);
                            }}
                        >
                            View
                        </button>
                        <button
                            type="button"
                            on:click={() => {
                                openEditBoard(index);
                            }}
                        >
                            Edit
                        </button>
                        <button
                            type="button"
                            on:click={() => {
                                deleteBoard(index);
                            }}
                        >
                            Delete
                        </button>
                    </td>
                {/if}
            </tr>
        {/each}
    </tbody>
</table>

<datalist id="grant-list">
    <option value="admin">Admin</option>
    <option value="manager">Manager</option>
    <option value="user_active">Regular user</option>
    <option value="user_hold">Pending user</option>
    <option value="user_banned">Banned user</option>
    <option value="guest">Guest</option>
</datalist>

<datalist id="board-types">
    <option value="board">Board</option>
    <option value="gallery">Gallery</option>
</datalist>

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
