var m = require("mithril")

var AppRoot = require("./views/AppRoot")
var LoginForm = require("./views/LoginForm")
var FirstUserForm = require("./views/FirstUserForm")
var NewUserForm = require("./views/NewUserForm")
var LeagueView = require("./views/LeagueView")
var LeaguesView = require("./views/LeaguesView")
var AddLeagueForm = require("./views/AddLeagueForm")

m.route(document.body, "/", {
    "/": AppRoot,
    "/leagues": LeaguesView,
    "/firstuser": FirstUserForm,
    "/leagues/add": AddLeagueForm,
    "/league/:id": LeagueView,
    "/login": LoginForm,
    "/user/new": NewUserForm,
})
