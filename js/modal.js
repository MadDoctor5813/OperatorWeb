function initDataModal(json) {
	var $modal = $('#data-modal');

	$modal.find('.category').text(json.category);
	$modal.find('.date').text(formatDate(json.initTime));
	$modal.find('.time').text(formatTime(json.initTime));
	$modal.find('.street').text(json.street);
	$modal.find('.city').text(json.city + ', ' + json.province);
	$modal.find('.postal-code').text(json.postalCode);
	$modal.find('.details').text(json.details);
	$modal.find('.description').text(json.description);
	
	// display image
	
	$modal.find('.response').text(json.response);
	$modal.find('.notes').text(json.notes);
	$modal.find('.level #inlineRadio' + json.level).prop('checked', true);
	$modal.find('.status #optionsRadios' + json.status).prop('checked', true);

	// display map

	$modal.modal('show');
}

function clearModal() {
    var $modal = $('#data-modal');

    $modal.find('.category').text('');
    $modal.find('.date').text('');
	$modal.find('.time').text('');
	$modal.find('.street').text('');
	$modal.find('.city').text('');
	$modal.find('.postal-code').text('');
	$modal.find('.details').text('');
	$modal.find('.description').text('');
	
	// clear image
	
	$modal.find('.response').text('');
	$modal.find('.notes').text('');
	$modal.find('.level #inlineRadio' + json.level).prop('checked', true); // **** EDIT ****
	$modal.find('.status #optionsRadios' + json.status).prop('checked', true);

	// clear map
}

function emergencyEventHandler() {
	$('#data-modal').on('shown.bs.modal', function() {
		console.log('modal open');
		
	    var currentCenter = map.getCenter();  // get current center before resizing
		google.maps.event.trigger(map, 'resize');
		map.setCenter(currentCenter); // re-set previous center
	});

	$('#data-modal').on('hidden.bs.modal', function() {
		console.log('modal closed');
	});
}

var map;
var markers = [];
var bounds;
var infowindow;
var flightPath;
var flightPlanCoordinates = [];

function initMap() {
	map = new google.maps.Map(document.getElementById('map'), {
		center: {lat: 43.471867, lng: -80.5415358},
		zoom: 16
	});
}