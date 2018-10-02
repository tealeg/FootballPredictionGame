var m = require("mithril")
var User = require("../models/User")

var LogOutButton = {
    view: function(){
	return m("form", {
	    onsubmit: function(e) {
		e.preventDefault()
		User.logOut().then(
		    function (response) {
			window.location.href = "/"
		    }
		).catch(secure)
	    }, 
	}, [m("input", {type: "submit", value: "Log Out"})])
    }
}

module.exports = LogOutButton
