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

//todo; temporally fix for datepicker and vuejs causes empty value.
function setPostData(arrayObj) {
    var postData = {};
    $.each(arrayObj, function(i, val) {
        if (val.frmValue === true) {
            postData[val.frmName] = "on";
        } else if ((typeof val.frmFormat !== "undefined") && (val.frmFormat == "date-time")) {
            postData[val.frmName] = $("#f_" + val.frmName).val();
        } else {
            postData[val.frmName] = val.frmValue;
        }
    });
    return postData;
}

function setFormUtils() {
    $('.form-datetime').datetimepicker({
        format: "yyyy-mm-dd hh:ii"
    });
    $.getJSON("/admin/image/upload/url", function(data) {
        var uploadUrl = data.uploadurl;
        var options = {
            uploadUrl: uploadUrl,
            uploadFieldName: 'file',
            downloadFieldName: 'filename',
            allowedTypes: [
                'image/jpeg',
                'image/png',
                'image/jpg',
                'image/gif'
            ],
            progressText: '![Uploading file...]()',
            urlText: "![file]({filename})\n\n",
            errorText: "Error uploading file",
            extraParams: {},
            extraHeaders: {},
            onReceivedFile: function(file) {},
            onUploadedFile: function(json) {},
            customErrorHandler: function() {
                return true;
            },
            customUploadHandler: function(file) {
                return true;
            },
            dataProcessor: function(data) {
                //refresh upload url on uploaded
                var that = this;
                $.getJSON("/admin/image/upload/url", function(refreshData) {
                    that.uploadUrl = refreshData.uploadurl;
                });
                return data;
            }
        };
        $('textarea').inlineattach(options);
    });

}

$(document).ready(function() {

    Vue.directive("disable", function(value) {
        if ((value == "url") && (this.el.value !== "")) {
            this.el.setAttribute("disabled", "disabled");
        }
    });
    Vue.filter('dateFormat', function(value) {
        value = value.replace(/T/, " ");
        var localTime = moment.utc(value.slice(0, 16)).toDate();
        localTime = moment(localTime).format('YYYY-MM-DD HH:mm (Z)');
        return localTime;
    });
    Vue.filter('dateFormatUTC', function(value) {
        value = value.replace(/T/, " ");
        return value.slice(0, 16);
    });
    var formApp = {};
    formApp.drawForm = function() {
        var hashArr = location.hash.split("/");
        this.modelName = hashArr[1];
        this.crudMethod = hashArr[2];
        this.entityKey = hashArr[3];
        this.perPage = hashArr[4];
        this.sortResults = [];
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
                if (typeof val.format !== "undefined") {
                    tmpObject.frmFormat = val.format;
                }
                formElementArr.push(tmpObject);
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
                        requestUrl += "?offset=" + offset;
                    }
                    $.getJSON(requestUrl, function(listData) {
                        var nextPage = false;
                        var previousPage = false;
                        if (listData.has_next) {
                            nextPage = "/admin/form/#/" + that.modelName + "/list/" + listData.next_offset;
                        }
                        if (!isNaN(offset)) {
                            if (offset - listData.per_page !== 0) {
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
                                previous: previousPage,
                                offset: offset,
                                per_page: listData.per_page
                            },
                            methods: {
                                deleteEntity: function(e) {
                                    console.log(e.targetVM);
                                    if (window.confirm('Delete this entity?')) {
                                        $.del("/admin/rest/" + that.modelName + "/" + e.targetVM.$data.url, {}, function(delResponse) {
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
                case "sort":
                    var perPage = parseInt(that.perPage); 
                    if (isNaN(perPage)) {
                        perPage = 20;
                    }
                    var requestUrl = "/admin/rest/" + that.modelName+"?per_page="+perPage*2;
                    var offset = parseInt(that.entityKey);
                    if (!isNaN(offset)) {
                        requestUrl += "&offset=" + offset;
                    }
                    $.getJSON(requestUrl, function(listData) {
                        var listVue = new Vue({
                            el: "#formContainer",
                            template: "#modelListSortable",
                            data: {
                                items: listData.items,
                                modelName: listData.model_name,
                                offset: offset
                            },
                            methods: {
                                submitOrderUpdate: function(e) {
                                    console.log(that.sortResults)
                                    $.each(that.sortResults, function(i, val) {
                                        if (val.neworder != val.oldorder) {
                                            var putData = {
                                                pageorder: val.neworder
                                            };
                                            $.put("/admin/rest/" + that.modelName + "/" + val.url, putData, function(putResponse) {
                                                console.log(putResponse);
                                                if (putResponse.message == "updated") {
                                                    location.href = "/admin/form/#/" + that.modelName + "/list";
                                                    location.reload(true);
                                                } else {
                                                    $("#postAlert").html('<div class="alert alert-danger" role="alert">error posting data</div>');
                                                }
                                            });
                                        }
                                    });
                                }
                            }
                        });
                        var orderArray = [];
                        $.each(listData.items, function(i, val) {
                            orderArray.push(parseInt(val.pageorder));
                        });
                        var sortEl = document.querySelector(".sortable");
                        new Sortable(sortEl,{
                            onUpdate:function(e){
                                 var newOrderArray = [];
                                $.each($("li.list-group-item"), function(i, val) {
                                    newOrderArray.push({
                                        url: val.getAttribute("data-url"),
                                        neworder: orderArray[i],
                                        oldorder: parseInt(val.getAttribute("data-order"))
                                    });
                                });
                                console.log(newOrderArray);
                                that.sortResults = newOrderArray;
                                $("#SaveOrder").removeAttr("disabled");
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
                        filters: {
                            marked: marked
                        },
                        methods: {
                            submitUpdate: function(e) {
                                e.preventDefault();
                                console.log(this.$data);
                                postData = setPostData(this.$data.items);
                                $.post("/admin/rest/" + that.modelName, postData, function(data) {
                                    if (data.message == "created") {
                                        location.href = "/admin/form/#/" + that.modelName + "/list";
                                        location.reload(true);
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
                    // $('.form-datetime').datetimepicker({
                    //         format: "yyyy-mm-dd hh:ii"
                    // });
                    setFormUtils();
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
                            filters: {
                                marked: marked
                            },
                            methods: {
                                submitUpdate: function(e) {
                                    e.preventDefault();
                                    console.log(this.$data);
                                    putData = setPostData(this.$data.items);
                                    $.put("/admin/rest/" + that.modelName + "/" + that.entityKey, putData, function(putResponse) {
                                        if (putResponse.message == "updated") {
                                            location.href = "/admin/form/#/" + that.modelName + "/list";
                                            location.reload(true);
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
                        setFormUtils();
                        // $('.form-datetime').datetimepicker({
                        //     format: "yyyy-mm-dd hh:ii"
                        // });
                    });
            }
        });


    };
    formApp.drawForm();

    window.onhashchange = formApp.drawForm;
});
