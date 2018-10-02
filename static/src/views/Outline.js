var m = require("mithril")

var LogOutButton = require("../views/LogOutButton")

Outline = {
    view: function(vnode) {
	return m(".app-root", [
	    m(".container", [
		m(".row", [
		    m(".col-12", [
			m("header.app-root-header", [
			    m("a", {href: "/"},
			      "Fooball Prediction App")
			]),
		    ]),
		]),
		m(".row",
		  [vnode.children]),
		m(".row", [
		    m(".col-9"),
		    m(".col-3", [m(LogOutButton)])
		])
	    ])
	])
    }
}

module.exports = Outline
