function notFound() {
    $('body').text('Not found!');
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
                    window.location.pathname = '/sign-in';
                    
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

function formatDate(date) {
    var year = date.substr(0, 4);
    var month = date.substr(4, 2);
    var day = date.substr(6, 2);

    switch (month) {
        case '01':
            month = 'January';
            break;
        case '02':
            month = 'February';
            break;
        case '03':
            month = 'March';
            break;
        case '04':
            month = 'April';
            break;
        case '05':
            month = 'May';
            break;
        case '06':
            month = 'June';
            break;
        case '07':
            month = 'July';
            break;
        case '08':
            month = 'August';
            break;
        case '09':
            month = 'September';
            break;
        case '10':
            month = 'October';
            break;
        case '11':
            month = 'November';
            break;
        case '12':
            month = 'December';
            break;
        default:
            console.log('Format date case not matched.');
    }

    return month + " " + day + ", " + year;
}

function formatTime(time) {
    var hour = time.substr(8, 2);
    var minute = time.substr(10, 2);
    var second = time.substr(12, 2);
    var am = 'AM';

    if (hour < 10) {
        hour = hour.substr(1, 1);
    }
    else if (hour == 12) {
        am = 'PM';
    }
    else if (hour > 12) {
        hour = hour - 12;
        am = 'PM';
    }

    return hour + ':' + minute + ':' + second + ' ' + am;
}