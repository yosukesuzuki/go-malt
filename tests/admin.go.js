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

casper.run();
