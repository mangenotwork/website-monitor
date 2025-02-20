const { createApp, ref } = Vue;
const common = new Utils;

const app = createApp({
    data() {
        return {
            mail: Mail,
            webSite: AddWebSite,
            websiteList: {
                api: function (){ return "/api/website/list"; },
                page: 1,
                data: [],
            },
            deleteWebsite: {
                api: function (){ return "/api/website/delete/" + this.hostId; },
                hostId: "",
                hostName: "",
            },
            point: {
                hostId: "",
                hostUri: "",
                apiAdd: function (){ return "/api/website/point/add/"+this.hostId; },
                apiList: function () { return "/api/website/point/list/"+this.hostId; },
                apiDel: function () { return "/api/website/point/del/"+this.hostId; },
                apiClear: function () { return "/api/website/point/clear/"+this.hostId; },
                param: {
                    uri:"",
                },
                uriList: [],
                nowUri: "",
            },
            websiteInfo: {
                host: "",
                hostId: "",
                api: function() { return "/api/website/info/"+this.hostId; },
                refresh: function () { return "/api/website/info/refresh?host="+this.host+"&id="+this.hostId;},
                data: {
                    base: {
                        host: "",
                        hostID: "",
                        monitorRate: 0,
                        contrastUrl: "",
                        contrastTime: 0,
                        ping: "",
                        pingTime: 0,
                        notes: "",
                        date: ""
                    },
                    info:{
                        title: "",
                        description: "",
                        keywords: "",
                        icon: "",
                        DNS: {
                            ips: "",
                            cname: "",
                            isCDN: false,
                        },
                        IPAddr: [],
                        server: "",
                        contentEncoding: "",
                        contentLanguage: "",
                        SSLCertificateInfo: {
                            url: "",
                            effectiveTime: "",
                            dnsName: "",
                            ocspServer: "",
                            crlDistributionPoints: "",
                            issuer: "",
                            issuingCertificateURL: "",
                            publicKeyAlgorithm: "",
                            subject: "",
                            version: "",
                            signatureAlgorithm: ""
                        },
                        whois: {
                            root: "",
                            rse: ""
                        },
                        IPC: {
                            host: "",
                            company: "",
                            nature: "",
                            ipc: "",
                            websiteName: "",
                            websiteIndex: "",
                            auditDate: "",
                            restrictAccess: ""
                        }
                    },
                    alarmRule: {
                        websiteSlowResponseTime: 0,
                        websiteSlowResponseCount: 0,
                        SSLCertificateExpire: 0,
                        notTDK: false,
                        badLink: false,
                        extLinkChange: false
                    },
                    scanCheckUp: {
                        uriDepth: 0,
                        scanRate: 0,
                        scanExtLinks: false,
                        scanBadLink: false,
                        scanTDK: false
                    },
                }
            },
            alertList: {
                api: "/api/alert/list",
                list: [],
                len: 0,
            },
            monitorErrList: {
                api: "/api/monitor/err/list",
                clear: "/api/monitor/err/clear",
                list: [],
                len: 0,
            },
            isOk: "",
            monitorLog: {
                hostId: "",
                logApi: function (){ return "/api/website/log/" + this.hostId; },
                loadApi: function (){ return "/api/website/log/" + this.hostId + "?day=" + this.log; },
                logListApi: function () { return "/api/website/log/list/" + this.hostId; },
                logUpload: function () { return "/api/website/log/upload/" + this.hostId + "?day=" + this.log; },
                data: {},
                logList: [],
                log: "",
            },
            editWebsiteConf: {
                hostId: "",
                confApi: function () { return "/api/website/conf/" + this.hostId; },
                editApi: function () { return "/api/website/edit/" + this.hostId; },
                base: {},
                alarmRule: {},
                scanCheckUp: {},
            },
            chartData: {
                hostId: "",
                api: function () {
                  return "/api/website/chart/" + this.hostId + "?day=" + this.selectDay + "&uri=" + this.selectUriType;
                },
                dayApi: function () { return "/api/website/log/list/" + this.hostId; },
                list: [],
                host: "",
                selectDay: "",
                selectUriType: "",
                dayList: [],
                uriType: [{name:"健康监测",value:"Health"}, {name:"随机抽查",value:"Random"},{name:"指定监测点",value:"Point"}],
                uriTypeName: {"Health":"健康监测", "Random":"随机抽查", "Point":"指定监测点"},
            },
            websiteAlert: {
                hostId: "",
                api: function () { return "/api/alert/wbesite/" + this.hostId; },
                list: [],
                len: 0,
                del: function (alertId) { return "/api/alert/del/" + alertId; },
                clear: function () { return "/api/alert/clear/" + this.hostId; }
            },
            websiteUrl: {
                hostId: "",
                api: function () { return "/api/website/urls/" + this.hostId; },
                data: {},
            },
            monitor: {
                api: "/api/monitor/list",
                list: [],
            },
            objFidWebsiteInfo: {
                id: FidWebsiteInfo,
                host: "",
                api: function (){ return "/api/tool/website/collectInfo?host=" + this.host; },
                rse: {},
            },
            history: {
                add: function (id, value) { return "/api/tool/history?toolID="+id+"&value="+value; },
                get: function (id) { return "/api/tool/history?toolID="+id; },
                clear: function (id) {return "/api/tool/history/clear?toolID="+id; },
                data: [],
            },
            objFidNsLookUp: {
                id: FidNsLookUp,
                host: "",
                api: function (){ return "/api/tool/nsLookUp/all?host=" + this.host; },
                rse: {},
            },
            objFidWhois: {
                id: FidWhois,
                host: "",
                api: function (){ return "/api/tool/whois?host=" + this.host; },
                rse: {},
            },
        }
    },
    created:function(){
        let t = this;
        t.webSite.getList();
        t.mail.hasSet();
        t.mail.getInfo();
        t.getMonitorList();
        t.getAlertList();

        t.timer = window.setInterval(() => {
            t.webSite.getList();
            t.getMonitorList();
            t.getAlertList();
        }, 10000);

    },
    destroyed:function () {
        let t = this;
        // window.clearInterval(t.timer);
    },
    methods: {
        getAlertList: function () {
            let t = this;
            common.AjaxGet(t.alertList.api, function (data) {
                t.alertList.list = data.data.list;
                t.alertList.len = t.alertList.list.length;
            });
        },

        refreshAlertList: function () {
            let t = this;
            t.getAlertList();
        },

        getMonitorList: function () {
            let t = this;
            common.AjaxGet(t.monitor.api, function (data) {
                t.monitor.list = data.data;
            });
        },

        refreshMonitorList: function () {
            let t = this;
            t.getMonitorList();
        },

        getMonitorErrList: function () {
            let t = this;
            common.AjaxGet(t.monitorErrList.api, function (data){
                t.monitorErrList.list = data.data;
                t.monitorErrList.len = t.monitorErrList.list.length;
            });
        },

        monitorErrClear: function () {
            let t = this;
            t.isOk = "monitorErrClear";
            $("#isOkModal").modal("show");
        },

        monitorErrClearSubmit: function () {
            let t = this;
            common.AjaxGet(t.monitorErrList.clear, function (data){
                common.ToastShow(data.msg);
                t.getMonitorErrList();
                $("#isOkModal").modal('toggle');
            });
        },

        gotoList: function (pg) {
            let t = this;
            t.websiteList.page = pg;
            t.getList();
        },

        setUriPoint: function (item) {
            let t = this;
            console.log("setUriPoint item... ")
            console.log(item)
            t.point.hostUri = item.host + "/";
            t.point.hostId = item.hostID;
            t.getUriPoint();
            $("#setUriModal").modal('show');
        },

        getUriPoint: function () {
            let t = this;
            common.AjaxGet(t.point.apiList(), function (data){
                t.point.uriList = []
                if (data.code === 0) {
                    t.point.uriList = data.data.url;
                }
            });
        },

        addUriPoint: function () {
            let t = this;
            if (t.point.nowUri === "") {
                common.ToastShow("请输入URI");
                return
            }
            t.point.param.uri = t.point.hostUri + t.point.nowUri;
            common.AjaxPost(t.point.apiAdd(), t.point.param, function (data){
                common.ToastShow(data.msg);
                t.getUriPoint();
            });
        },

        gotoUriPoint: function (hostId, uri) {
            let t = this;
            t.point.hostId = hostId;
            t.point.param.uri = uri;
            common.AjaxPost(t.point.apiAdd(), t.point.param, function (data){
                common.ToastShow(data.msg);
                t.getUriPoint();
            });
        },

        delUriPoint: function (uri) {
            let t = this;
            t.point.param.uri = uri;
            common.AjaxPost(t.point.apiDel(), t.point.param, function (data){
                common.ToastShow(data.msg);
                t.getUriPoint();
            });
        },

        openWebsiteInfo: function (id) {
            let t = this;
            t.websiteInfo.hostId = id;
            common.AjaxGet(t.websiteInfo.api(), function (data){
                console.log(data);
                if (data.code === 0) {
                    t.websiteInfo.data = data.data;
                    $("#websiteInfoModal").modal('show');
                } else {
                    common.ToastShow(data.msg);
                }
            });
        },

        refreshWebsiteInfo: function (host, id) {
            let t = this;
            t.websiteInfo.host = host;
            t.websiteInfo.hostId = id;
            common.AjaxGetNotAsync(t.websiteInfo.refresh(), function (data){
                if (data.code === 0) {
                    t.openWebsiteInfo(id);
                } else {
                    common.ToastShow(data.msg);
                }
            });
        },

        gotoWebsite: function (item) {
            window.open(item.host, '_blank');
        },

        logShow: function (item) {
            let t = this;
            t.monitorLog.hostId = item.hostID;
            common.AjaxGetNotAsync(t.monitorLog.logListApi(), function (data){
                t.monitorLog.logList = data.data;
                t.monitorLog.log = t.monitorLog.logList[0];
            })
            common.AjaxGet(t.monitorLog.logApi(), function (data){
                if (data.code === 0) {
                    t.monitorLog.data = data.data;
                    $("#logModal").modal('show');
                } else {
                    common.ToastShow(data.msg);
                }
            });
        },

        loadLog: function () {
            let t = this;
            common.AjaxGet(t.monitorLog.loadApi(), function (data){
                if (data.code === 0) {
                    t.monitorLog.data = data.data;
                    $("#logModal").modal('show');
                } else {
                    common.ToastShow(data.msg);
                }
            });
        },

        uploadLog: function () {
            let t = this;
            window.location = t.monitorLog.logUpload();
        },

        deleteWebsiteOpen: function (item) {
            let t = this;
            t.deleteWebsite.hostId = item.hostID;
            t.deleteWebsite.hostName = item.host;
            t.isOk = "deleteWebsite";
            $("#isOkModal").modal("show");
        },

        deleteWebsiteSubmit: function () {
            let t = this;
            common.AjaxGet(t.deleteWebsite.api(), function (data){
                common.ToastShow(data.msg);
                //t.getMonitorErrList();
                $("#isOkModal").modal('toggle');
                location.reload();
            });
        },

        openEditWebsiteConf: function (hostId) {
            let t = this;
            t.editWebsiteConf.hostId = hostId;
            common.AjaxGetNotAsync(t.editWebsiteConf.confApi(), function (data) {
                t.editWebsiteConf.base = data.data.base;
                t.editWebsiteConf.alarmRule = data.data.alarmRule;
                t.editWebsiteConf.scanCheckUp = data.data.scanCheckUp;
                console.log(t.editWebsiteConf.param)
                $("#setAlertModal").modal("show");
            })

        },

        openEditWebsiteConfAtInfo: function (hostId) {
            let t = this;
            t.editWebsiteConf.hostId = hostId;
            common.AjaxGetNotAsync(t.editWebsiteConf.confApi(), function (data) {
                t.editWebsiteConf.base = data.data.base;
                t.editWebsiteConf.alarmRule = data.data.alarmRule;
                t.editWebsiteConf.scanCheckUp = data.data.scanCheckUp;
                console.log(t.editWebsiteConf.param)
                $("#websiteInfoModal").modal('toggle');
                $("#setAlertModal").modal("show");
            })
        },

        editWebsiteConfSubmit: function () {
            let t = this;
            let param = {
                "host": t.editWebsiteConf.base.host,
                "notes": t.editWebsiteConf.base.notes,
                "monitorRate": Number(t.editWebsiteConf.base.monitorRate),
                "contrastUrl": t.editWebsiteConf.base.contrastUrl,
                "contrastTime": Number(t.editWebsiteConf.base.contrastTime),
                "ping": t.editWebsiteConf.base.ping,
                "pingTime": Number(t.editWebsiteConf.base.pingTime),
                "uriDepth": Number(t.editWebsiteConf.scanCheckUp.uriDepth),
                "scanRate": Number(t.editWebsiteConf.scanCheckUp.scanRate),
                "scanBadLink": t.editWebsiteConf.scanCheckUp.scanBadLink,
                "scanTDK": t.editWebsiteConf.scanCheckUp.scanTDK,
                "scanExtLinks": t.editWebsiteConf.scanCheckUp.scanExtLinks,
                "websiteSlowResponseTime": Number(t.editWebsiteConf.alarmRule.websiteSlowResponseTime),
                "websiteSlowResponseCount": Number(t.editWebsiteConf.alarmRule.websiteSlowResponseCount),
                "SSLCertificateExpire": Number(t.editWebsiteConf.alarmRule.SSLCertificateExpire),
            }
            common.AjaxPost(t.editWebsiteConf.editApi(), param, function (data){
                common.ToastShow(data.msg);
                $("#setAlertModal").modal('toggle');
            })
        },

        openWebsiteUrl: function (item) {
            let t =this;
            t.websiteUrl.hostId = item.hostID;
            common.AjaxGet(t.websiteUrl.api(), function (data){
                t.websiteUrl.data = data.data;
                $("#urlInfoModal").modal("show");
            });
        },

        copy: function () {
            let t = this;
            let clipboard = new ClipboardJS('.copy');
            clipboard.on('success', e => {
                common.ToastShow("复制成功!");
                e.clearSelection();
            });
            clipboard.on('error', e => {
                common.ToastShow("复制失败！请重试或者手动复制内容!");
            });
        },

        openChart: function (item){
            let t = this;
            t.chartData.hostId = item.hostID;
            t.chartData.host = item.host;
            t.chartData.selectUriType = t.chartData.uriType[0].value;
            common.AjaxGetNotAsync(t.chartData.dayApi(), function (data) {
                t.chartData.dayList = data.data;
                t.chartData.selectDay = t.chartData.dayList[0];
            });
            t.$nextTick(() => {
                    t.DrawChart();
                }
            );
            $("#chartModal").modal("show");
        },

        loadingChart: function () {
            let t = this;
            t.$nextTick(() => {
                    t.DrawChart();
                }
            );
        },

        DrawChart: function () {
            let t = this;
            common.AjaxGetNotAsync(t.chartData.api(), function (data) {
                t.chartData.list = data.data;
            });
            var option = {
                tooltip: {
                    trigger: 'axis',
                    position: function (pt) {
                        return [pt[0], '10%'];
                    }
                },
                title: {
                    left: 'center',
                    text: t.chartData.host+"[" + t.chartData.uriTypeName[t.chartData.selectUriType] + "]",
                },
                toolbox: {
                    feature: {
                        saveAsImage: {}
                    }
                },
                xAxis: {
                    type: 'time',
                    boundaryGap: false
                },
                yAxis: {
                    type: 'value',
                    boundaryGap: [0, '100%'],
                    name:"(单位:ms)",
                },
                dataZoom: [
                    {
                        type: 'inside',
                        start: 0,
                        end: 100
                    },
                    {
                        start: 0,
                        end: 100
                    }
                ],
                series: [
                    {
                        name: '响应时间(ms)',
                        type: 'line',
                        smooth: true,
                        symbol: 'none',
                        areaStyle: {},
                        data: t.chartData.list
                    }
                ]
            };
            let myChart = echarts.init(document.getElementById('myChart'), '', {
                renderer: 'canvas',
                useDirtyRect: false
            });
            // 使用刚指定的配置项和数据显示图表。
            myChart.setOption(option);
        },

        openAlert: function (id) {
            let t = this;
            t.websiteAlert.hostId = id;
            common.AjaxGet(t.websiteAlert.api(), function (data){
               t.websiteAlert.list = data.data.list;
               t.websiteAlert.len = t.websiteAlert.list.length;
            });
            $("#alertModal").modal("show");
        },

        delAlert: function (id, hostId) {
            let t = this;
            common.AjaxGet(t.websiteAlert.del(id), function (data) {
                common.ToastShow(data.msg);
                if (data.code === 0) {
                    t.openAlert(hostId);
                }
            });
        },

        submitFidWebsiteInfo: function () {
            let t = this;
            if (t.objFidWebsiteInfo.host === "") {
                common.ToastShow("请输入Host！");
                return
            }
            $("#FidWebsiteInfoRse").hide();
            $("#FidWebsiteInfoLoading").show();
            common.AjaxGet(t.objFidWebsiteInfo.api(), function (data){
                t.objFidWebsiteInfo.rse = data.data;
                $("#FidWebsiteInfoRse").show();
                $("#FidWebsiteInfoLoading").hide();
                t.setHistory(t.objFidWebsiteInfo.id, t.objFidWebsiteInfo.host);
                t.getHistory(t.objFidWebsiteInfo.id);
            });
        },

        gotoFidWebsiteInfo: function (value) {
            let t = this;
            t.objFidWebsiteInfo.host = value;
            t.submitFidWebsiteInfo();
        },

        clearFidWebsiteInfoHistory: function () {
            let t = this;
            t.clearHistory(t.objFidWebsiteInfo.id);
            t.getHistory(t.objFidWebsiteInfo.id);
        },

        getHistory: function (id) {
            let t = this;
            common.AjaxGet(t.history.get(id), function (data){
                t.history.data = data.data;
            });
        },

        setHistory: function (id, value) {
            let t = this;
            common.AjaxPost(t.history.add(id, value), "", function (data){
                console.log(data)
            });
        },

        clearHistory: function (id) {
            let t = this;
            common.AjaxGet(t.history.clear(id), function (data){
                t.history.data = data.data;
            });
        },

        open: function (fid) {
            let t = this;
            t.getHistory(fid);
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

        submitFidNsLookUp: function () {
            let t = this;
            if (t.objFidNsLookUp.host === "") {
                common.ToastShow("请输入Host！");
                return
            }
            $("#FidNsLookUpRse").hide();
            $("#FidNsLookUpLoading").show();
            common.AjaxGet(t.objFidNsLookUp.api(), function (data){
                t.objFidNsLookUp.rse = data.data;
                $("#FidNsLookUpLoading").hide();
                $("#FidNsLookUpRse").show();
                t.setHistory(t.objFidNsLookUp.id, t.objFidNsLookUp.host);
                t.getHistory(t.objFidNsLookUp.id);
            });
        },

        gotoFidNsLookUp: function (value) {
            let t = this;
            t.objFidNsLookUp.host = value;
            t.submitFidNsLookUp();
        },

        clearFidNsLookUpHistory: function () {
            let t = this;
            t.clearHistory(t.objFidNsLookUp.id);
            t.getHistory(t.objFidNsLookUp.id);
        },

        submitFidWhois: function () {
            let t = this;
            if (t.objFidWhois.host === "") {
                common.ToastShow("请输入Host！");
                return
            }
            common.AjaxGet(t.objFidWhois.api(), function (data){
                t.objFidWhois.rse = data.data;
                t.setHistory(t.objFidWhois.id, t.objFidWhois.host);
                t.getHistory(t.objFidWhois.id);
            });
        },

        gotoFidWhois: function (value) {
            let t = this;
            t.objFidWhois.host = value;
            t.submitFidWhois();
        },

        clearFidWhoisHistory: function () {
            let t = this;
            t.clearHistory(t.objFidWhois.id);
            t.getHistory(t.objFidWhois.id);
        },

        initTippy: function () {
            tippy('.websiteInfoBtn',{
                content: "查看网站信息",
            });
            tippy('.websiteAlertBtn',{
                content: "打开报警信息",
            });
            tippy('.gotoWebsiteBtn',{
                content: "新页签打开该网站",
            });
            tippy('.editWebsiteConfBtn',{
                content: "修改网站监测配置",
            });
            tippy('.websiteUrlBtn',{
                content: "查看网站扫描的Url",
            });
            tippy('.uriPointBtn',{
                content: "设置网站监测点(Url)",
            });
            tippy('.logShowBtn',{
                content: "查看监测日志",
            });
            tippy('.openChartBtn',{
                content: "查看监测图",
            });
            tippy('.websiteOpenBtn',{
                content: "删除该网站监测",
            });
        },

    },
    computed: {
    },
    mounted(){
        this.initTippy()
    },
    updated() {
        this.initTippy()
    },
});
app.mount('#app');
