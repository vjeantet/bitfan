+++
title = "Codecs"
description = ""
weight = 64
+++

Codecs are essentially stream decoder or encoder depending on where it operates, as part of an input or output.
They have their own options to handle charset, formating, etc...

This example codec, used in stdout or http or file or .... will encode your event into a utf-8 json with "   " as indentation
```
codec => json {
    charset => "UTF-8"
    indent => "    "
}

```

{{%children style="h2" description="true"%}}