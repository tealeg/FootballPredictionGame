var m = require("mithril")
var LogOut = require("../models/LogOut")

var LogOutButton = {
    view: function(){
	return m("form", {
	    onsubmit: function(e) {
		e.preventDefault()
		LogOut.logOut().then(
		    function (response) {
			window.location.href = "/"
		    }
		).catch(secure)
	    }, 
	}, [m("input", {type: "submit", value: "Log Out"})])
    }
}

module.exports = LogOutButton
