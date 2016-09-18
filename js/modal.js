function initDataModal(json) {
	var $modal = $('#data-modal');

	// date and time
	$modal.find('.category').text(json.category);
	$modal.find('.date').text(formatDate(json.initTime));
	$modal.find('.time').text(formatTime(json.initTime));

	// address
	var size = json.locations.length;
	var cityProvince = '';
	if (json.locations[size - 1].street != '') {
		$modal.find('.street').text(location.street);
	}
	if (json.locations[size - 1].city != '') {
		cityProvince = json.locations[size - 1].city;
	}
	if (json.locations[size - 1].province != '') {
		if (cityProvince == '') { // no city
			cityProvince = json.locations[size - 1].province;
		}
		else { // yes city
			cityProvince += ', ' + json.locations[size - 1].province;
		}	
	}
	$modal.find('.city').text(cityProvince);
	if (json.locations[size - 1].postalCode != '') {
		$modal.find('.postal-code').text(json.locations[size - 1].postalCode);
	}
	
	// user info: details and description
	if (json.details != '') {
		$modal.find('.details').text(json.details);
	}
	else {
		$modal.find('.details').text('No details have been recorded.');
	}
	if (json.description != '') {
		$modal.find('.description').text(json.description);
	}
	else {
		$modal.find('.description').text('No description has been recorded.');
	}
	
	// display image
	if (json.imageName != '') {
    	$modal.find('.image img').attr('src', '/img/' + json.imageName + ".jpg");
    }
    else {
    	$modal.find('.image span').text('No image has been uploaded.');
    }
	
	// admin info: response and notes
	if (json.response != '') {
		$modal.find('.response').text(json.response);
	}
	else {
		$modal.find('.response').text('No response has been recorded.');
	}
	if (json.notes != '') {
		$modal.find('.notes').text(json.notes);
	}
	else {
		$modal.find('.notes').text('No notes have been recorded.');
	}

	// level and status
	$modal.find('.level #inlineRadio' + json.level).prop('checked', true);
	$modal.find('.status #optionsRadios' + json.status).prop('checked', true);

	// display map
	deleteMarkers();
	flightPlanCoordinates = [];
	bounds = new google.maps.LatLngBounds();

	var counter = 1;
	$.each(json.locations, function(i, location) {
		var latLng = new google.maps.LatLng(location.latitude, location.longitude);
        
        bounds.extend(latLng);
        flightPlanCoordinates.push(latLng);

        if (i % 5 == 0) { // for every fifth location
        	if (i == 0) { // first location
        		addMarker(counter.toString(), location, latLng);
        		counter++;
        	}
        	if (i >= 5) { // all other locations
        		if (json.locations[i].latitude == json.locations[i - 5].latitude
        			&& json.locations[i].longitude == json.locations[i - 5].longitude) { // if coordinates are the same as coordinates of previous marker, do not set new marker
        			// console.log(i + ' ' + json.locations[i].latitude + ' ' + json.locations[i - 5].latitude);
        			// console.log(i + ' ' + json.locations[i].longitude + ' ' + json.locations[i - 5].longitude);
        		}
        		else {
        			addMarker(counter.toString(), location, latLng);
        			counter++;
        		}
			}
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
	$modal.find('.image img').attr('src', '');
	$modal.find('.image span').text('');
	
	$modal.find('.response').text('');
	$modal.find('.notes').text('');
	$modal.find('.level input').prop('checked', false);
	$modal.find('.status input').prop('checked', false);

	// clear map
	deleteMarkers();
	flightPath.setMap(null);
	flightPlanCoordinates = [];

	$modal.removeClass('initialized');
	clearInterval(emergencyTimerID);
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

	$('#data-modal').on('dblclick', '.text-edit', function() {
		var $this = $(this);
		var text = $this.text();
		var $textareaHTML = $('<textarea class="form-control"></textarea>').val(text);
		
		$this.text('').append($textareaHTML);
		$this.find('textarea').focus();
	});
	$('#data-modal').on('click', '.text-edit', function(e) {
		e.stopPropagation();
	});
	$('#data-modal').on('blur', '.text-edit', function() {
		var $this = $(this);
		var text = $this.find('textarea').val();

		$this.text(text);
	});

	$('#data-modal').on('click', 'button.save', function() {
		updateEmergency();
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