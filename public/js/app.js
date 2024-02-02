(function (window) {
	'use strict';

	function invokeToggle(id) {
		const form =  document.getElementById("toggleForm");
		const field = document.getElementById("toggleFormField");
		field.value = id;
		return form.submit();
	}

	function startEdit(id) {
		console.log("start edit?")
		const label = document.getElementById("label-" + id)
		const li = document.getElementById("li-" + id);
		const edit = document.getElementById("edit-" + id);
		edit.value = label.textContent;
		li.classList.add("editing")
		edit.focus();
	}

	function invokeEdit(id) {
		console.log("invokeEdit", id)
		const li = document.getElementById("li-" + id);
		if (!li.classList.contains("editing")) {
			// the onblur event is triggering because the user
			// pressed Escape
			return;
		}
		const form =  document.getElementById("editForm");
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

})(window);
