new Vue({
    el: '#app',
    data: {
	adminUserExists: null,
    },
    watch: { adminUserExists: 'handleAdminExists' },
    created: function(){ this.fetchData()},
    methods: {
	fetchData: function () {
	    axios.get('/user/admin/exists.json').then(response => this.adminUserExists = response)
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
