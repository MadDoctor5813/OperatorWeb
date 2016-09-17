function initDataModal(json) {
	var $modal = $('#data-modal');

	$modal.find('.date').text(formatDate(json.date));
	$modal.find('.time').text(formatTime(json.initTime));
	$modal.find('.street').text(json.street);
	$modal.find('.city').text(json.city + ', ' + json.province);
	$modal.find('.postal-code').text(json.postal-code);
	$modal.find('.details').text(json.details);
	$modal.find('.description').text(json.description);
	
	// display image
	
	$modal.find('.response').text(json.response);
	$modal.find('.notes').text(json.notes);
	$modal.find('.level #inlineRadio' + json.level).prop('checked', true);
	$modal.find('.status #optionsRadios' + json.status).prop('checked', true);

	// display map
}

function clearModal() {    
    if ($('body').data('autosave-timer')) {
        clearInterval($('body').data('autosave-timer'));
        $('body').removeData('autosave-timer');
    }

    var $modal = $('#data-modal');

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