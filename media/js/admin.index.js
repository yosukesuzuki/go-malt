$(document).ready(function(){
    $.getJSON("/admin/rest/models",function(data){
        $.each(data.models,function(i,val){
            console.log(val);
        });
    });
});
