$(document).ready(function(){
    var hashArr = location.hash.split("/");
    $.getJSON("/admin/rest/schema/"+hashArr[1],function(data){
        var vue = new Vue({
            el: '#formContainer',
            data: {
                formtitle: data.schema.title 
            }
        });  
    });
});
