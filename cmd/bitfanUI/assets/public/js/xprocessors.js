$(document).ready(function() {
    $.views.settings.delimiters("[[", "]]");
    console.log("xprocessors !!!");

    // Handle Enter on forms input
    // Handle Form Submit and Send Ajax
    $("#bitfan-xprocessor-form,#bitfan-xprocessor-form-content").submit(function() {
        console.log($("#bitfan-xprocessor-form,#bitfan-xprocessor-form-content").serialize());
        $.ajax({
            type: "POST",
            data: $("#bitfan-xprocessor-form,#bitfan-xprocessor-form-content").serialize(),
            url: "",
            success: function(data) {
            	notie.alert({ type: "success", text: 'Saved !' })
                console.log(data) ;
            },
            error: function(output) {
            	console.log(output) ;
				notie.alert({ type: 'error', text: 'There was an error while processing!<br>' + output.responseText });
            }
        });
		

        // $.ajax({
        //     type: 'post',
        //     data: 'json',
        //     url: e.attr("href"),
        //     beforeSend: function() {
        //         e.addClass("disabled")
        //         e.siblings().addClass("disabled")
        //     },
        //     success: function(output) {
        //         console.log(output);
        //         var notoggle = e.attr('notoggle');
        //         if (typeof notoggle == typeof undefined || notoggle == false) {
        //             e.toggleClass("hidden")
        //             e.siblings().toggleClass("hidden")
        //         }
        //         e.removeClass("disabled")
        //         e.siblings().removeClass("disabled")


        //         // var tmpl = $.templates("#flash-template"); // Get compiled template
        //         // $('section.flash').append(tmpl.render({message: output}));   

        //         notie.alert({ type: "success", text: 'Success !' })

        //         // $(template).fadeOut( "slow" );
        //     },
        //     error: function(output) {
        //         console.log(output);
        //         notie.alert({ type: 'error', text: 'There was an error while processing!<br>' + output.responseText });
        //         e.removeClass("disabled")
        //         e.siblings().removeClass("disabled")
        //         return false;
        //     }
        // });
        // console.log($("#bitfan-xprocessor-form,#bitfan-xprocessor-form-content").serializeArray());
        return false;
    });
});