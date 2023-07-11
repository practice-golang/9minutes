<script>
    export let data;
    const { columns } = data;

    let editIDX = -1;
    let showNewCOL = false;
    const newCOL = {
        "display-name": "",
        "column-code": "",
        "column-type": "",
        "column-name": "",
    };
    const editCOL = {
        "display-name": "",
        "column-code": "",
        "column-type": "",
        "column-name": "",
    };

    function closeNewColumn() {
        newCOL["display-name"] = "";
        newCOL["column-code"] = "";
        newCOL["column-type"] = "";
        newCOL["column-name"] = "";

        showNewCOL = false;
    }

    function saveNewColumn() {
        closeNewColumn();
    }

    function closeEditColumn() {
        editCOL["display-name"] = "";
        editCOL["column-code"] = "";
        editCOL["column-type"] = "";
        editCOL["column-name"] = "";

        editIDX = -1;
    }
    function saveEditColumn() {
        closeEditColumn();
    }

    function deleteColumn(idx) {
        console.log(columns[idx]);
    }
</script>

<h1>Admin / User columns</h1>

<h1>User fields</h1>

<button
    type="button"
    on:click={() => {
        showNewCOL = true;
    }}>Add column</button
>

<table id="column-list-container">
    <thead>
        <tr>
            <td>
                <input
                    type="checkbox"
                    name="select-all"
                    placeholder="Select all"
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
            <td>&nbsp;</td>
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
                <input
                    type="text"
                    name="column-type"
                    bind:value={newCOL["column-type"]}
                    onchange="restrictDatalist(this)"
                    list="column-types"
                    placeholder="Column type"
                />
            </td>
            <td>
                <input
                    type="text"
                    name="column-name"
                    bind:value={newCOL["column-name"]}
                    placeholder="Column name"
                />
            </td>
            <td>&nbsp;</td>
            <td>
                <button type="button" on:click={closeNewColumn}>
                    Cancel
                </button>
                <button type="button" on:click={saveNewColumn}> Save </button>
            </td>
        </tr>
    {/if}

    <tbody id="column-list-body">
        {#each columns as col, idx}
            {#if editIDX != idx}
                <tr lr-if="columnEditIndex != $index">
                    {#if idx <= 6}
                        <td />
                    {:else}
                        <td>
                            <input
                                type="checkbox"
                                name="select$index"
                                placeholder="Select"
                            />
                        </td>
                    {/if}
                    <td>{col["display-name"]}</td>
                    <td>{col["column-code"]}</td>
                    <td>{col["column-type"]}</td>
                    <td>{col["column-name"]}</td>
                    <td>{col["sort-order"]}</td>
                    {#if idx <= 6}
                        <td />
                    {:else}
                        <td>
                            <button
                                type="button"
                                on:click={() => {
                                    editIDX = idx;
                                }}
                            >
                                Edit
                            </button>
                            <button
                                type="button"
                                lr-click="deleteColumn($index)"
                            >
                                Delete
                            </button>
                        </td>
                    {/if}
                </tr>
            {:else}
                <tr lr-if="columnEditIndex == $index">
                    <td />
                    <td>
                        <input
                            type="text"
                            name="display-name"
                            value={col["display-name"]}
                            placeholder="Display name"
                        />
                    </td>
                    <td>
                        <input
                            type="text"
                            name="column-code"
                            value={col["column-code"]}
                            placeholder="Column code"
                        />
                    </td>
                    <td>
                        <input
                            type="text"
                            name="column-type"
                            value={col["column-type"]}
                            onchange="restrictDatalist(this)"
                            list="column-types"
                            placeholder="Column type"
                            disabled
                        />
                    </td>
                    <td>
                        <input
                            type="text"
                            name="column-name"
                            value={col["column-name"]}
                            placeholder="Column name"
                        />
                    </td>
                    <td>
                        <input
                            type="hidden"
                            name="sort-order"
                            value={col["sort-order"]}
                        />
                        <span>{col["sort-order"]}</span>
                    </td>
                    <td>
                        <button
                            type="button"
                            on:click={() => {
                                editIDX = -1;
                            }}
                        >
                            Cancel
                        </button>
                        <button type="button" onclick="updateColumn()">
                            Save
                        </button>
                    </td>
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
