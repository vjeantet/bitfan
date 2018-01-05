var websocketOUT;
var websocketIN;
var UUID;
var keyPressTimeoutId = 0;
var editorInput ;
var editorInputRaw ;
var editorFilter ;
var editorOutput;
var editorOutputRaw ;



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

    pgEditor.commands.addCommand({
        name: 'Help',
        bindKey: {
            win: 'Ctrl-H',
            mac: 'Command-H',
            sender: 'editor|cli'
            },
        exec: function(env, args, request) {
            console.log(pgEditor.container.id) ;
            console.log("cursor at row "+ pgEditor.getCursorPosition().row + ", column : "+pgEditor.getCursorPosition().column);
            $('#exampleModal').on('shown.bs.modal', function () {
                $('#bitbar-input').focus();
            }) ; 
            $('#exampleModal').modal({
              keyboard: true
            });
            $(document).on("keydown", "#bitbar-input", function(event) { 
                if(event.which==38 || event.which==40){
                    console.log("UP DOWN") ;
                    event.preventDefault();
                }
                if(event.which == 13) {
                    console.log("ENTER") ;
                }
            });
        }
    });

    pgEditor.setOption("firstLineNumber", firstLineNumber)
    return pgEditor
}

$(document).ready(function() {
    UUID = guid();

    // ######### EDITORS 
    // OUTPUT CONFIGURATION
    editorOutput = PgNewEditor("output-configuration", 8, "logstash","monokai",true)

    editorOutputRaw = PgNewEditor("output-raw", 1, "json","eclipse",false)
    
    // FILTER CONFIGURATION
    editorFilter = PgNewEditor("filter-configuration", 6,"logstash","monokai",true)
    editorFilter.getSession().on('change', function() {
        editorOutput.setOption("firstLineNumber", editorFilter.getSession().getLength() + 6 + 1);
    });



    // Autocompletion - autocomplete processor names and their respectives options
    // TODO:  to autocomplete reliably, api needs :
    //   * all configuration before the cursor position
    //   * API server side will calculate  
    //     * what kind of token
    //       * is currently written
    //       * OR is to be written
    //     * the list of previous token (processor name, option name)
    //   * with this token list, bitfan can propose 
    //      * some words.
    //      * help, description, documentation.
    // * OR
    //   * section type
    //   * processor context name (null, mutate, grok, etc....)
    //   * option context name
    //   * 
    // var langTools = ace.require("ace/ext/language_tools");
    // editorFilter.setOptions({enableBasicAutocompletion: true});
    // // uses http://rhymebrain.com/api.html
    // var bitfanCompleter = {
    //     getCompletions: function(editor, session, pos, prefix, callback) {
    //         if (prefix.length === 0) { callback(null, []); return }
    //         console.log(pos,prefix ); 
    //         $.getJSON(
    //             "http://rhymebrain.com/talk?function=getRhymes&word=" + prefix,
    //             function(wordList) {
    //                 // wordList like [{"word":"flow","freq":24,"score":300,"flags":"bc","syllables":"1"}]
    //                 callback(null, wordList.map(function(ea) {
    //                     return {name: ea.word, value: ea.word, score: ea.score, meta: "rhyme"}
    //                 }));
    //             })
    //     }
    // }
    // langTools.addCompleter(rhymeCompleter);

    // INPUT CONFIGURATION
    editorInput = PgNewEditor("input-configuration", 2,"logstash","monokai",true)
    editorInput.getSession().on('change', function() {
        editorFilter.setOption("firstLineNumber", 1+editorInput.getSession().getLength()+2)
        editorOutput.setOption("firstLineNumber", editorFilter.getOption("firstLineNumber")+editorFilter.getSession().getLength()+1)
    });
    
    // INPUT RAW
    editorInputRaw = PgNewEditor("input-raw", 1,"json","eclipse",false)
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
            $("#bitfan-playground-form button[name='sendEvent']").click();
        }
    });
    // editorFilter.commands.addCommand({
    //     name: 'Help',
    //     bindKey: {
    //         win: 'Ctrl-H',
    //         mac: 'Command-H',
    //         sender: 'editor|cli'
    //         },
    //     exec: function(env, args, request) {
    //         alert("GO") ;
    //     }
    // });
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
        if (websocketIN != null){
            websocketIN.send($("#section-input-raw").val());
        } 
    });
    // When codec selection change Then play playground
    $("#section-input-codec").on('change', function(e) { //use on if jQuery 1.7+
        play()
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
        if (
            LogMessage.Data.pipeline_uuid == "playground-" + UUID ||
            LogMessage.Message.indexOf(UUID)>0
            ) {
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


    // Shortcuts
    // define a handler
    // var delta = 500;
    // var lastKeypressTime = 0;
    // function KeyHandler(event)
    // {
    //    if ( event.key == 'g' )
    //    {
    //       var thisKeypressTime = new Date();
    //       if ( thisKeypressTime - lastKeypressTime <= delta )
    //       {
    //         console.log("GO !");
    //         // optional - if we'd rather not detect a triple-press
    //         // as a second double-press, reset the timestamp
    //         thisKeypressTime = 0;
    //       }
    //       lastKeypressTime = thisKeypressTime;
    //    }
    // }
    // // register the handler 
    // document.addEventListener('keyup', KeyHandler, false);

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
                    editorOutputRaw.getSession().setValue(event.data) ;
                    $("#section-output-raw-content").addClass("success");
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
    $("#bitfan-playground-form button[name='sendEvent']").show();
    $("#section-output-raw-content").removeClass("error");
    $("#section-output-raw-content").removeClass("success");
}

function playError(errorTxt) {
    $("#bitfan-playground-form button[name='sendEvent']").hide();
    $("#bitfan-playground-form #bitfan-http-input-url").hide();
    $("#section-output-raw-content").addClass("error");
    $("#section-output-raw-content").removeClass("success");

        var logmessagetmpl = $.templates("#logmessage-template");
        LogMessage = {
            Message: errorTxt,
            Level: 2,
            Data : {} 
        }
            $('#logs').append(logmessagetmpl.render({
                ev: LogMessage,
                timeString: moment(Date.now()).format('LTS'),
                // eventHTML: syntaxHighlightIfEvent(LogMessage.Data.event),
            }));
            $('#logs').scrollTop($('#logs')[0].scrollHeight);



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