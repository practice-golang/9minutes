<script>
    export let data;
</script>

<h1>Admin / Board list</h1>

<h1>Boards</h1>

<div>
    <label for="search">Search:</label>
    <input
        type="text"
        id="search"
        onkeyup="pressEnter()"
        placeholder="Search for..."
    />
    <button type="button" onclick="search()">Search</button>
</div>

<button type="button" onclick="openAdd()">Create</button>

<table id="boards-list-container">
    <thead>
        <tr>
            <!-- <td><input type="checkbox" name="select-all" placeholder="Select all" /></td> -->
            <th>Index</th>
            <th>Name</th>
            <th>Code</th>
            <th>Type</th>
            <th>Board table</th>
            <th>Comment table</th>
            <th>Grant read</th>
            <th>Grant write</th>
            <th>Grant comment</th>
            <th>Grant upload</th>
            <th>Control</th>
        </tr>
    </thead>
    <tr id="add-board">
        <!-- <td>&nbsp;</td> -->
        <td>&nbsp;</td>
        <td>
            <input
                type="text"
                name="board-name"
                value=""
                placeholder="Board name"
            />
        </td>
        <td>
            <input
                type="text"
                name="board-code"
                onchange="setTableNameByCode(this)"
                value=""
                placeholder="Board code"
            />
        </td>
        <td>
            <input
                type="text"
                name="board-type"
                onchange="restrictDatalist(this)"
                list="board-types"
                value=""
                autocomplete="off"
                placeholder="Board type"
            />
        </td>
        <td>
            <input
                type="text"
                name="board-table"
                value=""
                placeholder="Board table"
                disabled
            />
        </td>
        <td>
            <input
                type="text"
                name="comment-table"
                value=""
                placeholder="Comment table"
                disabled
            />
        </td>
        <td>
            <input
                type="text"
                name="grant-read"
                onchange="restrictDatalist(this)"
                list="grant-list"
                value=""
                autocomplete="off"
                placeholder="Grant read"
            />
        </td>
        <td>
            <input
                type="text"
                name="grant-write"
                onchange="restrictDatalist(this)"
                list="grant-list"
                value=""
                autocomplete="off"
                placeholder="Grant write"
            />
        </td>
        <td>
            <input
                type="text"
                name="grant-comment"
                onchange="restrictDatalist(this)"
                list="grant-list"
                value=""
                autocomplete="off"
                placeholder="Grant comment"
            />
        </td>
        <td>
            <input
                type="text"
                name="grant-upload"
                onchange="restrictDatalist(this)"
                list="grant-list"
                value=""
                autocomplete="off"
                placeholder="Grant upload"
            />
        </td>
        <td>
            <button type="button" onclick="closeAdd()">Cancel</button>
            <button type="button" onclick="addBoard()">Save</button>
        </td>
    </tr>
    <tbody id="boards-list-body" lr-loop="boardsList">
        <tr lr-if="boardEditIndex != $index">
            <!-- <td><input type="checkbox" name="select$index" placeholder="Select" /></td> -->
            <td>idx</td>
            <td>board-name</td>
            <td>board-code</td>
            <td>board-type</td>
            <td>board-table</td>
            <td>comment-table</td>
            <td>grant-read</td>
            <td>grant-write</td>
            <td>grant-comment</td>
            <td>grant-upload</td>
            <td>
                <button type="button" lr-click="moveToListView($index)">
                    View
                </button>
                <button type="button" lr-click="openEdit($index)">Edit</button>
                <button type="button" lr-click="deleteBoard($index)">
                    Delete
                </button>
            </td>
        </tr>
        <tr lr-if="boardEditIndex == $index">
            <!-- <td>&nbsp;</td> -->
            <td>
                <input
                    type="hidden"
                    name="idx"
                    value="idx"
                    placeholder="Index"
                />
                <span>idx</span>
            </td>
            <td>
                <input
                    type="text"
                    name="board-name"
                    value={data.boardName}
                    placeholder="Board name"
                />
            </td>
            <td>
                <input
                    type="text"
                    name="board-code"
                    value={data.boardCode}
                    placeholder="Board code"
                />
            </td>
            <td>
                <input
                    type="text"
                    name="board-type"
                    list="board-types"
                    value={data.boardType}
                    autocomplete="off"
                    placeholder="Board type"
                    onchange="restrictDatalist(this)"
                />
            </td>
            <td>
                <input
                    type="text"
                    name="board-table"
                    value={data.boardTable}
                    placeholder="Board table"
                    disabled
                />
            </td>
            <td>
                <input
                    type="text"
                    name="comment-table"
                    value={data.commentTable}
                    placeholder="Comment table"
                    disabled
                />
            </td>
            <td>
                <input
                    type="text"
                    name="grant-read"
                    list="grant-list"
                    value={data.grantRead}
                    autocomplete="off"
                    placeholder="Grant read"
                    onchange="restrictDatalist(this)"
                />
            </td>
            <td>
                <input
                    type="text"
                    name="grant-write"
                    list="grant-list"
                    value={data.grantWrite}
                    autocomplete="off"
                    placeholder="Grant write"
                    onchange="restrictDatalist(this)"
                />
            </td>
            <td>
                <input
                    type="text"
                    name="grant-comment"
                    list="grant-list"
                    value={data.grantComment}
                    autocomplete="off"
                    placeholder="Grant comment"
                    onchange="restrictDatalist(this)"
                />
            </td>
            <td>
                <input
                    type="text"
                    name="grant-upload"
                    list="grant-list"
                    value={data.grantUpload}
                    autocomplete="off"
                    placeholder="Grant upload"
                    onchange="restrictDatalist(this)"
                />
            </td>
            <td>
                <button type="button" onclick="closeEdit()">Cancel</button>
                <button type="button" onclick="updateBoard()">Save</button>
            </td>
        </tr>
    </tbody>
</table>

<datalist id="grant-list">
    <option value="admin">Admin</option>
    <option value="manager">Manager</option>
    <option value="regular_user">Regular user</option>
    <option value="pending_user">Pending user</option>
    <option value="banned_user">Banned user</option>
    <option value="guest">Guest</option>
</datalist>

<datalist id="board-types">
    <option value="board">Board</option>
    <option value="gallery">Gallery</option>
</datalist>

<div id="pages-container">
    <div lr-loop="pages">
        <span lr-if="$index == 0 && pages[0].page > 1">&laquo;</span>
        <span lr-if="$index == 0 && pages[0].page > 1">&lt;</span>

        <b lr-if="page == boardsData['current-page']">
            <a href={data.link} rel="external">
                {data.page}
            </a>
        </b>
        <a
            lr-if="page != boardsData['current-page']"
            href={data.link}
            rel="external"
        >
            {data.page}
        </a>

        <span
            lr-if="page < boardsData['total-page'] && $index == (pages.length - 1)"
        >
            &gt;
        </span>
        <span
            lr-if="page < boardsData['total-page'] && $index == (pages.length - 1)"
        >
            &raquo;
        </span>
    </div>
</div>

<style>
    table,
    th,
    td {
        border: 1px solid black;
    }

    #add-board {
        display: none;
    }
</style>
