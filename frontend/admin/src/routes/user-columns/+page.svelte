<script>
    import { invalidateAll } from "$app/navigation";

    export let data;
    $: columns = data.columns;

    const lastDefaultColIDX = 7;

    $: selectAll = selectedIndices.length == columns.length - lastDefaultColIDX;
    let selectedIndices = [];

    let editINDEX = -1;
    let showNewCOL = false;
    let newCOL = {};
    let editCOL = {};

    const colTYPES = {
        text: "Text",
        "number-integer": "Number Integer",
        "number-real": "Number Real",
    };

    // https://stackoverflow.com/a/73866413/8964990 , https://svelte.dev/repl/e54cf61806b0474a803769daaeb12f1b?version=3.50.1
    function toggleSelectAll() {
        if (selectAll) {
            selectedIndices = [];
        } else {
            selectedIndices = columns
                .filter((col) => col["idx"] > lastDefaultColIDX)
                .map((col) => col["idx"]);
        }
    }

    function closeNewColumn() {
        newCOL = {};
        showNewCOL = false;
    }

    async function saveNewColumn() {
        const uri = "/api/admin/user-columns";
        const r = await fetch(uri, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(newCOL),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        closeNewColumn();
        invalidateAll();
    }

    function openEditColumn(index) {
        editINDEX = index;
        editCOL = {};
        for (const k in columns[index]) {
            editCOL[k] = columns[index][k];
        }
    }

    function closeEditColumn() {
        editCOL = {};
        editINDEX = -1;
    }

    async function updateEditColumn() {
        const uri = "/api/admin/user-columns";
        const r = await fetch(uri, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify([editCOL]),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        closeEditColumn();
        invalidateAll();
    }

    async function deleteColumn(index) {
        const colIDX = columns[index]["idx"];

        const uri = "/api/admin/user-columns";
        const r = await fetch(uri, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify([{ idx: parseInt(colIDX) }]),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        selectedIndices = [];
        invalidateAll();
    }

    async function deleteSelectedColumns() {
        if (selectedIndices.length == 0) {
            alert("No columns selected");
            return;
        }

        const colIndices = [];
        for (let i = 0; i < selectedIndices.length; i++) {
            colIndices.push({ idx: parseInt(selectedIndices[i]) });
        }

        const uri = "/api/admin/user-columns";
        const r = await fetch(uri, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(colIndices),
        });

        if (!r.ok) {
            alert(await r.text());
        }

        selectedIndices = [];
        invalidateAll();
    }
</script>

<h1>Admin / User columns</h1>

<h1>User fields</h1>

<button
    type="button"
    on:click={() => {
        showNewCOL = true;
    }}
>
    Add column
</button>

<button type="button" on:click={deleteSelectedColumns}>
    Delete selected columns
</button>

<table id="column-list-container">
    <thead>
        <tr>
            <td>
                <input
                    type="checkbox"
                    checked={selectAll}
                    on:click={toggleSelectAll}
                />
            </td>
            <th>Display name</th>
            <th>Column code</th>
            <th>Column type</th>
            <th>Column name</th>
            <th>Sort Order</th>
            <th>Control</th>
        </tr>
    </thead>

    {#if showNewCOL}
        <tr id="add-column">
            <td />
            <td>
                <input
                    type="text"
                    name="display-name"
                    bind:value={newCOL["display-name"]}
                    placeholder="Display name"
                />
            </td>
            <td>
                <input
                    type="text"
                    name="column-code"
                    bind:value={newCOL["column-code"]}
                    placeholder="Column code"
                />
            </td>
            <td>
                <select name="column-type" bind:value={newCOL["column-type"]}>
                    {#each Object.entries(colTYPES) as [key, name]}
                        <option value={key}>{name}</option>
                    {/each}
                </select>
            </td>
            <td>
                <input
                    type="text"
                    name="column-name"
                    bind:value={newCOL["column-name"]}
                    placeholder="Column name"
                />
            </td>
            <td />
            <td>
                <button type="button" on:click={closeNewColumn}>
                    Cancel
                </button>
                <button type="button" on:click={saveNewColumn}> Save </button>
            </td>
        </tr>
    {/if}

    <tbody id="column-list-body">
        {#each columns as col, index}
            {#if editINDEX == index}
                <tr>
                    <td />
                    <td>
                        <input
                            type="text"
                            name="display-name"
                            bind:value={editCOL["display-name"]}
                            placeholder="Display name"
                        />
                    </td>
                    <td>
                        <input
                            type="text"
                            name="column-code"
                            bind:value={editCOL["column-code"]}
                            placeholder="Column code"
                        />
                    </td>
                    <td>{colTYPES[col["column-type"]]}</td>
                    <td>
                        <input
                            type="text"
                            name="column-name"
                            bind:value={editCOL["column-name"]}
                            placeholder="Column name"
                        />
                    </td>
                    <td>
                        <input
                            type="number"
                            name="sort-order"
                            bind:value={editCOL["sort-order"]}
                        />
                    </td>
                    <td>
                        <button type="button" on:click={closeEditColumn}>
                            Cancel
                        </button>
                        <button type="button" on:click={updateEditColumn}>
                            Save
                        </button>
                    </td>
                </tr>
            {:else}
                <tr>
                    {#if index <= 6}
                        <td />
                    {:else}
                        <td>
                            <input
                                type="checkbox"
                                bind:group={selectedIndices}
                                value={col["idx"]}
                            />
                        </td>
                    {/if}
                    <td>{col["display-name"]}</td>
                    <td>{col["column-code"]}</td>
                    <td>{colTYPES[col["column-type"]]}</td>
                    <td>{col["column-name"]}</td>
                    <td>{col["sort-order"]}</td>
                    {#if index <= 6}
                        <td />
                    {:else}
                        <td>
                            <button
                                type="button"
                                on:click={() => openEditColumn(index)}
                            >
                                Edit
                            </button>
                            <button
                                type="button"
                                on:click={() => deleteColumn(index)}
                            >
                                Delete
                            </button>
                        </td>
                    {/if}
                </tr>
            {/if}
        {/each}
    </tbody>
</table>

<datalist id="column-types">
    <option value="text">Text</option>
    <option value="number-integer">Number Integer</option>
    <option value="number-real">Number Real</option>
</datalist>
