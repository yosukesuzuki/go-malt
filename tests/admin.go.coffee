baseURL = "http://localhost:8080"
casper.start()
casper.open baseURL + "/admin/"
casper.then ->
  @echo @getTitle()
  currentUrl = @getCurrentUrl()
  if currentUrl.match(/login/)
    @click "#admin"
    @thenClick "#submit-login"

# test of /admin/
casper.then ->
  @echo @getTitle()
  @test.assertHttpStatus 200

# test of /admin/rest/models
casper.thenOpen baseURL + "/admin/rest/models", ->
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.models.length, 2, "total count of models should be 2"

# test of /admin/rest/schema/adminpage
casper.thenOpen baseURL + "/admin/rest/schema/adminpage", ->
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())

#this.test.assertEqual(jsonData.fields.length,9,'total count of models should be 9');
#this.test.assertEqual(jsonData.model_name,'adminpage','title should be adminpage');

# add entity by post
casper.thenOpen baseURL + "/admin/rest/adminpage",
  method: "post"
  data:
    title: "title1"
    url: "url1"
    pageorder: 1
    content: "foobar"
    displaypage: "on"
, ->
  @echo "POST request has been sent."
  @test.assertHttpStatus 201
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "created", "return created message"

# add entity by post
casper.thenOpen baseURL + "/admin/rest/adminpage",
  method: "post"
  data:
    title: "title0"
    url: "url0"
    pageorder: 0
    content: "foobar"
    displaypage: ""
, ->
  @echo "POST request has been sent."
  @test.assertHttpStatus 201
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "created", "return created message"

casper.wait 1000, ->
  @echo "I've waited for a second."

casper.thenOpen baseURL + "/admin/rest/adminpage", ->
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.items[0].displaypage, false, "displaypage of the first entity should be false"
  @test.assertEqual jsonData.items[0].title, "title0", "title of the first entity should be title0"
  @test.assertEqual jsonData.items[0].url, "url0", "url of the first entity should be url0"
  @test.assertEqual jsonData.items[0].pageorder, 0, "content of the first entity should be 0"
  @test.assertEqual jsonData.items[0].content, "foobar", "title of the first entity should be foobar"
  @test.assertEqual jsonData.items[1].displaypage, true, "displaypage of the first entity should be true"
  @test.assertEqual jsonData.items[1].title, "title1", "title of the first entity should be title1"
  @test.assertEqual jsonData.items[1].url, "url1", "url of the first entity should be url1"
  @test.assertEqual jsonData.items[1].pageorder, 1, "content of the first entity should be 1"
  @test.assertEqual jsonData.items[1].content, "foobar", "title of the first entity should be foobar"


casper.thenOpen baseURL + "/admin/rest/adminpage/url0",
  method: "delete"
, ->
  @echo "DELETE request has been sent."
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "deleted", "return deleted message"

casper.thenOpen baseURL + "/admin/rest/adminpage/url1",
  method: "delete"
, ->
  @echo "DELETE request has been sent."
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "deleted", "return deleted message"

# add entity by post
casper.thenOpen baseURL + "/admin/rest/adminpage",
  method: "post"
  data:
    title: "title1"
    url: "url1"
    pageorder: 1
    content: "foobar"
    displaypage: "on"
, ->
  @echo "POST request has been sent."
  @test.assertHttpStatus 201
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "created", "return created message"

casper.thenOpen baseURL + "/admin/rest/adminpage",
  method: "post"
  data:
    title: "title0"
    url: "url0"
    pageorder: 0
    content: "foobar"
    displaypage: ""
, ->
  @echo "POST request has been sent."
  @test.assertHttpStatus 201
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "created", "return created message"

casper.thenOpen baseURL + "/admin/form/#/adminpage/list", ->
  @test.assertHttpStatus 200

casper.waitForSelector "tr td:first-child", ->
  tempArr = @getElementsInfo("tr td:first-child")
  @test.assertEqual tempArr.length isnt 0, true, "check table list"

casper.then ->
  @click "#createEntity"

casper.waitForSelector "form.form-horizontal", ->
  @fill "form.form-horizontal",
    displaypage: true
    title: "title999"
    url: "url999"
    pageorder: 999
    content: "content999"
    externalurl: "link999"
  , true

casper.waitForSelector "tr td:first-child", ->
  tempArr = @getElementsInfo("tr td:first-child")
  @test.assertEqual tempArr.length isnt 0, true, "check table list"
  @mouseEvent "click", "tr td:nth-child(3) a:first-child"

casper.waitForSelector "form.form-horizontal", ->
  @fill "form.form-horizontal",
    displaypage: false
  , true

casper.waitForSelector "tr td:first-child", ->
  @mouseEvent "click", "tr td:nth-child(3) a:first-child"

casper.waitForSelector "form.form-horizontal", ->
  @test.assertEqual false, @getFormValues("form.form-horizontal").displaypage, "check data update"

i = 0
casper.repeat 30, ->
  i++
  @thenOpen baseURL + "/admin/rest/adminpage",
    method: "post"
    data:
      title: "title" + i
      url: "url" + i
      pageorder: i
      content: "foobar"
      displaypage: "on"
  , ->
    @echo "POST request has been sent."
    @test.assertHttpStatus 201
    jsonData = JSON.parse(@getPageContent())
    @test.assertEqual jsonData.message, "created", "return created message"

