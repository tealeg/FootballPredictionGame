var m = require("mithril")
var User = require("../models/User")
var secure = require("../models/secure")
var NewUserButton = require("../views/NewUserButton")

var LoginForm = {
    errors: [],
    view: function() {
	return m(".login-form", [
	    m("ul.errors", LoginForm.errors.map(
		function(err) {
		    return m("li", err)
		}
	    )),
	    m("form", {
		onsubmit: function(e) {
		    e.preventDefault()
		    User.login().then(
			function (response) {
			    if (response.Errors.length > 0) {
				LoginForm.errors = response.Errors
			    } else {
				window.location.href = "/"
			    }
			}
		    ).catch(secure).catch(
			function(err) {
			    if ("Errors" in err) {
				LoginForm.errors = err.Errors
			    } else {
				if ("message" in err) {
				    LoginForm.errors.push(err.message)
				}
			    }
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
		]),
	    ]),
	    m(NewUserButton),	    	    	    
	])
    }
}


module.exports = LoginForm
