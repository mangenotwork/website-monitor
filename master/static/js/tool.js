const { createApp, ref } = Vue;
import common from './public.js'
const FidWebsiteQuery = 1;
const FidIpQuery = 2;
const FidGetExternal = 3;
const FidDeadLink = 4;
const FidGetICP = 5;
const FidPing = 6;
const FidNsLookUp = 7;
const app = createApp({
    data() {
        return {
            func: [
                {fid: FidWebsiteQuery, title: "网站服务器查询", content: "查询网站的服务器IP地址，IP属地，服务器类型."},
                {fid: FidIpQuery, title: "IP属地查询", content: "查询IP属地."},
                {fid: FidGetExternal, title: "网站外链抓取", content: "网站外链抓取."},
                {fid: FidDeadLink, title: "网站死链检查", content: "网站死链检查."},
                {fid: FidGetICP, title: "备案查询", content: "网站备案查询."},
                {fid: FidPing, title: "Ping", content: "执行ping功能."},
                {fid: FidNsLookUp, title: "DNS查询", content: "执行nslookup功能."},
            ]
        }
    },
    created:function(){
        let t = this;
    },
    destroyed:function () {
    },
    methods: {
        open: function (fid) {
            let t = this;
            switch (fid) {
                case FidWebsiteQuery :
                    console.log("FidWebsiteQuery...")
                    break;
                case FidIpQuery :
                    console.log("FidIpQuery...")
                    break;
                case FidGetExternal :
                    console.log("FidGetExternal...")
                    break;
                case FidDeadLink :
                    console.log("FidDeadLink...")
                    break;
                case FidGetICP :
                    console.log("FidGetICP...")
                    break;
                case FidPing :
                    console.log("FidPing...")
                    $("#toolPingModal").modal("show");
                    break;
                case FidNsLookUp :
                    console.log("FidNsLookUp...")
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