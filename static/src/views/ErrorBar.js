var m = require("mithril")

function renderErrorBar(err) {
    return m(".row", [
	m(".col-12", [
	    m(".error", err)
	])
    ])
}

ErrorBar = {
    errors: [],
    view: function() {
	return ErrorBar.errors.map(renderErrorBar)
    }
}

module.exports = ErrorBar
