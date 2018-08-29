var m = require("mithril")

var User = {
    current: {},
    save: function() {
	return m.request({
	    method: "PUT",
	    url: "/user/new.json",
	    data: User.current,
	    withCredentials: true
	})
    },
    login: function() {
	return m.request({
	    method: "POST",
	    url: "/authenticate",
	    data: User.current,
	    withCredentials: true,
	})
    }
}

module.exports = User
