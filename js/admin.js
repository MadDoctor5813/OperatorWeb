function initDataTable(json) {
	$('#data-table').DataTable({
		'dataSrc': '',
		'destroy': true,
		'info': false,
		'order': [[5, 'asc']],
		'ordering': true,
		'paging': false,
		'processing': false,
		'searching': false,
		'serverSide': false,
		'stateSave': true,

		'ajax': function(data, callback, settings) {
			var prevData = $('#data-table').data('data');

			clearContainer();
			callback(json);

			if ($('body').hasClass('initialized')) { // same page
				if (JSON.stringify(prevData) == JSON.stringify(data)) { // same data
					console.log('same page, same data');
				}
				else { // new data
					console.log('same page, new data');

					$('#data-table').data('data', data);
				}
			}
			else { // new page
				if (jQuery.isEmptyObject(json)) { // no data
					console.log('new page, no data');
				}
				else { // yes data
					console.log('new page, yes data');
				}


				$('body').addClass('initialized');
				$('#data-table').data('data', data);
			}
		},

		'columnDefs': [
            {
                'data': 'category',
                'targets': 0,
                'title': 'Category',
                'orderable': true
            },
            {
                'data': 'street',
                'targets': 1,
                'title': 'Street',
                'orderable': true
            },
            {
                'data': 'city',
                'targets': 2,
                'title': 'City',
                'orderable': true
            },
            {
                'data': 'province',
                'targets': 3,
                'title': 'Province',
                'orderable': true
            },
            {
                'data': 'postalCode',
                'targets': 4,
                'title': 'Postal Code',
                'orderable': true
            },
            {
                'data': 'initTime',
                'targets': 5,
                'title': 'Time',
                'orderable': true
            },
            {
                'data': 'level',
                'targets': 6,
                'title': 'Level',
                'orderable': true
            }
        ],

        'drawCallback': function(settings) {
            var api = this.api();

            api.rows().every(function() {
                var data = this.data();
                $(this.node()).attr('id', data.emergencyID);
            });
        }
	});
}