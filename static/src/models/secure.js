var m = require("mithril")

function adminUserExists() {
    let exists = false
    m.request({
	    method: "GET",
	    url: "/user/admin/exists.json",
	    async: false,
	}).then(
	    function(result){
		exists = result
	    }
	).catch(
	    function(err) {
		console.log(err)
		exists = false
	    }
	)
    return exists
}

module.exports  = function(error) {
    if (error.message == "Cookie check failed") {
	if (adminUserExists()) {
	    window.location.href = "/#!/login"	    
	}
	window.location.href = "/#!/firstuser"
    }
}
