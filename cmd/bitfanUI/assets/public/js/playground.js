var websocketOUT;
var websocketIN;
var UUID;
var keyPressTimeoutId = 0;
var editorInput;
var editorInputRaw;
var editorFilter;
var editorOutput;
var editorOutputRaw;
var autoStartPlayGround = false;




function PgNewEditor(name, firstLineNumber, syntaxName, themeName, playWithKeyPress) {
    var pgEditor = ace.edit("section-" + name + "-content");
    pgEditor.setAutoScrollEditorIntoView(true);
    pgEditor.$blockScrolling = Infinity
    pgEditor.setTheme("ace/theme/" + themeName);
    pgEditor.getSession().setMode("ace/mode/" + syntaxName);
    var textarea = $('textarea[name="section-' + name + '"]').hide();
    pgEditor.getSession().setValue(textarea.val());

    pgEditor.getSession().on('change', function() {
        textarea.val(pgEditor.getSession().getValue());
        pgEditor.getSession().clearAnnotations();
        if (autoStartPlayGround) {
            if (playWithKeyPress == true) {
                clearTimeout(keyPressTimeoutId); // doesn't matter if it's 0
                keyPressTimeoutId = setTimeout(play, 500);
            }
        }
    });

    pgEditor.commands.addCommand({
        name: 'Help',
        bindKey: {
            win: 'Ctrl-Shift-P',
            mac: 'Command-Shift-P',
            sender: 'editor|cli'
        },
        exec: function(env, args, request) {
            // console.log(pgEditor.container.id);
            // console.log("cursor at row " + pgEditor.getCursorPosition().row + ", column : " + pgEditor.getCursorPosition().column);
            // context = editorID, column, row, selection
            var context = {
                aceEditor: pgEditor
            };

            zone = "playground"
            switch (name) {
                case "input-configuration":
                    zone = "doc input";
                    break;
                case "filter-configuration":
                    zone = "doc filter";
                    break;
                case "output-configuration":
                    zone = "doc output";
                    break;
            }

            bitbar.show(zone + " ", context)
        }
    });

    pgEditor.commands.addCommand({
        name: 'sendEvent',
        bindKey: {
            win: 'Ctrl-S',
            mac: 'Command-S',
            sender: 'editor|cli'
        },
        exec: function(env, args, request) {
            // console.log(env)
            // console.log(args)
            // console.log(request)
            play();
        }
    });

    pgEditor.setOption("firstLineNumber", firstLineNumber)
    return pgEditor
}

function bitfanProcessorSelectOption(c, item) {
    c.aceEditor.session.insert(c.aceEditor.selection.getCursor(), item.value)
    c.aceEditor.focus();
}

function bitfanProcessorListptions(c) {
    $.ajax({
        type: 'GET',
        dataType: "json",
        processData: false,
        url: 'http://' + baseApiHost + '/api/v2/docs/processors/' + c.selected[1],
        success: function(processor_doc) {
            var items = []
            let proc = processor_doc[c.selected[1]]
            // console.log(proc.Options)
            var optionsLen = proc.Options.Options.length;
            for (var i = 0; i < optionsLen; i++) {
                let option = proc.Options.Options[i]
                // console.log(option);

                //Ignore common options
                if (option.Type == "processors.CommonOptions") {
                    continue
                }

                // When no alias set, use Name
                var labelStr = option.Alias
                if (labelStr == "") {
                    labelStr = (option.Name.replace(/\.?([A-Z])/g, function(x, y) { return "_" + y.toLowerCase() }).replace(/^_/, ""));
                }

                var labelRequired = ""
                if (option.Required == true) {
                    labelRequired = " - (required)"
                }

                value = bitfanOptionText(option, true)

                items.push({
                    onSelect: bitfanProcessorSelectOption,
                    id: labelStr,
                    label: labelStr + labelRequired,
                    value: value,
                    help: option.Doc,
                })

            }
            bitbar.showWith(items, "", c)
        },
        error: function(output) {
            // console.log("error getting processor documentations");
            return false;
        }
    });
}

