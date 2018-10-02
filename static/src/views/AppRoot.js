var m = require("mithril")
var Leagues = require("../models/Leagues")
var Outline = require("./Outline")

AppRoot = {
    view: function() {
	return m(Outline, [
	    m(".col-12", [
		m(".container", [
		    m(".col-12", [
			m("a", {href: "/leagues", oncreate: m.route.link}, "Admin")
		    ])
		])
	    ])
	])
    }
}

module.exports = AppRoot
