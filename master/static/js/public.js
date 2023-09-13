class Utils {
    AjaxGet (url, func) {
        $.ajax({
            type: "get",
            url: url,
            data: "",
            dataType: 'json',
            success: function(data){
                func(data);
            },
            error: function(xhr,textStatus) {
                console.log(xhr, textStatus);
            }
        });
    }

    AjaxGetNotAsync(url, func) {
        $.ajax({
            type: "get",
            url: url,
            data: "",
            dataType: 'json',
            async: false,
            success: function(data){
                func(data);
            },
            error: function(xhr,textStatus) {
                console.log(xhr, textStatus);
            }
        });
    }

    AjaxPost(url, param, func) {
        $.ajax({
            type: "post",
            url: url,
            data: JSON.stringify(param),
            dataType: 'json',
            success: function(data){
                func(data);
            },
            error: function(xhr,textStatus) {
                console.log(xhr, textStatus);
            }
        });
    }

    ToastShow(msg) {
        const toastLiveExample = $('#liveToast')
        const toast = new bootstrap.Toast(toastLiveExample)
        $("#liveToastMsg").append(msg)
        console.log("ToastShow...")
        toast.show()
    }

    MailToListJoin(toList) {
        if (Array.isArray(toList)) {
            return toList.join(";");
        }
        return toList
    }

}

export default new Utils()


