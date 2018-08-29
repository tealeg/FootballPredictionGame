var m = require("mithril")

var Leagues = {
    list: [],
    loadList: function() {
	return m.request({
	    method: "GET",
	    url: "/leagues.json",
	    withCredentials: true,
	}).then(
	    function(result) {
		Leagues.list = result.data
	    }
	)
    },
    new: [],
    save: function() {
	return m.request({
	    method: "PUT",
	    url: "/leagues/new.json",
	    data: Leagues.new,
	    withCredentials: true,
	})
    }
}

module.exports = Leagues
