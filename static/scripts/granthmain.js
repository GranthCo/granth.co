

//var granthMainApp = angular.module('granthMainApp', ['ui.bootstrap', 'ngRoute', 'ui.router']);
var granthMainApp = angular.module('granthMainApp', ['ui.bootstrap', 'ngRoute']);


granthMainApp.config(['$routeProvider', '$locationProvider', 
	function($routeProvider, $locationProvider) {
		
		$routeProvider.when('/', {
			templateUrl: 'templates/home.html',
			controller: 'HomeController'
		}).when('/:page', {
			templateUrl: 'templates/page.html',
			controller: 'ReplyController'
		}).when('/:languages/:page', {
			templateUrl: 'templates/page.html',
			controller: 'ReplyController'
		});
		
		$locationProvider.html5Mode(false);
}]);


granthMainApp.controller('AlertController', 
	['$scope', 

	function($scope) {
	  	//$scope.alerts = [
	  	//  { type: 'danger', msg: 'Error quering database' },
	  	//  { type: 'success', msg: 'Well done! You successfully read this important alert message.' }
	  	//];

	  	$scope.addAlert = function() {
	    	$scope.alerts.push({msg: "Another alert!"});
	  	};

	  	$scope.addDatabaseAlert = function() {
	  		$scope.alerts.push({msg: "Error quering databse!"});
	  	}

	  	$scope.closeAlert = function(index) {
	    	$scope.alerts.splice(index, 1);
	  	};	
	}]);





granthMainApp.controller('HomeController', 
	['$rootScope', '$scope', '$location', function($rootScope, $scope, $location){

		mixpanel.track("Main Page");
		$scope.raagData = [{raag:'Jap Ji Sahib', page:'1'},
					{raag:'Rehraas Sahib', page:'8'},
					{raag:'Kirtan Sohila', page:'12'},
					{raag:'Siri Raag', page:'14'},
					{raag:'Raag Majh', page:'94'},
					{raag:'Raag Gauri', page:'151'},
					{raag:'Raag Aasaa', page:'347'},
					{raag:'Raag Gujari', page:'489'},
					{raag:'Raag Devgandhari', page:'527'},
					{raag:'Raag Bihagara', page:'537'},
					{raag:'Raag Vadhans', page:'557'},
					{raag:'Raag Sorath', page:'595'},
					{raag:'Raag Dhanasari', page:'660'},
					{raag:'Raag Jaitsiri', page:'696'},
					{raag:'Raag Todi', page:'711'},
					{raag:'Raag Bairari', page:'719'},
					{raag:'Raag Tilang', page:'721'},
					{raag:'Raag Suhi', page:'728'},
					{raag:'Raag Bilaval', page:'795'},
					{raag:'Raag Gond', page:'859'},
					{raag:'Raag Ramkali', page:'876'},
					{raag:'Raag Nat Narayan', page:'975'},
					{raag:'Raag Mali Gaura', page:'984'},
					{raag:'Raag Maru', page:'889'},
					{raag:'Raag Tukhari', page:'1107'},
					{raag:'Raag Kedara', page:'1118'},
					{raag:'Raag Bhairao', page:'1125'},
					{raag:'Raag Basant', page:'1168'},
					{raag:'Raag Sarang', page:'1197'},
					{raag:'Raag Malaar', page:'1254'},
					{raag:'Raag Kanara', page:'1294'},
					{raag:'Raag Kalian', page:'1319'},
					{raag:'Raag Parbhati', page:'1327'},
					{raag:'Raag Jaijavanti', page:'1352'},
					{raag:'Raag Mala', page:'1429'}];
		
		  
	}]);

granthMainApp.controller('PaginationController2', 
	['$rootScope', '$scope', '$location', function($rootScope, $scope, $location){

		  	$scope.totalItems = 1430;
		  	$scope.itemsPerPage = 1;
		  	$scope.numPages = 1430;

		  	$scope.maxSize = 10;
		  
		  	$rootScope.setPage = function (pageNo) {
		    	$scope.currentPage = pageNo;
		  	};

		  	$scope.pageChanged = function(page) {
				var currentPath = $location.path();
				var split = currentPath.substring(0, currentPath.lastIndexOf("/"));
				$location.path(split + "/" + page);   
		  	};
		  	$rootScope.pageChanged = $scope.pageChanged;
		  	$rootScope.shabadChanged = function(shabad) {
		  		$scope.pageChanged('h' + shabad);
		  	};
		  	$rootScope.lineChanged = function(line) {
		  		$scope.pageChanged(line + '-' + line);
		  	};

	}]);


granthMainApp.controller('ReplyController', 
	['$scope', '$http', '$routeParams', '$rootScope', 
	function ($scope, $http, $routeParams, $rootScope) {

		var languages = $routeParams.languages;
		if(languages === undefined) {
			languages = "english";
		}
		var page = $routeParams.page;
		if(page === undefined) {
			//redirect to a default page?
			page = 1;

		} 

		$http({method: 'GET', url: '/rest/' + page + "?lang=" + languages}).
		  	success(function(data, status, headers, config) {
		    	// this callback will be called asynchronously
		    	// when the response is available
		    	$scope.reply = data;
		    	$rootScope.setPage($scope.reply.Lines[0].Page);    
		  	}).
		  	error(function(data, status, headers, config) {
		    	// called asynchronously if an error occurs
		    	// or server returns response with an error status.
		    	alert("some error " + status +  '/rest/' + page + "?lang=" + languages);
		    	$scope.AlertController.addDatabaseAlert();
		  	});     
	}]);




granthMainApp.controller('DropdownController', 
	['$scope', function($scope) {
	  
	  $scope.items = [
	    "The first choice!",
	    "And another choice for you.",
	    "but wait! A third!"];
	}]);





