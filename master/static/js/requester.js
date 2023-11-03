const { createApp, ref } = Vue;
const common = new Utils;
const app = createApp({
    data() {
        return {
            mail: Mail,
            webSite: AddWebSite,
            bodyType: "json",
        }
    },
    created:function(){
        let t = this;
        t.mail.getInfo();
    },
    destroyed:function () {
    },
    methods: {
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
    },
    computed: {
    },
    mounted:function(){
    },
});
app.mount('#app');