const { createApp, ref } = Vue;
const common = new Utils;
const app = createApp({
    data() {
        return {
            mail: Mail,
            addWebSite: AddWebSite,
            func: funcData,
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
                    console.log("FidWebsiteQuery...")
                    break;
                // ip信息查询
                case FidIp :
                    console.log("FidIp...")
                    break;
                // 查询dns
                case FidNsLookUp :
                    console.log("FidNsLookUp...")
                    break;
                // Whois查询
                case FidWhois :
                    console.log("FidWhois...")
                    break;
                // 查询备案
                case FidICP :
                    console.log("FidICP...")
                    break;
                // 在线ping
                case FidPing :
                    console.log("FidPing...")
                    break;
                // 获取证书
                case FidSSL :
                    console.log("FidSSL...")
                    break;
                // 网站信息获取
                case FidWebsiteInfo :
                    console.log("FidWebsiteInfo...")
                    break;
            }
        },
    },
    computed: {
    },
    mounted:function(){
    },
});
app.mount('#app');