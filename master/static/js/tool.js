const { createApp, ref } = Vue;
const common = new Utils;
const app = createApp({
    data() {
        return {
            mail: Mail,
            addWebSite: AddWebSite,
            func: funcData,
            objFidWebsiteTDKI: {
                url: "",
                api: function (){ return "/api/tool/website/tdki?url=" + this.url; },
                rse: {},
            },
            objFidIp:{
                ip: "",
                api: function (){ return "/api/tool/ip?ip=" + this.ip; },
                rse: "",
            },
            objFidNsLookUp: {
                host: "",
                api: function (){ return "/api/tool/nsLookUp/all?host=" + this.host; },
                rse: {},
            },
            objFidWhois: {
                host: "",
                api: function (){ return "/api/tool/whois?host=" + this.host; },
                rse: {},
            },
            objFidICP: {
                host: "",
                api: function (){ return "/api/tool/icp?host=" + this.host; },
                rse: {},
            },
            objFidPing: { // TODO...
                host: "",
                api: function (){ return "/api/tool/icp?host=" + this.host; },
                rse: {},
            },
            objFidSSL: {
                url: "",
                api: function (){ return "/api/tool/certificate?url=" + this.url; },
                rse: {},
            },
            objFidWebsiteInfo: {
                url: "",
                api: function (){ return "/api/tool/website/collectInfo?url=" + this.url; },
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
        open: function (fid) {
            let t = this;
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
            });
        },

        submitFidIp: function () {
            let t = this;
            if (t.objFidIp.ip === "") {
                common.ToastShow("请输入IP地址！");
                return
            }
            common.AjaxGet(t.objFidIp.api(), function (data){
                t.objFidIp.rse = data.data;
            });
        },

        submitFidNsLookUp: function () {
            let t = this;
            $("#FidNsLookUpRse").hide();
            $("#FidNsLookUpLoading").show();
            common.AjaxGet(t.objFidNsLookUp.api(), function (data){
                t.objFidNsLookUp.rse = data.data;
                $("#FidNsLookUpLoading").hide();
                $("#FidNsLookUpRse").show();
            });
        },

        submitFidWhois: function () {
            let t = this;
            common.AjaxGet(t.objFidWhois.api(), function (data){
                t.objFidWhois.rse = data.data;
            });
        },

        submitFidICP: function () {
            let t = this;
            common.AjaxGet(t.objFidICP.api(), function (data){
                t.objFidICP.rse = data.data;
            });
        },

        submitFidPing: function () {
            let t = this;
            common.AjaxGet(t.objFidPing.api(), function (data){
                t.objFidPing.rse = data.data;
            });
        },

        submitFidSSL: function () {
            let t = this;
            common.AjaxGet(t.objFidSSL.api(), function (data){
                t.objFidSSL.rse = data.data;
            });
        },

        submitFidWebsiteInfo: function () {
            let t = this;
            common.AjaxGet(t.objFidWebsiteInfo.api(), function (data){
                t.objFidWebsiteInfo.rse = data.data;
            });
        }

    },
    computed: {
    },
    mounted:function(){
    },
});
app.mount('#app');