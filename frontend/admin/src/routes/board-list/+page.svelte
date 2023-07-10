<script>
    export let data;
</script>

<h1>Admin / Board list</h1>

<h1>User fields</h1>

<button type="button" onclick="openAdd()">Add column</button>

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
    <tr id="add-column">
        <td>&nbsp;</td>
        <td>
            <input
                type="text"
                name="display-name"
                value=""
                placeholder="Display name"
            />
        </td>
        <td>
            <input
                type="text"
                name="column-code"
                value=""
                placeholder="Column code"
            />
        </td>
        <td>
            <input
                type="text"
                name="column-type"
                value=""
                onchange="restrictDatalist(this)"
                list="column-types"
                placeholder="Column type"
            />
        </td>
        <td>
            <input
                type="text"
                name="column-name"
                value=""
                placeholder="Column name"
            />
        </td>
        <td>&nbsp;</td>
        <td>
            <button type="button" onclick="closeAdd()">Cancel</button>
            <button type="button" onclick="addColumn()">Save</button>
        </td>
    </tr>
    <tbody id="column-list-body" lr-loop="columnList">
        <tr lr-if="columnEditIndex != $index">
            <td>
                <input
                    type="checkbox"
                    name="select$index"
                    placeholder="Select"
                />
            </td>
            <td>
                {data.displayName}
            </td>
            <td>{data.columnCode}</td>
            <td>{data.columnType}</td>
            <td>{data.columnName}</td>
            <td>{data.sortOrder}</td>
            <td lr-if="$index <= 6">&nbsp;</td>
            <td lr-if="$index > 6">
                <button type="button" lr-click="openEdit($index)">Edit</button>
                <button type="button" lr-click="deleteColumn($index)">
                    Delete
                </button>
            </td>
        </tr>
        <tr lr-if="columnEditIndex == $index">
            <td>&nbsp;</td>
            <td>
                <input
                    type="text"
                    name="display-name"
                    value={data.displayName}
                    placeholder="Display name"
                />
            </td>
            <td>
                <input
                    type="text"
                    name="column-code"
                    value={data.columnCode}
                    placeholder="Column code"
                />
            </td>
            <td>
                <input
                    type="text"
                    name="column-type"
                    value={data.columnType}
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
                    value={data.columnName}
                    placeholder="Column name"
                />
            </td>
            <td>
                <input type="hidden" name="sort-order" value={data.sortOrder} />
                <span>{data.sortOrder}</span>
            </td>
            <td>
                <button type="button" onclick="closeEdit()">Cancel</button>
                <button type="button" onclick="updateColumn()">Save</button>
            </td>
        </tr>
    </tbody>
</table>

<datalist id="column-types">
    <option value="text">Text</option>
    <option value="number-integer">Number Integer</option>
    <option value="number-real">Number Real</option>
</datalist>
