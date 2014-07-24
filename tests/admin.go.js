//execute # casperjs test tests/admin.go.js

baseURL = 'http://localhost:8080';

casper.start();

casper.open(baseURL + '/admin/');

casper.then(function() {
    this.echo(this.getTitle());
    var currentUrl = this.getCurrentUrl();
    if (currentUrl.match(/login/)) {
        this.click('#admin');
        this.thenClick('#submit-login');
    }
});

// test of /admin/
casper.then(function() {
    this.echo(this.getTitle());
    this.test.assertHttpStatus(200);
});

// test of /admin/rest/models
casper.thenOpen(baseURL + '/admin/rest/models', function() {
    this.test.assertHttpStatus(200);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.models.length, 2, 'total count of models should be 2');
});

// test of /admin/rest/schema/adminpage
casper.thenOpen(baseURL + '/admin/rest/schema/adminpage', function() {
    this.test.assertHttpStatus(200);
    var jsonData = JSON.parse(this.getPageContent());
    //this.test.assertEqual(jsonData.fields.length,9,'total count of models should be 9');
    //this.test.assertEqual(jsonData.model_name,'adminpage','title should be adminpage');
});

// add entity by post
casper.thenOpen(baseURL + '/admin/rest/adminpage', {
    method: "post",
    data: {
        title: 'title1',
        url: 'url1',
        pageorder: 1,
        content: 'foobar',
        displaypage: 'on',
    }
}, function() {
    this.echo("POST request has been sent.")
    this.test.assertHttpStatus(201);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.message, 'created', 'return created message');
});

// add entity by post
casper.thenOpen(baseURL + '/admin/rest/adminpage', {
    method: "post",
    data: {
        title: 'title0',
        url: 'url0',
        pageorder: 0,
        content: 'foobar',
        displaypage: '',
    }
}, function() {
    this.echo("POST request has been sent.")
    this.test.assertHttpStatus(201);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.message, 'created', 'return created message');
});

casper.wait(1000, function() {
    this.echo("I've waited for a second.");
});

casper.thenOpen(baseURL + '/admin/rest/adminpage', function() {
    this.test.assertHttpStatus(200);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.items[0].displaypage, false, 'displaypage of the first entity should be false');
    this.test.assertEqual(jsonData.items[0].title, 'title0', 'title of the first entity should be title0');
    this.test.assertEqual(jsonData.items[0].url, 'url0', 'url of the first entity should be url0');
    this.test.assertEqual(jsonData.items[0].pageorder, 0, 'content of the first entity should be 0');
    this.test.assertEqual(jsonData.items[0].content, 'foobar', 'title of the first entity should be foobar');
    this.test.assertEqual(jsonData.items[1].displaypage, true, 'displaypage of the first entity should be true');
    this.test.assertEqual(jsonData.items[1].title, 'title1', 'title of the first entity should be title1');
    this.test.assertEqual(jsonData.items[1].url, 'url1', 'url of the first entity should be url1');
    this.test.assertEqual(jsonData.items[1].pageorder, 1, 'content of the first entity should be 1');
    this.test.assertEqual(jsonData.items[1].content, 'foobar', 'title of the first entity should be foobar');
});
/*
//something wrong with PUT request by casperjs
casper.thenOpen(baseURL+'/admin/rest/adminpage/url0', {
    method: "PUT",
    data: {
        title: 'title0-2',
        url: 'url0',
        pageorder: 0,
        content: 'foobar-2',
        displaypage: 'on',
    }
}, function() {
    this.echo("PUT request has been sent.")
    this.test.assertHttpStatus(200);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.message,'updated','return update message');
});

casper.thenOpen(baseURL+'/admin/rest/adminpage',function(){
    this.test.assertHttpStatus(200);
    var jsonData = JSON.parse(this.getPageContent());
//    this.test.assertEqual(jsonData.items[0].displaypage,true,'displaypage of the first entity should be false');
    this.test.assertEqual(jsonData.items[0].title,'title0-2','title of the first entity should be title0');
    this.test.assertEqual(jsonData.items[0].url,'url0','url of the first entity should be url0');
    this.test.assertEqual(jsonData.items[0].pageorder,0,'content of the first entity should be 0');
    this.test.assertEqual(jsonData.items[0].content,'foobar-2','title of the first entity should be foobar');
});
*/
casper.thenOpen(baseURL + '/admin/rest/adminpage/url0', {
    method: "delete",
}, function() {
    this.echo("DELETE request has been sent.")
    this.test.assertHttpStatus(200);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.message, 'deleted', 'return deleted message');
});

