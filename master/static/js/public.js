/* ==========================================================================
   Website Monitor · 公共脚本
   Utils：Ajax 与提示
   UI：主题、抽屉、确认框、徽标
   SiteDetail / AlertDetail：站点与报警详情抽屉状态
   AddWebSite：创建监测表单
   ========================================================================== */

class Utils {
    AjaxGet(url, func) {
        $.ajax({
            type: "get", url: url, data: "", dataType: 'json', async: true,
            success: function (data) { func(data); },
            error: function (xhr, textStatus) { console.log(xhr, textStatus); }
        });
    }

    AjaxPost(url, param, func) {
        $.ajax({
            type: "post", url: url,
            data: param === "" ? "" : JSON.stringify(param),
            dataType: 'json', async: true,
            success: function (data) { func(data); },
            error: function (xhr, textStatus) { console.log(xhr, textStatus); }
        });
    }

    ToastShow(msg) {
        UI.toast(msg);
    }
}

/* -------------------------------------------------------------------------- */

const UI = {
    toastTimer: null,

    toast(msg) {
        const t = document.getElementById('liveToast');
        if (!t) return;
        const m = document.getElementById('liveToastMsg');
        if (m) m.textContent = msg == null ? '' : String(msg);
        t.classList.add('show');
        clearTimeout(this.toastTimer);
        this.toastTimer = setTimeout(function () { t.classList.remove('show'); }, 2400);
    },

    copy(text) {
        if (navigator.clipboard) {
            navigator.clipboard.writeText(text).catch(function () { });
        }
        this.toast('已复制 ' + text);
    },

    setAccent(name) {
        document.documentElement.setAttribute('data-accent', name);
        try { localStorage.setItem('wm.accent', name); } catch (e) { }
        document.querySelectorAll('.accent-switch button').forEach(function (b) {
            b.classList.toggle('on', b.getAttribute('data-accent') === name);
        });
    },

    toggleDensity() {
        const root = document.documentElement;
        const next = root.getAttribute('data-density') === 'compact' ? 'cozy' : 'compact';
        this.setDensity(next);
        this.toast(next === 'compact' ? '已切换为紧凑模式' : '已切换为舒适模式');
    },

    setDensity(mode) {
        const m = mode === 'compact' ? 'compact' : 'cozy';
        document.documentElement.setAttribute('data-density', m);
        try { localStorage.setItem('wm.density', m); } catch (e) { }
        const btn = document.getElementById('densityBtn');
        if (!btn) return;
        const compact = m === 'compact';
        btn.title = compact ? '当前：紧凑 · 点击切换为舒适' : '当前：舒适 · 点击切换为紧凑';
        btn.setAttribute('aria-pressed', compact ? 'true' : 'false');
        btn.classList.toggle('on', compact);
    },

    setAlertBadge(n) {
        const dot = document.getElementById('topAlertDot');
        const nav = document.getElementById('navAlertCount');
        if (dot) dot.style.display = n > 0 ? '' : 'none';
        if (nav) {
            nav.textContent = n > 99 ? '99+' : String(n);
            nav.style.display = n > 0 ? '' : 'none';
        }
    },

    setSiteCount(n) {
        const el = document.getElementById('navSiteCount');
        if (!el) return;
        el.textContent = String(n);
        el.style.display = n > 0 ? '' : 'none';
    },

    setProbeStatus(online, total) {
        const el = document.getElementById('topProbeText');
        if (!el) return;
        el.textContent = total > 0 ? ('探针 ' + online + '/' + total + ' 在线') : '探针管理';
    },

    openDrawer(id) {
        this.closeDrawers();
        const el = document.getElementById(id);
        if (!el) return;
        el.classList.add('open');
        const mask = document.getElementById('uiMask');
        if (mask) mask.classList.add('open');
        document.body.style.overflow = 'hidden';
    },

    closeDrawers() {
        document.querySelectorAll('.drawer').forEach(function (d) { d.classList.remove('open'); });
        const mask = document.getElementById('uiMask');
        if (mask) mask.classList.remove('open');
        document.body.style.overflow = '';
    },

    confirm(title, desc, onOk, danger) {
        const mask = document.getElementById('uiConfirmMask');
        if (!mask) { if (window.confirm(title + '\n' + desc)) onOk(); return; }
        document.getElementById('uiConfirmTitle').textContent = title;
        document.getElementById('uiConfirmDesc').textContent = desc;
        const ok = document.getElementById('uiConfirmOk');
        ok.className = 'btn primary';
        ok.style.background = danger ? 'var(--danger)' : '';
        ok.textContent = danger ? '确认删除' : '确认';
        ok.onclick = function () {
            UI.closeConfirm();
            onOk();
        };
        mask.classList.add('open');
    },

    closeConfirm() {
        const mask = document.getElementById('uiConfirmMask');
        if (mask) mask.classList.remove('open');
    },

    fmtTime(ts) {
        if (!ts) return '—';
        let n = Number(ts);
        if (!isFinite(n) || n <= 0) return '—';
        if (n < 1e12) n = n * 1000;
        const d = new Date(n);
        const p = function (x) { return String(x).padStart(2, '0'); };
        return d.getFullYear() + '-' + p(d.getMonth() + 1) + '-' + p(d.getDate()) + ' ' +
            p(d.getHours()) + ':' + p(d.getMinutes()) + ':' + p(d.getSeconds());
    },

    ago(ts) {
        if (!ts) return '—';
        let n = Number(ts);
        if (!isFinite(n) || n <= 0) return '—';
        if (n > 1e12) n = n / 1000;
        let diff = Math.floor(Date.now() / 1000) - n;
        if (diff < 0) diff = 0;
        if (diff < 60) return diff + 's 前';
        if (diff < 3600) return Math.floor(diff / 60) + 'm 前';
        if (diff < 86400) return Math.floor(diff / 3600) + 'h 前';
        return Math.floor(diff / 86400) + 'd 前';
    },

    daysLeft(ts) {
        if (!ts) return '—';
        let n = Number(ts);
        if (!isFinite(n) || n <= 0) return '—';
        if (n > 1e12) n = n / 1000;
        const d = Math.floor((n - Date.now() / 1000) / 86400);
        return d > 0 ? d + ' 天' : '已过期';
    },

    shortHost(url) {
        return String(url || '').replace(/^https?:\/\//i, '').replace(/\/+$/, '');
    }
};

window.UI = UI;

/* -------------------------------------------------------------------------- */

document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') { UI.closeDrawers(); UI.closeConfirm(); }
    if (e.key === '/' && !/input|textarea|select/i.test(document.activeElement.tagName)) {
        const input = document.querySelector('.top-search input');
        if (input) { e.preventDefault(); input.focus(); }
    }
});

