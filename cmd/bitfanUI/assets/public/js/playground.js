var websocketOUT;
var websocketIN;
var UUID;
var keyPressTimeoutId = 0;

function PgNewEditor(name,firstLineNumber,syntaxName,themeName, playWithKeyPress) {
    var pgEditor = ace.edit("section-"+name+"-content");
    pgEditor.setAutoScrollEditorIntoView(true);
    pgEditor.setTheme("ace/theme/"+themeName);
    pgEditor.getSession().setMode("ace/mode/"+syntaxName);
    var textarea = $('textarea[name="section-'+name+'"]').hide();
    pgEditor.getSession().setValue(textarea.val());

    pgEditor.getSession().on('change', function() {
        textarea.val(pgEditor.getSession().getValue());
        pgEditor.getSession().clearAnnotations();
        if (playWithKeyPress == true) {
            clearTimeout(keyPressTimeoutId); // doesn't matter if it's 0
            keyPressTimeoutId = setTimeout(play, 500);
        }
    });

    pgEditor.setOption("firstLineNumber", firstLineNumber)
    return pgEditor
}

$(document).ready(function() {
    UUID = guid();


// ######### EDITORS
// OUTPUT CONFIGURATION
    var editorOutput = PgNewEditor("output-configuration", 8, "logstash","monokai",true)
// FILTER CONFIGURATION
    var editorFilter = PgNewEditor("filter-configuration", 6,"logstash","monokai",true)
    editorFilter.getSession().on('change', function() {
        editorOutput.setOption("firstLineNumber", editorFilter.getSession().getLength() + 6 + 1);
    });
// INPUT CONFIGURATION
    var editorInput = PgNewEditor("input-configuration", 2,"logstash","monokai",true)
    editorInput.getSession().on('change', function() {
        editorFilter.setOption("firstLineNumber", 1+editorInput.getSession().getLength()+2)
        editorOutput.setOption("firstLineNumber", editorFilter.getOption("firstLineNumber")+editorFilter.getSession().getLength()+1)
    });
// INPUT RAW
    var editorInputRaw = PgNewEditor("input-raw", 1,"json","eclipse",false)
    editorInputRaw.getSession().on('change', function() {
        editorFilter.setOption("firstLineNumber", 6)
        editorOutput.setOption("firstLineNumber", editorFilter.getOption("firstLineNumber") + editorFilter.getSession().getLength() + 1)
    });

    editorInputRaw.commands.addCommand({
        name: 'sendEvent',
        bindKey: {
            win: 'Ctrl-S',
            mac: 'Command-S',
            sender: 'editor|cli'
            },
        exec: function(env, args, request) {
        // websocketIN.send($("#section-input-raw").val());
        $("#bitfan-playground-form button[name='sendEvent']").click();
        }
    });
// ######### End editors






    // On form submit, play...ground    
    // $('#bitfan-playground-form').on('submit', function(e) { //use on if jQuery 1.7+
    //     e.preventDefault(); //prevent form from submitting
    //     play();
    //     return false;
    // });

    // #########
    // INPUT RAW EVENT
    // #########
    // When user click on "Send again" button Then send the raw event to pipeline using its websocket input
    $("#bitfan-playground-form button[name='sendEvent']").on('click', function(e) { //use on if jQuery 1.7+
        websocketIN.send($("#section-input-raw").val());
    });
    // When codec selection change Then play playground
    $("#section-input-codec").on('change', function(e) { //use on if jQuery 1.7+
        play();
    });


    // #########
    // GLOBAL
    // #########
    // When leaving page Then delete currently running pipeline
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
    // When user toggle any tab Then play playground 
    $('a[data-toggle="tab"]').on('shown.bs.tab', function(e) {
        play();
    });

    // #########
    // LOGS
    // #########
    // When page loaded Then connect to the logs websocke
    var websocketLOGS = new WebSocket("ws://" + baseApiHost + "/api/v2/logs");
    websocketLOGS.onopen = function(event) {
        console.log("LOGS : Connection established! ");
    }

    // When a log message comes Then display it
    var logmessagetmpl = $.templates("#logmessage-template");
    websocketLOGS.onmessage = function(event) {
        var LogMessage = JSON.parse(event.data);
        if (LogMessage.Data.pipeline_uuid == "playground-" + UUID) {
            $('#logs').append(logmessagetmpl.render({
                ev: LogMessage,
                timeString: moment(LogMessage.Time).format('LTS'),
                eventHTML: syntaxHighlightIfEvent(LogMessage.Data.event),
            }));
            $('#logs').scrollTop($('#logs')[0].scrollHeight);
        }
    };
    // When an error occurs on logs websocket Then alert user
    websocketLOGS.onerror = function(event) {
        notie.alert({ type: 'warning', stay: false, text: 'Problem due to some Error' });
    };
    // When websocket connexion closes Then alert user
    websocketLOGS.onclose = function(event) {
        notie.alert({ type: 'warning', stay: false, text: 'Connection Closed' });
    };

    // $('#frmChat').on("submit",function(event){
    //  event.preventDefault();
    //  $('#chat-user').attr("type","hidden");      
    //  var messageJSON = {
    //      chat_user: $('#chat-user').val(),
    //      chat_message: $('#chat-message').val()
    //  };
    //  websocket.send(JSON.stringify(messageJSON));
    // });

});

