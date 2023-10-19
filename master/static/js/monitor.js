const { createApp, ref } = Vue;
const common = new Utils;
const app = createApp({
    data() {
        return {
            mail: Mail,
            webSite: AddWebSite,
            monitor: {
                api: "/api/monitor/list",
                list: [],
            },
        }
    },
    created:function(){
        let t = this;
        t.mail.getInfo();
        t.getList();
    },
    destroyed:function () {
    },
    methods: {
        getList: function () {
            let t = this;
            common.AjaxGet(t.monitor.api, function (data) {
                t.monitor.list = data.data;
            });
        },
    },
    computed: {
    },
    mounted:function(){
    },
});
app.mount('#app');