document.addEventListener('DOMContentLoaded', function () {
    const accent = document.documentElement.getAttribute('data-accent') || 'indigo';
    document.querySelectorAll('.accent-switch button').forEach(function (b) {
        b.classList.toggle('on', b.getAttribute('data-accent') === accent);
    });
    UI.setDensity(document.documentElement.getAttribute('data-density') || 'cozy');
});

/* -------------------------------------------------------------------------- */

/*
 * ECharts 按需加载：图表库约 1MB，只在打开站点详情并需要绘图时加载，
 * 不参与首页首屏，避免阻塞 Vue 挂载与接口请求。
 */
let echartsLoading = null;

function loadECharts() {
    if (window.echarts) return Promise.resolve(window.echarts);
    if (!echartsLoading) {
        echartsLoading = new Promise(function (resolve, reject) {
            const s = document.createElement('script');
            s.src = '/static/js/echarts.min.js?v=5.4.3';
            s.async = true;
            s.onload = function () {
                if (window.echarts) resolve(window.echarts);
                else reject(new Error('ECharts 未就绪'));
            };
            s.onerror = function () {
                echartsLoading = null;
                reject(new Error('ECharts 加载失败'));
            };
            document.head.appendChild(s);
        });
    }
    return echartsLoading;
}

function newSiteDetail() {
    return {
        hostId: '', host: '', tab: 'overview',
        info: null, infoLoading: false,
        points: [], pointLoading: false, newPoint: '',
        urls: null, urlLoading: false, urlFilter: '',
        logs: [], logDays: [], logDay: '', logLoading: false,
        chart: null, chartData: [], chartLoading: false,
        stat: { count: 0, avg: 0, p95: 0, max: 0 },

        openAt: function (hostId, host) {
            this.hostId = hostId;
            this.host = host || '';
            this.tab = 'overview';
            this.info = null;
            this.points = [];
            this.urls = null;
            this.logs = [];
            this.logDays = [];
            this.logDay = '';
            this.chartData = [];
            this.stat = { count: 0, avg: 0, p95: 0, max: 0 };
            UI.openDrawer('siteDrawer');
            this.loadInfo();
            this.loadChart();
        },

        close: function () { UI.closeDrawers(); },

        copy: function (text) { UI.copy(text); },

        shortHost: function () { return UI.shortHost(this.host); },

        fmtDay: function (day) {
            const s = String(day || '');
            if (s.length === 8) return s.slice(0, 4) + '-' + s.slice(4, 6) + '-' + s.slice(6, 8);
            return s || '—';
        },

        loadInfo: function () {
            const self = this;
            if (!this.hostId) return;
            this.infoLoading = true;
            new Utils().AjaxGet('/api/website/info/' + this.hostId, function (data) {
                self.infoLoading = false;
                if (data.code === 0) {
                    self.info = data.data;
                    self.host = data.data.base.host;
                } else {
                    UI.toast(data.msg || '获取站点信息失败');
                }
            });
        },

        refreshInfo: function () {
            const self = this;
            if (!this.hostId || !this.info) return;
            UI.toast('正在重新采集站点信息，请稍候…');
            new Utils().AjaxGet('/api/website/info/refresh?host=' + encodeURIComponent(this.host) + '&id=' + this.hostId, function (data) {
                if (data.code === 0) { self.loadInfo(); UI.toast('采集完成'); }
                else { UI.toast(data.msg || '采集失败'); }
            });
        },

        switchTab: function (tab) {
            this.tab = tab;
            if (tab === 'points') this.loadPoints();
            if (tab === 'urls') this.loadUrls();
            if (tab === 'logs') this.loadLogs();
            if (tab === 'overview' && this.chart) {
                const self = this;
                setTimeout(function () { self.chart.resize(); }, 30);
            }
        },

        loadPoints: function () {
            const self = this;
            this.pointLoading = true;
            new Utils().AjaxGet('/api/website/point/list/' + this.hostId, function (data) {
                self.pointLoading = false;
                if (data.code === 0) {
                    self.points = (data.data && data.data.url) ? data.data.url : [];
                } else {
                    self.points = [];
                }
            });
        },

        addPoint: function () {
            const self = this;
            const uri = (this.newPoint || '').trim();
            if (!uri) { UI.toast('请输入要监测的 URI'); return; }
            new Utils().AjaxPost('/api/website/point/add/' + this.hostId, { uri: uri }, function (data) {
                UI.toast(data.msg || '');
                if (data.code === 0) { self.newPoint = ''; self.loadPoints(); }
            });
        },

        delPoint: function (uri) {
            const self = this;
            UI.confirm('删除监测点', '将从轮询中移除：' + uri, function () {
                new Utils().AjaxPost('/api/website/point/del/' + self.hostId, { uri: uri }, function (data) {
                    UI.toast(data.msg || '');
                    if (data.code === 0) self.loadPoints();
                });
            }, true);
        },

        loadUrls: function () {
            const self = this;
            this.urlLoading = true;
            new Utils().AjaxGet('/api/website/urls/' + this.hostId, function (data) {
                self.urlLoading = false;
                self.urls = data.code === 0 ? data.data : null;
            });
        },

        loadLogs: function () {
            const self = this;
            const load = function () {
                self.logDay = self.logDay || self.logDays[0] || '';
                self.logLoading = true;
                new Utils().AjaxGet('/api/website/log/' + self.hostId + (self.logDay ? ('?day=' + self.logDay) : ''), function (data) {
                    self.logLoading = false;
                    self.logs = data.code === 0 ? (data.data || []) : [];
                });
            };

            if (this.logDays.length === 0) {
                new Utils().AjaxGet('/api/website/log/list/' + this.hostId, function (data) {
                    if (data.code === 0) self.logDays = data.data || [];
                    load();
                });
                return;
            }
            load();
        },

        changeLogDay: function (day) {
            this.logDay = day;
            this.loadLogs();
        },

        dnsText: function () {
            const d = this.info && this.info.info ? this.info.info.DNS : null;
            return d && d.ips && d.ips.length ? d.ips.join(', ') : '—';
        },

        cnameText: function () {
            const d = this.info && this.info.info ? this.info.info.DNS : null;
            return d && d.cname ? d.cname : '—';
        },

        sslText: function () {
            const s = this.info && this.info.info ? this.info.info.SSLCertificateInfo : null;
            if (!s || !s.notAfter) return '—';
            return UI.fmtTime(s.notAfter).slice(0, 10) + ' · ' + UI.daysLeft(s.notAfter);
        },

        ipcText: function () {
            const i = this.info && this.info.info ? this.info.info.IPC : null;
            if (!i) return '—';
            return i.ipc || i.company || '—';
        },

        serverText: function () {
            const i = this.info && this.info.info ? this.info.info : null;
            return i && i.server ? i.server : '—';
        },

        titleText: function () {
            const i = this.info && this.info.info ? this.info.info : null;
            return i && i.title ? i.title : '—';
        },

        filteredUrls: function () {
            if (!this.urls || !this.urls.allUri) return [];
            const kw = (this.urlFilter || '').toLowerCase();
            const list = kw ? this.urls.allUri.filter(function (u) { return u.toLowerCase().indexOf(kw) !== -1; }) : this.urls.allUri;
            return list.slice(0, 200);
        },

        urlTotal: function () {
            return this.urls && this.urls.allUri ? this.urls.allUri.length : 0;
        },

        loadChart: function () {
            const self = this;
            this.chartLoading = true;
            new Utils().AjaxGet('/api/website/log/list/' + this.hostId, function (data) {
                if (data.code === 0) self.logDays = data.data || [];
                self.logDay = self.logDays[0] || '';
                const day = self.logDay;
                new Utils().AjaxGet('/api/website/chart/' + self.hostId + '?day=' + day + '&uri=Health', function (res) {
                    self.chartLoading = false;
                    if (res.code !== 0) { self.chartData = []; return; }
                    const d = res.data;
                    if (Array.isArray(d)) { // 兼容旧接口
                        self.chartData = d;
                        self.calcStat();
                    } else {
                        self.chartData = (d && d.list) ? d.list : [];
                        if (d && d.stat) self.stat = d.stat;
                        else self.calcStat();
                    }
                    setTimeout(function () { self.renderChart(); }, 30);
                });
            });
        },

        calcStat: function () {
            const arr = this.chartData.map(function (v) { return v[1]; }).filter(function (v) { return typeof v === 'number'; });
            if (arr.length === 0) { this.stat = { count: 0, avg: 0, p95: 0, max: 0 }; return; }
            const sorted = arr.slice().sort(function (a, b) { return a - b; });
            let sum = 0;
            for (let i = 0; i < arr.length; i++) sum += arr[i];
            const idx = Math.min(sorted.length - 1, Math.floor(sorted.length * 0.95));
            this.stat = {
                count: arr.length,
                avg: Math.round(sum / arr.length),
                p95: Math.round(sorted[idx]),
                max: Math.round(sorted[sorted.length - 1])
            };
        },

        renderChart: function (retry) {
            if (this.chartData.length === 0) return;
            const el = document.getElementById('siteChart');
            if (!el) {
                // chart-box 由 v-if 控制，数据刚写入时 DOM 可能还没更新，短暂重试
                if ((retry || 0) < 20) {
                    const self = this;
                    setTimeout(function () { self.renderChart((retry || 0) + 1); }, 50);
                }
                return;
            }
            if (!window.echarts) {
                const self = this;
                loadECharts().then(function () {
                    self.renderChart();
                }).catch(function () {
                    UI.toast('图表组件加载失败，请检查网络后重试');
                });
                return;
            }
            const accent = getComputedStyle(document.documentElement).getPropertyValue('--accent').trim() || '#4F46E5';
            if (this.chart && this.chart.getDom() !== el) {
                this.chart.dispose();
                this.chart = null;
            }
            if (!this.chart) this.chart = echarts.init(el, '', { renderer: 'canvas' });
            const threshold = this.info && this.info.alarmRule ? Number(this.info.alarmRule.websiteSlowResponseTime) : 0;
            const series = {
                name: '响应时间',
                type: 'line',
                smooth: false,
                symbol: 'none',
                data: this.chartData,
                lineStyle: { width: 1.6, color: accent },
                areaStyle: { color: accent, opacity: .06 },
                emphasis: { disabled: true }
            };
            if (threshold > 0) {
                series.markLine = {
                    silent: true, symbol: 'none', animation: false,
                    lineStyle: { color: '#C9C9D1', type: 'dashed' },
                    label: { formatter: '阈值 ' + threshold + 'ms', color: '#93939E', fontSize: 10, position: 'insideEndTop' },
                    data: [{ yAxis: threshold }]
                };
            }
            this.chart.setOption({
                animation: false,
                grid: { left: 6, right: 14, top: 18, bottom: 22, containLabel: true },
                tooltip: {
                    trigger: 'axis',
                    backgroundColor: '#18181B', borderWidth: 0,
                    textStyle: { color: '#fff', fontSize: 11 },
                    formatter: function (p) {
                        const d = new Date(p[0].value[0]);
                        const pad = function (x) { return String(x).padStart(2, '0'); };
                        return pad(d.getHours()) + ':' + pad(d.getMinutes()) + ':' + pad(d.getSeconds()) +
                            '<br/>响应 ' + p[0].value[1] + ' ms';
                    }
                },
                xAxis: {
                    type: 'time', boundaryGap: false,
                    axisLine: { show: false }, axisTick: { show: false },
                    axisLabel: { color: '#93939E', fontSize: 10, hideOverlap: true },
                    splitLine: { show: false }
                },
                yAxis: {
                    type: 'value',
                    axisLine: { show: false }, axisTick: { show: false },
                    axisLabel: { color: '#93939E', fontSize: 10 },
                    splitLine: { lineStyle: { color: '#F0F0F2', type: 'dashed' } }
                },
                series: [series]
            }, true);
            if (!this._resizeBound) {
                this._resizeBound = true;
                const self = this;
                window.addEventListener('resize', function () { if (self.chart) self.chart.resize(); });
            }
        },

        del: function () {
            const self = this;
            UI.confirm('删除站点监测', '将删除 ' + (this.host || '') + ' 的全部监测配置、采集 URL 与监测点，此操作不可撤销。', function () {
                new Utils().AjaxGet('/api/website/delete/' + self.hostId, function (data) {
                    UI.toast(data.msg || '');
                    if (data.code === 0) {
                        UI.closeDrawers();
                        setTimeout(function () { location.reload(); }, 400);
                    }
                });
            }, true);
        }
    };
}

