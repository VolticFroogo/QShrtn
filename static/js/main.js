$(document).ready(function(){
    $("#shorten-button").click(function(){
        var url = $("#url").val();
        var lower = url.toLowerCase();

        // You can't shorten qshr.tn links, so we'll assume they want to copy it.
        if (lower.includes("qshr.tn")) {
            // Show a toast with the new URL.
            M.toast({html: "Copied " + url + " to clipboard."});

            // Copy the URL to the clipboard.
            document.getElementById("url").select();
            document.execCommand("copy");

            return;
        }

        // Define our code outside of the AJAX scope.
        var Code = -1;

        // Send an AJAX request to create a new URL.
        $.ajax({
            url: "/new/",
            type: "POST",
            async: false,
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify({
                URL: url
            }),
            dataType: "json",
            success: function(r) {
                // Set the code for the code outside of the AJAX scope.
                Code = r.Code;

                switch(r.Code) {
                    // Success.
                    case 0:
                        // Set the URL field to our new URL.
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

        /*
            Comment about using AJAX synchronously which is deprecated:

            Because Firefox is stupid, the callbacks from an AJAX request
            even though technically originating from a click event are
            not click events anymore.

            I can't find any workaround for this except effectively
            hanging the page until the URL is back.

            :(
        */

        // If we successfully shortened the URL:
        if (Code === 0) {
            // Show a toast with the new URL.
            M.toast({html: "Copied " + url + " to clipboard."});

            // Copy the URL to the clipboard.
            document.getElementById("url").select();
            document.execCommand("copy");
        }
    });
});