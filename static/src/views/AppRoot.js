var m = require("mithril")
var Leagues = require("../models/Leagues")
var secure = require("../models/secure")

module.exports = {
    oninit: function(vnode) {
	Leagues.loadList().catch(secure)
    },
    view: function() {
	return m(".app-root", [
	    m("header.app-root-header", "Fooball Prediction App"),
	    m("aside.app-root-sidebar", [
		m("h3.leagues", "Leagues"),
		m(".leagues", 
		  Leagues.list.map(
		    function(league) {
			return m("a.league-list-item", {href: "/league/" + league.ID, oncreate: m.route.link}, league.Name)
		    }
		)),
		m("a.add-league-button", {href: "/leagues/add", oncreate: m.route.link}, "+")
	    ])
	])
    }
}
