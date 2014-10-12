$(document).ready(function() {
  $.getJSON("/rest/search?keyword=" + $("#keyword").val(), function(data) {
      var searchVue = new Vue({
          el: '#searchList',
          data: {
              items: data.items,
          }
      });
  });
});
