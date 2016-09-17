function initDataModal(json) {
	var $modal = $('#data-modal');

	$modal.find('.date').text(formatDate(json.date));
	$modal.find('.time').text(formatTime(json.initTime));
	$modal.find('.street').text(json.street);
	$modal.find('.city').text(json.city + ', ' + json.province);
	$modal.find('.postal-code').text(json.postal-code);
	
	if (json.details != '') {
		$modal.find('.details').text(json.details);
	}
	else {
		$modal.find('.details').text('No details have been recorded.');
	}
	$modal.find('.description').text(json.description);
	
	// display image
	
	$modal.find('.response').text(json.response);
	$modal.find('.notes').text(json.notes);
	$modal.find('.level #inlineRadio' + json.level).prop('checked', true);
	$modal.find('.status #optionsRadios' + json.status).prop('checked', true);

	// display map
}

function clearTable() {    
    if ($('body').data('autosave-timer')) {
        clearInterval($('body').data('autosave-timer'));
        $('body').removeData('autosave-timer');
    }

    $('#data-table tbody').empty();
}