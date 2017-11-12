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


  var changeApiURLAjax = function (e) {
    btn = $('#bitfan-location button[href]')
    input = $('#bitfan-location input')
    sendData = {url:input.val()}

    $.ajax({
      type: 'PUT',
      contentType: "application/json; charset=utf-8",
      data: JSON.stringify(sendData),
      dataType: 'json',
      url: btn.attr("href"),
      beforeSend: function(){
        
      },
      success: function (settings) {
          input.val(settings.url)
          window.location = "/pipelines"
      },
      error: function (output) {
        console.log(output);
          notie.alert({ type: 'error', text: 'There was an error while processing!<br>'+ output.responseText}) ;
          return false;
      }
    });

    return false;
  };

  $('#bitfan-location input').bind("enterKey",function(e){
    changeApiURLAjax();
  });
  $('#bitfan-location input').keyup(function(e){
    if(e.keyCode == 13)
    {
        $(this).trigger("enterKey");
    }
  });
  $('#bitfan-location button[href]').click(function() {
    changeApiURLAjax();
  });


});