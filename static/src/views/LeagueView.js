var m = require("mithril")
var ErrorBar = require("./ErrorBar")
var Outline = require("./Outline")
var League = require("../models/League")
var secure = require("../models/secure")

function renderSeasons(seasonList) {
    if (seasonList) {
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
    return null
}

function renderAddSeasonButton(lid) {
    return m("a", {
	href: "/league/" + lid + "/seasons/new" ,
	oncreate: m.route.link,
    }, "Add Season")

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
		    renderAddSeasonButton(League.current.id),
		])
	    ])
	])
    }
}

module.exports = LeagueView
