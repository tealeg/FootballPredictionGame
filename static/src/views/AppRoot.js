var m = require("mithril")
var Leagues = require("../models/Leagues")
var secure = require("../models/secure")
var LogOutButton = require("../views/LogOutButton")

AppRoot = {
    view: function() {
	return m(".app-root", [
	    m(".container", [
		m(".row", [
		    m(".col-12", [
			m("header.app-root-header", "Fooball Prediction App"),
		    ])
		]),

		m(".row", [
		    m(".col-3", [m("a", {href: "/leagues", oncreate: m.route.link}, "Admin")]),
		    m(".col-3", []),
		    m(".col-3", [			
			    m(LogOutButton)
		    ])
		])
	    ])
	])
    }
}

module.exports = AppRoot
