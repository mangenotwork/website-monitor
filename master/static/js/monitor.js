const { createApp } = Vue;
const common = new Utils;

const app = createApp({
    data() {
        return {
            mail: Mail,
            monitor: {
                api: "/api/monitor/list?detail=1",
                list: [],
            },
        }
    },

    computed: {
        online() { return this.monitor.list.filter(function (i) { return i.online === true; }).length; },
        total() { return this.monitor.list.length; },
    },

    created: function () {
        const t = this;
        t.mail.getInfo();
        t.getList();
        t.timer = window.setInterval(function () { t.getList(); }, 60000);
    },

    destroyed: function () {
        window.clearInterval(this.timer);
    },

    methods: {
        getList: function () {
            const t = this;
            common.AjaxGet(t.monitor.api, function (data) {
                if (data.code !== 0) return;
                t.monitor.list = data.data || [];
                UI.setProbeStatus(t.online, t.total);
            });
        },

        refresh: function () {
            this.getList();
            UI.toast("已刷新探针状态");
        },

        hostName: function (item) {
            if (item.osInfo && item.osInfo.hostName) return item.osInfo.hostName;
            return item.name || item.key || "—";
        },

        osLine: function (item) {
            if (!item.osInfo) return "—";
            const o = item.osInfo;
            const parts = [];
            if (o.osType || o.osArch) {
                const os = o.osType ? (o.osType.charAt(0).toUpperCase() + o.osType.slice(1)) : "";
                parts.push((os + " " + (o.osArch || "")).trim());
            }
            if (o.cpuCoreNumber) parts.push(o.cpuCoreNumber);
            return parts.join(" · ") || "—";
        },

        addrLine: function (item) {
            if (item.publicIP && item.publicIPAddr) return item.publicIPAddr;
            return item.addr || "—";
        },

        lastSeen: function (item) {
            return UI.ago(item.lastTime);
        },

        copy: function (text) {
            UI.copy(text);
        },
    },
});
app.mount('#app');
