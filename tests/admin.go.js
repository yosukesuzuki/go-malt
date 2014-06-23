//execute # casperjs test tests/admin.go.js

baseURL = 'http://localhost:8080';

casper.start();

casper.open(baseURL+'/admin/');

casper.then(function() {
    this.echo(this.getTitle());
    var currentUrl = this.getCurrentUrl();
    if(currentUrl.match(/login/)){
        this.click('#admin');
        this.thenClick('#submit-login');
    }
});

// test of /admin/
casper.then(function(){
    this.echo(this.getTitle());
    this.test.assertHttpStatus(200);
});

// test of /admin/rest/metadata
casper.thenOpen(baseURL+'/admin/rest/metadata',function(){
    this.test.assertHttpStatus(200);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.models.length,2,'total count of models should be 2');
});

// test of /admin/rest/metadata/adminpage
casper.thenOpen(baseURL+'/admin/rest/metadata/adminpage',function(){
    this.test.assertHttpStatus(200);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.fields.length,9,'total count of models should be 9');
    this.test.assertEqual(jsonData.model_name,'adminpage','title should be adminpage');
});

casper.thenOpen(baseURL+'/admin/rest/adminpage', {
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
    this.test.assertEqual(jsonData.message,'created','return created message');
});

casper.thenOpen(baseURL+'/admin/rest/adminpage', {
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
    this.test.assertEqual(jsonData.message,'created','return created message');
});

casper.wait(1000, function() {
    this.echo("I've waited for a second.");
});

casper.thenOpen(baseURL+'/admin/rest/adminpage',function(){
    this.test.assertHttpStatus(200);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.items[0].displaypage,false,'displaypage of the first entity should be false');
    this.test.assertEqual(jsonData.items[0].title,'title0','title of the first entity should be title0');
    this.test.assertEqual(jsonData.items[0].url,'url0','url of the first entity should be url0');
    this.test.assertEqual(jsonData.items[0].pageorder,0,'content of the first entity should be 0');
    this.test.assertEqual(jsonData.items[0].content,'foobar','title of the first entity should be foobar');
    this.test.assertEqual(jsonData.items[1].displaypage,true,'displaypage of the first entity should be true');
    this.test.assertEqual(jsonData.items[1].title,'title1','title of the first entity should be title1');
    this.test.assertEqual(jsonData.items[1].url,'url1','url of the first entity should be url1');
    this.test.assertEqual(jsonData.items[1].pageorder,1,'content of the first entity should be 1');
    this.test.assertEqual(jsonData.items[1].content,'foobar','title of the first entity should be foobar');
});

/*
something wrong with PUT request by casperjs
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
    this.test.assertEqual(jsonData.items[0].displaypage,true,'displaypage of the first entity should be false');
    this.test.assertEqual(jsonData.items[0].title,'title0-2','title of the first entity should be title0');
    this.test.assertEqual(jsonData.items[0].url,'url0','url of the first entity should be url0');
    this.test.assertEqual(jsonData.items[0].pageorder,0,'content of the first entity should be 0');
    this.test.assertEqual(jsonData.items[0].content,'foobar-2','title of the first entity should be foobar');
});
*/

/*
casper.thenOpen(baseURL+'/admin/rest/adminpage/url0', {
    method: "delete",
}, function() {
    this.echo("DELETE request has been sent.")
    this.test.assertHttpStatus(200);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.message,'deleted','return deleted message');
});

casper.thenOpen(baseURL+'/admin/rest/adminpage/url1', {
    method: "delete",
}, function() {
    this.echo("DELETE request has been sent.")
    this.test.assertHttpStatus(200);
    var jsonData = JSON.parse(this.getPageContent());
    this.test.assertEqual(jsonData.message,'deleted','return deleted message');
});
*/

casper.run();