function newAlertDetail() {
    return {
        data: null,
        onChange: null,

        open: function (item, onChange) {
            this.data = item;
            this.onChange = onChange || null;
            UI.openDrawer('alertDrawer');
        },

        close: function () { UI.closeDrawers(); },

        copy: function () {
            if (!this.data) return;
            UI.copy(JSON.stringify(this.data, null, 2));
        },

        shortHost: function () {
            return this.data ? UI.shortHost(this.data.host) : '';
        },

        severity: function () {
            if (!this.data) return 'neutral';
            if (this.data.alertType === 'timeout' || this.data.alertType === 'code' || this.data.alertType === 'err') return 'danger';
            return 'warn';
        },

        severityText: function () {
            if (!this.data) return '';
            const m = {
                timeout: '超时', code: '状态码', err: '请求失败',
                ssl: '证书', slow: '慢响应'
            };
            return m[this.data.alertType] || (this.data.alertType || '异常');
        },

        read: function () {
            const self = this;
            if (!this.data) return;
            new Utils().AjaxGet('/api/alert/read/' + this.data.alertId, function (data) {
                UI.toast(data.msg || '');
                if (data.code === 0) {
                    self.data.read = 1;
                    if (self.onChange) self.onChange('read');
                }
            });
        },

        del: function () {
            const self = this;
            if (!this.data) return;
            UI.confirm('删除报警', '删除后不可恢复：' + (this.data.uri || this.data.host), function () {
                new Utils().AjaxGet('/api/alert/del/' + self.data.alertId, function (data) {
                    UI.toast(data.msg || '');
                    if (data.code === 0) {
                        UI.closeDrawers();
                        if (self.onChange) self.onChange('del');
                    }
                });
            }, true);
        },

        gotoUri: function () {
            if (this.data && this.data.uri) window.open(this.data.uri, '_blank');
        }
    };
}

