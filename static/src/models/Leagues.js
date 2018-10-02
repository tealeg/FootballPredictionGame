var m = require("mithril")

var Leagues = {
    errors: [],
    list: [],
    loadList: function() {
	return m.request({
	    method: "GET",
	    url: "/leagues.json",
	    withCredentials: true,
	}).then(
	    function(result) {
		Leagues.list = result
	    }
	)
    },
    new: {},
    save: function() {
	return m.request({
	    method: "POST",
	    url: "/leagues/new.json",
	    data: Leagues.new,
	    withCredentials: true,
	})
    }
}

module.exports = Leagues
