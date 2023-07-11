<script>
    import { invalidateAll } from "$app/navigation";

    export let data;

    const columns = data.columns;
    $: users = data.users;
</script>

<h1>Admin / User list</h1>

<h1>Users list</h1>

<button type="button" onclick="openAdd()">Add user</button>

<label for="search">Search:</label>
<input
    type="text"
    id="search"
    onkeyup="pressEnter()"
    placeholder="Search for..."
/>
<button type="button" onclick="search()">Search</button>

<table id="users-list-container">
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
    <tr id="add-user">
        <td />
        <td />
        <td>
            <input type="text" name="userid" value="" placeholder="Userid" />
        </td>
        <td>
            <input
                type="password"
                name="password"
                value=""
                placeholder="Password"
            />
        </td>
        <td><input type="text" name="email" value="" placeholder="Email" /></td>
        <!-- <td><input type="text" name="phone" value="" placeholder="Phone" /></td> -->
        <td>
            <input
                type="text"
                name="grade"
                onchange="restrictDatalist(this)"
                list="grade-list"
                value=""
                placeholder="Grade"
            />
        </td>
        <td>
            <input
                type="text"
                name="approval"
                onchange="restrictDatalist(this)"
                list="yn-list"
                value=""
                autocomplete="off"
                placeholder="Approval"
            />
        </td>
        <td>
            <button type="button" onclick="closeAdd()">Cancel</button>
            <button type="button" onclick="addUser()">Save</button>
        </td>
    </tr>
    <tbody id="users-list-body">
        <tr>
            <td />
            <td>
                <input
                    type="hidden"
                    name="idx"
                    value={data.idx}
                    placeholder="Index"
                />
                <span>{data.idx}</span>
            </td>
            <td>
                <input
                    type="text"
                    name="userid"
                    value={data.userid}
                    placeholder="Userid"
                />
            </td>
            <td>
                <input
                    type="password"
                    name="password"
                    value=""
                    placeholder="Password"
                />
            </td>
            <td>
                <input
                    type="text"
                    name="email"
                    value={data.email}
                    placeholder="Email"
                />
            </td>
            <!-- <td><input type="text" name="phone" value="{data.phone}" placeholder="Phone" /></td> -->
            <td>
                <input
                    type="text"
                    name="grade"
                    onchange="restrictDatalist(this)"
                    list="grade-list"
                    value={data.grade}
                    placeholder="Grade"
                />
            </td>
            <td>
                <input
                    type="text"
                    name="approval"
                    onchange="restrictDatalist(this)"
                    list="yn-list"
                    value={data.approval}
                    autocomplete="off"
                    placeholder="Approval"
                />
            </td>
            <td>
                <button type="button" onclick="closeEdit()">Cancel</button>
                <button type="button" onclick="updateUser()">Save</button>
            </td>
        </tr>
        {#each users as user}
            <tr>
                <td>
                    <input type="checkbox" />
                </td>
                {#each columns as col}
                    <td>{user[col["column-code"]]}</td>
                {/each}
                <td>
                    <button type="button" lr-click="openEdit($index)">
                        Edit
                    </button>
                    <button type="button" lr-click="deleteUser($index)">
                        Delete
                    </button>
                </td>
            </tr>
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

<datalist id="yn-list">
    <option value="Y">Y</option>
    <option value="N">N</option>
</datalist>

<div id="pages-container">
    <div lr-loop="pages">
        <span lr-if="$index == 0 && pages[0].page > 1">&laquo;</span>
        <span lr-if="$index == 0 && pages[0].page > 1">&lt;</span>

        <b lr-if="page == usersData['current-page']">
            <a href={data.link} rel="external">
                {data.page}
            </a>
        </b>
        <a
            lr-if="page != usersData['current-page']"
            href={data.link}
            rel="external"
        >
            {data.page}
        </a>

        <span
            lr-if="page < usersData['total-page'] && $index == (pages.length - 1)"
            >&gt;</span
        >
        <span
            lr-if="page < usersData['total-page'] && $index == (pages.length - 1)"
            >&raquo;</span
        >
    </div>
</div>
