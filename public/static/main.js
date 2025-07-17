function handleGetItemsByChannelId(event) {
    event.preventDefault();
    const id = document.getElementById('getChannelId').value;
    window.location.href = `/channels/${id}/`;
}

function handleDeleteChannel(event) {
    event.preventDefault();
    const id = document.getElementById('deleteChannelId').value;
    fetch(`/channels/${id}/`, {
        method: 'DELETE'
    }).then(() => window.location.reload());
}

function handlePutChannel(event) {
    event.preventDefault();
    const data = {
        id: document.getElementById('putChannelId').value.trim(),
        title: document.getElementById('putChannelTitle').value,
        language: document.getElementById('putChannelLanguage').value,
        description: document.getElementById('putChannelDescription').value,
    };

    fetch(`/channels/${data.id}/`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    }).then(() => window.location.reload());
}

function handleGetItemById(event) {
    event.preventDefault();
    const id = document.getElementById('getItemId').value;
    window.location.href = `/items/${id}/`;
}

function handleDeleteItem(event) {
    event.preventDefault();
    const id = document.getElementById('deleteItemId').value;
    fetch(`/items/${id}/`, {
        method: 'DELETE'
    }).then(() => window.location.reload());
}

function handlePutItem(event) {
    event.preventDefault();
    const data = {
        id: document.getElementById('putItemId').value.trim(),
        title: document.getElementById('putItemTitle').value,
        description: document.getElementById('putItemDescription').value,
        pub_date: document.getElementById('putItemPubDate').value,
    };

    fetch(`/items/${data.id}/`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    }).then(() => window.location.reload());
}
