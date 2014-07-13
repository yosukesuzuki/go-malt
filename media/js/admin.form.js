$(document).ready(function(){
    var hashArr = location.hash.split("/");
    $.getJSON("/admin/rest/schema/"+hashArr[1],function(data){
        var crudVue = new Vue({
            el: '#crudContainer',
            data: {
                formtitle: data.schema.title,
                id: data.schema.id
            }
        });
        var crudMethod = hashArr[2];
        switch(crudMethod){
            case "list":
                $.getJSON("/admin/rest/"+hashArr[1],function(listData){
                    var listVue = new Vue({
                        el: "#formContainer",
                        template: "#modelListTable",
                        data: {
                            items: listData.items,
                            modelName: listData.model_name
                        }
                    });
                });
        }
    });
    Vue.filter('dateFormat', function (value) {
        return value.slice(0,16)
    });
});
