'use strict';

/* Filters */

var filters = angular.module('app.filters', []);

filters.filter('userpic', function() {
	return function(user) {
		if (!_.isEmpty(user)) {
			return user.picture.data.url;
		}
	};
});


filters.filter('strLimit', function() {
	return function(str, limit) {
		if (str.length > limit) {
			return str.substr(0, limit-3) + " ..."
		};
		return str
	};
});