function syntaxHighlightIfEvent(data) {
    if (data) {
        return syntaxHighlight(data)
    }
}

function play() {
    var input_mode = $('#pan-input .nav-tabs .active').attr("bitfan-section-type");
    var input_value = $("#section-input-" + input_mode).val()
    var input_codec = $("#section-input-codec").val()
    var filter_mode = $('#pan-filter .nav-tabs .active').attr("bitfan-section-type");
    var filter_value = $("#section-filter-" + filter_mode).val()
    var output_mode = $('#pan-output .nav-tabs .active').attr("bitfan-section-type");
    var output_value = $("#section-output-" + output_mode).val()

    var dataObject = {
        'uuid': "playground-" + UUID,
        'input_value': input_value,
        'input_mode': input_mode,
        'input_codec': input_codec,
        'filter_value': filter_value,
        'filter_mode': filter_mode,
        'output_value': output_value,
        'output_mode': output_mode,
    };

    // console.table(dataObject);

    $.ajax({
        type: 'PUT',
        contentType: "application/json; charset=utf-8",
        data: JSON.stringify(dataObject),
        dataType: 'json',
        url: window.location.href,
        beforeSend: function() {},
        success: function(settings) {
            // console.log(settings)
            // console.log("success");
            playErrorReset();

            if (settings.wsout != "") {
                new_uri = "ws://" + settings.apiHost + settings.wsout;
                websocketOUT = new WebSocket(new_uri);
                websocketOUT.onopen = function(event) {
                    // console.log("Connection is established!");
                }
                websocketOUT.onmessage = function(event) {
                    // console.log(event.data);
                    $("#bitfan-playground-form div[name='output']").html(syntaxHighlight(event.data));
                    $("#bitfan-playground-form div[name='output']").addClass("success");
                };
                websocketOUT.onerror = function(event) {
                    // notie.alert({ type: 'warning', stay: false, text: 'Problem due to some Error' });
                    console.log("error on wsout");
                    console.log(event);
                };
                websocketOUT.onclose = function(event) {

                };
            }

            if (settings.wsin != "") {
                new_uri = "ws://" + settings.apiHost + settings.wsin;
                websocketIN = new WebSocket(new_uri);
                websocketIN.onopen = function(event) {
                    websocketIN.send(dataObject.input_value);
                };
            }
        },
        error: function(output) {
            console.log("error playing");

            playError(output.responseJSON);
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


// Utils

// getModeByFileExtension returns the mode path for a given file extension
function getModeByFileExtension(path) {
    var modelist = ace.require("ace/ext/modelist");
    return modelist.getModeForPath(path).mode;
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