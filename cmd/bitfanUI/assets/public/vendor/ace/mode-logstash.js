ace.define("ace/mode/logstash_highlight_rules",["require","exports","module","ace/lib/oop","ace/lib/lang","ace/mode/text_highlight_rules"], function(require, exports, module) {
"use strict";

var oop = require("../lib/oop");
var TextHighlightRules = require("./text_highlight_rules").TextHighlightRules;

var LogstashHighlightRules = function() {
    // regexp must not have capturing parentheses. Use (?:) instead.
    // regexps are ordered -> the first match is used

    this.$rules = {
        start: [{
            token: "storage.type.logstash",
            regex: /^(?:input|filter|codec|output)/,
            comment: "classes: inputs, codecs, filters and outputs"
        }, {
            token: [
                "keyword.operator.logstash",
                "text.logstash",
                "text.logstash",
                "entity.name.function.logstash",
                "text.logstash"
            ],
            regex: /(?:(and|or)(\s+)(\[)(\w*)(\]))+/,
            comment: "complex if/else if statements"
        }, {
            token: "keyword.operator.logstash",
            regex: /==|!=|<|>|<=|>=|=~|!~|in|not in|and|or|nand|xor|!/,
            comment: "Operators"
        }, {
            token: [
                "entity.name.function.logstash",
                "entity.name.function.logstash",
                "entity.name.function.logstash"
            ],
            regex: /(%{)(\w*)(})/,
            comment: "Groked Field"
        }, {
            token: "string.text.logstash",
            regex: /".+?[^"]*"/,
            comment: "String values"
        }, {
            token: [
                "keyword.control.logstash",
                "text.logstash",
                "text.logstash",
                "entity.name.function.logstash",
                "text.logstash",
                "text.logstash",
                "keyword.operator.logstash"
            ],
            regex: /(if|else if)(\s+)(\[)(\w*)(\])(\s*)(==|!=|<|>|<=|>=|=~|!~|in|not in|!)/,
            comment: "if/else if statements"
        }, {
            token: [
                "keyword.control.logstash",
                "text.logstash",
                "text.logstash"
            ],
            regex: /(else)(\s+)({)/,
            comment: "else statements"
        }, {
            token: [
                "text.logstash",
                "entity.name.function.logstash",
                "text.logstash",
                "string.text.logstash",
                "text.logstash",
                "text.logstash",
                "variable.text.logstash",
                "text.logstash",
                "keyword.operator.logstash",
                "text.logstash"
            ],
            regex: /^(\s*)(\w+)(\s*)("?.+?[^"]*"?)(\s*{)((?:\s*)?)((?:\w+)?)((?:\s*)?)((?:=>)?)((?:\s*)?)/,
            comment: "functions: types of inputs, codecs, filters and outputs"
        }, {
            token: "comment.line.number-sign.logstash",
            regex: /#.+$/,
            comment: "Comments"
        }, {
            token: [
                "keyword.text.logstash",
                "variable.text.logstash",
                "keyword.text.logstash",
                "keyword.operator.logstash",
                "keyword.text.logstash",
                "constant.numeric.logstash"
            ],
            regex: /^((?:\s*)?)(\w+)((?:\s*)?)(=>)((?:\s*)?)(\d+)/,
            comment: "Variables: Number values"
        }, {
            token: [
                "keyword.text.logstash",
                "variable.text.logstash",
                "keyword.text.logstash",
                "keyword.operator.logstash",
                "keyword.text.logstash"
            ],
            regex: /^((?:\s*)?)(\w+)((?:\s*)?)(=>)((?:\s*)?)/,
            comment: "Variables: String values"
        }]
    }
    
    this.normalizeRules();
};

oop.inherits(LogstashHighlightRules, TextHighlightRules);

exports.LogstashHighlightRules = LogstashHighlightRules;
});

