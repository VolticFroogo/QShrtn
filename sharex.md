# [ShareX](https://getsharex.com/)
[ShareX](https://getsharex.com/) is a great screen capture and file sharing tool for Windows
(which is also [open source](https://github.com/ShareX/ShareX)).
It also supports automatically shortening URLs after performing a task using a URL shortener such as qshr.tn.

## Add qshr.tn to ShareX
To setup qshr.tn in ShareX, open the main window, click destinations, then custom uploader settings.

Copy the configuration below to your clipboard.

```json
{
  "Version": "12.4.1",
  "DestinationType": "URLShortener",
  "RequestMethod": "POST",
  "RequestURL": "https://qshr.tn/new/",
  "Body": "JSON",
  "Data": "{\"url\":\"$input$\"}",
  "URL": "https://qshr.tn/$json:id$"
}
```

Now click import, and select from clipboard.

If you have other custom uploaders already added, make sure to set the custom URL shortener to qshr.tn.

## Enable URL shortening
Now qshr.tn has been added, you will need to enable it as you see fit.

First, to actually use qshr.tn over other URL shorteners in ShareX, click destinations, URL shortener, and select custom URL shortener.

To enable automatic URL shortening, click after upload tasks, and enable shorten URL.

You can customise this further, but I'm sure you could figure out how to do that just by tinkering.
