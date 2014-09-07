# default setting == localdev
env = "localdev"

settings = {
  localdev : "http://localhost:8080"
  gaedev : "https://dev-goappstarter.appspot.com"
}

if casper.cli.has('env')
  env = casper.cli.get('env')

baseURL = settings[env]
casper.start()
casper.open baseURL + "/admin/"
casper.then ->
  if env is "localdev"
    @echo @getTitle()
    currentUrl = @getCurrentUrl()
    if currentUrl.match(/login/)
      @click "#admin"
      @thenClick "#submit-login"
  else if env is "gaedev"
    #fill id and password passed from cli
    user_id = casper.cli.get("user_id")
    user_password = casper.cli.get("user_password")
    @fill "form#gaia_loginform",
      Email: user_id
      Passwd:user_password 
    , true

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
  @test.assertEqual jsonData.items[0].content, "foobar", "title of the first entity should be foobar"
  @test.assertEqual jsonData.items[1].displaypage, true, "displaypage of the first entity should be true"
  @test.assertEqual jsonData.items[1].title, "title1", "title of the first entity should be title1"
  @test.assertEqual jsonData.items[1].url, "url1", "url of the first entity should be url1"
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
    content: "foobar"
    displaypage: ""
, ->
  @echo "POST request has been sent."
  @test.assertHttpStatus 201
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "created", "return created message"

casper.thenOpen baseURL + "/admin/rest/adminpage/url0",
  # put method does not work correctly in casperjs,no data is sent
  method: "post"
  data:
    pageorder: "0"
, ->
  @echo "page order update request has been sent."
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "updated", "return created message"

casper.thenOpen baseURL + "/admin/rest/adminpage/url0", ->
  @echo "check if page order is updated"
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.item.pageorder, 0, "page order should be 0"
  @test.assertEqual jsonData.item.title, "title0", "title should be title0"
  @test.assertEqual jsonData.item.url, "url0", "url should be url0"
  @test.assertEqual jsonData.item.content, "foobar", "content should be foobar"

casper.thenOpen baseURL + "/admin/rest/adminpage/url0",
  # put method does not work correctly in casperjs,no data is sent
  # so use post instead of put here
  method: "post"
  data:
    title: "title0-2"
    url: "url0-2"
, ->
  @echo "partial update request has been sent."
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "updated", "return created message"

casper.thenOpen baseURL + "/admin/rest/adminpage/url0", ->
  @echo "check partial update"
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.item.pageorder, 0, "page order should be 0"
  @test.assertEqual jsonData.item.title, "title0-2", "title should be title0-2"
  @test.assertEqual jsonData.item.url, "url0", "url should be url0"
  @test.assertEqual jsonData.item.content, "foobar", "content should be foobar"

casper.thenOpen baseURL + "/admin/form/#/adminpage/list", ->
  @test.assertHttpStatus 200

casper.waitForSelector "tr td:first-child", ->
  tempArr = @getElementsInfo("tr td:first-child")
  @test.assertEqual tempArr.length isnt 0, true, "check table list"

casper.then ->
  @click "#createEntity"

casper.waitForSelector "form.form", ->
  @fill "form.form",
    displaypage: true
    title: "title999"
    url: "url999"
    content: "content999"
    externalurl: "link999"
  , true

casper.waitForSelector "tr td:first-child", ->
  tempArr = @getElementsInfo("tr td:first-child")
  @test.assertEqual tempArr.length isnt 0, true, "check table list"
  @mouseEvent "click", "tr td:nth-child(3) a:first-child"

casper.waitForSelector "form.form", ->
  @fill "form.form",
    displaypage: false
  , true

casper.waitForSelector "tr td:first-child", ->
  @mouseEvent "click", "tr td:nth-child(3) a:first-child"

casper.waitForSelector "form.form", ->
  @test.assertEqual false, @getFormValues("form.form").displaypage, "check data update"

i = 0
casper.repeat 30, ->
  i++
  @thenOpen baseURL + "/admin/rest/adminpage",
    method: "post"
    data:
      title: "title" + i
      url: "url" + i
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

casper.waitForSelector "#nextButton", ->
  @test.assertEqual @getCurrentUrl(), baseURL + "/admin/form/#/adminpage/list", "previous button works"
  @click "#sortMode"

