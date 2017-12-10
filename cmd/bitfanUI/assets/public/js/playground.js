var websocketOUT;
var websocketIN;


$(document).ready(function() {

    UUID = guid();



    $('#bitfan-playground-form').on('submit', function(e) { //use on if jQuery 1.7+
        e.preventDefault(); //prevent form from submitting
        play();
        return false;
    });

    $("#bitfan-playground-form button[name='sendEvent']").on('click', function(e) { //use on if jQuery 1.7+
        websocketIN.send($("#bitfan-playground-form textarea[name='event']").val());
    });
    
    $("#bitfan-playground-form select[name='event_type']").on('change', function(e) { //use on if jQuery 1.7+
        play();
    });


    $(window).on('beforeunload', function() {
        var dataObject = {
            'event': "",
            'event_type': "",
            'filter': "",
            'uuid': "playground-" + UUID,
        };

        $.ajax({
            url: window.location.href,
            type: 'DELETE',
            contentType: "application/json; charset=utf-8",
            dataType: 'json',
            data: JSON.stringify(dataObject),
            success: function(result) {

            }
        });
    });

});



function play() {
    var dataObject = {
        'event': $("#bitfan-playground-form textarea[name='event']").val(),
        'event_type': $("#bitfan-playground-form select[name='event_type']").val(),
        'filter': $("#bitfan-playground-form textarea[name='filter']").val(),
        'uuid': "playground-" + UUID,
    };

    $.ajax({
        type: 'PUT',
        contentType: "application/json; charset=utf-8",
        data: JSON.stringify(dataObject),
        dataType: 'json',
        url: window.location.href,
        beforeSend: function() {},
        success: function(settings) {
            console.log(settings)
            console.log("success");
            playErrorReset();

            httpin_url = "http://" + settings.apiHost + settings.httpin;
            $("#bitfan-playground-form #bitfan-http-input-url").show();
            var httpInTmpl = $.templates("#httpin-template");
            $('#bitfan-http-input-url').html(httpInTmpl.render({
                ev: {"url":httpin_url}, 
            }));


            new_uri = "ws://" + settings.apiHost + settings.wsout;
            websocketOUT = new WebSocket(new_uri);
            websocketOUT.onopen = function(event) {
                console.log("Connection is established!");
            }
            websocketOUT.onmessage = function(event) {
                // var Data = JSON.parse(event.data);
                console.log(event.data);
                // $("#bitfan-playground-form textarea[name='output']").val(event.data);
                $("#bitfan-playground-form div[name='output']").html(syntaxHighlight(event.data));
                $("#bitfan-playground-form div[name='output']").addClass("success");
            };
            websocketOUT.onerror = function(event) {
                // notie.alert({ type: 'warning', stay: false, text: 'Problem due to some Error' });
                console.log(event);
            };
            websocketOUT.onclose = function(event) {

            };

            new_uri = "ws://" + settings.apiHost + settings.wsin;
            websocketIN = new WebSocket(new_uri);
            websocketIN.onopen = function(event) {
                websocketIN.send(dataObject.event);
            };
        },
        error: function(output) {
            console.log(output);
            playError(output.responseText);
            return false;
        }
    });

}






function playErrorReset() {
    $("#playground-error").text("");
    $("#bitfan-playground-form button[name='sendEvent']").show();
    $("#playground-error").removeClass("error");
    $("#bitfan-playground-form div[name='output']").removeClass("error");
    $("#bitfan-playground-form div[name='output']").removeClass("success");
}

function playError(errorTxt) {
    $("#playground-error").text(errorTxt);
    $("#bitfan-playground-form button[name='sendEvent']").hide();
    $("#bitfan-playground-form #bitfan-http-input-url").hide();
    $("#playground-error").addClass("error");
    $("#bitfan-playground-form div[name='output']").addClass("error");
    $("#bitfan-playground-form div[name='output']").removeClass("success");
}


function guid() {
    function s4() {
        return Math.floor((1 + Math.random()) * 0x10000)
            .toString(16)
            .substring(1);
    }
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
        s4() + '-' + s4() + s4() + s4();
}