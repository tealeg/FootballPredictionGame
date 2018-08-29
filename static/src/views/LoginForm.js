var m = require("mithril")
var User = require("../models/User")

var LoginForm = {
    errors: [],
    view: function() {
	return m(".login-form", [
	    m("form", {
		onsubmit: function(e) {
		    e.preventDefault()
		    User.login().then(
			window.location.href = "/"
		    ).catch(secure).catch(
			function(err) {
			    LoginForm.errors.push(err)
			}
		    )
		}
	    }, [
		m("fieldset", [
		    m(".container", [
			m(".row", [
			    m(".col-12", [
				m("legend", "Please enter your login details")
			    ])
			]),
			m(".row", [
			    m(".col-3", [
				m("label", {for: "username"}, "Username")
			    ]),
			    m(".col-9", [
				m("input", {
				    type: "text",
				    name: "username",
				    oninput: m.withAttr("value", function(value){
					User.current.UserName = value
				    })
				})
			    ])
			]),
			m(".row", [
			    m(".col-3", [
				m("label", {for: "password"}, "password")
			    ]),
			    m(".col-9", [
				m("input", {
				    type: "password",
				    name: "password",
				    oninput: m.withAttr("value", function(value){
					User.current.Password = value
				    })
				})
			    ])
			]),
			m(".row", [
			    m(".col-9"),
			    m(".col-3", [
				m("input", {
				    type: "submit",
				    value: "Log In",
				})
			    ])
			])
			
		    ])
		])
	    ])
	])
    }
}


module.exports = LoginForm
