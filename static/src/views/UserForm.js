var m = require("mithril")
var User = require("../models/User")

var UserForm = {
    user: User,
    errors: [],
    view: function () {
	return m("fieldset", [
	    m(".container", [
		UserForm.errors.map(
		    function(err) {
			m(".row", [
			    m(".col-12", [
				m(".error", err),
			    ])
			])
		    }
		),
		m(".row", [
		    m(".col-12", [
			m("legend", "No adminstration account has been created. Please create one now.")
		    ])
		]),
		m(".row", [
		    m(".col-3", [m("label", {for: "forename"}, "Forename")]),
		    m(".col-9", [m("input", {
			type: "text",
			name: "forename",
			oninput: m.withAttr("value", function(value){
			    User.current.forename = value
			}),
			value: User.current.forename,
		    })]),
		]),
		m(".row", [
		    m(".col-3", m("label", {for: "surname"}, "Surname")),
		    m(".col-9", m("input", {
			type: "text",
			name: "surname",
			oninput: m.withAttr("value", function(value) {
			    User.current.surname = value
			}),
			value: User.current.surname,
		    })),
		]),
		m(".row", [
		    m(".col-3", m("label", {for: "username"}, "Username")),
		    m(".col-9", m("input", {
			type: "text",
			name: "username",
			oninput: m.withAttr("value", function(value) {
			    User.current.username = value
			}),
			value: User.current.username,
		    })),
		]),
		m(".row", [
		    m(".col-3", m("label", {for: "password"}, "Password")),
		    m(".col-9", m("input", {
			type: "password",
			name: "password",
			oninput: m.withAttr("value", function(value){
			    User.current.password = value
			}),
			value: User.current.password,
		    })),
		]),
		m(".row", [
		    m(".col-9", []),
		    m(".col-3", [
			m("input", {type: "submit", value: "Create Account"})
		    ])
		])
	    ])
	])
    }
}


module.exports = UserForm
