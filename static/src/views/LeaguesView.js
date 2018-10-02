
var m = require("mithril")

var Leagues = require("../models/Leagues")
var secure = require("../models/secure")
var ErrorBar = require("./ErrorBar")
var Outline = require("./Outline")

function renderLeagues(leagueList) {
    return m("ul", leagueList.map(
	function(league) {
	    return m("li", [
		m("a", {
		    href: "/league/" + league.ID,
		    oncreate: m.route.link
		}, league.Name)
	    ])
	}
    ))
}

function renderAddLeagueButton() {
    return m("a", {
	href: "/leagues/add",
	oncreate: m.route.link,
    }, "Add League")

}

LeaguesView = {
    oninit: function(vnode) {
	Leagues.loadList().catch(secure).catch(
	    function(err) {
		ErrorBar.errors.push(err)
	    }
	)
    },

    view: function() {
	return m(Outline, [
	    m(".col-12", [
		m(".container", [
		    m(ErrorBar),
		    renderLeagues(Leagues.list),
		    renderAddLeagueButton(),
		])
	    ])
	])
    }
}

module.exports = LeaguesView
