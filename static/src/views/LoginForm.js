var m = require("mithril")

module.exports = {
    view: function() {
	
	return m(".login-form", [
	    m("form.form", [
		m("fieldset.fieldset", [
		    m(".container", {class: "container"}, [
			m(".row", {class: "row"}, [
			    m(".legend-col", {class: "col-12"}, [
				m("legend.legend", "Please enter your login details")
			    ])
			])
		    ])
		])
	    ])
	])
    }
}
