<script>
    import { onMount, onDestroy, beforeUpdate, afterUpdate } from "svelte";
    import { invalidateAll } from "$app/navigation";
    import { page } from "$app/stores";

    import moment from "moment";
    // import "moment/dist/locale/ko";

    export let data;

    onMount(() => {});

    afterUpdate(() => {});
</script>

<h1>Read</h1>


<div class="content-container">
    <h1>$TITLE$</h1>
    <p>Author: $AUTHOR_NAME$ / Views: $VIEWS$</p>
    <div id="content">$CONTENT$</div>

    <div id="filelist-container">
        <ul lr-loop="filelist">
            <li><a href="{{.link}}" download="{{.filename}}" target="_blank">{{.filename}}</a></li>
        </ul>
    </div>

    <button type="button" onclick="moveToEdit()">Edit</button>
    <button type="button" onclick="deleteContent()">Delete</button>
    <button type="button" onclick="moveToList()">Back to list</button>

    <hr />

    <div id="comments-container">
        Comments:
        <div lr-loop="comments">
            <div class="comment-item">
                {{.authorName}}
                <button type="button" lr-click="deleteComment($index)">X</button>
            </div>
            <div>{{.content}}</div>
        </div>
    </div>

    <div id="pages-container">
        <div lr-loop="pages">
            <span lr-if="$index == 0 && pages[0].page > 1">&laquo;</span>
            <span lr-if="$index == 0 && pages[0].page > 1">&lt;</span>

            <b lr-if="page == commentsData['current-page']">{{.page}}</b>
            <a lr-if="page != commentsData['current-page']" lr-click="fetchComments('{{.page}}')">{{.page}}</a>

            <span lr-if="$index == (pages.length - 1) && pages[0].page < pages[pages.length - 1]">&gt;</span>
            <span lr-if="$index == (pages.length - 1) && pages[0].page < pages[pages.length - 1]">&raquo;</span>
        </div>
    </div>

    <label for="comment-area">Write a comment:</label>
    <div>
        <textarea id="comment-area"></textarea>
    </div>
    <button type="button" onclick="writeComment()">Save comment</button>
</div>