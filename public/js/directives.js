'use strict';

/* Directives */

var directives = angular.module('app.directives', []);

directives.directive('focusMe', function($timeout, $parse) {
  return {
    link: function(scope, element, attrs) {
      var model = $parse(attrs.focusMe);
      scope.$watch(model, function(value) {
        if(value === true) {
          $timeout(function() {
            element[0].focus();
          });
        }
      });
      // on blur event: set attribute value to 'false'
      element.bind('blur', function() {
         scope.$apply(model.assign(scope, false));
      });
    }
  };
});