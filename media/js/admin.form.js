//extend jquery ajax from http://d.hatena.ne.jp/anatoo/20090527/1243415081
$.extend({
    "put": function(url, data, success, error) {
        error = error || function() {};
        return $.ajax({
            "url": url,
            "data": data,
            "success": success,
            "type": "PUT",
            "cache": false,
            "error": error,
            "dataType": "json"
        });
    },
    "del": function(url, data, success, error) {
        error = error || function() {};
        return $.ajax({
            "url": url,
            "data": data,
            "success": success,
            "type": "DELETE",
            "cache": false,
            "error": error,
            "dataType": "json"
        });
    }
});

$(document).ready(function() {
    Vue.filter('dateFormat', function(value) {
        return value.slice(0, 16)
    });
    var formApp = {};
    formApp.drawForm = function() {
        var hashArr = location.hash.split("/");
        this.modelName = hashArr[1];
        this.crudMethod = hashArr[2];
        this.entityKey = hashArr[3];
        var that = this;
        $.getJSON("/admin/rest/schema/" + that.modelName, function(data) {
            var crudVue = new Vue({
                el: '#crudContainer',
                data: {
                    formtitle: data.schema.title,
                    id: data.schema.id
                }
            });
            var formElementArr = [];
            $.each(data.schema.properties, function(i, val) {
                val.name = i;
                var tmpObject = {
                    frmName: val.name,
                    frmTitle: val.title,
                    frmType: val.type,
                    frmFieldOrder: val.fieldOrder,
                    frmValue: ""
                };
                if (typeof val.maxLength !== "undefined") {
                    tmpObject.frmMaxLength = val.maxLength;
                }
                formElementArr.push(tmpObject)
            });
            formElementArr.sort(function(a, b) {
                var x = a.frmFieldOrder;
                var y = b.frmFieldOrder;
                if (x > y) return 1;
                if (x < y) return -1;
                return 0;
            });
            var crudMethod = that.crudMethod;

            switch (crudMethod) {
                case "list":
                    var offset = parseInt(that.entityKey);
                    var requestUrl = "/admin/rest/" + that.modelName;
                    if (!isNaN(offset)) {
                        requestUrl += "?offset=" + offset
                    }
                    $.getJSON(requestUrl, function(listData) {
                        var nextPage = false;
                        var previousPage = false;
                        if (listData.has_next) {
                            nextPage = "/admin/form/#/" + that.modelName + "/list/" + listData.next_offset;
                        }
                        if (!isNaN(offset)) {
                            if (offset - listData.per_page != 0) {
                                previousPage = "/admin/form/#/" + that.modelName + "/list/" + (offset - listData.per_page);
                            } else {
                                previousPage = "/admin/form/#/" + that.modelName + "/list";
                            }
                        }
                        var listVue = new Vue({
                            el: "#formContainer",
                            template: "#modelListTable",
                            data: {
                                items: listData.items,
                                modelName: listData.model_name,
                                next: nextPage,
                                previous: previousPage
                            },
                            methods: {
                                deleteEntity: function(e) {
                                    console.log(e.targetVM);
                                    if(window.confirm('Delete this entity?')){
                                        $.del("/admin/rest/" + that.modelName + "/" + e.targetVM.$data.url,{},function(delResponse) {
                                            if (delResponse.message == "deleted") {
                                                //location.href = "/admin/form/#/" + that.modelName + "/list/success";
                                                //e.targetVM.$data;
                                                console.log("deleted");
                                            } else {
                                                $("#postAlert").html('<div class="alert alert-danger" role="alert">error deleting data</div>');
                                            }
                                        });
                                        e.targetVM.$remove(e.targetVM.$data);
                                    }
                                }
                            }
                        });
                    });
                    break;
                case "create":
                    var createVue = new Vue({
                        el: "#formContainer",
                        template: "#modelForm",
                        data: {
                            items: formElementArr,
                        },
                        methods: {
                            submitUpdate: function(e) {
                                e.preventDefault();
                                console.log(this.$data);
                                postData = {}
                                $.each(this.$data.items, function(i, val) {
                                    if (val.frmValue == true) {
                                        postData[val.frmName] = "on";
                                    } else {
                                        postData[val.frmName] = val.frmValue;
                                    }
                                });
                                $.post("/admin/rest/" + that.modelName, postData, function(data) {
                                    if (data.message == "created") {
                                        location.href = "/admin/form/#/" + that.modelName + "/list";
                                    } else {
                                        $("#postAlert").html('<div class="alert alert-danger" role="alert">error posting data</div>');
                                    }
                                });
                            },
                            cancel: function(e) {
                                location.href = "/admin/form/#/" + that.modelName + "/list";
                            }
                        }
                    });
                    break;
                case "update":
                    $.getJSON("/admin/rest/" + that.modelName + "/" + that.entityKey, function(data) {
                        $.each(formElementArr, function(i, val) {
                            val.frmValue = data.item[val.frmName];
                            formElementArr[i] = val;
                        });
                        var updateVue = new Vue({
                            el: "#formContainer",
                            template: "#modelForm",
                            data: {
                                items: formElementArr,
                            },
                            methods: {
                                submitUpdate: function(e) {
                                    e.preventDefault();
                                    console.log(this.$data);
                                    putData = {}
                                    $.each(this.$data.items, function(i, val) {
                                        if (val.frmValue == true) {
                                            putData[val.frmName] = "on";
                                        } else {
                                            putData[val.frmName] = val.frmValue;
                                        }
                                    });
                                    $.put("/admin/rest/" + that.modelName + "/" + that.entityKey, putData, function(putResponse) {
                                        if (putResponse.message == "updated") {
                                            location.href = "/admin/form/#/" + that.modelName + "/list/success";
                                        } else {
                                            $("#postAlert").html('<div class="alert alert-danger" role="alert">error posting data</div>');
                                        }
                                    });
                                },
                                cancel: function(e) {
                                    location.href = "/admin/form/#/" + that.modelName + "/list";
                                }
                            }
                        });
                    });
            }
        });
    };
    formApp.drawForm();
    window.onhashchange = formApp.drawForm;
});
