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

# add content by post
casper.thenOpen baseURL + "/admin/rest/adminpage",
  method: "post"
  data:
    title: "About this site."
    url: "about"
    content: """
  This site is the showcase of GAE Go Starter.

  GAE, Google App Engine, is one of PaaS platform to build a web application.

  GAE has some excellent features for developers, for example, auto scaling and maintainance free key value store, taskque and so on.

  Go is the best perfermance language for web application, and GAE with go launchs new instances faster than python and java.

  This starter template is collection of sample implementations of GAE/Go, so developers can develop new applications based on this app.

  And this Starter app has content management functions, so it is easy to use this application as a corporate page or blog engine.
"""
    displaypage: "on"
    pageorder: "0"
, ->
  @echo "POST request has been sent."
  @test.assertHttpStatus 201
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "created", "return created message"

casper.thenOpen baseURL + "/admin/rest/adminpage",
  method: "post"
  data:
    title: "CMS included"
    url: "cms"
    content: """
  Malt has simple content management system.

  Markdown editor works on browser and you can drag & drop image files to insert.

  Draft content is supported. Save content which is under construction.


"""
    displaypage: "on"
, ->
  @echo "POST request has been sent."
  @test.assertHttpStatus 201
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "created", "return created message"

casper.thenOpen baseURL + "/admin/rest/adminpage",
  method: "post"
  data:
    title: "Search Engine"
    url: "searchengine"
    content: """
  Content which is published is automatically saved Full-Text Search index.

  So you can provide search engine function by default to your users.

"""
    displaypage: "on"
, ->
  @echo "POST request has been sent."
  @test.assertHttpStatus 201
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "created", "return created message"

casper.thenOpen baseURL + "/admin/rest/adminpage",
  method: "post"
  data:
    title: "Deploy in a minute"
    url: "deploy"
    content: """
  It's very very very easy to deploy this appliction.

  You have to type just 'goapp deploy'.

  You can get a scalable web application in a minute.
"""
    displaypage: "on"
, ->
  @echo "POST request has been sent."
  @test.assertHttpStatus 201
  jsonData = JSON.parse(@getPageContent())
  @test.assertEqual jsonData.message, "created", "return created message"

casper.run()
