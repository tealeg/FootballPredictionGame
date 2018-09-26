var m = require("mithril")

var AppRoot = require("./views/AppRoot")
var LoginForm = require("./views/LoginForm")
var FirstUserForm = require("./views/FirstUserForm")
var NewUserForm = require("./views/NewUserForm")
var AddLeagueForm = require("./views/AddLeagueForm")

m.route(document.body, "/", {
    "/": AppRoot,
    "/firstuser": FirstUserForm,
    "/leagues/add": AddLeagueForm,
    "/login": LoginForm,
    "/user/new": NewUserForm,
})