/* -------------------------------------------------------------------------- */

let Mail = {
    init: { api: "/api/mail/init", data: false },
    param: {
        host: "smtp.qq.com",
        port: 25,
        from: "",
        authCode: "",
        toList: "",
        open: false,
    },
    infoApi: "/api/mail/info",
    setConfApi: "/api/mail/conf",
    sendApi: "/api/mail/sendTest",

    toListJoin: function (toList) {
        return Array.isArray(toList) ? toList.join(";") : toList;
    },

    hasSet: function () {
        const self = this;
        new Utils().AjaxGet(this.init.api, function (data) {
            self.init.data = data.data;
        });
    },

    getInfo: function () {
        const self = this;
        new Utils().AjaxGet(this.infoApi, function (data) {
            if (data.code === 0 && data.data) {
                self.param = data.data;
            }
        });
    },

    setConf: function () {
        const self = this;
        this.param.toList = this.toListJoin(this.param.toList);
        this.param.port = Number(this.param.port);
        new Utils().AjaxPost(this.setConfApi, this.param, function (data) {
            UI.toast(data.msg || '');
            if (data.code === 0) {
                const el = document.getElementById('mailSetModal');
                if (el) bootstrap.Modal.getOrCreateInstance(el).hide();
                self.getInfo();
            }
        });
    },

    send: function () {
        const self = this;
        this.param.toList = this.toListJoin(this.param.toList);
        this.param.port = Number(this.param.port);
        new Utils().AjaxPost(this.sendApi, this.param, function (data) {
            UI.toast(data.msg || '');
            if (data.code === 0) self.getInfo();
        });
    }
};

