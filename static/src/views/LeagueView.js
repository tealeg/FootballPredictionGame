var m = require("mithril")
var ErrorBar = require("./ErrorBar")
var Outline = require("./Outline")
var League = require("../models/League")
var secure = require("../models/secure")

function renderSeasons(seasonList) {
    return m("ul", seasonList.map(
	function(season) {
	    return m("li", [
		m("a", {
		    href: "/league/" + LeagueView.leauge.ID + "/season/" + season.ID,
		    oncreate: m.route.link,
		})
	    ])
	}
    ));
}

LeagueView = {
    oninit: function(vnode){
	League.load(vnode.attrs.id).catch(secure).catch(
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
		    m("h3", League.current.name),
		    renderSeasons(League.current.seasons),
		])
	    ])
	])
    }
}

module.exports = LeagueView
