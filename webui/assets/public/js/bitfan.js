$(document).ready(function () {
  $.views.settings.delimiters("[[", "]]");

  var pipelineActionAjax = function (e) {
  $.ajax({
      type: 'get',
      dataType: 'json',
      url: e.attr("href"),
      beforeSend: function(){
        e.addClass("disabled")
        e.siblings().addClass("disabled")
      },
      success: function (output) {
          console.log(output);
          var notoggle = e.attr('notoggle');
          if (typeof notoggle == typeof undefined || notoggle == false) {
            e.toggleClass("hidden")
            e.siblings().toggleClass("hidden")
          }
          e.removeClass("disabled")
          e.siblings().removeClass("disabled")

          
          // var tmpl = $.templates("#flash-template"); // Get compiled template
          // $('section.flash').append(tmpl.render({message: output}));   

          notie.alert({ type:"success", text: 'Success !' })       

          // $(template).fadeOut( "slow" );
      },
      error: function (output) {
        console.log(output);
          notie.alert({ type: 'error', text: 'There was an error while processing!<br>'+ output.responseText}) ;
          e.removeClass("disabled")
          e.siblings().removeClass("disabled")
          return false;
      }
  });
  return false;
  };

  $('#pipeline-actions button[href]').click(function() {
    if ($(this).hasClass('disabled')) {
        return false;
    }    

    pipelineActionAjax($(this));
  });
});