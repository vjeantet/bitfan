+++
title = "pop3"
description = "Read mail from a pop3 server"
weight = 10
+++

{{% processordetails pop3processor %}}

## Produced event examples
{{%expand%}}
```
{
  "delivered-to": {
    "bitfan@free.fr": "",
  },
  "cc": {
    "valere.jeantet+cc@gmail.com": "",
  },
  "from": {
    "valere.jeantet@gmail.com": "Valere Jeantet",
  },
  "@timestamp":  2018-01-14 00:22:57 Local,
  "uid":         "0acv1f1810845a5aa63e020051224870",
  "subject":     "MSGBI",
  "html":        "<div dir=\"ltr\"><img src=\"cid:160f19327978e9f50f81\" alt=\"Capture d’écran 2018-01-12 à 12.11.52.png\" class=\"\" style=\"max-width: 100%;\"><br></div>\r\n",
  "return-path": {
    "valere.jeantet@gmail.com": "",
  },
  "to": {
    "bitfan@free.fr":           "bitfan@free.fr",
    "valere.jeantet@gmail.com": "Valere Jeantet",
  },
  "message": "[image: Capture d’écran 2018-01-12 à 12.11.52.png]\r\n",
  "headers": {
    "Return-Path": []string{
      "<valere.jeantet@gmail.com>",
    },
    "Dkim-Signature": []string{
      "v=1; a=rsa-sha256; c=relaxed/relaxed; d=gmail.com; s=XXX; h=mime-version:from:date:message-id:subject:to:cc; bh=ODEQpVa/5Qcw==",
    },
    "X-Received": []string{
      "by 10.XXX.55.XX with SMTP id z37mr30058367uad.XXX.1515881XXX; Sat, 13 Jan 2018 14:11:25 -0800 (PST)",
    },
    "Content-Type": []string{
      "multipart/mixed; boundary=\"001a11409cc4725d6c0562afa937\"",
    },
    "Received": []string{
      "from mail-SSS-f17X.google.com (mx25XXX.priv.proxad.net [172.XXX.XXX.95]) by SSSS-g26.priv.proxad.net (Postfix) with ESMTP id 1C055A00634 for <bitfan@free.fr>; Sat, 13 Jan 2018 23:11:27 +0100 (CET)",
      "from mail-ua0-XXXX.google.com ([209.XXX.XXX.179]) by mx1-XXX.XXX.fr (MXproxy) with ESMTPS for bitfan@free.fr (version=TLSv1/SSLv3 cipher=AES128-XXX-SHA256 bits=128); Sat, 13 Jan 2018 23:11:28 +0100 (CET)",
      "by mail-ua0-SSSS.google.com with SMTP id x10so6315721ual.8 for <bitfan@free.fr>; Sat, 13 Jan 2018 14:11:26 -0800 (PST)",
    },
    "X-Gm-Message-State": []string{
      "AKwxytfekxmrRWAc21hNb05lq kKc+spWodYAlFaCMPvuFfWOAjybzg==",
    },
    "Mime-Version": []string{
      "1.0",
    },
    "From": []string{
      "Valere Jeantet <valere.jeantet@gmail.com>",
    },
    "Message-Id": []string{
      "<CACE-i0C-T5D75xd1gJp1Q9imHw3rNwktuxc+OwxpyQ@mail.gmail.com>",
    },
    "To": []string{
      "\"bitfan@free.fr\" <bitfan@free.fr>, Valere Jeantet <valere.jeantet@gmail.com>",
    },
    "Delivered-To": []string{
      "bitfan@free.fr",
    },
    "X-Proxad-Sc": []string{
      "state=HAM score=0",
    },
    "X-Proxad-Cause": []string{
      "(null)",
    },
    "X-Google-Smtp-Source": []string{
      "ACJfBosF8RqtE9BCxZvUZbC0WnryH5nlfIn12xmVhWoP3KazNLz4RfB82gEdZofHxtKo=",
    },
    "Subject": []string{
      "Message from me",
    },
    "X-Google-Dkim-Signature": []string{
      "v=1; a=rsa-sha256; c=relaxed/relaxed; d=1e100.net; s=20561025; h=x-gm-message-state:mime-version:from:date:message-id:subject:to:cc; bh=ODEQpVa/ym6uT9lFT7X pdvw==",
    },
    "Date": []string{
      "Sat, 13 Jan 2018 22:11:13 +0000",
    },
    "Cc": []string{
      "valere.jeantet+cc@gmail.com",
    },
  },
  "text":   "[image: Capture d’écran 2018-01-12 à 12.11.52.png]\r\n",
  "sentAt": 2018-01-13 22:11:13 ,
}
```
{{%expand%}}