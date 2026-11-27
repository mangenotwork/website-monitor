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
                total: 0,
                unread: 0,
                counts: { all: 0, unread: 0, timeout: 0, code: 0 },
                hosts: [],
                page: 1,
                size: 100,
                hasMore: true,
                loading: false,
            },
            seg: segs.indexOf(segParam) !== -1 ? segParam : "all",
            kw: "",
            hostFilter: params.get("host") || "",
            reqSeq: 0,
            kwTimer: null,
        }
    },

    computed: {
        unread() { return this.alerts.unread; },
        hostCount() { return this.alerts.hosts.length; },
        hosts() { return this.alerts.hosts; },
        counts() { return this.alerts.counts; },
        // 服务端已经完成筛选与排序，这里只做别名，模板无需改动
        filtered() { return this.alerts.list; },
    },

    watch: {
        kw: function () {
            const self = this;
            clearTimeout(this.kwTimer);
            this.kwTimer = setTimeout(function () { self.getList(true); }, 300);
        },
    },

    created: function () {
        this.mail.getInfo();
        this.getList(true);
    },

    methods: {
        // reset=true 重新从第 1 页加载；false 追加下一页（下拉到底触发）
        getList: function (reset) {
            const t = this;

            if (reset) {
                t.reqSeq++;              // 让在途的旧请求失效
                t.alerts.page = 1;
                t.alerts.hasMore = true;
                t.alerts.loading = false;
            } else if (t.alerts.loading || !t.alerts.hasMore) {
                return;
            }

            t.alerts.loading = true;
            const seq = ++t.reqSeq;

            const q = new URLSearchParams();
            q.set("page", String(t.alerts.page));
            q.set("size", String(t.alerts.size));
            if (t.seg !== "all") q.set("seg", t.seg);
            if (t.hostFilter) q.set("host", t.hostFilter);
            if (t.kw) q.set("kw", t.kw);

            common.AjaxGet(t.alerts.api + "?" + q.toString(), function (data) {
                if (seq !== t.reqSeq) return;
                t.alerts.loading = false;
                if (data.code !== 0) { UI.toast(data.msg || ""); return; }

                const d = data.data || {};
                const list = d.list || [];
                t.alerts.list = reset ? list : t.alerts.list.concat(list);
                t.alerts.total = d.count || 0;
                t.alerts.unread = d.unread || 0;
                t.alerts.counts = d.counts || { all: 0, unread: 0, timeout: 0, code: 0 };
                t.alerts.hosts = d.hosts || [];
                t.alerts.hasMore = !!d.hasMore;
                if (t.alerts.hasMore) t.alerts.page += 1;
                UI.setAlertBadge(t.alerts.unread);
            });
        },

        loadMore: function () {
            this.getList(false);
        },

        onScroll: function (e) {
            const el = e.target;
            if (!el || this.alerts.loading || !this.alerts.hasMore) return;
            if (el.scrollTop + el.clientHeight >= el.scrollHeight - 80) {
                this.getList(false);
            }
        },

        reload: function () {
            this.getList(true);
        },

        open: function (item) {
            const self = this;
            const wasUnread = item.read === 0;
            this.alertDetail.open(item, function (reason) {
                if (reason === 'del' || self.seg === 'unread') {
                    self.getList(true);
                    return;
                }
                // 其余筛选下就地更新，避免标记已读后跳回第一页
                if (reason === 'read' && wasUnread) {
                    self.alerts.unread = Math.max(0, self.alerts.unread - 1);
                    self.alerts.counts.unread = self.alerts.unread;
                    UI.setAlertBadge(self.alerts.unread);
                }
            });
        },

        setSeg: function (seg) {
            if (this.seg === seg) return;
            this.seg = seg;
            this.syncUrl();
            this.getList(true);
        },

        onHostChange: function () {
            this.syncUrl();
            this.getList(true);
        },

        syncUrl: function () {
            const u = new URL(location.href);
            if (this.seg === "all") u.searchParams.delete("seg");
            else u.searchParams.set("seg", this.seg);
            if (this.hostFilter) u.searchParams.set("host", this.hostFilter);
            else u.searchParams.delete("host");
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

        markAllRead: function () {
            const self = this;
            if (this.alerts.unread === 0) { UI.toast("没有未读报警"); return; }
            UI.confirm("全部标记为已读", "共 " + this.alerts.unread + " 条未读报警，将全部标记为已读。", function () {
                common.AjaxGet("/api/alert/all/read", function (data) {
                    UI.toast(data.msg || "已全部标记为已读");
                    self.getList(true);
                });
            });
        },
    },
});
app.mount('#app');
