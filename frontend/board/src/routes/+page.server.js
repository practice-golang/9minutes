export async function load() {
    let message = 'Hello from the server at build time';

    const r = await fetch('https://jsonplaceholder.typicode.com/posts');
    if (r.ok) {
        const data = await r.json();
        message = data[0].body;
    }

    return {
        message: message,
    };
}