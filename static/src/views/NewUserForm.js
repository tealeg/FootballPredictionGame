var m = require("mithril")
var UserForm = require("../views/UserForm")

var NewUserForm = {
    view: function () {
	return m(".new-user-form", [
	    m("form.form",
	      {
		  onsubmit: function(e) {
		      e.preventDefault()
		      UserForm.user.current.isAdmin = false
		      UserForm.user.save().then(
			  function(response){ 
			      window.location.href = "/"
			  }
		      ).catch(secure).catch(
			  function(err) {
			      NewUserForm.errors.push(err)
			  }
		      )
		  }
	      }
	      , m(UserForm)
	     )
	])
    }
}

module.exports = NewUserForm