casper.waitForSelector "ul.sortable", ->
  @test.assertEqual @getCurrentUrl(), baseURL + "/admin/form/#/adminpage/sort/NaN/20", "sort button works"

casper.thenOpen baseURL + "/admin/form/#/adminpage/list", ->
  @echo "back to list"
 
casper.waitForSelector "tr td:nth-child(4) button", ->
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
  @test.assertEqual jsonData.items[0].content, "foobar", "title of the first entity should be foobar"
  @test.assertEqual jsonData.items[0].displaytime.replace("T", " ").slice(0, 16), "2014-07-25 12:01", "title of the first entity should be foobar"

casper.thenOpen baseURL + "/admin/rest/article/url0",
  # put method does not work correctly in casperjs,no data is sent
  method: "post"
  data:
    pageorder: "0"
, ->
  @echo "page order update request has been sent."
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "updated", "return created message"

casper.thenOpen baseURL + "/admin/rest/article/url0", ->
  @echo "check if page order is updated"
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.item.pageorder, 0, "page order should be 0"
  @test.assertEqual jsonData.item.title, "title0", "title should be title0"
  @test.assertEqual jsonData.item.url, "url0", "url should be url0"
  @test.assertEqual jsonData.item.content, "foobar", "content should be foobar"

casper.thenOpen baseURL + "/admin/rest/article/url0",
  # put method does not work correctly in casperjs,no data is sent
  # so use post instead of put here
  method: "post"
  data:
    title: "title0-2"
    url: "url0-2"
, ->
  @echo "partial update request has been sent."
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "updated", "return created message"

casper.thenOpen baseURL + "/admin/rest/article/url0", ->
  @echo "check partial update"
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.item.pageorder, 0, "page order should be 0"
  @test.assertEqual jsonData.item.title, "title0-2", "title should be title0-2"
  @test.assertEqual jsonData.item.url, "url0", "url should be url0"
  @test.assertEqual jsonData.item.content, "foobar", "content should be foobar"

casper.thenOpen baseURL + "/admin/rest/article/url0",
  # put method does not work correctly in casperjs,no data is sent
  # so use post instead of put here
  method: "post"
  data:
    draft: "on"
    content: "foobar draft"
, ->
  @echo "draft put request has been sent."
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "updated", "return created message"

casper.thenOpen baseURL + "/admin/rest/article/url0", ->
  @echo "check draft update"
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.is_draft, true, "draft flg is true"

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

casper.waitForSelector "form.form", ->
  @fill "form.form",
    displaypage: true
    title: "title999"
    url: "url999"
    displaytime: "2014-07-25 12:01"
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

casper.waitForSelector "#nextButton", ->
  @test.assertEqual @getCurrentUrl(), baseURL + "/admin/form/#/article/list", "previous button works"
  @click "#sortMode"

casper.waitForSelector "ul.sortable", ->
  @test.assertEqual @getCurrentUrl(), baseURL + "/admin/form/#/article/sort/NaN/20", "sort button works"

casper.thenOpen baseURL + "/admin/form/#/article/list", ->
  @echo "back to list"

casper.waitForSelector "tr td:nth-child(4) button", ->
  @mouseEvent "click", "tr td:nth-child(4) button"

casper.setFilter "page.confirm", (message) ->
  self.received = message
  @echo "message to confirm : " + message
  true

casper.thenOpen baseURL + "/admin/image/upload/url", ->
  @test.assertHttpStatus 200
  jsonData = JSON.parse(@getPageContent())
  @test.assertMatch jsonData.uploadurl,/^\/_ah\/upload/,"test upload url"


# add content by post
casper.thenOpen baseURL + "/admin/rest/adminpage",
  @echo "create about page"
  method: "post"
  data:
    title: "About this site."
    url: "about"
    content: """This site is the showcase of GAE Go Starter.
      GAE, Google App Engine, is one of PaaS platform to build a web application.
      GAE has some excellent features for developers, for example, auto scaling and maintainance free key value store, taskque and so on.
      Go is the best perfermance language for web application, and GAE with go launchs new instances fastest than python and java.
      This starter template is collection of sample implementations of GAE/Go, so developers can develop new applications based on this app.
      And this Starter app has content management functions, so it is easy to use this application as a corporate page or blog engine.
      """
    displaypage: "on"
, ->
  @echo "POST request has been sent."
  @test.assertHttpStatus 201
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "created", "return created message"


casper.run()
