var LoginForm = new Vue({
    el: "#login-form",
    data: {
	errors: [],
	username: null,
	password: null,
    },
    methods: {
	login: function(e) {
	    e.preventDefault()
	    this.errors = []
	    let self = this
	    axios.post("/authenticate", {
		username: this.username,
		password: this.password
	    }).then(
		function(response) {
		    window.location.href = "/app.html"
		}
	    ).catch(
		function(error) {
		    self.errors.push(error.message)
		}
	    )
	}
    }
    
})
