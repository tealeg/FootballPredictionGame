var m = require("mithril")

var User = {
    current: {},
    save: function() {
	return m.request({
	    method: "POST",
	    url: "/user/new.json",
	    data: User.current,
	})
    },
    login: function() {
	return m.request({
	    method: "PUT",
	    url: "/authenticate",
	    data: User.current,
	})
    },
    logOut: function() {
	return m.request({
	    method: "GET",
	    url: "/logout",
	})
    }

}

module.exports = User
