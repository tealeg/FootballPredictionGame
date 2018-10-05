var m = require("mithril")
var Leagues = require("../models/Leagues")
var Outline = require("./Outline")
var User = require("../models/User")
var secure = require("../models/secure")
var ErrorBar = require("./ErrorBar")

function renderAdminLink(adminMode) {
    if (adminMode) {
	return m(".col-12", [
	    m(".container", [
		m(".col-12", [
		    m("a", {href: "/leagues", oncreate: m.route.link}, "Admin")
		])
	    ])
	])
    }
}

AppRoot = {
    adminMode: false,
    oninit: function(vnode){
	return User.isAdmin().then(
	    function(result){
		AppRoot.adminMode = result
	    }).catch(secure).catch(
		function(err) {
		    ErrorBar.errors.push(err)
		}
	    )
    },
    view: function() {
	return m(Outline, [
	    m(".col12", [
		m(".container", [
		    m(ErrorBar),
		    renderAdminLink(AppRoot.adminMode),
		])
	    ])

	])
    }
}

module.exports = AppRoot
