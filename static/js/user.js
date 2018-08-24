var UserForm = new Vue({
    el: "#user-form",
    data: {
	errors: [],
	forename: null,
	surname: null,
	username: null,
	password: null,
    },
    methods: {
	createUser: function(isAdmin) {
	    let self = this
	    this.errors = []
	    axios.post('/user/new.json', {
	    	forename: this.forename,
	    	surname: this.surname,
	    	username: this.username,
	    	password: this.password,
	    	isadmin: isAdmin,
	    }).then(
	    	function(response) {
	    	    self.errors = response.data.Errors
	    	}
	    ).catch(
	    	function(error) {
	    	    self.errors.push(error.message)
	    	}
	    )

	},
	createAdminUser: function(e) {
	    e.preventDefault()
	    return this.createUser(true)
	},
	createNormalUser: function(e) {
	    e.preventDefault()
	    return this.createUser(false)
	}
    }
	
})
