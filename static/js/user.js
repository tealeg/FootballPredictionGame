var app = new Vue({
    el: "#user-form",
    data: {
	errors: [],
	forename: null,
	surname: null,
	username: null,
	password: null,
	isadmin: null,
    },
    methods: {
	createUser: function(isAdmin) {
	    this.isadmin = isAdmin
	    return function (e) {
		let self = this
		this.errors = []
		axios.post('/user/new.json', {
	    	    forename: this.forename,
	    	    surname: this.surname,
	    	    username: this.username,
	    	    password: this.password,
	    	    isadmin: this.isadmin,
		}).then(
	    	    function(response) {
	    		self.errors = response.data.Errors
	    	    }
		).catch(
	    	    function(error) {
	    		self.errors.push(error.message)
	    	    }
		    
		)
		e.preventDefault()

	    }
	}
	
    }
})
