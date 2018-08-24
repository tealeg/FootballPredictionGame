new Vue({
    el: '#app',
    data: {
	adminUserExists: null,
    },
    watch: { adminUserExists: 'handleAdminExists' },
    created: function(){ this.fetchData()},
    methods: {
	fetchData: function () {
	    let self = this
	    axios.get('/user/admin/exists.json').then(
		function(response) {
		    self.adminUserExists = response.data
		}
	    ).catch(
		function(error) {
		    console.log(error)
		}
	    )
	},
	handleAdminExists: function() {
	    if (this.adminUserExists == true) {
		window.location.href = "/login.html"
	    } else {
		window.location.href = "/firstuser.html"
	    }
	    
	}
    }
})
