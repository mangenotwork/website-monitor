const { createApp, ref } = Vue;
const common = new Utils;
const app = createApp({
    data() {
        return {
            mail: Mail,
            addWebSite: AddWebSite,
            func: funcData,
            history: {
                add: function (id, value) { return "/api/tool/history?toolID="+id+"&value="+value; },
                get: function (id) { return "/api/tool/history?toolID="+id; },
                clear: function (id) {return "/api/tool/history/clear?toolID="+id; },
                data: [],
            },
            objFidWebsiteTDKI: {
                id:  FidWebsiteTDKI,
                url: "",
                api: function (){ return "/api/tool/website/tdki?url=" + this.url; },
                rse: {},
            },
            objFidIp:{
                id: FidIp,
                ip: "",
                api: function (){ return "/api/tool/ip?ip=" + this.ip; },
                rse: "",
            },
            objFidNsLookUp: {
                id: FidNsLookUp,
                host: "",
                api: function (){ return "/api/tool/nsLookUp/all?host=" + this.host; },
                rse: {},
            },
            objFidWhois: {
                id: FidWhois,
                host: "",
                api: function (){ return "/api/tool/whois?host=" + this.host; },
                rse: {},
            },
            objFidICP: {
                id: FidICP,
                host: "",
                api: function (){ return "/api/tool/icp?host=" + this.host; },
                rse: {},
            },
            objFidPing: {
                id: FidPing,
                ip: "",
                api: function (){ return "/api/tool/ping?ip=" + this.ip; },
                rse: {},
            },
            objFidSSL: {
                id: FidSSL,
                url: "",
                api: function (){ return "/api/tool/certificate?url=" + this.url; },
                rse: {},
            },
            objFidWebsiteInfo: {
                id: FidWebsiteInfo,
                host: "",
                api: function (){ return "/api/tool/website/collectInfo?host=" + this.host; },
                rse: {},
            },
        }
    },
    created:function(){
        let t = this;
        t.mail.getInfo();
    },
    destroyed:function () {
    },
    methods: {
        getHistory: function (id) {
            let t = this;
            common.AjaxGet(t.history.get(id), function (data){
                t.history.data = data.data;
            });
        },
        setHistory: function (id, value) {
            let t = this;
            common.AjaxPost(t.history.add(id, value), "", function (data){
                console.log(data)
            });
        },
        clearHistory: function (id) {
            let t = this;
            common.AjaxGet(t.history.clear(id), function (data){
                t.history.data = data.data;
            });
        },
        open: function (fid) {
            let t = this;
            t.getHistory(fid);
            switch (fid) {
                // 获取网站的T, D, K, 图标
                case FidWebsiteTDKI :
                    $("#FidWebsiteTDKIModal").modal("show");
                    break;
                // ip信息查询
                case FidIp :
                    $("#FidIpModal").modal("show");
                    break;
                // 查询dns
                case FidNsLookUp :
                    $("#FidNsLookUpModal").modal("show");
                    break;
                // Whois查询
                case FidWhois :
                    $("#FidWhoisModal").modal("show");
                    break;
                // 查询备案
                case FidICP :
                    $("#FidICPModal").modal("show");
                    break;
                // 在线ping
                case FidPing :
                    $("#FidPingModal").modal("show");
                    break;
                // 获取证书
                case FidSSL :
                    $("#FidSSLModal").modal("show");
                    break;
                // 网站信息获取
                case FidWebsiteInfo :
                    $("#FidWebsiteInfoModal").modal("show");
                    break;
            }
        },

        submitFidWebsiteTDKI: function () {
            let t = this;
            if (t.objFidWebsiteTDKI.url === "") {
                common.ToastShow("请输入链接！");
                return
            }
            common.AjaxGet(t.objFidWebsiteTDKI.api(), function (data){
                t.objFidWebsiteTDKI.rse = data.data;
                t.setHistory(t.objFidWebsiteTDKI.id, t.objFidWebsiteTDKI.url);
                t.getHistory(t.objFidWebsiteTDKI.id);
            });
        },

        gotoFidWebsiteTDKI: function (value) {
            let t = this;
            t.objFidWebsiteTDKI.url = value;
            t.submitFidWebsiteTDKI();
        },

        clearFidWebsiteTDKIHistory: function () {
            let t = this;
            t.clearHistory(t.objFidWebsiteTDKI.id);
            t.getHistory(t.objFidWebsiteTDKI.id);
        },

        submitFidIp: function () {
            let t = this;
            if (t.objFidIp.ip === "") {
                common.ToastShow("请输入IP地址！");
                return
            }
            common.AjaxGet(t.objFidIp.api(), function (data){
                t.objFidIp.rse = data.data;
                t.setHistory(t.objFidIp.id, t.objFidIp.ip);
                t.getHistory(t.objFidIp.id);
            });
        },

        gotoFidIp: function (value) {
            let t = this;
            t.objFidIp.ip = value;
            t.submitFidIp();
        },

        clearFidIpHistory: function () {
            let t = this;
            t.clearHistory(t.objFidIp.id);
            t.getHistory(t.objFidIp.id);
        },

        submitFidNsLookUp: function () {
            let t = this;
            if (t.objFidNsLookUp.host === "") {
                common.ToastShow("请输入Host！");
                return
            }
            $("#FidNsLookUpRse").hide();
            $("#FidNsLookUpLoading").show();
            common.AjaxGet(t.objFidNsLookUp.api(), function (data){
                t.objFidNsLookUp.rse = data.data;
                $("#FidNsLookUpLoading").hide();
                $("#FidNsLookUpRse").show();
                t.setHistory(t.objFidNsLookUp.id, t.objFidNsLookUp.host);
                t.getHistory(t.objFidNsLookUp.id);
            });
        },

        gotoFidNsLookUp: function (value) {
            let t = this;
            t.objFidNsLookUp.host = value;
            t.submitFidNsLookUp();
        },

        clearFidNsLookUpHistory: function () {
            let t = this;
            t.clearHistory(t.objFidNsLookUp.id);
            t.getHistory(t.objFidNsLookUp.id);
        },

        submitFidWhois: function () {
            let t = this;
            if (t.objFidWhois.host === "") {
                common.ToastShow("请输入Host！");
                return
            }
            common.AjaxGet(t.objFidWhois.api(), function (data){
                t.objFidWhois.rse = data.data;
                t.setHistory(t.objFidWhois.id, t.objFidWhois.host);
                t.getHistory(t.objFidWhois.id);
            });
        },

        gotoFidWhois: function (value) {
            let t = this;
            t.objFidWhois.host = value;
            t.submitFidWhois();
        },

        clearFidWhoisHistory: function () {
            let t = this;
            t.clearHistory(t.objFidWhois.id);
            t.getHistory(t.objFidWhois.id);
        },

        submitFidICP: function () {
            let t = this;
            if (t.objFidICP.host === "") {
                common.ToastShow("请输入Host！");
                return
            }
            common.AjaxGet(t.objFidICP.api(), function (data){
                t.objFidICP.rse = data.data;
                t.setHistory(t.objFidICP.id, t.objFidICP.host);
                t.getHistory(t.objFidICP.id);
            });
        },

        gotoFidICP: function (value) {
            let t = this;
            t.objFidICP.host = value;
            t.submitFidICP();
        },

        clearFidICPHistory: function () {
            let t = this;
            t.clearHistory(t.objFidICP.id);
            t.getHistory(t.objFidICP.id);
        },

        submitFidPing: function () {
            let t = this;
            if (t.objFidPing.ip === "") {
                common.ToastShow("请输入IP！");
                return
            }
            $("#FidPingRse").hide();
            $("#FidPingLoading").show();
            common.AjaxGet(t.objFidPing.api(), function (data){
                t.objFidPing.rse = data.data;
                $("#FidPingRse").show();
                $("#FidPingLoading").hide();
                t.setHistory(t.objFidPing.id, t.objFidPing.ip);
                t.getHistory(t.objFidPing.id);
            });
        },

        gotoFidPing: function (value) {
            let t = this;
            t.objFidPing.ip = value;
            t.submitFidPing();
        },

        clearFidPingHistory: function () {
            let t = this;
            t.clearHistory(t.objFidPing.id);
            t.getHistory(t.objFidPing.id);
        },

        submitFidSSL: function () {
            let t = this;
            if (t.objFidSSL.url === "") {
                common.ToastShow("请输入Url！");
                return
            }
            common.AjaxGet(t.objFidSSL.api(), function (data){
                t.objFidSSL.rse = data.data;
                t.setHistory(t.objFidSSL.id, t.objFidSSL.url);
                t.getHistory(t.objFidSSL.id);
            });
        },

        gotoFidSSL: function (value) {
            let t = this;
            t.objFidSSL.url = value;
            t.submitFidSSL();
        },

        clearFidSSLHistory: function () {
            let t = this;
            t.clearHistory(t.objFidSSL.id);
            t.getHistory(t.objFidSSL.id);
        },

        submitFidWebsiteInfo: function () {
            let t = this;
            if (t.objFidWebsiteInfo.host === "") {
                common.ToastShow("请输入Host！");
                return
            }
            common.AjaxGet(t.objFidWebsiteInfo.api(), function (data){
                t.objFidWebsiteInfo.rse = data.data;
                t.setHistory(t.objFidWebsiteInfo.id, t.objFidWebsiteInfo.host);
                t.getHistory(t.objFidWebsiteInfo.id);
            });
        },

        gotoFidWebsiteInfo: function (value) {
            let t = this;
            t.objFidWebsiteInfo.host = value;
            t.submitFidWebsiteInfo();
        },

        clearFidWebsiteInfoHistory: function () {
            let t = this;
            t.clearHistory(t.objFidWebsiteInfo.id);
            t.getHistory(t.objFidWebsiteInfo.id);
        },

    },
    computed: {
    },
    mounted:function(){
    },
});
app.mount('#app');