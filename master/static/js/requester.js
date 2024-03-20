const { createApp, ref } = Vue;
const common = new Utils;
const app = createApp({
    data() {
        return {
            mail: Mail,
            webSite: AddWebSite,
            bodyType: "json",
            param: {
                reqId: "",
                name : "新建请求",
                note : "",
                method : "GET",
                url : "",
                header : {},
                bodyType : "",
                bodyJson : "",
                bodyFromData : {},
                bodyXWWWFrom : {},
                bodyXml : "",
                bodyText : "",
            },
            resp:{},
            history: {
                api: "/api/requester/history/list",
                list: [],
            },
            nowReqList: {
                api: "/api/requester/list",
                list: [],
                len: 0,
            },
            createTabApi: "/api/requester/create/tab",
            getNowApi: function (id) { return "/api/requester/data/"+id; },
            closeTabApi: function (id) { return "/api/requester/close/tab/"+id; },
            deleteHistoryApi: function (id) { return  "/api/requester/history/delete/"+id; },
        }
    },
    created:function(){
        let t = this;
        t.mail.getInfo();

        t.getHistory();
        t.getNowReqList();
    },
    destroyed:function () {
    },
    methods: {
        getHistory: function () {
            let t = this;
            common.AjaxGet(t.history.api, function (data) {
               if (data.code === 0) {
                   t.history.list = data.data;
               }
            });
        },
        getNowReqList: function () {
            let t = this;
            common.AjaxGet(t.nowReqList.api, function (data) {
                if (data.code === 0) {
                    t.nowReqList.list = data.data;
                    t.nowReqList.len = t.nowReqList.list.length;
                    if (t.nowReqList.len === 0) {
                       t.createTab();
                    }else {
                        t.param.reqId = t.nowReqList.list[0].id;
                        t.getNow();
                    }
                }
            });
        },
        openTab: function (id) {
            let t = this;
            t.param.reqId = id;

            for (item in t.nowReqList.list) {
                if (t.nowReqList.list[item].id === id ) {
                    t.nowReqList.list[item].isNow = true;
                } else {
                    t.nowReqList.list[item].isNow = false;
                }
            }
            t.getNow();
        },
        closeTab: function (id) {
            let t = this;
            common.AjaxGet(t.closeTabApi(id), function (data){
                if (data.code === 0) {
                    for (item in t.nowReqList.list) {
                        if (t.nowReqList.list[item].id === id ) {
                            t.nowReqList.list.splice(item, 1);
                        }
                    }
                    if (t.nowReqList.list.length === 0) {
                        t.getNowReqList();
                    }
                }
            });
        },
        createTab: function () {
            let t = this;
            common.AjaxGet(t.createTabApi, function (data){
                if (data.code === 0) {
                    t.getNowReqList()
                    $('.gdtiao').scrollLeft(0);
                }
            });
        },
        getNow: function () {
          let t = this;
          common.AjaxGet(t.getNowApi(t.param.reqId), function (data) {
              if (data.code === 0) {
                  console.log(data.data);
                  t.resp = data.data;
                  t.param.name = data.data.name;
                  t.param.note = data.data.note;
                  t.param.method = data.data.method;
                  t.param.url = data.data.url;
                  t.param.header = data.data.header;
                  // t.param.bodyType = data.data.bodyType;
                  // t.param.bodyJson = data.data.bodyJson;
                  // t.param.bodyFromData = data.data.bodyFromData;
                  // t.param.bodyXWWWFrom = ""
                  // t.param.bodyXml = ""
                  // t.param.bodyText = "",
              }
          });
        },
        deleteHistory: function (id) {
            let t = this;
            common.AjaxGet(this.deleteHistoryApi(id), function (data){
                if (data.code === 0) {
                    t.getHistory();
                }
            });
        },
        closeSet: function () {
            console.log("closeSet");
            $("#setHeaderTable").hide();
            $("#setQueryTable").hide();
            $("#setBodyTable").hide();
            $("#openHeaderLi").removeClass("select_activity");
            $("#openQueryLi").removeClass("select_activity");
            $("#openBodyLi").removeClass("select_activity");
        },
        openHeaderTable: function () {
            let t =this;
            console.log("openHeaderTable");
            t.closeSet();
            $("#setHeaderTable").show();
            $("#openHeaderLi").addClass("select_activity");
        },
        openQueryTable: function () {
            let t =this;
            console.log("openQueryTable");
            t.closeSet();
            $("#setQueryTable").show();
            $("#openQueryLi").addClass("select_activity");
        },
        openBodyTable: function () {
            let t =this;
            console.log("openBodyTable");
            t.closeSet();
            $("#setBodyTable").show();
            $("#openBodyLi").addClass("select_activity");
        },
        closeBodyMain: function () {
            $("#bodyJson").hide();
            $("#bodyFromData").hide();
            $("#bodyXwwwFrom").hide();
            $("#bodyXml").hide();
            $("#bodyText").hide();
        },
        openBodyMain: function (type) {
            let t =this;
            t.bodyType = type;
            t.closeBodyMain();
            switch (type) {
                case "json":
                    console.log("json");
                    $("#bodyJson").show();
                    break;
                case "from-data":
                    console.log("from-data");
                    $("#bodyFromData").show();
                    break;
                case "x-www-form-urlencoded":
                    console.log("x-www-form-urlencoded");
                    $("#bodyXwwwFrom").show();
                    break;
                case "xml":
                    console.log("xml");
                    $("#bodyXml").show();
                    break;
                case "plain":
                    console.log("plain");
                    $("#bodyText").show();
                    break;
            }
        },
        closeRse: function () {
            $("#rseDiv").hide();
            $("#rseRpHeader").hide();
            $("#rseHeader").hide();
            $("#rseCookie").hide();
            $("#openRse").removeClass("select_activity");
            $("#openRseRpHeader").removeClass("select_activity");
            $("#openRseHeader").removeClass("select_activity");
            $("#openRseCookie").removeClass("select_activity");
        },
        openRseDiv: function () {
            let t = this;
            t.closeRse();
            $("#rseDiv").show();
            $("#openRse").addClass("select_activity");
        },
        openRseRpHeaderDiv: function () {
            let t = this;
            t.closeRse();
            $("#rseRpHeader").show();
            $("#openRseRpHeader").addClass("select_activity");
        },
        openRseHeaderDiv: function () {
            let t = this;
            t.closeRse();
            $("#rseHeader").show();
            $("#openRseHeader").addClass("select_activity");
        },
        openRseCookieDiv: function () {
            let t = this;
            t.closeRse();
            $("#rseCookie").show();
            $("#openRseCookie").addClass("select_activity");
        },
        openGlobalParamModal: function () {
            let t = this;
            $("#globalParamModal").modal("show");
        },
        openCookieManageModal: function () {
            let t = this;
            $("#cookieManageModal").modal("show");
        },
        openApiNoteModal: function () {
            let t = this;
            $("#apiNoteModal").modal("show");
        },
        openCodeModal: function () {
            let t = this;
            $("#codeModal").modal("show");
        },
        openGotoDirModal: function () {
            let t = this;
            $("#gotoDirModal").modal("show");
        },
        openDirModal: function () {
            let t = this;
            $("#dirModal").modal("show");
        },
        // 发送
        execute: function () {
            let t = this;
            console.log(t.param);
            common.AjaxPost("/api/requester/execute", t.param, function (data){
                console.log(data);
                if (data.code === 0 ){
                    t.resp = data.data;
                    t.getNowReqList();
                    t.getHistory();
                }
            });
        },
    },
    computed: {
    },
    mounted:function(){
    },
});
app.mount('#app');