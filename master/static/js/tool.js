const { createApp } = Vue;
const common = new Utils;

const TOOLS = [
    { id: 1, group: "站点分析", name: "TDK 与图标", icon: "layers", key: "url", placeholder: "https://www.example.com", desc: "获取网页的 Title、Description、Keywords 与 favicon" },
    { id: 8, group: "站点分析", name: "站点信息", icon: "globe", key: "host", placeholder: "www.example.com", desc: "采集标题、DNS、IP 属地、证书、备案与 Whois" },
    { id: 2, group: "网络", name: "IP 属地", icon: "zap", key: "ip", placeholder: "223.5.5.5", desc: "查询 IP 的地理位置与运营商" },
    { id: 3, group: "网络", name: "DNS 查询", icon: "server", key: "host", placeholder: "www.example.com", desc: "使用本地与公共 DNS 解析域名，对比解析结果与耗时" },
    { id: 6, group: "网络", name: "在线 Ping", icon: "activity", key: "ip", placeholder: "223.5.5.5", desc: "连续发起 4 次 ICMP 探测" },
    { id: 4, group: "域名与证书", name: "Whois", icon: "search", key: "host", placeholder: "example.com", desc: "查询域名注册信息与 IANA 根服务器响应" },
    { id: 5, group: "域名与证书", name: "ICP 备案", icon: "book", key: "host", placeholder: "example.com", desc: "查询网站备案主体与备案号" },
    { id: 7, group: "域名与证书", name: "SSL 证书", icon: "shield", key: "url", placeholder: "https://www.example.com", desc: "获取证书有效期、颁发者、算法与备用域名" },
];

const API = {
    1: function (v) { return "/api/tool/website/tdki?url=" + encodeURIComponent(v); },
    2: function (v) { return "/api/tool/ip?ip=" + encodeURIComponent(v); },
    3: function (v) { return "/api/tool/nsLookUp/all?host=" + encodeURIComponent(v); },
    4: function (v) { return "/api/tool/whois?host=" + encodeURIComponent(v); },
    5: function (v) { return "/api/tool/icp?host=" + encodeURIComponent(v); },
    6: function (v) { return "/api/tool/ping?ip=" + encodeURIComponent(v); },
    7: function (v) { return "/api/tool/certificate?url=" + encodeURIComponent(v); },
    8: function (v) { return "/api/tool/website/collectInfo?host=" + encodeURIComponent(v); },
};

const app = createApp({
    data() {
        return {
            mail: Mail,
            tools: TOOLS,
            active: 1,
            params: { 1: "", 2: "", 3: "", 4: "", 5: "", 6: "", 7: "", 8: "" },
            result: {},
            loading: false,
            error: "",
            historyList: [],
        }
    },

    computed: {
        current() {
            return this.tools.find(function (t) { return t.id === this.active; }, this) || this.tools[0];
        },
        groups() {
            const names = [];
            const map = {};
            this.tools.forEach(function (t) {
                if (!map[t.group]) { map[t.group] = []; names.push(t.group); }
                map[t.group].push(t);
            });
            return names.map(function (n) { return { name: n, tools: map[n] }; });
        },
    },

    created: function () {
        this.mail.getInfo();
        this.loadHistory(this.active);
    },

    methods: {
        select: function (id) {
            this.active = id;
            this.error = "";
            this.loading = false;
            this.loadHistory(id);
        },

        submit: function () {
            const self = this;
            const id = this.active;
            const value = (this.params[id] || "").trim();
            if (!value) { UI.toast("请输入查询内容"); return; }
            if (this.loading) return;

            this.loading = true;
            this.error = "";
            common.AjaxGet(API[id](value), function (data) {
                self.loading = false;
                if (data.code !== 0) {
                    self.error = data.msg || "查询失败";
                    return;
                }
                self.result[id] = data.data;
                self.setHistory(id, value);
            });
        },

        hasResult: function () {
            const r = this.result[this.active];
            return r !== undefined && r !== null && r !== "";
        },

        setHistory: function (id, value) {
            const self = this;
            common.AjaxPost("/api/tool/history?toolID=" + id + "&value=" + encodeURIComponent(value), "", function () {
                self.loadHistory(id);
            });
        },

        loadHistory: function (id) {
            const self = this;
            common.AjaxGet("/api/tool/history?toolID=" + id, function (data) {
                if (data.code === 0) self.historyList = data.data || [];
            });
        },

        pick: function (value) {
            this.params[this.active] = value;
            this.submit();
        },

        clearHistory: function () {
            const self = this;
            common.AjaxGet("/api/tool/history/clear?toolID=" + this.active, function () {
                self.historyList = [];
                UI.toast("已清空历史");
            });
        },

        ipText: function (d) {
            if (!d) return "—";
            if (d.ips && d.ips.length) return d.ips.join(", ");
            return "—";
        },

        ms: function (v) {
            return Math.round(Number(v) || 0);
        },
    },
});
app.mount('#app');
