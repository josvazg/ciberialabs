function MainCtrl($scope) {

	var addresses = $scope.addresses = [];

	$scope.newBtcAddress = '';
	$scope.newBtcLabel = '';

	$scope.addressCount = $scope.addresses.length;

	$scope.$watch('addresses', function() {
		$scope.addressCount = $scope.addresses.length;
		console.log($scope.addressCount);
	}, true);

	$scope.addBtcAddress = function() {
		var newBtcAddress = $scope.newBtcAddress.trim();
		if (!newBtcAddress.length) {
			return;
		}

		var newBtcLabel = $scope.newBtcLabel.trim();

		addresses.push({
			name : newBtcAddress,
			label : newBtcLabel
		});

		$scope.newBtcAddress = '';
		$scope.newBtcLabel = '';

		$('#addBtcAddressModal').modal('hide');

	};

	$scope.removeBtcAddress = function(address) {
		addresses.splice(addresses.indexOf(address), 1);

	}

}
