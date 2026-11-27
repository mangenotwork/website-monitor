const { createApp } = Vue;
const common = new Utils;

const app = createApp({
    data() {
        return {
            mail: Mail,
            webSite: AddWebSite,
            site: newSiteDetail(),
            alertDetail: newAlertDetail(),
            alertList: {
                api: "/api/alert/list?size=5",
                list: [],
                unread: 0,
            },
            monitor: {
                api: "/api/monitor/list",
                list: [],
            },
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
        siteTotal() { return this.webSite.list.length; },
        siteMonitoring() { return this.webSite.list.filter(function (i) { return i.state === 1; }).length; },
        siteAlertSites() { return this.webSite.list.filter(function (i) { return i.state === 2; }).length; },
        siteWaiting() { return this.webSite.list.filter(function (i) { return i.state === 0; }).length; },
        unreadAlerts() { return this.alertList.unread; },
        recentAlerts() {
            return this.alertList.list.slice().sort(function (a, b) {
                return String(b.date).localeCompare(String(a.date));
            });
        },
        probeOnline() { return this.monitor.list.filter(function (i) { return i.online === true; }).length; },
        probeTotal() { return this.monitor.list.length; },
        lastAlertTime() {
            if (this.alertList.list.length === 0) return "—";
            return this.hm(this.alertList.list[0].date || "");
        },
        monitorRate() {
            if (this.siteTotal === 0) return 0;
            return Math.round(this.siteMonitoring / this.siteTotal * 100);
        },
    },

    created: function () {
        const t = this;
        t.getList();
        t.getAlertList();
        t.getMonitorList();
        t.mail.hasSet();
        t.mail.getInfo();

        t.timer = window.setInterval(function () {
            t.getList();
            t.getAlertList();
        }, 30000);

        t.monitorTimer = window.setInterval(function () {
            t.getMonitorList();
        }, 60000);
    },

    destroyed: function () {
        window.clearInterval(this.timer);
        window.clearInterval(this.monitorTimer);
    },

    methods: {
        getList: function () {
            this.webSite.getList();
        },

        getAlertList: function () {
            const t = this;
            common.AjaxGet(t.alertList.api, function (data) {
                if (data.code !== 0) return;
                t.alertList.list = (data.data && data.data.list) ? data.data.list : [];
                t.alertList.unread = (data.data && data.data.unread) || 0;
                UI.setAlertBadge(t.unreadAlerts);
            });
        },

        getMonitorList: function () {
            const t = this;
            common.AjaxGet(t.monitor.api, function (data) {
                if (data.code !== 0) return;
                t.monitor.list = data.data || [];
                UI.setProbeStatus(t.probeOnline, t.probeTotal);
            });
        },

        refreshAll: function () {
            this.getList();
            this.getAlertList();
            this.getMonitorList();
            UI.toast("已刷新");
        },

        openSite: function (item) {
            this.site.openAt(item.hostID, item.host);
        },

        openSitePoints: function (item) {
            this.site.openAt(item.hostID, item.host);
            this.site.switchTab("points");
        },

        openAlertItem: function (item) {
            const self = this;
            this.alertDetail.open(item, function () { self.getAlertList(); });
        },

        alertsOfSite: function (hostId) {
            location.href = "/alert?host=" + hostId;
        },

        openEditWebsiteConf: function (hostId) {
            const t = this;
            t.editWebsiteConf.hostId = hostId;
            common.AjaxGet(t.editWebsiteConf.confApi(), function (data) {
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
                    t.getList();
                    if (t.site.info) t.site.loadInfo();
                }
            });
        },

        gotoWebsite: function (item) {
            window.open(item.host, "_blank");
        },

        deleteWebsite: function (item) {
            const t = this;
            UI.confirm("删除站点监测", "将删除 " + item.host + " 的全部监测配置与监测点，此操作不可撤销。", function () {
                common.AjaxGet("/api/website/delete/" + item.hostID, function (data) {
                    UI.toast(data.msg || "");
                    if (data.code === 0) setTimeout(function () { location.reload(); }, 400);
                });
            }, true);
        },

        hm: function (date) {
            if (!date) return "—";
            return String(date).length >= 16 ? String(date).slice(11, 16) : String(date);
        },

        sevClass: function (item) {
            const msg = item.msg || "";
            if (msg.indexOf("证书") !== -1) return "warn";
            if (msg.indexOf("响应") !== -1 && msg.indexOf("超过") !== -1) return "warn";
            return (item.alertType === "timeout" || item.alertType === "code" || item.alertType === "err") ? "" : "warn";
        },

        osText: function (item) {
            if (!item || !item.osInfo) return "—";
            const o = item.osInfo;
            const parts = [];
            if (o.osType || o.osArch) parts.push((o.osType || "") + " " + (o.osArch || ""));
            if (o.cpuCoreNumber) parts.push(o.cpuCoreNumber);
            return parts.join(" · ") || "—";
        },

        probeTag: function (item) {
            if (!item.online) return "离线";
            if (item.osInfo && item.osInfo.osType) {
                const t = item.osInfo.osType;
                return t.charAt(0).toUpperCase() + t.slice(1);
            }
            return "在线";
        },

        probeAddr: function (item) {
            const ip = item.publicIP || item.IP || "";
            const addr = item.publicIPAddr || item.addr || "";
            if (ip && addr) return ip + " · " + addr;
            return ip || addr || "—";
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
