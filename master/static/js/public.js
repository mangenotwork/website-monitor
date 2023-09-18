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
}

let Mail= {
    init: {
        api: "/api/mail/init",
        data: true,
    },
    conf: {
        api: "/api/mail/conf",
        param: {
            host: "smtp.qq.com",
            port: 25,
            from: "",
            authCode: "",
            toList: "",
        },
    },
    infoApi: "/api/mail/info",
    sendApi: "/api/mail/sendTest",
    common: new Utils(),
    toListJoin: function(toList) {
        if (Array.isArray(toList)) {
            return toList.join(";");
        }
        return toList
    },
    hasSet: function () {
        let t = this;
        console.log("getMail...")
        t.common.AjaxGet(t.init.api, function (data){
            t.init.data = data.data;
        })
    },
    setConf: function () {
        let t = this;
        t.conf.param.toList = t.toListJoin(t.conf.param.toList);
        t.conf.param.port = Number(t.conf.param.port)
        t.common.AjaxPost(t.conf.api, t.conf.param, function (data){
            if (data.code === 0) {
                $("#mailSetModal").modal('toggle');
                t.getMail();
            }
            common.ToastShow(data.msg);
        });
    },
    getInfo: function () {
        let t = this;
        t.common.AjaxGet(t.infoApi, function (data){
            t.conf.param = data.data;
        });
    },
    send: function () {
        let t = this;
        t.conf.param.toList = t.toListJoin(t.conf.param.toList);
        t.conf.param.port = Number(t.conf.param.port)
        t.common.AjaxPost(t.sendApi, t.conf.param, function (data){
            common.ToastShow(data.msg);
        });
    },
}

let AddWebSite = {
    api: "/api/website/add",
        param : {
        host: "",
            healthUri: "",
            rate: 10,
            alarmResTime: 3000,
            uriDepth: 2,
            uriUpdateRate: 24,
    }
}