var m = require("mithril")

var LogOut = {
    logOut: function() {
	return m.request({
	    method: "POST",
	    url: "/logout",
	})
    }
}

module.exports = LogOut
