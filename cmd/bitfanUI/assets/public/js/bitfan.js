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


function syntaxHighlight(json) {
  if (typeof json != 'string') {
        json = JSON.stringify(json, undefined, '\t');
    }

    json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
        var cls = 'number';
        if (/^"/.test(match)) {
            if (/:$/.test(match)) {
                cls = 'key';
            } else {
                cls = 'string';
            }
        } else if (/true|false/.test(match)) {
            cls = 'boolean';
        } else if (/null/.test(match)) {
            cls = 'null';
        }
        return '<span class="' + cls + '">' + match + '</span>';
    });
}

// Base64 Encode/Decode
var Base64 = {
    characters: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=" ,

    encode: function( string )
    {
        var characters = Base64.characters;
        var result     = '';

        var i = 0;
        do {
            var a = string.charCodeAt(i++);
            var b = string.charCodeAt(i++);
            var c = string.charCodeAt(i++);

            a = a ? a : 0;
            b = b ? b : 0;
            c = c ? c : 0;

            var b1 = ( a >> 2 ) & 0x3F;
            var b2 = ( ( a & 0x3 ) << 4 ) | ( ( b >> 4 ) & 0xF );
            var b3 = ( ( b & 0xF ) << 2 ) | ( ( c >> 6 ) & 0x3 );
            var b4 = c & 0x3F;

            if( ! b ) {
                b3 = b4 = 64;
            } else if( ! c ) {
                b4 = 64;
            }

            result += Base64.characters.charAt( b1 ) + Base64.characters.charAt( b2 ) + Base64.characters.charAt( b3 ) + Base64.characters.charAt( b4 );

        } while ( i < string.length );

        return result;
    } ,

    decode: function( string )
    {
        var characters = Base64.characters;
        var result     = '';

        var i = 0;
        do {
            var b1 = Base64.characters.indexOf( string.charAt(i++) );
            var b2 = Base64.characters.indexOf( string.charAt(i++) );
            var b3 = Base64.characters.indexOf( string.charAt(i++) );
            var b4 = Base64.characters.indexOf( string.charAt(i++) );

            var a = ( ( b1 & 0x3F ) << 2 ) | ( ( b2 >> 4 ) & 0x3 );
            var b = ( ( b2 & 0xF  ) << 4 ) | ( ( b3 >> 2 ) & 0xF );
            var c = ( ( b3 & 0x3  ) << 6 ) | ( b4 & 0x3F );

            result += String.fromCharCode(a) + (b?String.fromCharCode(b):'') + (c?String.fromCharCode(c):'');

        } while( i < string.length );

        return result;
    }
};