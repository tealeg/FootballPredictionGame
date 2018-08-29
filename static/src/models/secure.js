var m = require("mithril")

function adminUserExists() {
    let exists = false
    return m.request({
	    method: "GET",
	    url: "/user/admin/exists.json",
	})
}

module.exports  = function(error) {
    if (error.message == "Cookie check failed") {
	return adminUserExists().then(
	    function(result){
		if (result == true) {
		    window.location.href = "/#!/login"
		} else {
		    window.location.href = "/#!/firstuser"
		}
	    }
	)
    }
    throw error
}
