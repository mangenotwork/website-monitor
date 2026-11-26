const { createApp } = Vue;
const common = new Utils;

const app = createApp({
    data() {
        const params = new URLSearchParams(location.search);
        const segParam = params.get("seg");
        const segs = ["all", "unread", "timeout", "code"];
        return {
            mail: Mail,
            alertDetail: newAlertDetail(),
            alerts: {
                api: "/api/alert/list",
                read: function (id) { return "/api/alert/read/" + id; },
                del: function (id) { return "/api/alert/del/" + id; },
                list: [],
            },
            seg: segs.indexOf(segParam) !== -1 ? segParam : "all",
            kw: "",
            hostFilter: params.get("host") || "",
        }
    },

    computed: {
        unread() { return this.alerts.list.filter(function (i) { return i.read === 0; }).length; },
        hostCount() {
            const s = {};
            this.alerts.list.forEach(function (i) { s[i.hostId] = i.host; });
            return Object.keys(s).length;
        },
        hosts() {
            const map = {};
            this.alerts.list.forEach(function (i) { map[i.hostId] = i.host; });
            return Object.keys(map).map(function (k) { return { id: k, host: map[k] }; });
        },
        counts() {
            const list = this.alerts.list;
            return {
                all: list.length,
                unread: list.filter(function (i) { return i.read === 0; }).length,
                timeout: list.filter(function (i) { return i.alertType === "timeout"; }).length,
                code: list.filter(function (i) { return i.alertType === "code"; }).length,
            };
        },
        filtered() {
            const seg = this.seg;
            const host = this.hostFilter;
            const kw = (this.kw || "").toLowerCase();
            const list = this.alerts.list.filter(function (i) {
                if (seg === "unread" && i.read !== 0) return false;
                if (seg === "timeout" && i.alertType !== "timeout") return false;
                if (seg === "code" && i.alertType !== "code") return false;
                if (host && i.hostId !== host) return false;
                if (kw) {
                    const text = ((i.host || "") + " " + (i.uri || "") + " " + (i.msg || "")).toLowerCase();
                    if (text.indexOf(kw) === -1) return false;
                }
                return true;
            });
            list.sort(function (a, b) { return String(b.date).localeCompare(String(a.date)); });
            return list;
        },
    },

    created: function () {
        this.mail.getInfo();
        this.getList();
    },

    methods: {
        getList: function () {
            const t = this;
            common.AjaxGet(t.alerts.api, function (data) {
                if (data.code !== 0) { UI.toast(data.msg || ""); return; }
                t.alerts.list = (data.data && data.data.list) ? data.data.list : [];
                UI.setAlertBadge(t.unread);
            });
        },

        open: function (item) {
            const self = this;
            this.alertDetail.open(item, function () { self.getList(); });
        },

        setSeg: function (seg) {
            this.seg = seg;
            const u = new URL(location.href);
            if (seg === "all") u.searchParams.delete("seg");
            else u.searchParams.set("seg", seg);
            try { history.replaceState(null, "", u); } catch (e) { }
        },

        typeClass: function (item) {
            const t = this.typeText(item);
            if (t === "证书" || t === "慢响应") return "warn";
            return "danger";
        },

        typeText: function (item) {
            const msg = item.msg || "";
            if (msg.indexOf("证书") !== -1) return "证书";
            if (msg.indexOf("响应") !== -1 && msg.indexOf("超过") !== -1) return "慢响应";
            const map = { timeout: "超时", code: "状态码", err: "请求失败" };
            return map[item.alertType] || "异常";
        },

        timeText: function (date) {
            if (!date) return "—";
            const s = String(date);
            return s.length >= 19 ? s.slice(5, 19) : s;
        },

        shortHost: function (text) {
            return UI.shortHost(text);
        },

        markAllRead: async function () {
            const self = this;
            const ids = this.alerts.list.filter(function (i) { return i.read === 0; }).map(function (i) { return i.alertId; });
            if (ids.length === 0) { UI.toast("没有未读报警"); return; }
            UI.confirm("全部标记为已读", "共 " + ids.length + " 条报警，将逐条标记为已读。", async function () {
                for (let i = 0; i < ids.length; i++) {
                    try { await fetch("/api/alert/read/" + ids[i]); } catch (e) { }
                }
                self.getList();
                UI.toast("已全部标记为已读");
            });
        },
    },
});
app.mount('#app');
