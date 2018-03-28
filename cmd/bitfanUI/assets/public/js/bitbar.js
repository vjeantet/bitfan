// UMD style ?
// UMD (Universal Module Definition) for bitbar
;
(function(root, UMD) {
    if (typeof define === 'function' && define.amd) define([], UMD)
    else if (typeof module === 'object' && module.exports) module.exports = UMD()
    else root.bitbar = UMD()
})(this, function UMD() {

    function isFunction(functionToCheck) {
        var getType = {};
        return functionToCheck && getType.toString.call(functionToCheck) === '[object Function]';
    }

    function bitbarNew(instanceOptions) {

        initWith([])
        var bitbar = {
            hide: function() {
                $('#bitbar').modal('hide');
            },
            // Display bitbar with default items
            show: function(label, context) {
                // console.log("bitbar - GO - label=" + label);
                if (this.items.length == 0) {
                    // console.log("no bitbar item set");
                    return
                } else {
                    this.showWith(this.items, label, context)
                }
            },
            // Display bitbar with given items
            showWith: function(items, label, context) {
                if (label == null) { label = "" }
                // console.log("bitbar - GO With keys - label=" + label);

                if (context.selected == null) {
                    context.selected = [];
                }
                this.context = context
                this.citems = items

                // set input value with label
                $('#bitbar-input').val(label);
                this.search(label);

                var lastActiveElement = document.activeElement;
                
                // Display modal and focus on input
                $('#bitbar').on('shown.bs.modal', function() {
                    $('#bitbar-input').focus();
                });
                $('#bitbar').on('hidden.bs.modal', function() {
                    lastActiveElement.focus();
                });
                $('#bitbar').modal({
                    keyboard: true
                });
            },
            // Do search term
            search: function(txt) {
                const term = getSearchLower(txt);

                // console.log('bitbar - search - term = "' + term + '"')

                // let objects = [ 'Favorite Color', 'Google Chrome',  'Launch Chrome' }]
                if (txt != null && txt != "") {
                    this.renderResults(fuzzysort.go(term, this.citems, { keys: ['label'] }));
                } else {
                    this.renderAll();
                }

            },
            // render all items
            renderAll: function() {
                // Display all
                var html = ''
                for (var i = 0; i < this.citems.length; i++) {
                    const result = this.citems[i]
                    if (i == 0) {
                        html += this.renderItem(result.id, result.label, "active");
                    } else {
                        html += this.renderItem(result.id, result.label, "");
                    }
                }
                if (this.citems.length > 0) {
                    this.focus(this.citems[0].id)
                }
                $('#bitbar-results').html(html);
            },
            // render results
            renderResults: function(results) {
                // console.log("result count = " + results.length);
                var html = ''
                for (var i = 0; i < results.length; i++) {
                    const result = results[i]
                    if (i == 0) {
                        html += this.renderItem(result.obj.id, fuzzysort.highlight(result[0]), "active");
                    } else {
                        html += this.renderItem(result.obj.id, fuzzysort.highlight(result[0]), "");
                    }
                }
                if (results.length > 0) {
                    this.focus(results[0].obj.id)    
                }
                
                $('#bitbar-results').html(html);
            },
            renderItem: function(id, label, className) {
                return '<li bid="' + id + '" class="' + className + '">' + label + '</li>';
            },
            focus: function(id){
                // console.log("focus "+id)
                var result = $.grep(this.citems, function(e) { return e.id == id; });
                if (result.length == 0) {
                    // console.log("id " + id + " not found")
                    return
                } else if (result.length > 1) {
                    // multiple items found
                    // console.log("crazy ids found for " + id)
                    return
                }

                if (result[0].help != null && result[0].help != "") {
                    $("#bitbar-console div").html("<pre style='color:white'>"+result[0].help+"</pre>")
                    $("#bitbar-console").show()
                }else{
                    $("#bitbar-console").hide()
                }

                if (isFunction(result[0].onFocus)) {
                    result[0].onFocus(this.context,result[0])
                } else {
                    // console.log("no focus callback for item " + id)
                }
                
            },
            // Action Entry selection
            select: function(id) {
                // console.log("selected entry id = " + id)
                // hide modal
                this.hide();
                // find the "onSelect" callback
                var result = $.grep(this.citems, function(e) { return e.id == id; });
                if (result.length == 0) {
                    // console.log("id " + id + " not found")
                } else if (result.length == 1) {
                    if (isFunction(result[0].onSelect)) {
                        this.context.selected.unshift(id)
                        result[0].onSelect(this.context,result[0])
                    } else {
                        // console.log("no callback for item " + id)
                    }
                } else {
                    // multiple items found
                    // console.log("crazy ids found for " + id)
                }
            },

            new: bitbarNew,
            items: [],
            citems: [],
            context: {},
        }

        return bitbar;
    }

    function getSearchLower(txt) { return txt.toLowerCase() }

    function scrollToVisible(selected) {
        if (
            selected.position().top + selected.height() < $('#bitbar-results').position().top + $('#bitbar-results').height() &
            selected.position().top >= $('#bitbar-results').position().top
        ) {return} 

        // hidden on top
        if (selected.position().top >= $('#bitbar-results').position().top) {
            $('#bitbar-results').scrollTop(
                $('#bitbar-results').scrollTop() + selected.position().top + selected.height() - $('#bitbar-results').height() -$('#bitbar-results').position().top
            );
        }

        // hidden bottom
        if (selected.position().top + selected.height() < $('#bitbar-results').position().top + $('#bitbar-results').height()){            
            $('#bitbar-results').scrollTop(
                $('#bitbar-results').scrollTop() + (selected.position().top - $('#bitbar-results').position().top)
            );
        }
    }

    // Init Bitbar with data
    function initWith(keys) {
        //      Register shortcut CTRL+SHIFT+P (optional)
        //      Create the modal screen on DOM
        //      Prevent UP DOWN and ENTER on input
        $(document).on("blur", "#bitbar-input", function(event) {
            $("#bitbar-input").focus();
        });
        $(document).on("keydown", "#bitbar-input", function(event) {
            if (event.which == 9) {
                event.preventDefault();
                bitbar.hide();
            }
            if (event.which == 38) {
                // console.log("UP");
                event.preventDefault();

                active = $('#bitbar-results > li.active')
                if (active.length > 0) { // if found

                    if (active.is(':first-child')) {
                        // if first
                    } else {
                        // if not first
                        active.prev().addClass("active");
                        active.removeClass("active");
                        bitbar.focus(active.prev().attr("bid"));    
                        scrollToVisible(active.prev());
                    }
                } else {
                    // console.log("todo : use the last visible element");
                }
            }
            if (event.which == 40) {
                // console.log("DOWN");
                // find active li in #bitbar-results
                active = $('#bitbar-results > li.active')

                if (active.length > 0) { // if found
                    // console.log("active li found");
                    // if the active one is the last one
                    if (active.is(':last-child')) {
                        // console.log("active is the last one");
                    } else {
                        // console.log("active is not the last one");
                        // select next sibbling element
                        active.next().addClass("active");
                        active.removeClass("active");
                        bitbar.focus(active.next().attr("bid"));    
                        scrollToVisible(active.next())
                    }
                } else { // try to select first
                    // console.log("no active li found");
                    $('#bitbar-results > li:first').addClass("active");
                }

            }
            if (event.which == 13) {
                // console.log("ENTER");
                active = $('#bitbar-results > li.active')
                bitbar.select(active.attr("bid"));
            }
        });
        // On input change do Search
        $(document).on("input", "#bitbar-input", function(event) {
            // $('#bitbar-input').on("input",function(e){
            // console.log("bitbar - input : " + event.target.value);
            bitbar.search(event.target.value);
        });

        $(document).on("mousedown", "#bitbar-results li", function(event) {
            $('#bitbar-results > li.active').removeClass("active")
            $(event.target).addClass("active")
        });
        $(document).on("click", "#bitbar-results li", function(event) {
            bitbar.select($(event.target).attr("bid"));
        });


    }

    return bitbarNew()
}) // UMD