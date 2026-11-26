const { createApp } = Vue;
const common = new Utils;

const app = createApp({
    data() {
        const params = new URLSearchParams(location.search);
        const segParam = params.get("seg");
        const segs = ["all", "monitoring", "alert", "waiting"];
        return {
            mail: Mail,
            webSite: AddWebSite,
            site: newSiteDetail(),
            kw: params.get("kw") || "",
            seg: segs.indexOf(segParam) !== -1 ? segParam : "all",
            selected: [],
            busy: false,
            editWebsiteConf: {
                hostId: "",
                confApi: function () { return "/api/website/conf/" + this.hostId; },
                editApi: function () { return "/api/website/edit/" + this.hostId; },
                base: {},
                alarmRule: {},
                scanCheckUp: {},
            },
        }
    },

    computed: {
        counts() {
            const list = this.webSite.list;
            return {
                all: list.length,
                monitoring: list.filter(function (i) { return i.state === 1; }).length,
                alert: list.filter(function (i) { return i.state === 2; }).length,
                waiting: list.filter(function (i) { return i.state === 0; }).length,
            };
        },
        filtered() {
            const seg = this.seg;
            const kw = (this.kw || "").toLowerCase();
            return this.webSite.list.filter(function (i) {
                if (seg === "monitoring" && i.state !== 1) return false;
                if (seg === "alert" && i.state !== 2) return false;
                if (seg === "waiting" && i.state !== 0) return false;
                if (kw) {
                    const text = ((i.host || "") + " " + (i.notes || "")).toLowerCase();
                    if (text.indexOf(kw) === -1) return false;
                }
                return true;
            });
        },
        selectedCount() { return this.selected.length; },
        allSelected() {
            if (this.filtered.length === 0) return false;
            const self = this;
            return this.filtered.every(function (i) { return self.selected.indexOf(i.hostID) !== -1; });
        },
    },

    created: function () {
        this.mail.getInfo();
        this.webSite.getList();
    },

    methods: {
        setSeg: function (seg) {
            this.seg = seg;
            const u = new URL(location.href);
            if (seg === "all") u.searchParams.delete("seg");
            else u.searchParams.set("seg", seg);
            try { history.replaceState(null, "", u); } catch (e) { }
        },

        statusText: function (item) {
            if (item.state === 2) return "报警 " + item.alertCount;
            if (item.state === 1) return "正常";
            return "等待探针";
        },

        statusClass: function (item) {
            if (item.state === 2) return "danger";
            if (item.state === 1) return "ok";
            return "wait";
        },

        openSite: function (item) {
            this.site.openAt(item.hostID, item.host);
        },

        toggle: function (item) {
            const i = this.selected.indexOf(item.hostID);
            if (i === -1) this.selected.push(item.hostID);
            else this.selected.splice(i, 1);
        },

        toggleAll: function (e) {
            if (e.target.checked) {
                this.selected = this.filtered.map(function (i) { return i.hostID; });
            } else {
                this.selected = [];
            }
        },

        clearSel: function () {
            this.selected = [];
        },

        bulkDelete: function () {
            const self = this;
            const ids = this.selected.slice();
            if (ids.length === 0) return;
            UI.confirm("批量删除监测", "将删除选中的 " + ids.length + " 个站点及其全部监测配置，此操作不可撤销。", async function () {
                self.busy = true;
                for (let i = 0; i < ids.length; i++) {
                    try { await fetch("/api/website/delete/" + ids[i]); } catch (e) { }
                }
                self.busy = false;
                UI.toast("已删除 " + ids.length + " 个站点");
                setTimeout(function () { location.reload(); }, 500);
            }, true);
        },

        bulkPause: function () {
            UI.toast("批量暂停需要后端提供暂停接口，暂未开放");
        },

        exportCsv: function () {
            const rows = [["站点", "备注", "状态", "频率(s)", "对照组", "Ping", "报警数"]];
            this.filtered.forEach(function (i) {
                rows.push([i.host, i.notes || "", i.state === 2 ? "报警" : (i.state === 1 ? "正常" : "等待探针"),
                    i.monitorRate, i.contrastUrl, i.ping, i.alertCount]);
            });
            const csv = rows.map(function (r) {
                return r.map(function (v) {
                    const s = String(v == null ? "" : v);
                    return '"' + s.replace(/"/g, '""') + '"';
                }).join(",");
            }).join("\r\n");
            const blob = new Blob(["\ufeff" + csv], { type: "text/csv;charset=utf-8" });
            const a = document.createElement("a");
            a.href = URL.createObjectURL(blob);
            a.download = "sites-" + new Date().toISOString().slice(0, 10) + ".csv";
            a.click();
            URL.revokeObjectURL(a.href);
            UI.toast("已导出 " + this.filtered.length + " 条");
        },

        openEditWebsiteConf: function (hostId) {
            const t = this;
            t.editWebsiteConf.hostId = hostId;
            common.AjaxGetNotAsync(t.editWebsiteConf.confApi(), function (data) {
                if (data.code !== 0) { UI.toast(data.msg || "获取配置失败"); return; }
                t.editWebsiteConf.base = data.data.base;
                t.editWebsiteConf.alarmRule = data.data.alarmRule;
                t.editWebsiteConf.scanCheckUp = data.data.scanCheckUp;
                bootstrap.Modal.getOrCreateInstance(document.getElementById("setAlertModal")).show();
            });
        },

        editWebsiteConfSubmit: function () {
            const t = this;
            const param = {
                host: t.editWebsiteConf.base.host,
                notes: t.editWebsiteConf.base.notes,
                monitorRate: Number(t.editWebsiteConf.base.monitorRate),
                contrastUrl: t.editWebsiteConf.base.contrastUrl,
                contrastTime: Number(t.editWebsiteConf.base.contrastTime),
                ping: t.editWebsiteConf.base.ping,
                pingTime: Number(t.editWebsiteConf.base.pingTime),
                uriDepth: Number(t.editWebsiteConf.scanCheckUp.uriDepth),
                scanRate: Number(t.editWebsiteConf.scanCheckUp.scanRate),
                scanBadLink: t.editWebsiteConf.scanCheckUp.scanBadLink,
                scanTDK: t.editWebsiteConf.scanCheckUp.scanTDK,
                scanExtLinks: t.editWebsiteConf.scanCheckUp.scanExtLinks,
                websiteSlowResponseTime: Number(t.editWebsiteConf.alarmRule.websiteSlowResponseTime),
                websiteSlowResponseCount: Number(t.editWebsiteConf.alarmRule.websiteSlowResponseCount),
                SSLCertificateExpire: Number(t.editWebsiteConf.alarmRule.SSLCertificateExpire),
            };
            common.AjaxPost(t.editWebsiteConf.editApi(), param, function (data) {
                UI.toast(data.msg || "");
                if (data.code === 0) {
                    bootstrap.Modal.getOrCreateInstance(document.getElementById("setAlertModal")).hide();
                    t.webSite.getList();
                    if (t.site.info) t.site.loadInfo();
                }
            });
        },

        copy: function (text) {
            UI.copy(text);
        },

        shortHost: function (text) {
            return UI.shortHost(text);
        },
    },
});
app.mount('#app');
