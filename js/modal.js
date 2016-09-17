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
    $modal.find('.image img').attr('src', '/img/' + json.imageName);
	
	$modal.find('.response').text(json.response);
	$modal.find('.notes').text(json.notes);
	$modal.find('.level #inlineRadio' + json.level).prop('checked', true);
	$modal.find('.status #optionsRadios' + json.status).prop('checked', true);

	// display map
	bounds = new google.maps.LatLngBounds();

	$.each(json.locations, function(i, location) {
		var latLng = new google.maps.LatLng(location.latitude, location.longitude);
        
        bounds.extend(latLng);
        flightPlanCoordinates.push(latLng);

        if (i % 5 == 0) {
        	var index = (i / 5 + 1).toString();
        	addMarker(index, location, latLng);
        }
	});

	flightPath = new google.maps.Polyline({
		path: flightPlanCoordinates,
		geodesic: true,
		strokeColor: '#FF0000',
		strokeOpacity: 1.0,
		strokeWeight: 2
    });
    flightPath.setMap(map);

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
	$modal.find('.level input').prop('checked', false);
	$modal.find('.status input').prop('checked', false);

	// clear map
	deleteMarkers();
	flightPath.setMap(null);
	flightPlanCoordinates = [];
}

function emergencyEventHandler() {
	$('#data-modal').on('shown.bs.modal', function() {
		console.log('modal open');
		zoomAndCenter();
	});

	$('#data-modal').on('hidden.bs.modal', function() {
		console.log('modal closed');
		clearModal();
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

	infowindow = new google.maps.InfoWindow();

	map.addListener('click', function() {
		infowindow.close();
    });
}

function zoomAndCenter() {
	google.maps.event.trigger(map, 'resize');
	map.setCenter(bounds.getCenter());
	map.fitBounds(bounds);

	if (map.getZoom() > 16) { // set maximum zoom
		map.setZoom(16);
	}
}

// Adds a marker to the map and push to the array.
function addMarker(index, location, latLng) {
	var marker = new google.maps.Marker({
		position: latLng,
		map: map,
		label: index
	});
	markers.push(marker);

	var infowindowHTML = '\
		<div id="infowindow">\
			<div class="time">' + formatTime(location.time) + '</div>\
			<div class="street">' + location.street + '</div>\
			<div class="city province">' + location.city + ', ' + location.province + '</div>\
			<div class="postal-code">' + location.postalCode + '</div>\
		</div>';

	google.maps.event.addListener(marker, 'click', function() {
        infowindow.close();
        infowindow.setContent(infowindowHTML);
        infowindow.open(map, marker);
    });
}

// Sets the map on all markers in the array.
function setMapOnAll(map) {
	for (var i = 0; i < markers.length; i++) {
		markers[i].setMap(map);
	}
}

// Removes the markers from the map, but keeps them in the array.
function clearMarkers() {
	setMapOnAll(null);
}

// Shows any markers currently in the array.
function showMarkers() {
	setMapOnAll(map);
}

// Deletes all markers in the array by removing references to them.
function deleteMarkers() {
	clearMarkers();
	markers = [];
}