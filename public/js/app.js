'use strict';

function startEdit(id) {
	window.location.replace(window.location.pathname + "?edit=" + id)
}

function abortEdit() {
	window.location.replace("/")
}

function invokeEdit(id) {
	console.log("invokeEdit", id)
	const form = document.getElementById("editForm");
	const idField = document.getElementById("editFormItemId");
	const titleField = document.getElementById("editFormItemTitle");
	const edit = document.getElementById("edit-" + id);
	idField.value = id;
	titleField.value = edit.value
	return form.submit();
}

function handleEditKeyDown(event) {
	if (event.key === 'Enter') {
		invokeEdit(event.target.id.split("-")[1])
	}
	if (event.key === 'Escape') {
		abortEdit(event.target.id.split("-")[1])
	}
}