#change offset value to match to perPage value in admin.go
casper.thenOpen baseURL + "/admin/form/#/adminpage/list", ->
  @test.assertHttpStatus 200
  @test.assertExists "#nextButton", "found next button"
  @click "#nextButton"
  @test.assertEqual @getCurrentUrl(), baseURL + "/admin/form/#/adminpage/list/20", "next button works"

casper.waitForSelector "#previousButton", ->
  @test.assertEqual @exists("#nextButton"), false, "not found next button"
  @click "#previousButton"
  @test.assertEqual @getCurrentUrl(), baseURL + "/admin/form/#/adminpage/list", "previous button works"

casper.waitForSelector "tr td:first-child", ->
  @mouseEvent "click", "tr td:nth-child(4) button"

casper.setFilter "page.confirm", (message) ->
  self.received = message
  @echo "message to confirm : " + message
  true

# test of /admin/rest/schema/adminpage
casper.thenOpen baseURL + "/admin/rest/schema/article", ->
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())

#this.test.assertEqual(jsonData.fields.length,9,'total count of models should be 9');
#this.test.assertEqual(jsonData.model_name,'adminpage','title should be adminpage');

# add entity by post
casper.thenOpen baseURL + "/admin/rest/article",
  method: "post"
  data:
    title: "title1"
    url: "url1"
    displaytime: "2014-07-25 12:00"
    pageorder: 1
    content: "foobar"
    tagstring: "tag"
    displaypage: "on"
, ->
  @echo "POST request has been sent."
  @test.assertHttpStatus 201
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "created", "return created message"

# add entity by post
casper.thenOpen baseURL + "/admin/rest/article",
  method: "post"
  data:
    title: "title0"
    url: "url0"
    displaytime: "2014-07-25 12:01"
    pageorder: 0
    content: "foobar"
    tagstring: "tag0"
    displaypage: ""
, ->
  @echo "POST request has been sent."
  @test.assertHttpStatus 201
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "created", "return created message"

casper.wait 1000, ->
  @echo "I've waited for a second."

casper.thenOpen baseURL + "/admin/rest/article", ->
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.items[0].displaypage, false, "displaypage of the first entity should be false"
  @test.assertEqual jsonData.items[0].title, "title0", "title of the first entity should be title0"
  @test.assertEqual jsonData.items[0].url, "url0", "url of the first entity should be url0"
  @test.assertEqual jsonData.items[0].pageorder, 0, "content of the first entity should be 0"
  @test.assertEqual jsonData.items[0].content, "foobar", "title of the first entity should be foobar"
  @test.assertEqual jsonData.items[0].displaytime.replace("T", " ").slice(0, 16), "2014-07-25 12:01", "title of the first entity should be foobar"

casper.thenOpen baseURL + "/admin/rest/article/url0",
  method: "delete"
, ->
  @echo "DELETE request has been sent."
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "deleted", "return deleted message"

casper.thenOpen baseURL + "/admin/rest/article/url1",
  method: "delete"
, ->
  @echo "DELETE request has been sent."
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "deleted", "return deleted message"

casper.thenOpen baseURL + "/admin/form/#/article/list", ->
  @test.assertHttpStatus 200

casper.then ->
  @click "#createEntity"

casper.waitForSelector "form.form-horizontal", ->
  @fill "form.form-horizontal",
    displaypage: true
    title: "title999"
    url: "url999"
    displaytime: "2014-07-25 12:01"
    pageorder: 999
    content: "content999"
    tagstring: "tag0"
    externalurl: "link999"
  , false
  @mouseEvent "click", "#submitButton"

j = 0
casper.repeat 30, ->
  j++
  @thenOpen baseURL + "/admin/rest/article",
    method: "post"
    data:
      title: "title" + j
      url: "url" + j
      displaytime: "2014-07-25 12:01"
      pageorder: j
      content: "foobar"
      tagstring: "tag" + j
      displaypage: "on"
  , ->
    @echo "POST request has been sent."
    @test.assertHttpStatus 201
    jsonData = JSON.parse(@getPageContent())
    @test.assertEqual jsonData.message, "created", "return created message"

#change offset value to match to perPage value in admin.go
casper.thenOpen baseURL + "/admin/form/#/article/list", ->
  @test.assertHttpStatus 200
  @test.assertExists "#nextButton", "found next button"
  @click "#nextButton"
  @test.assertEqual @getCurrentUrl(), baseURL + "/admin/form/#/article/list/20", "next button works"

casper.waitForSelector "#previousButton", ->
  @test.assertEqual @exists("#nextButton"), false, "not found next button"
  @click "#previousButton"
  @test.assertEqual @getCurrentUrl(), baseURL + "/admin/form/#/article/list", "previous button works"

casper.waitForSelector "tr td:first-child", ->
  @mouseEvent "click", "tr td:nth-child(4) button"

casper.setFilter "page.confirm", (message) ->
  self.received = message
  @echo "message to confirm : " + message
  true

casper.thenOpen baseURL + "/admin/image/upload/url", ->
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())

casper.run()
