var m = require("mithril")

var AppRoot = require("./views/AppRoot")
var LoginForm = require("./views/LoginForm")
var FirstUserForm = require("./views/FirstUserForm")

m.route(document.body, "/", {
    "/": AppRoot,
    "/login": LoginForm,
    "/firstuser": FirstUserForm,
})
