$(document).ready(function() {
    $.views.settings.delimiters("[[", "]]");
    console.log("env !!!");

    function addVarItem(varObj) {
        var tmpl = $.templates("#env-item"); // Get compiled template
        $('#env-items').append(tmpl.render({ "env": varObj }));

        // Bind Delete
        $("#env-" + varObj.uuid + " button[name='delete']").click(function(e) {
            e.preventDefault();
            console.log($(this).attr("href"))
            var uuid = $(this).attr("uuid")
            $.ajax({
                type: 'delete',
                url: 'http://' + baseApiHost + '/api/v2/env/' + uuid,
                success: function(output) {
                    console.log("#env-" + uuid);
                    $("#env-" + uuid).remove()
                    // notie.alert({ type: "success", text: 'Success !' })
                },
                error: function(output) {
                    console.log(output);
                    var errorMessage = JSON.parse(output);
                    console.log(errorMessage)
                    notie.alert({ type: 'error', text: 'There was an error while processing!<br>' + errorMessage.error });
                    return false;
                }
            });
            return false
        });
    }


    $.ajax({
        type: 'get',
        // url: window.location.href,
        // data: JSON.stringify(sendData),
        url: 'http://' + baseApiHost + '/api/v2/env',
        success: function(envVars) {
            console.log(envVars);
            $.each(envVars, function(i, obj) {
                console.log(obj.uuid)
                addVarItem(obj)
            });
        },
        error: function(output, textStatus, errorThrown) {
            console.log(textStatus)
                console.log(errorThrown)

                 if (output.readyState == 0) {
                // HTTP error (can be checked by XMLHttpRequest.status and XMLHttpRequest.statusText)
                    notie.alert({ type: 'warning', stay:false, text: 'Connection Closed' }) ;
                    return
                }

            console.log(output);
            output.responseText
            notie.alert({ type: 'error', text: 'There was an error while processing!<br>' + output.responseText });
            return false;
        }
    });



    // Bind form to ajax request
    // when sucess 
    //    add a new env-item
    //    clear form
    // when ko alert
    $("#add_env").submit(function(e) {
        // values = $(this).serializeArray()
        // console.log(values);
        console.log('http://' + baseApiHost + '/api/v2/env');

        var name = $(this).find('input[name="name"]').val()
        var value = $(this).find('input[name="value"]').val()
        var secret = $(this).find('input[name="secret"]').is(":checked")
        var sendData = { "name": name, "value": value, "secret": secret };

        $.ajax({
            type: 'post',
            data: 'json',
            // url: window.location.href,
            data: JSON.stringify(sendData),
            url: 'http://' + baseApiHost + '/api/v2/env',
            beforeSend: function() {
                $(e.target).attr("disabled", true)
                $(e.target).children().attr("disabled", true)
            },
            success: function(envVar) {
                console.log(envVar);
                $(e.target).trigger("reset"); // reset form        
                $(e.target).removeAttr("disabled")
                $(e.target).children().removeAttr("disabled")
                addVarItem(envVar)
            },
            error: function(output, textStatus, errorThrown) {
                console.log(output);
                var errorObj = output.responseJSON
                var errorMessage = output.responseText
                if (errorObj != null) {
                    errorMessage = errorObj.error
                }
                notie.alert({ type: 'error', text: errorMessage });
                $(e.target).removeAttr("disabled")
                $(e.target).children().removeAttr("disabled")
                return false;
            }
        });

        return false;
    });




})