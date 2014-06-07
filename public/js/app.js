'use strict';


// Declare app level module which depends on filters, and services
var app = angular.module('app', [
	'ngCookies',
	'ui.router',
	'restangular',
	'ui.bootstrap',
	'app.filters',
	'app.services',
	'app.directives',
	'app.controllers',
	'facebook'
]);

app.config(function($stateProvider, $urlRouterProvider, FacebookProvider) {
	// Initializing facebook provider
	FacebookProvider.init(env.appId);

	// For any unmatched url, redirect to /
  	$urlRouterProvider.otherwise("/");
	$stateProvider.state('home', {
			url: "/",
			templateUrl: "view/home.html",
			controller: "HomeController"
		});
	
	app.$stateProvider = $stateProvider;
});

app.run(function ($rootScope, $log, $http, $cookies, Facebook, Restangular) {
	// Alert's system
	$rootScope.alerts = [];
	$rootScope.addAlert = function(msg, type) {
		// Type could be danger, success or warning
		if (!type) {
			type = 'warning'
		}
		$rootScope.alerts.push({type: type, msg: msg});
	};
	$rootScope.closeAlert = function(index) {
		$rootScope.alerts.splice(index, 1);
	};

	// Helper function to divide array in a 2d array
	// Return an array with max of elements in each row
	$rootScope.getRows = function(array, maxElem) {
		var rows = [];
		var i, j, temparray, chunk = maxElem;
		for (i=0,j=array.length; i<j; i+=chunk) {
			temparray = array.slice(i, i+chunk);
			rows.push(temparray);
		}
		return rows;
	};
	// Helper function to divide array in a 2d array
	// Divide the inner array in x pieces
	$rootScope.getCols = function(array, pieces) {
		var rows = [];
		var i, j, temparray, chunk = array.length/pieces;
		for (i=0,j=array.length; i<j; i+=chunk) {
			temparray = array.slice(i, i+chunk);
			rows.push(temparray);
		}
		return rows;
	};

	// Cofiguring our Rest client Restangular

	// Calls to api will be done in /api path
	Restangular.setBaseUrl('/api');


	// Setting the Api error response handler
	Restangular.setErrorInterceptor(function (response) {
		if (response.status === 419) {
			// The users credentials has expired
			// If he had credentials, let's remove them
			document.execCommand("ClearAuthenticationCache");
			delete $cookies.credentials;
			Restangular.setDefaultHeaders({Authorization: ""});
			$http.defaults.headers.common.Authorization = "";
			// Alerting the user
			$rootScope.addAlert("Voce nao esta mais logado no sistema, logue-se novamente.");
		} else if (response.status === 401) {
			// The current user is not logged
			$rootScope.addAlert("Voce nao esta logado no sistema.", "danger");
		} else if (response.status === 404) {
			$rootScope.addAlert("O recurso requisitado nao existe.", "danger");
		} else {
			$rootScope.addAlert("Ops ocorreu um erro.", "danger");
			$log.log("Ops ocorreu um erro.....");
			$log.log(response);
		}
		// JUST FOR DEBBUNGING, LETS SHOW THE ERROR MESSAGE SENT BY SERVER
		if (response.data.error) {
			$rootScope.addAlert(response.data.error.message);
		};
		// return false; // stop the promise chain
	});

	// Adding the function to logout the current user
	$rootScope.logout = function () {
	    document.execCommand("ClearAuthenticationCache");
	    delete $cookies.credentials;
	    Restangular.setDefaultHeaders({Authorization: ""});
	    $http.defaults.headers.common.Authorization = "";
	    delete $rootScope.me;
		$rootScope.addAlert("Voce se deslogou do sistema.");
    }

	// Setting the current user session
	if ($cookies.credentials) {
		// Set the Authorization header saved in this session
		$http.defaults.headers.common.Authorization = $cookies.credentials;
	    Restangular.setDefaultRequestParams({Authorization: $cookies.credentials});
	    // Getting the actual user
	    //Restangular.one("me").get().then(function (user) {
	    //	$rootScope.me = user; // Prevent to create empty object using .$object
	    //});
	}


  	// Hidden collapse navbar on click in a link
  	// REMOVE COLAPSE MENU !!!
	$('.nav a').on('click', function(){
		$(".navbar-toggle").click();
		//$log.log("Clicked in navbar, close navbar");
	});




      




	// Define user empty data :/
	$rootScope.user = {};

	// Defining user logged status
	$rootScope.logged = false;

	// Defining if Facebook is ready
	$rootScope.facebookReady = false;

	/**
	* Watch for Facebook to be ready.
	* There's also the event that could be used
	*/
	$rootScope.$watch(
		function() {
			return Facebook.isReady();
		},
		function(newVal) {
			if (newVal){
				Facebook.getLoginStatus(function(response) {
					if (response.status == 'connected') {
						$log.log("AUTO-LOGIN OK");
						$log.log(response);
						// We need to wait for user to facebook became ready
						$rootScope.connect(response);
					}
					else {
						// User isn't logged, facebook is ready
						$rootScope.facebookReady = true;
					}
				});
			}
		}
	);


	/**
	* Login
	*/
	$rootScope.login = function() {
		Facebook.login(function(response) {
			if (response.status == 'connected') {
				$log.log("LOGIN OK");
				$log.log(response);
				$rootScope.connect(response);
			}
		}, {
			scope: 'email,user_friends,read_stream', 
			return_scopes: true
		});
	};

	/**
	* Get me in the facebook servers
	*/
	$rootScope.me = function() {
		Facebook.api('/me', {fields: 'id,email,name,gender,link,timezone,verified,picture,permissions'}, function(response) {
			/**
			* Using $rootScope.$apply since this happens outside angular framework.
			*/
			console.log("LOGGED");
			console.log(response);
			$rootScope.$apply(function() {
				$rootScope.user = response;
				$rootScope.logged = true;
			});
			// We've got the user, facebook is ready
			$rootScope.facebookReady = true;
		});
	};

	/**
	* Logout
	*/
	$rootScope.logout = function() {
		Facebook.logout(function() {
			$rootScope.$apply(function() {
				$rootScope.user   = {};
				$rootScope.logged = false;  
			});

			// Clearing the access default params
		    Restangular.setDefaultRequestParams({accessToken: ""});
		});

	    
	}

	/**
	* Connect the user with the system
	*/
	$rootScope.connect = function(response) {
		
		// Setting the access default params
	    Restangular.setDefaultRequestParams({accessToken: response.authResponse.accessToken});

		// Sending this connected user to the turma server
		Restangular.service('me').post({
			accessToken: response.authResponse.accessToken,
			userID: response.authResponse.userID,
		}).then(function(response) {
			$log.log("User sent to server, look the response");
			$log.log(response);
		});

		// Get user in the facebook servers
		$rootScope.me();

	};

});

app.controller('AppController', function($scope, $log, $timeout, Restangular) {
	//$scope.timeFormat = "d-MM-yy 'às' HH:mm";
	$scope.showSearch = false;
	$scope.isSearching = false;

});