ace.define("ace/mode/folding/cstyle",["require","exports","module","ace/lib/oop","ace/range","ace/mode/folding/fold_mode"], function(require, exports, module) {
"use strict";

var oop = require("../../lib/oop");
var Range = require("../../range").Range;
var BaseFoldMode = require("./fold_mode").FoldMode;

var FoldMode = exports.FoldMode = function(commentRegex) {
    if (commentRegex) {
        this.foldingStartMarker = new RegExp(
            this.foldingStartMarker.source.replace(/\|[^|]*?$/, "|" + commentRegex.start)
        );
        this.foldingStopMarker = new RegExp(
            this.foldingStopMarker.source.replace(/\|[^|]*?$/, "|" + commentRegex.end)
        );
    }
};
oop.inherits(FoldMode, BaseFoldMode);

(function() {
    
    this.foldingStartMarker = /([\{\[\(])[^\}\]\)]*$|^\s*(\/\*)/;
    this.foldingStopMarker = /^[^\[\{\(]*([\}\]\)])|^[\s\*]*(\*\/)/;
    this.singleLineBlockCommentRe= /^\s*(\/\*).*\*\/\s*$/;
    this.tripleStarBlockCommentRe = /^\s*(\/\*\*\*).*\*\/\s*$/;
    this.startRegionRe = /^\s*(\/\*|\/\/)#?region\b/;
    this._getFoldWidgetBase = this.getFoldWidget;
    this.getFoldWidget = function(session, foldStyle, row) {
        var line = session.getLine(row);
    
        if (this.singleLineBlockCommentRe.test(line)) {
            if (!this.startRegionRe.test(line) && !this.tripleStarBlockCommentRe.test(line))
                return "";
        }
    
        var fw = this._getFoldWidgetBase(session, foldStyle, row);
    
        if (!fw && this.startRegionRe.test(line))
            return "start"; // lineCommentRegionStart
    
        return fw;
    };

    this.getFoldWidgetRange = function(session, foldStyle, row, forceMultiline) {
        var line = session.getLine(row);
        
        if (this.startRegionRe.test(line))
            return this.getCommentRegionBlock(session, line, row);
        
        var match = line.match(this.foldingStartMarker);
        if (match) {
            var i = match.index;

            if (match[1])
                return this.openingBracketBlock(session, match[1], row, i);
                
            var range = session.getCommentFoldRange(row, i + match[0].length, 1);
            
            if (range && !range.isMultiLine()) {
                if (forceMultiline) {
                    range = this.getSectionRange(session, row);
                } else if (foldStyle != "all")
                    range = null;
            }
            
            return range;
        }

        if (foldStyle === "markbegin")
            return;

        var match = line.match(this.foldingStopMarker);
        if (match) {
            var i = match.index + match[0].length;

            if (match[1])
                return this.closingBracketBlock(session, match[1], row, i);

            return session.getCommentFoldRange(row, i, -1);
        }
    };
    
    this.getSectionRange = function(session, row) {
        var line = session.getLine(row);
        var startIndent = line.search(/\S/);
        var startRow = row;
        var startColumn = line.length;
        row = row + 1;
        var endRow = row;
        var maxRow = session.getLength();
        while (++row < maxRow) {
            line = session.getLine(row);
            var indent = line.search(/\S/);
            if (indent === -1)
                continue;
            if  (startIndent > indent)
                break;
            var subRange = this.getFoldWidgetRange(session, "all", row);
            
            if (subRange) {
                if (subRange.start.row <= startRow) {
                    break;
                } else if (subRange.isMultiLine()) {
                    row = subRange.end.row;
                } else if (startIndent == indent) {
                    break;
                }
            }
            endRow = row;
        }
        
        return new Range(startRow, startColumn, endRow, session.getLine(endRow).length);
    };
    this.getCommentRegionBlock = function(session, line, row) {
        var startColumn = line.search(/\s*$/);
        var maxRow = session.getLength();
        var startRow = row;
        
        var re = /^\s*(?:\/\*|\/\/|--)#?(end)?region\b/;
        var depth = 1;
        while (++row < maxRow) {
            line = session.getLine(row);
            var m = re.exec(line);
            if (!m) continue;
            if (m[1]) depth--;
            else depth++;

            if (!depth) break;
        }

        var endRow = row;
        if (endRow > startRow) {
            return new Range(startRow, startColumn, endRow, line.length);
        }
    };

}).call(FoldMode.prototype);

});

ace.define("ace/mode/logstash",["require","exports","module","ace/lib/oop","ace/mode/text","ace/mode/logstash_highlight_rules"], function(require, exports, module) {
"use strict";

var oop = require("../lib/oop");
var TextMode = require("./text").Mode;
var LogstashHighlightRules = require("./logstash_highlight_rules").LogstashHighlightRules;
// TODO: pick appropriate fold mode
var FoldMode = require("./folding/cstyle").FoldMode;
var Mode = function() {
    this.HighlightRules = LogstashHighlightRules;
    this.foldingRules = new FoldMode();
    this.$behaviour = this.$defaultBehaviour;
};
oop.inherits(Mode, TextMode);

(function() {
       
    this.lineCommentStart = "#";
    
    this.$id = "ace/mode/logstash";
}).call(Mode.prototype);

exports.Mode = Mode;
});




