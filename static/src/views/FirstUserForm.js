var m = require("mithril")
var UserForm = require("./UserForm")

var FirstUserForm = {
    errors: [],
    view: function() {
	return m(".first-user-form", [
	    m("form.form",
	      {
		  onsubmit: function(e) {
		      e.preventDefault()
		      UserForm.user.current.isAdmin = true
		      UserForm.user.save().then(
			  function(response) {
			      window.location.href = "/"
			  }
		      ).catch(secure).catch(
			  function(err) {
			      FirstUserForm.errors.push(err)
			  }
		      )
		  }
	      }
	      ,[m(UserForm)])
	])
    }
}

module.exports = FirstUserForm
