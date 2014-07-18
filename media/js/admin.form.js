$(document).ready(function(){
    Vue.filter('dateFormat', function (value) {
        return value.slice(0,16)
    });
    var formApp = {};
    formApp.drawForm = function(){
        var hashArr = location.hash.split("/");
        this.modelName = hashArr[1];
        this.crudMethod = hashArr[2];
        this.entityKey = hashArr[3];
        var that = this;
        $.getJSON("/admin/rest/schema/"+that.modelName,function(data){
            var crudVue = new Vue({
                el: '#crudContainer',
                data: {
                    formtitle: data.schema.title,
                    id: data.schema.id
                }
            });
            var formElementArr = [];
            $.each(data.schema.properties,function(i,val){
                val.name = i;
                var tmpObject = {frmName:val.name,frmTitle:val.title,frmType:val.type,frmFieldOrder:val.fieldOrder,frmValue:""};
                if(typeof val.maxLength !== "undefined"){
                    tmpObject.frmMaxLength = val.maxLength;
                }
                formElementArr.push(tmpObject)
            });
            formElementArr.sort(function(a, b){
                var x = a.frmFieldOrder;
                var y = b.frmFieldOrder;
                if (x > y) return 1;
                if (x < y) return -1;
                return 0;
            });
            var crudMethod = that.crudMethod;

            switch(crudMethod){
                case "list":
                    $.getJSON("/admin/rest/"+that.modelName,function(listData){
                        var listVue = new Vue({
                            el: "#formContainer",
                            template: "#modelListTable",
                            data: {
                                items: listData.items,
                                modelName: listData.model_name
                            }
                        });
                    });
                    if((typeof entityKey !== "undefined") && (entityKey == "success")){
                        $("#postAlert").html('<div class="alert alert-success" role="alert"><button type="button" class="close" data-dismiss="alert"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button>success</div>');
                        $("#postAlert").css("display","block");
                    }
                    break;
                case "create":
                    var createVue = new Vue({
                    el: "#formContainer",
                    template: "#modelForm",
                    data: {
                        items: formElementArr,
                    },
                    methods: {
                        submitUpdate: function (e) {
                            e.preventDefault();
                            console.log(this.$data);
                            postData = {}
                            $.each(this.$data.items,function(i,val){
                                if(val == true){
                                    postData[val.frmName] = "on";
                                }else{
                                    postData[val.frmName] = val.frmValue;
                                }
                            });
                            $.post("/admin/rest/"+that.modelName,postData,function(data){
                                if(data.message == "created"){
                                    location.href = "/admin/form/#/"+that.modelName+"/list/success";
                                }else{
                                    $("#postAlert").html('<div class="alert alert-danger" role="alert">error posting data</div>');
                                }
                            });
                        },
                        cancel: function (e) {
                            location.href = "/admin/form/#/"+that.modelName+"/list";
                        }
                    }
                    });
                    break;
                case "update":
                    var updateVue = new Vue({
                        el: "#formContainer",
                        template: "#modelForm",
                        data: {
                                items: formElementArr,
                        }
                    });
            }
        });
    };
    formApp.drawForm();
    window.onhashchange = formApp.drawForm;
});