function bitfanOptionText(option, withDoc) {
    // When no alias set, use Name
    var labelStr = option.Alias
    if (labelStr == "") {
        labelStr = (option.Name.replace(/\.?([A-Z])/g, function(x, y) { return "_" + y.toLowerCase() }).replace(/^_/, ""));
    }

    var value = labelStr + " => ";

    switch (option.Type) {
        case "hash":
            if (option.DefaultValue != null) {
                value += option.DefaultValue;
            } else {
                value += "{}";
            }
            break;
        case "array":
            if (option.DefaultValue != null) {
                value += option.DefaultValue;
            } else {
                value += "[]";
            }
            break;
        case "string":
            if (option.DefaultValue != null && option.DefaultValue != "") {
                value += option.DefaultValue;
            } else {
                value += '""';
            }
            break;
        case "int":
            if (option.DefaultValue != null) {
                value += option.DefaultValue;
            } else {
                value += "123";
            }
            break;
        case "bool":
            if (option.DefaultValue != null) {
                value += option.DefaultValue;
            } else {
                value += "false";
            }
            break;
        case "time.Duration":
            if (option.DefaultValue != null) {
                value += option.DefaultValue;
            } else {
                value += "";
            }
            break;
        case "location":
            if (option.DefaultValue != null) {
                value += option.DefaultValue;
            } else {
                value += '"" # a go template, as string or as path ';
            }
            break;
        case "interval":
            if (option.DefaultValue != null) {
                value += option.DefaultValue;
            } else {
                value += '"" # cron spec like * * * * *';
            }
            break;
    }

    finalTxt = ""

    if (withDoc) {
        if (option.Doc != "") {
            finalTxt += "  # " + option.Doc.replace(/[\n\r]/g, "\n  # ") + "\n"
        }
        if (option.Required == true) {
            finalTxt += "  # @Required option !\n"
        }
        if (option.ExampleLS != "") {
            finalTxt += "  # @Example : " + option.ExampleLS.replace(/[\n\r]/g, "\n  # ") + "\n"
        }
        finalTxt += "  " + value + "\n\n"
    } else {
        finalTxt += "  " + value
        if (option.Required == true) {
            finalTxt += " # required"
        }
        finalTxt += "\n"
    }

    return finalTxt;
}

function bitfanProcessorInsertTemplate(c) {
    $.ajax({
        type: 'GET',
        dataType: "json",
        processData: false,
        url: 'http://' + baseApiHost + '/api/v2/docs/processors/' + c.selected[1],
        success: function(processor_doc) {
            var items = []
            let key = c.selected[1]
            var finalTxt = ""
            let proc = processor_doc[key]

            if (key.startsWith("input") || key.startsWith("output")) {
                // key = key.replace("input_","")
                key = key.replace("output_", "").replace("input_", "")
            }
            finalTxt = key + " {\n"

            withDoc = false
            if (c.selected[0] == "insert_full") {
                withDoc = true
            }

            var optionsLen = proc.Options.Options.length;
            if (optionsLen == 0 && withDoc == true) {
                const regex = /\n/gm;
                const subst = `  # `;
                // The substituted value will be contained in the result variable
                finalTxt += "  # " + proc.Doc.replace(regex, subst) + "\n"
            }
            for (var i = 0; i < optionsLen; i++) {
                let option = proc.Options.Options[i]

                //Ignore common options
                if (option.Type == "processors.CommonOptions") {
                    continue
                }

                finalTxt += bitfanOptionText(option, withDoc)
            }
            finalTxt = finalTxt + "}\n"

            c.aceEditor.session.insert(c.aceEditor.selection.getCursor(), finalTxt)
            c.aceEditor.focus();

        },
        error: function(output) {
            // console.log("error getting processor documentations");
            return false;
        }
    });
}

function bitfanProcessorMenu(c) {
    keys = [
        // { onSelect: bitfanProcessorShowDoc, id: 'show_doc', label: "read the doc of " + c.selected[0] },
        { onSelect: bitfanProcessorListptions, id: 'list_options', label: "list options" },
        { onSelect: bitfanProcessorInsertTemplate, id: 'insert_full', label: "insert full blueprint" },
        { onSelect: bitfanProcessorInsertTemplate, id: 'insert_min', label: "insert only options" },
    ];
    bitbar.showWith(keys, "", c)
}



