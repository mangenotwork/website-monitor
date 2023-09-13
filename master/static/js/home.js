const { createApp, ref } = Vue;
import common from './public.js'
const app = createApp({
    data() {
        return {
            title: "",
            addWebSite: {
                api: "/api/website/add",
                param : {
                    host: "",
                    healthUri: "",
                    rate: 10,
                    alarmResTime: 3000,
                    uriDepth: 2,
                    uriUpdateRate: 24,
                }
            },
            msg: "",
            websiteList: {
                api: function (){ return "/api/website/list"; },
                page: 1,
                data: [],
            },
            deleteWebsite: {
                api: function (){ return "/api/website/delete/" + this.id; },
                id: "",
                hostName: "",
            },
            hasMail: {
                api: "/api/mail/init",
                data: false,
            },
            mailConf: {
                api: "/api/mail/conf",
                param: {
                    host: "smtp.qq.com",
                    port: 25,
                    from: "",
                    authCode: "",
                    toList: "",
                },
            },
            mailInfo: {
                api: "/api/mail/info",
            },
            mailSend: {
                api: "/api/mail/sendTest",
            },
            point: {
                hostId: "",
                hostUri: "",
                apiAdd: function (){ return "/api/point/add/"+this.hostId; },
                apiList: function () { return "/api/point/list/"+this.hostId; },
                apiDel: function () { return "/api/point/del/"+this.hostId; },
                param: {
                    uri:"",
                },
                uriList: [],
                nowUri: "",
            },
            websiteInfo: {
                hostId: "",
                api: function() { return "/api/website/info/"+this.hostId; },
                data: {}
            },
            alertList: {
                api: "/api/alert/list",
                clear: "/api/alert/clear",
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
                api: function (){ return "/api/monitor/log/" + this.hostId; },
                logListApi: function () { return "/api/log/list/" + this.hostId; },
                logUpload: function () { return "/api/log/upload/" + this.hostId + "?day=" + this.log; },
                data: {},
                logList: [],
                log: "",
            },
            editWebsiteConf: {
                api: "/api/website/edit",
                param: {
                    hostId: "",
                    rate: 10,
                    alarmResTime: 3000,
                    uriDepth: 2,
                },
                host: "",
            },
            chartData: {
                hostId: "",
                api: function () {
                  return "/api/website/chart/" + this.hostId + "?day=" + this.selectDay + "&uri=" + this.selectUriType;
                },
                dayApi: function () { return "/api/log/list/" + this.hostId; },
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
                api: function () { return "/api/website/alert/" + this.hostId; },
                list: [],
                len: 0,
                del: function (date) { return "/api/website/alert/del/" + this.hostId + "?date="+date; }
            }
        }
    },
    created:function(){
        let t = this;
        t.getList();
        t.getMail();
        t.getMailInfo();
        t.getAlertList();
        t.getMonitorErrList();
        t.timer = window.setInterval(() => {
            t.getList();
        }, 10000);
    },
    destroyed:function () {
        let t = this;
        window.clearInterval(t.timer);
    },
    methods: {
        addWebSiteMonitor: function () {
            let t = this;
            t.addWebSite.param.rate = Number(t.addWebSite.param.rate);
            t.addWebSite.param.alarmResTime = Number(t.addWebSite.param.alarmResTime);
            t.addWebSite.param.uriDepth = Number(t.addWebSite.param.uriDepth);
            common.AjaxPost(t.addWebSite.api, t.addWebSite.param, function (data){
                if (data.code === 0) {
                    $("#addHostModal").modal('toggle');
                    t.getList();
                }
                common.ToastShow(data.msg);
            });
        },
        test:function () {
            console.log("test...")
        },
        getList: function () {
            let t = this;
            common.AjaxGet(t.websiteList.api(), function (data){
                t.websiteList.data = data.data;
            });
        },
        getAlertList: function () {
            let t = this;
            common.AjaxGet(t.alertList.api, function (data) {
                t.alertList.list = data.data;
                t.alertList.len = t.alertList.list.length;
            });
        },
        alertClear: function () {
            let t = this;
            t.isOk = "alertClear";
            $("#isOkModal").modal("show");
        },
        alertClearSubmit: function () {
            let t = this;
            common.AjaxGet(t.alertList.clear, function (data){
                common.ToastShow(data.msg);
                t.getAlertList();
                $("#isOkModal").modal('toggle');
            });
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
        getMail: function () {
            let t = this;
            common.AjaxGet(t.hasMail.api, function (data){
                t.hasMail.data = data.data;
            })
        },
        setMailConf: function () {
            let t = this;
            t.mailConf.param.toList = common.MailToListJoin(t.mailConf.param.toList);
            t.mailConf.param.port = Number(t.mailConf.param.port)
            common.AjaxPost(t.mailConf.api, t.mailConf.param, function (data){
                if (data.code === 0) {
                    $("#mailSetModal").modal('toggle');
                    t.getMail();
                }
                common.ToastShow(data.msg);
            });
        },
        getMailInfo: function () {
            let t = this;
            common.AjaxGet(t.mailInfo.api, function (data){
                t.mailConf.param.host = data.data.Host;
                t.mailConf.param.port = data.data.Port;
                t.mailConf.param.from = data.data.From;
                t.mailConf.param.authCode = data.data.AuthCode;
                t.mailConf.param.toList = data.data.ToList;
            });
        },
        mailSendTest: function () {
            let t = this;
            t.mailConf.param.toList = common.MailToListJoin(t.mailConf.param.toList);
            t.mailConf.param.port = Number(t.mailConf.param.port)
            common.AjaxPost(t.mailSend.api, t.mailConf.param, function (data){
                common.ToastShow(data.msg);
            });
        },
        setUriPoint: function (item) {
            let t = this;
            t.point.hostUri = item.HealthUri + "/";
            t.point.hostId = item.ID;
            t.getUriPoint();
            $("#setUriModal").modal('show');
        },
        getUriPoint: function () {
            let t = this;
            common.AjaxGet(t.point.apiList(), function (data){
                t.point.uriList = []
                if (data.code === 0) {
                    t.point.uriList = data.data;
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
                if (data.code === 0) {
                    t.websiteInfo.data = data.data;
                    $("#websiteInfoModal").modal('show');
                } else {
                    common.ToastShow(data.msg);
                }
            });
        },
        logShow: function (id) {
            let t = this;
            t.monitorLog.hostId = id;
            common.AjaxGetNotAsync(t.monitorLog.logListApi(), function (data){
                t.monitorLog.logList = data.data.DayList;
                t.monitorLog.log = t.monitorLog.logList[0];
            })
            common.AjaxGet(t.monitorLog.api(), function (data){
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
            t.deleteWebsite.hostId = item.ID;
            t.deleteWebsite.hostName = item.Host;
            t.isOk = "deleteWebsite";
            $("#isOkModal").modal("show");
        },
        deleteWebsiteSubmit: function () {
            let t = this;
            common.AjaxGet(t.deleteWebsite.api(), function (data){
                common.ToastShow(data.msg);
                t.getMonitorErrList();
                $("#isOkModal").modal('toggle');
                t.getList();
            });
        },
        openEditWebsiteConf: function (item) {
            let t = this;
            t.editWebsiteConf.host = item.Host;
            t.editWebsiteConf.param.hostId = item.ID;
            t.editWebsiteConf.param.rate = item.Rate;
            t.editWebsiteConf.param.alarmResTime = item.AlarmResTime;
            t.editWebsiteConf.param.uriDepth = item.UriDepth;
            console.log(item)
            console.log(t.editWebsiteConf)
            $("#setAlertModal").modal("show");
        },
        editWebsiteConfSubmit: function () {
            let t = this;
            t.editWebsiteConf.param.rate = Number(t.editWebsiteConf.param.rate);
            t.editWebsiteConf.param.alarmResTime = Number(t.editWebsiteConf.param.alarmResTime);
            t.editWebsiteConf.param.uriDepth = Number(t.editWebsiteConf.param.uriDepth);
            common.AjaxPost(t.editWebsiteConf.api, t.editWebsiteConf.param, function (data){
                common.ToastShow(data.msg);
                $("#setAlertModal").modal('toggle');
            })
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
            t.chartData.hostId = item.ID;
            t.chartData.host = item.Host;
            t.chartData.selectUriType = t.chartData.uriType[0].value;
            common.AjaxGetNotAsync(t.chartData.dayApi(), function (data) {
                t.chartData.dayList = data.data.DayList;
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
               t.websiteAlert.list = data.data;
               t.websiteAlert.len = t.websiteAlert.list.length;
            });
            $("#alertModal").modal("show");
        },
        delAlert: function (date) {
            let t = this;
            common.AjaxGet(t.websiteAlert.del(date), function (data) {
                common.ToastShow(data.msg);
                if (data.code === 0) {
                    t.openAlert(data.data);
                }
            });
        }
    },
    computed: {
    },
    mounted:function(){
    },
});
app.mount('#app');