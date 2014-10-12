system = require 'system'

# load global settings
settings = require '../helpers/settings'

casper.test.begin 'test search', 3, (test) ->

  # test of /admin/rest/schema/adminpage
  casper.thenOpen settings.baseURL() + "/rest/search?keyword=foobar", ->
    test.assertHttpStatus 200
    jsonData = JSON.parse(@getPageContent())

  # add entity by post
  casper.thenOpen settings.baseURL() + "/",->
    test.assertHttpStatus 200
    @fill "form.navbar-form",
      keyword: "foobar"
    , true

  casper.waitForUrl "/search", ->
    test.assertHttpStatus 200

  casper.run ->
    do test.done