$(document).ready(function() {

    var urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has('with') && urlParams.get('with') != null) {
        contentUUID = urlParams.get('with');
        $.ajax({
            type: 'GET',
            dataType: "json",
            processData: false,
            url: 'http://' + baseApiHost + '/api/v2/assets/' + contentUUID,
            success: function(asset) {
                contentValueString = Base64.decode(asset.Value);
                try {
                    const myRegexp = /input[^{]*{([\S\s.]*)}[^}]*filter[^{]*{([\S\s.]*)}[^}]*output[^{]*{([\S\s.]*)}[^}]*/gm;
                    var match = myRegexp.exec(contentValueString);
                    if (match != null && match.length == 4) {
                        editorInput.getSession().setValue(match[1]);
                        editorFilter.getSession().setValue(match[2]);
                        editorOutput.getSession().setValue(match[3]);
                    }
                } catch (e) {

                }

            },
            error: function(output) {
                // console.log("error getting base asset");
                return false;
            }
        });
    }



    UUID = guid();

    // ######## BITBAR load with processors docs
    $.ajax({
        type: 'GET',
        dataType: "json",
        processData: false,
        url: 'http://' + baseApiHost + '/api/v2/docs/processors',
        success: function(processors_docs) {
            for (var key in processors_docs) {
                var labelStr = "doc filter " + key
                if (key.startsWith("input") || key.startsWith("output")) {
                    labelStr = "doc " + key.replace("_", " ");
                }

                switch (processors_docs[key].Behavior) {
                    case "producer":
                        labelStr = "doc input " + key;
                        break;
                    case "transformer":
                        labelStr = "doc filter " + key;
                        break;
                    case "consumer":
                        labelStr = "doc output " + key;
                        break;
                }

                bitbar.items.push({
                    onSelect: bitfanProcessorMenu,
                    id: key,
                    label: labelStr,
                    help: processors_docs[key].Doc,
                })

                if (processors_docs[key].Behavior == "producer") {
                    bitbar.items.push({
                        onSelect: bitfanProcessorMenu,
                        id: key,
                        label: "doc filter " + key,
                        help: processors_docs[key].Doc,
                    })
                }


            }

        },
        error: function(output) {
            // console.log("error getting processor documentations");
            return false;
        }
    });



    bitbar.items.push({
        onSelect: function(c, i) { play() },
        id: "playground-play",
        label: "playground: play / replay",
        help: "Stop currently running playground's pipeline, Then Start it again",
    }, {
        onSelect: function(c, i) { stop() },
        id: "playground-stop",
        label: "playground: stop",
        help: "Stop currently running playground's pipeline",
    })

    // Init Bitbar
    // Load Data from api/docs
    // Create labels
    // bitbar.initWith(keys)

    // ######### EDITORS 
    // OUTPUT CONFIGURATION
    editorOutput = PgNewEditor("output-configuration", 8, "logstash", "monokai", true)

    editorOutputRaw = PgNewEditor("output-raw", 1, "json", "eclipse", false)

    // FILTER CONFIGURATION
    editorFilter = PgNewEditor("filter-configuration", 6, "logstash", "monokai", true)
    editorFilter.getSession().on('change', function() {
        editorOutput.setOption("firstLineNumber", editorFilter.getSession().getLength() + 6 + 4);
    });

    // INPUT CONFIGURATION
    editorInput = PgNewEditor("input-configuration", 2, "logstash", "monokai", true)
    editorInput.getSession().on('change', function() {
        editorFilter.setOption("firstLineNumber", 1 + editorInput.getSession().getLength() + 2)
        editorOutput.setOption("firstLineNumber", editorFilter.getOption("firstLineNumber") + editorFilter.getSession().getLength() + 1)
    });

    // INPUT RAW
    editorInputRaw = PgNewEditor("input-raw", 1, "json", "eclipse", false)
    editorInputRaw.getSession().on('change', function() {
        editorFilter.setOption("firstLineNumber", 6)
        editorOutput.setOption("firstLineNumber", editorFilter.getOption("firstLineNumber") + editorFilter.getSession().getLength() + 1)
    });


    var dragging = false;

    $('#pg-container .dragbar').mousedown(function(e) {
        e.preventDefault();
        window.dragging = true;


        var wrapper = $(e.target).parent()
        var my_editor = wrapper.find(".editor_wrap div");
        var top_offset = my_editor.offset().top;
        // Set editor opacity to 0 to make transparent so our wrapper div shows
        my_editor.css('opacity', 0);
        // handle mouse movement

        $(document).mousemove(function(e) {
            var actualY = e.pageY;
            var eheight = actualY - top_offset;

            // for each editor_wrap
            wrapper.parent().find(".editor_wrap").css('height', eheight);

            // Only one
            wrapper.find(".dragbar").css('opacity', 0.15);
        });
    });

    $(document).mouseup(function(e) {
        if (window.dragging) {
            window.dragging = false;
            $(document).unbind('mousemove');
            // For each Wrapper
            var wrapper = $(e.target).parent()
            var my_editor = wrapper.find(".editor_wrap div");
            // Set dragbar opacity back to 1
            $(e.target).css('opacity', 1);
            my_editor.css('opacity', 1);
            // Trigger resize() one each ace editor 
            wrapper.parent().find(".editor_wrap > div").each(function(index) {
                ace.edit($(this).attr("id")).resize();
            })
        }
    });



    // #########
    // INPUT RAW EVENT
    // #########
    // When user click on "Send again" button Then send the raw event to pipeline using its websocket input
    $("#bitfan-playground-form button[name='sendEvent']").on('click', function(e) { //use on if jQuery 1.7+
        if (websocketIN != null) {
            websocketIN.send($("#section-input-raw").val());
        }
    });
    // When codec selection change Then play playground
    $("#section-input-codec").on('change', function(e) { //use on if jQuery 1.7+
        if (autoStartPlayGround) {
            play()
        }
    });

    // #########
    // GLOBAL
    // #########
    // When leaving page Then delete currently running pipeline
    $(window).on('beforeunload', function() {
        stop();
    });
    // When user toggle any tab Then play playground 
    $('a[data-toggle="tab"]').on('shown.bs.tab', function(e) {
        if (autoStartPlayGround) {
            play();
        }
    });

    // #########
    // LOGS
    // #########
    // When page loaded Then connect to the logs websocke
    var websocketLOGS = new WebSocket("ws://" + baseApiHost + "/api/v2/logs");
    websocketLOGS.onopen = function(event) {
        // console.log("LOGS : Connection established! ");
    }

    // When a log message comes Then display it
    var logmessagetmpl = $.templates("#logmessage-template");
    websocketLOGS.onmessage = function(event) {
        var LogMessage = JSON.parse(event.data);
        if (
            LogMessage.Data.pipeline_uuid == "playground-" + UUID ||
            LogMessage.Message.indexOf(UUID) > 0
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

    // #########
    // PLAY ACTION BUTTONS
    // #########
    $('#play-actions button').click(function() {
        if ($(this).hasClass('disabled')) {
            return false;
        }

        if ($(this).attr("id") == "playground-stop") {
            stop();
        }

        if ($(this).attr("id") == "playground-play" || $(this).attr("id") == "playground-replay") {
            play();
        }
    });

    $('#play-options-autostart').change(function() {
        autoStartPlayGround = $(this).is(":checked")
    });


});



function stop() {
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

    $("#bitfan-playground-form button[name='sendEvent']").hide();
    $('#playground-play').show();
    $('#playground-stop').hide();
    $('#playground-replay').hide();
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
                    editorOutputRaw.getSession().setValue(event.data);
                    $("#section-output-raw-content").addClass("success");
                };
                websocketOUT.onerror = function(event) {
                    // notie.alert({ type: 'warning', stay: false, text: 'Problem due to some Error' });
                    // console.log("error on wsout");
                    // console.log(event);
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
            // console.log("error playing");

            playError(output.responseJSON);
            return false;
        }
    });

}






function playErrorReset() {

    $('#playground-play').hide();
    $('#playground-stop').show();
    $('#playground-replay').show();

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
        Data: {}
    }
    $('#logs').append(logmessagetmpl.render({
        ev: LogMessage,
        timeString: moment(Date.now()).format('LTS'),
        // eventHTML: syntaxHighlightIfEvent(LogMessage.Data.event),
    }));
    $('#logs').scrollTop($('#logs')[0].scrollHeight);



}

function syntaxHighlightIfEvent(data) {
    if (data) {
        return syntaxHighlight(data)
    }
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