casper.thenOpen(baseURL + '/admin/rest/adminpage/url1', {
    method: "delete",
}, function() {
    this.echo("DELETE request has been sent.")
    this.test.assertHttpStatus(200);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.message, 'deleted', 'return deleted message');
});

// add entity by post
casper.thenOpen(baseURL + '/admin/rest/adminpage', {
    method: "post",
    data: {
        title: 'title1',
        url: 'url1',
        pageorder: 1,
        content: 'foobar',
        displaypage: 'on',
    }
}, function() {
    this.echo("POST request has been sent.")
    this.test.assertHttpStatus(201);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.message, 'created', 'return created message');
});

// add entity by post
casper.thenOpen(baseURL + '/admin/rest/adminpage', {
    method: "post",
    data: {
        title: 'title0',
        url: 'url0',
        pageorder: 0,
        content: 'foobar',
        displaypage: '',
    }
}, function() {
    this.echo("POST request has been sent.")
    this.test.assertHttpStatus(201);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.message, 'created', 'return created message');
});

casper.thenOpen(baseURL + '/admin/form/#/adminpage/list', function() {
    this.test.assertHttpStatus(200);
});

casper.waitForSelector('tr td:first-child', function() {
    var tempArr = this.getElementsInfo('tr td:first-child');
    this.test.assertEqual(tempArr.length !== 0, true, 'check table list');
});

casper.then(function() {
    this.click('#createEntity');
});

casper.waitForSelector('form.form-horizontal', function() {
    this.fill('form.form-horizontal', {
        'displaypage': true,
        'title': 'title999',
        'url': 'url999',
        'pageorder': 999,
        'content': 'content999',
        'externalurl': 'link999',
    }, true)
});

casper.waitForSelector('tr td:first-child', function() {
    var tempArr = this.getElementsInfo('tr td:first-child');
    this.test.assertEqual(tempArr.length !== 0, true, 'check table list');
    this.mouseEvent('click', 'tr td:nth-child(3) a:first-child');
});

casper.waitForSelector('form.form-horizontal', function() {
    this.fill('form.form-horizontal', {
        'displaypage': false,
    }, true);
})

casper.waitForSelector('tr td:first-child', function() {
    this.mouseEvent('click', 'tr td:nth-child(3) a:first-child');
});

casper.waitForSelector('form.form-horizontal', function() {
    this.test.assertEqual(false, this.getFormValues('form.form-horizontal').displaypage, 'check data update');
})

var i = 0;
casper.repeat(30, function() {
    i++;
    this.thenOpen(baseURL + '/admin/rest/adminpage', {
        method: "post",
        data: {
            title: 'title' + i,
            url: 'url' + i,
            pageorder: i,
            content: 'foobar',
            displaypage: 'on',
        }
    }, function() {
        this.echo("POST request has been sent.")
        this.test.assertHttpStatus(201);
        var jsonData = JSON.parse(this.getPageContent());
        this.test.assertEqual(jsonData.message, 'created', 'return created message');
    });
});

//change offset value to match to perPage value in admin.go
casper.thenOpen(baseURL + '/admin/form/#/adminpage/list', function() {
    this.test.assertHttpStatus(200);
    this.test.assertExists('#nextButton', 'found next button');
    this.click('#nextButton');
    this.test.assertEqual(this.getCurrentUrl(), baseURL + '/admin/form/#/adminpage/list/20', 'next button works')
});

casper.waitForSelector('#previousButton', function() {
    this.test.assertEqual(this.exists('#nextButton'), false, 'not found next button');
    this.click('#previousButton');
    this.test.assertEqual(this.getCurrentUrl(), baseURL + '/admin/form/#/adminpage/list', 'previous button works')
});

casper.waitForSelector('tr td:first-child', function() {
    this.mouseEvent('click', 'tr td:nth-child(4) button');

});

casper.setFilter('page.confirm', function(message) {
    self.received = message;
    this.echo("message to confirm : " + message);
    return true;
});

/*
casper.wait(1000, function() {
    this.echo("I've waited for a second.");
});

casper.then(function(){
    var lists = this.getElementsInfo('tbody tr');
    this.test.assertEqual(lists.length,19,'list was reduced');
});
*/

casper.run();
