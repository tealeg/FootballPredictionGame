var m = require("mithril")

var NewUserButton = {
    view: function(){
	return m(".container", [
	    m(".row", [
		m(".col-12", [
		    m("button", {
			onclick: function(e) {
			    window.location.href = "/#!/user/new"
			}
		    }, "Create Account")
		])
	    ])
	])
    }
}

module.exports = NewUserButton


