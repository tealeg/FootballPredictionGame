var m = require("mithril")

var secure = require("./secure")


var League = {
    current: {},
    errors: [],
    seasons: [],
    load: function(id) {
	return m.request({
	    method: "GET",
	    url: "league/" + id,
	    withCredentials: true,
	}).then(function(result) {
            League.current = result
	    // League.seasons = loadSeasons(League.current)
        })
    },
    // loadSeasons: function(league) {
    // 	return m.request({
    // 	    method: "GET",
    // 	    url: ""
    // 	})
    // },
    
}

module.exports = League
