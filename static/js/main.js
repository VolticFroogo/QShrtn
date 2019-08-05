$(document).ready(function(){
    $(".sidenav").sidenav();

    var clipboard = new ClipboardJS("#copy-button");

    clipboard.on("success", function(e) {
        if (e.text !== "") {
            M.toast({html: "Copied " + e.text + " to your clipboard."});
        }
    });

    clipboard.on("error", function() {
        M.toast({html: "Copy failed, press CTRL+C to copy."});
    });

    $("#shorten-button").click(function(){
        var url = $("#url").val();
        url = url.toLowerCase();

        if (url.includes("qshr.tn")) {
            M.toast({html: "You can not shorten qshr.tn links."});
            return;
        }

        $.ajax({
            url: "/new/",
            type: "POST",
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify({
                URL: url
            }),
            dataType: "json",
            success: function(r) {
                switch(r.Code) {
                    // Success.
                    case 0:
                        M.toast({html: "Shortened to qshr.tn/" + r.ID});
                        $("#url").val("https://qshr.tn/" + r.ID);
                        break;

                    // Internal server error.
                    case 1:
                        M.toast({html: "Internal server error."});
                        break;

                    // Forbidden domain.
                    case 2:
                        M.toast({html: "You can not shorten qshr.tn links."});
                        break;
                }
            }
        });
    });
});