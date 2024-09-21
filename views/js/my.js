function get_auth_token() {
    token = window.sessionStorage.getItem("auth_token")
    if (token == null) {
        refresh_token = getCookie("refresh_token");
        if (refresh_token != null) {
            $.ajax({
                type: "POST",
                url: "/token",
                data: { "refresh_token": refresh_token },
                success: function (result) {
                    token = result;
                    window.sessionStorage.setItem("auth_token", result);
                }
            }).fail(function (result, result1, result2) {
            });
        }
    }
    return token;
};

function getCookie(name) {
    var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
    if (arr = document.cookie.match(reg))
        return arr[2];
    else
        return null;
}