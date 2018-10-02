var m = require("mithril")
var Leagues = require("../models/Leagues")
var secure = require("../models/secure")

var AddLeagueForm = {
    errors: [],
    view: function(){
	return m("form.add-league-form",
		 {
		     onsubmit: function(e) {
			 e.preventDefault()
			 Leagues.save().then(
			     function(response) {
				 window.location.href = "/league/" + response.ObjID
			     }
			 ).catch(secure).catch(
			     function(err) {
				 AddLeagueForm.errors.push(err)
			     }
			 )
		     }
		},
		 [
		     m("fieldset", [
			 m(".container", [
			     AddLeagueForm.errors.map(
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
				     m("legend", "New League")
				 ])
				 
			     ]),
			     m(".row", [
				 m(".col-3", [m("label", {for: "name"}, "League Name")]),
				 m(".col-9", [m("input", {
				     type: "text",
				     name: "name",
				     placeholder: "league name",
				     oninput: m.withAttr("value", function(value){
					 Leagues.new.Name = value
				     }),
				 })])
			     ]),
			     m(".row", [
				 m(".col-9", []),
				 m(".col-3", [
				     m("input", {
					 type: "submit",
					 value: "Add League"
				     })
				 ])
			     ])
			 ])
		     ]),
		 ])
    }
}

module.exports = AddLeagueForm

