'use strict';

function invokeToggle(itemId) {
	const form = document.getElementById("toggleForm");
	const field = document.getElementById("toggleFormField");
	field.value = itemId;
	return form.submit();
}

function invokeDestroy(itemId) {
	const form = document.getElementById("destroyForm");
	const field = document.getElementById("destroyFormField");
	field.value = itemId;
	return form.submit();
}

function startEdit(id) {
	console.log("start edit?")
	window.location.replace("/?edit=" + id)
}

function invokeEdit(id) {
	console.log("invokeEdit", id)
	const li = document.getElementById("li-" + id);
	if (!li.classList.contains("editing")) {
		// the onblur event is triggering because the user
		// pressed Escape
		return;
	}
	const form = document.getElementById("editForm");
	const idField = document.getElementById("editFormItemId");
	const titleField = document.getElementById("editFormItemTitle");
	const edit = document.getElementById("edit-" + id);
	idField.value = id;
	titleField.value = edit.value
	return form.submit();
}

function abortEdit(id) {
	console.log("abortEdit", id)
	const li = document.getElementById("li-" + id);
	li.classList.remove("editing")
}

function stopEdit(id) {
	invokeEdit(id)
}

function handleEditKeyDown(event) {
	if (event.key === 'Enter') {
		invokeEdit(event.target.id.split("-")[1])
	}
	if (event.key === 'Escape') {
		abortEdit(event.target.id.split("-")[1])
	}
}
