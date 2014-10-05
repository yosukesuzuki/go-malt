system = require 'system'

# load global settings
settings = require '../helpers/settings'

casper.test.begin 'test article', 102, (test) ->

  # test of /admin/rest/schema/adminpage
  casper.thenOpen settings.baseURL() + "/admin/rest/schema/article", ->
    test.assertHttpStatus 200
    jsonData = JSON.parse(@getPageContent())


  # add entity by post
  casper.thenOpen settings.baseURL() + "/admin/rest/article",
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
    test.assertHttpStatus 201
    jsonData = JSON.parse(@getPageContent())
    test.assertEqual jsonData.message, "created", "return created message"

  # add entity by post
  casper.thenOpen settings.baseURL() + "/admin/rest/article",
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
    test.assertHttpStatus 201
    jsonData = JSON.parse(@getPageContent())
    test.assertEqual jsonData.message, "created", "return created message"

  casper.wait 1000, ->
    @echo "I've waited for a second."

  casper.thenOpen settings.baseURL() + "/admin/rest/article", ->
    test.assertHttpStatus 200
    jsonData = JSON.parse(@getPageContent())
    test.assertEqual jsonData.items[0].displaypage, false, "displaypage of the first entity should be false"
    test.assertEqual jsonData.items[0].title, "title0", "title of the first entity should be title0"
    test.assertEqual jsonData.items[0].url, "url0", "url of the first entity should be url0"
    test.assertEqual jsonData.items[0].content, "foobar", "title of the first entity should be foobar"
    test.assertEqual jsonData.items[0].displaytime.replace("T", " ").slice(0, 16), "2014-07-25 12:01", "title of the first entity should be foobar"

  casper.thenOpen settings.baseURL() + "/admin/rest/article/url0",
    # put method does not work correctly in casperjs,no data is sent
    method: "post"
    data:
      pageorder: "0"
  , ->
    @echo "page order update request has been sent."
    test.assertHttpStatus 200
    jsonData = JSON.parse(@getPageContent())
    test.assertEqual jsonData.message, "updated", "return created message"

  casper.thenOpen settings.baseURL() + "/admin/rest/article/url0", ->
    @echo "check if page order is updated"
    test.assertHttpStatus 200
    jsonData = JSON.parse(@getPageContent())
    test.assertEqual jsonData.item.pageorder, 0, "page order should be 0"
    test.assertEqual jsonData.item.title, "title0", "title should be title0"
    test.assertEqual jsonData.item.url, "url0", "url should be url0"
    test.assertEqual jsonData.item.content, "foobar", "content should be foobar"

  casper.thenOpen settings.baseURL() + "/admin/rest/article/url0",
    # put method does not work correctly in casperjs,no data is sent
    # so use post instead of put here
    method: "post"
    data:
      title: "title0-2"
      url: "url0-2"
  , ->
    @echo "partial update request has been sent."
    test.assertHttpStatus 200
    jsonData = JSON.parse(@getPageContent())
    test.assertEqual jsonData.message, "updated", "return created message"

  casper.thenOpen settings.baseURL() + "/admin/rest/article/url0", ->
    @echo "check partial update"
    test.assertHttpStatus 200
    jsonData = JSON.parse(@getPageContent())
    test.assertEqual jsonData.item.pageorder, 0, "page order should be 0"
    test.assertEqual jsonData.item.title, "title0-2", "title should be title0-2"
    test.assertEqual jsonData.item.url, "url0", "url should be url0"
    test.assertEqual jsonData.item.content, "foobar", "content should be foobar"

  casper.thenOpen settings.baseURL() + "/admin/rest/article/url0",
    # put method does not work correctly in casperjs,no data is sent
    # so use post instead of put here
    method: "post"
    data:
      draft: "on"
      content: "foobar draft"
  , ->
    @echo "draft put request has been sent."
    test.assertHttpStatus 200
    jsonData = JSON.parse(@getPageContent())
    test.assertEqual jsonData.message, "updated", "return created message"

  casper.thenOpen settings.baseURL() + "/admin/rest/article/url0", ->
    @echo "check draft update"
    test.assertHttpStatus 200
    jsonData = JSON.parse(@getPageContent())
    test.assertEqual jsonData.is_draft, true, "draft flg is true"

  casper.thenOpen settings.baseURL() + "/admin/rest/article/url0",
    method: "delete"
  , ->
    @echo "DELETE request has been sent."
    test.assertHttpStatus 200
    jsonData = JSON.parse(@getPageContent())
    test.assertEqual jsonData.message, "deleted", "return deleted message"

  casper.thenOpen settings.baseURL() + "/admin/rest/article/url1",
    method: "delete"
  , ->
    @echo "DELETE request has been sent."
    test.assertHttpStatus 200
    jsonData = JSON.parse(@getPageContent())
    test.assertEqual jsonData.message, "deleted", "return deleted message"

  casper.thenOpen settings.baseURL() + "/admin/form/#/article/list", ->
    test.assertHttpStatus 200

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
    @thenOpen settings.baseURL() + "/admin/rest/article",
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
      test.assertHttpStatus 201
      jsonData = JSON.parse(@getPageContent())
      test.assertEqual jsonData.message, "created", "return created message"

  #change offset value to match to perPage value in admin.go
  casper.thenOpen settings.baseURL() + "/admin/form/#/article/list", ->
    test.assertHttpStatus 200
    test.assertExists "#nextButton", "found next button"
    @click "#nextButton"
    test.assertEqual @getCurrentUrl(), settings.baseURL() + "/admin/form/#/article/list/20", "next button works"

  casper.waitForSelector "#previousButton", ->
    test.assertEqual @exists("#nextButton"), false, "not found next button"
    @click "#previousButton"

  casper.waitForSelector "#nextButton", ->
    test.assertEqual @getCurrentUrl(), settings.baseURL() + "/admin/form/#/article/list", "previous button works"
    @click "#sortMode"

  casper.waitForSelector "ul.sortable", ->
    test.assertEqual @getCurrentUrl(), settings.baseURL() + "/admin/form/#/article/sort/NaN/20", "sort button works"

  casper.thenOpen settings.baseURL() + "/admin/form/#/article/list", ->
    @echo "back to list"

  casper.waitForSelector "tr td:nth-child(4) button", ->
    @mouseEvent "click", "tr td:nth-child(4) button"

  casper.setFilter "page.confirm", (message) ->
    self.received = message
    @echo "message to confirm : " + message
    true

  casper.thenOpen settings.baseURL() + "/admin/image/upload/url", ->
    test.assertHttpStatus 200
    jsonData = JSON.parse(@getPageContent())
    test.assertMatch jsonData.uploadurl,/^\/_ah\/upload/,"test upload url"

  casper.run ->
    do test.done
