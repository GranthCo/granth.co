

//var granthMainApp = angular.module('granthMainApp', ['ui.bootstrap', 'ngRoute', 'ui.router']);
var granthMainApp = angular.module('granthMainApp', ['ui.bootstrap', 'ngRoute']);


granthMainApp.controller('AccordionDemoCtrl', 
	['$scope', function($scope){

  $scope.oneAtATime = true;

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


granthMainApp.controller('TabsDemoCtrl', 
	['$scope', function($scope) {
	  $scope.tabs = [
	    { title:"Dynamic Title 1", content:"Dynamic content 1" },
	    { title:"Dynamic Title 2", content:"Dynamic content 2", disabled: true }
	  ];


	  $scope.navType = 'pills';
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
	['$rootScope', '$scope', '$location', '$window', '$route', '$log', 
		function($rootScope, $scope, $location, $window, $route, $log){

		  	$scope.totalItems = 1430;
		  	$scope.itemsPerPage = 1;
		  	$scope.numPages = 1430;

		  	$scope.maxSize = 10;
		  
		  	$rootScope.setPage = function (pageNo) {
		  		// Get the page number
		    	$scope.currentPage = pageNo;
		    	// But also build up the current string
		    	var languages = $window.location.search;
		    	$log.info("setPage: " + languages);

		    	var queryString = languages.substring(1,languages.length).split("&");
		    	$log.info("qs" + queryString);
		    	for(var i = 0; i < queryString.length; i++) {
		    		var pairs = queryString[i].split("=");
		    		if(pairs[0] == "lang") {
		    			
		    			$rootScope.languageSelection = pairs[1].split(",");
		    			$log.info("languageSelection " + $rootScope.languageSelection);
		    		}
		    	}

		    	if($rootScope.languageSelection === undefined || $rootScope.languageSelection.length ==0 ) {
		    		$rootScope.languageSelection = ['english'];
		    	}


		  	};

		  	$rootScope.pageChanged = function(page) {
				var currentPath = $window.location.search; 
				$log.info("path " + $window.location.pathname);
				var base = $window.location.pathname.substring(0, $window.location.pathname.lastIndexOf("/"));

				$window.location.href = base + "/" + page + currentPath ;

		  	};
			
			$rootScope.setLanguages = function(languages) {
				var languageString = languages.join(",");
				$log.info("setLanguagge " + $window.location.pathname + "?lang=" + languageString);

				$window.location.href = $window.location.pathname + "?lang=" + languageString;

		  	};


		  	/*$rootScope.setLanguages = function(languages) {
		  		$log.info(languages + ":" $window.location);

		  	};*/
		  	$rootScope.shabadChanged = function(shabad) {
		  		$rootScope.pageChanged('h' + shabad);
		  	};
		  	$rootScope.lineChanged = function(line) {
		  		$rootScope.pageChanged(line + '-' + line);
		  	};

	}]);

// Please note that $modalInstance represents a modal window (instance) dependency.
// It is not the same as the $modal service used above.

granthMainApp.controller('LanguageInstanceController', 
	['$rootScope', '$modalInstance', 'items', '$log', 
	function($rootScope, $modalInstance, items, $log) {


  		$rootScope.toggleSelection = function (languageName) {
  			$log.info("Toggle selection called with " + languageName);
		    var idx = $rootScope.languageSelection.indexOf(languageName);

		    // is currently selected
		    if (idx > -1) {
		      $rootScope.languageSelection.splice(idx, 1);
		    }

		    // is newly selected
		    else {
		      $rootScope.languageSelection.push(languageName);
		    }
		  };


  		$rootScope.ok = function () {
  			$log.info('oksy');
    		$modalInstance.close($rootScope.languageSelection);
    		$rootScope.setLanguages($rootScope.languageSelection);
  		};

  		$rootScope.cancel = function () {
  			$log.info('close ');
    		$modalInstance.dismiss('cancel');
  		};
	}]);

granthMainApp.controller('LanguageController', 
	['$rootScope', '$modal', '$log', 
	function($rootScope, $modal, $log) {

  		$rootScope.languages = ["afrikaans",
								"albanian",
								"arabic",
								"belarusian",
								"bulgarian",
								"catalan",
								"chineseSimplified",
								"chineseTraditional",
								"croatian",
								"czech",
								"danish",
								"dutch",
								"english",
								"estonian",
								"filipino",
								"finnish",
								"french",
								"galician",
								"german",
								"greek",
								"haitian",
								"hebrew",
								"hindi",
								"hungarian",
								"icelandic",
								"indonesian",
								"irish",
								"italian",
								"japanese",
								"korean",
								"latvian",
								"lithuanian",
								"macedonian",
								"malay",
								"maltese",
								"norwegian",
								"persian",
								"polish",
								"portuguese",
								"romanian",
								"russian",
								"serbian",
								"slovak",
								"slovenian",
								"spanish",
								"swahili",
								"swedish",
								"thai",
								"turkish",
								"ukrainian",
								"vietnamese",
								"welsh",
								"yiddish"];
 
	  	$rootScope.open = function () {

	    	var modalInstance = $modal.open({
	      		templateUrl: 'myModalContent.html',
	      		controller: 'LanguageInstanceController',
	      		resolve: {
	        		items: function () {
	          			return $rootScope.items;
	        		}	
	      		}
	    	});

	    	$log.info('Got the model instance');

	    	modalInstance.result.then(function (selectedItem) {
	      		$rootScope.selectedLanguages = selectedItem;
	    	}, function () {
	      		$log.info('Modal dismissed at: ' + new Date());
	    	});
	  	};
	}]);




