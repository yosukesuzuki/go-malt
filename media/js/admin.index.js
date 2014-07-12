$(document).ready(function(){
    $.getJSON("/admin/rest/models",function(data){
        var vue = new Vue({
            el: '#models',
            data: {
                models: data.models 
            }
        });
    });
});
