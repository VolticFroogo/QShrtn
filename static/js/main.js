const ResponseCode = {
    Success: 0,
    InternalServerError: 1,
    ForbiddenDomain: 2,
    TakenID: 3,
    InvalidURL: 4
};

$(document).ready(function(){
    var customURLVisible = false;

    // When the custom URL button is clicked, toggle the visibility of everything.
    $("#custom-url-button").click(function(){
        if (customURLVisible) {
            $("#custom-url").hide();
            $("#custom-url-button .up").hide();
            $("#custom-url-button .down").show();
        } else {
            $("#custom-url").show();
            $("#custom-url-button .up").show();
            $("#custom-url-button .down").hide();
        }

        customURLVisible = !customURLVisible;
    });

    // When the "scroll down" button is pressed, scroll down to download
    $("#footer").click(function(){
        $("#download")[0].scrollIntoView();
    });

    $("#shorten-button").click(function(){
        var url = $("#url").val();
        var id = $("#custom-url").val();
        var lower = url.toLowerCase();

        // You can't create recursive links, so we'll assume they want to copy it.
        if (lower.includes(window.location.hostname)) {
            // Show a toast with the new URL.
            M.toast({html: "Copied " + url + " to clipboard."});

            // Copy the URL to the clipboard.
            document.getElementById("url").select();
            document.execCommand("copy");

            return;
        }

        // Tell the user if their URL is too long.
        if (url.length > 2048) {
            // Show a toast with the new URL.
            M.toast({html: "URL is too long, maximum size is 2048."});

            return;
        }

        // Tell the user if their ID is too long.
        if (id.length > 8) {
            M.toast({html: "Custom URL is too long, maximum size is 8."});

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
            data: JSON.stringify(customURLVisible ?
                {   // If custom URL is visible, include an ID.
                    url: url,
                    id: id
                }:
                {   // Otherwise, just send the URL.
                    url: url
                }
            ),
            dataType: "json",
            success: function(r) {
                // Set the code for the code outside of the AJAX scope.
                Code = r.code;

                switch (r.code) {
                    // Success.
                    case ResponseCode.Success:
                        // Set the URL field to our new URL.
                        $("#url").val(window.location.protocol + "//" + window.location.hostname + "/" + r.id);

                        break;

                    // Internal server error.
                    case ResponseCode.InternalServerError:
                        M.toast({html: "Internal server error: " + r.error});
                        break;

                    // Forbidden domain.
                    case ResponseCode.ForbiddenDomain:
                        M.toast({html: "You can not shorten " + window.location.hostname + " links."});
                        break;

                    // Taken ID.
                    case ResponseCode.TakenID:
                        M.toast({html: "That custom link has already been taken."});
                        break;

                    // Invalid URL.
                    case ResponseCode.InvalidURL:
                        M.toast({html: "That is not a valid URL."});
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
            // Copy the URL to the clipboard.
            document.getElementById("url").select();
            document.execCommand("copy");

            // Show a toast with the new URL.
            M.toast({html: "Copied " + $("#url").val() + " to clipboard."});
        }
    });
});
