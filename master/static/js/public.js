class Utils {
    AjaxGet (url, func) {
        $.ajax({
            type: "get",
            url: url,
            data: "",
            dataType: 'json',
            async: true,
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
            async: true,
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
        $("#liveToastMsg").empty(msg);
        $("#liveToastMsg").append(msg);
        toast.show()
    }
}

// Mail
let Mail= {
    init: {
        api: "/api/mail/init",
        data: true,
    },
    conf: {
        api: "/api/mail/conf",
    },
    param: {
        host: "smtp.qq.com",
        port: 25,
        from: "",
        authCode: "",
        toList: "",
        open: 0,
    },
    paramHost: "smtp.qq.com",
    paramPort: 25,
    paramFrom: "",
    paramAuthCode: "",
    paramToList: "",
    setConfApi: "/api/mail/conf",
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
        Mail.common.AjaxGet(Mail.init.api, function (data){
            Mail.init.data = data.data;
        })
    },
    setConf: function () {
        let t = this;
        Mail.param.toList = Mail.toListJoin(Mail.param.toList);
        Mail.param.port = Number(Mail.param.port);
        Mail.common.AjaxPost(Mail.setConfApi, Mail.param, function (data){
            if (data.code === 0) {
                $("#mailSetModal").modal('toggle');
                Mail.getMail();
            }
            common.ToastShow(data.msg);
        });
    },
    getInfo: function () {
        let t = this;
        Mail.common.AjaxGet(Mail.infoApi, function (data){
            if (data.code === 0 ){
                Mail.param = data.data;
            }
        });
    },
    send: function () {
        let t = this;
        Mail.param.toList = Mail.toListJoin(Mail.param.toList);
        Mail.param.port = Number(Mail.param.port);
        Mail.common.AjaxPost(Mail.sendApi, Mail.param, function (data){
            common.ToastShow(data.msg);
        });
    },
}

// add website
let AddWebSite = {
    api: "/api/website/add",
    apiGet: "/api/website/list",
    param : {},
    list : [],
    common: new Utils(),
    next1Param: {
        hostProtocol: "https://",
        hostUrl: "",
        notes: "",
        monitorRate: 15,
        contrastUrl: "https://www.baidu.com",
        contrastTime: 1000,
        ping: "223.5.5.5",
        pingTime: 1000,
    },
    next2Param: {
        uriDepth: 2,
        scanRate: 24,
        scanBadLink: true,
        scanTDK: false,
        scanExtLinks: false,
    },
    next3Param: {
        websiteSlowResponseTime: 3000,
        websiteSlowResponseCount: 3,
        SSLCertificateExpire: 14,
    },
    addNext1: function () {
        console.log(AddWebSite.next1Param);
        if (AddWebSite.next1Param.hostUrl === "") {
            common.ToastShow("请输入检测网站的URL");
            return
        }
        if (AddWebSite.next1Param.monitorRate === 0) {
            common.ToastShow("检测频率建议大于1s");
            return
        }
        if (AddWebSite.next1Param.contrastUrl === "") {
            common.ToastShow("对照组不能为空");
            return
        }
        if (AddWebSite.next1Param.contrastTime === "") {
            common.ToastShow("对照组慢响应时间不低于100ms");
            return
        }
        if (AddWebSite.next1Param.ping === "") {
            common.ToastShow("Ping地址不能为空");
            return
        }
        if (AddWebSite.next1Param.pingTime === "") {
            common.ToastShow("Ping慢响应时间不低于100ms");
            return
        }
        $("#next1Forms").hide();
        $("#next2Forms").show();
        $("#next1").removeClass("addTipNow");
        $("#next2").addClass("addTipNow");
        $("#next1But").hide();
        $("#next2But").show();
    },
    addReturn1: function () {
        $("#next1Forms").show();
        $("#next2Forms").hide();
        $("#next1").addClass("addTipNow");
        $("#next2").removeClass("addTipNow");
        $("#next1But").show();
        $("#next2But").hide();
    },
    addNext2: function () {
        console.log(AddWebSite.next2Param);
        if (AddWebSite.next2Param.uriDepth <1) {
            common.ToastShow("扫描网站深度不能小于1");
            return
        }
        if (AddWebSite.next2Param.scanRate <1) {
            common.ToastShow("扫描网站频率不能小于1");
            return
        }
        $("#next2Forms").hide();
        $("#next3Forms").show();
        $("#next2").removeClass("addTipNow");
        $("#next3").addClass("addTipNow");
        $("#next2But").hide();
        $("#next3But").show();
    },
    addReturn2: function () {
        let t = this;
        $("#next2Forms").show();
        $("#next3Forms").hide();
        $("#next2").addClass("addTipNow");
        $("#next3").removeClass("addTipNow");
        $("#next2But").show();
        $("#next3But").hide();
    },
    addSubmit: function () {
        console.log(AddWebSite.next3Param);
        if (AddWebSite.next3Param.websiteSlowResponseTime <100) {
            common.ToastShow("网站响应慢不能小于100ms");
            return
        }
        if (AddWebSite.next3Param.websiteSlowResponseCount <1) {
            common.ToastShow("网站响应慢发送邮件连次数不能小于1次");
            return
        }
        if (AddWebSite.next3Param.SSLCertificateExpire <1) {
            common.ToastShow("网站证书过期触发报警天数不能小于1天");
            return
        }
        AddWebSite.param.host = AddWebSite.next1Param.hostProtocol + AddWebSite.next1Param.hostUrl;
        AddWebSite.param.notes = AddWebSite.next1Param.notes;
        AddWebSite.param.monitorRate = Number(AddWebSite.next1Param.monitorRate);
        AddWebSite.param.contrastUrl = AddWebSite.next1Param.contrastUrl;
        AddWebSite.param.contrastTime = Number(AddWebSite.next1Param.contrastTime);
        AddWebSite.param.ping = AddWebSite.next1Param.ping;
        AddWebSite.param.pingTime = Number(AddWebSite.next1Param.pingTime);
        AddWebSite.param.uriDepth = Number(AddWebSite.next2Param.uriDepth);
        AddWebSite.param.scanRate = Number(AddWebSite.next2Param.scanRate);
        AddWebSite.param.scanBadLink = AddWebSite.next2Param.scanBadLink;
        AddWebSite.param.scanTDK = AddWebSite.next2Param.scanTDK;
        AddWebSite.param.scanExtLinks = AddWebSite.next2Param.scanExtLinks;
        AddWebSite.param.websiteSlowResponseTime = Number(AddWebSite.next3Param.websiteSlowResponseTime);
        AddWebSite.param.websiteSlowResponseCount = Number(AddWebSite.next3Param.websiteSlowResponseCount);
        AddWebSite.param.SSLCertificateExpire = Number(AddWebSite.next3Param.SSLCertificateExpire);
        console.log(AddWebSite.param);
        common.AjaxPost(AddWebSite.api, AddWebSite.param, function (data){
            console.log(data);
            if (data.code === 0) {
                $("#addHostModal").modal('toggle');
                AddWebSite.next1Param.hostUrl = ""
                AddWebSite.next1Param.notes = ""
                location.reload();
            }
            common.ToastShow(data.msg);
        });

    },
    getList: function () {
        common.AjaxGet(AddWebSite.apiGet, function (data){
            console.log(data);
            if (data.code === 0) {
                AddWebSite.list = data.data;
            } else {
                common.ToastShow(data.msg);
            }
        });
    }
}