/* -------------------------------------------------------------------------- */

let AddWebSite = {
    api: "/api/website/add",
    apiGet: "/api/website/list",
    param: {},
    list: [],
    next1Param: {
        hostProtocol: "https://",
        hostUrl: "",
        notes: "",
        monitorRate: 15,
        contrastUrl: "https://www.baidu.com",
        contrastTime: 1000,
        ping: "223.5.5.5",
        pingTime: 1000,
    },
    next2Param: {
        uriDepth: 2,
        scanRate: 24,
        scanBadLink: true,
        scanTDK: false,
        scanExtLinks: false,
    },
    next3Param: {
        websiteSlowResponseTime: 3000,
        websiteSlowResponseCount: 3,
        SSLCertificateExpire: 14,
    },

    addSubmit: function () {
        const p1 = this.next1Param;
        const p2 = this.next2Param;
        const p3 = this.next3Param;

        if (!p1.hostUrl) { UI.toast('请输入站点地址'); return; }
        if (!p1.monitorRate || Number(p1.monitorRate) < 5) { UI.toast('监测频率不能小于 5 秒'); return; }
        if (!p1.contrastUrl) { UI.toast('对照组地址不能为空'); return; }
        if (!p1.ping) { UI.toast('Ping 地址不能为空'); return; }
        if (Number(p3.websiteSlowResponseTime) < 100) { UI.toast('慢响应阈值不能小于 100ms'); return; }
        if (Number(p3.websiteSlowResponseCount) < 1) { UI.toast('连续次数不能小于 1'); return; }
        if (Number(p3.SSLCertificateExpire) < 1) { UI.toast('证书阈值不能小于 1 天'); return; }

        this.param = {
            host: p1.hostProtocol + p1.hostUrl,
            notes: p1.notes,
            monitorRate: Number(p1.monitorRate),
            contrastUrl: p1.contrastUrl,
            contrastTime: Number(p1.contrastTime),
            ping: p1.ping,
            pingTime: Number(p1.pingTime),
            uriDepth: Number(p2.uriDepth),
            scanRate: Number(p2.scanRate),
            scanBadLink: p2.scanBadLink,
            scanTDK: p2.scanTDK,
            scanExtLinks: p2.scanExtLinks,
            websiteSlowResponseTime: Number(p3.websiteSlowResponseTime),
            websiteSlowResponseCount: Number(p3.websiteSlowResponseCount),
            SSLCertificateExpire: Number(p3.SSLCertificateExpire),
        };

        new Utils().AjaxPost(this.api, this.param, function (data) {
            UI.toast(data.msg || '');
            if (data.code === 0) {
                const el = document.getElementById('addHostModal');
                if (el) bootstrap.Modal.getOrCreateInstance(el).hide();
                AddWebSite.next1Param.hostUrl = '';
                AddWebSite.next1Param.notes = '';
                setTimeout(function () { location.reload(); }, 500);
            }
        });
    },

    getList: function () {
        const self = this;
        new Utils().AjaxGet(this.apiGet, function (data) {
            if (data.code === 0) {
                self.list = data.data;
                UI.setSiteCount(data.data.length);
            } else {
                UI.toast(data.msg || '');
            }
        });
    }
};
