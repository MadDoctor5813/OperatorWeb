function notFound() {
    $('body').text('Not found!');
}

function eventHandler() {
    $('body').on('click', '[data-uw-action]', function(e) {
        e.preventDefault();
        
        var $this = $(this);
        var action = $this.data('uw-action');
        var target = $this.data('uw-target');
        switch (action) {
            case 'sidebar-open':
                $this.addClass('toggled');
                $('#' + target).addClass('toggled');
                
                var selection = document.querySelector('#dynamic-sidebar-style');
                if (!selection) {
                    $('head').append('<style type="text/css" id="dynamic-sidebar-style"></style>');
                    selection = document.querySelector('#dynamic-sidebar-style');
                }
                
                var mainHeight = $(window).height() - $('.page-top').outerHeight(true);
                selection.textContent = '@media (max-width: 991px) { #page-content-wrapper.toggled { height: ' + mainHeight + 'px; min-height: 0} }';
                
                $('#page-content-wrapper').addClass('toggled');
                $('#page-content-wrapper').append('<div data-uw-action="sidebar-close" data-uw-target="' + target + '" class="sidebar-backdrop" onClick=""></div>');
                
                break;
            case 'sidebar-close':
                $('[data-uw-action="sidebar-open"]').removeClass('toggled');
                $('#' + target).removeClass('toggled');
                
                var selection = document.querySelector('#dynamic-sidebar-style');
                if (selection) {
                    selection.textContent = '';
                }
                
                $('#page-content-wrapper').removeClass('toggled');
                $('.sidebar-backdrop').remove();
                
                break;
        }
    });
}

function handleAjaxError(jqXHr, textStatus) {
    var message = '';
 
    switch (textStatus) {
        case 'notmodified':
            message = 'Not Modified';
            break;
        case 'parsererror':
            message = 'Parser Error';
            break;
        case 'timeout':
            message = 'Time Out';
            break;
        default:
            switch (jqXHr.status) {
                case 398: // error
                    if (jqXHr.responseJSON) {
                        message = jqXHr.responseJSON.message;
                    }
                    else {
                        message = '398 Error';
                    }
                    
                    break;
                case 401: // unauthorized
                    if (jqXHr.responseJSON) {
                        message = jqXHr.responseJSON.message;
                    }
                    else {
                        message = '401 Unauthorized';
                    }
                    
                    break;
                case 403: // forbidden
                    if (jqXHr.responseJSON) {
                        message = jqXHr.responseJSON.message;
                    }
                    else {
                        message = '403 Forbidden';
                    }
                    window.location.pathname = '/login';
                    
                    break;
                case 404: // not found
                    if (jqXHr.responseJSON) {
                        message = jqXHr.responseJSON.message;
                    }
                    else {
                        message = '404 Not Found';
                    }
                    
                    break;
                case 500: // internal server error
                    if (jqXHr.responseJSON) {
                        message = jqXHr.responseJSON.message;
                    }
                    else {
                        message = '500 Internal Server Error';
                    }
                    
                    break;
                case 503: // service unavailable
                    if (jqXHr.responseJSON) {
                        message = jqXHr.responseJSON.message;
                    }
                    else {
                        message = '503 Service Unavailable';
                    }
                    
                    break;
                default:
                    message = 'Error';
            }
    }
    
    if (message) {
        console.log('Error: ' + message);
        displayAlertMessage(message);
    }
}

function displayAlertMessage(message) {
    var alertMessageHTML = '\
        <div class="alert alert-warning alert-dismissible fade in" role="alert">\
            <button type="button" class="close default" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>' + message + '\
        </div>';
    
    $('body').append(alertMessageHTML);
}