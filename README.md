# wiregock
Small and very fast and stable implementation of Wiremock with Goland and MongoDB based of fiber lib.

## Request Matching

Stub matching and verification queries can use the following request attributes:

* URL
* HTTP Method
* Query parameters
* Form parameters
* Headers
* Basic authentication (a special case of header matching)
* Cookies
* Request body
* Traceparent

### HTTP methods

* **ANY** all methods are accepted
* **GET**
* **HEAD**
* **OPTIONS**
* **TRACE**
* **PUT**
* **DELETE**
* **POST**
* **PATCH**
* **CONNECT**

### Request mapping

* **urlPath**, **url** equality matching on path and query 
* **urlPattern** regex matching on path and query
* **method** HTTP method. To accept all, use **ANY**
* **headers**
* **queryParameters**
* **cookies**
* **bodyPatterns**
* **basicAuthCredentials**
* **matchingType** accept only **ALL** (default) params or **ANY** of params

### Comparation

* **equalTo** exact equality
* **binaryEqualTo** Unlike the above equalTo operator, this compares byte arrays (or their equivalent base64 representation).
* **contains** string contains the value
* **matches** compare by RegExp
* **wildcards** compare with wildcards (**\***, **?**)
* **equalToJson** if the attribute (most likely the request body in practice) is valid JSON and is a semantic match for the expected value.
* **equalToXml** if the attribute value is valid XML and is semantically equal to the expected XML document
* **matchesXPath** XPath matcher for XML objects.
* **ignoreArrayOrder** ignore order of array items
* **ignoreExtraElements** ignore extra elements of array items
* **matchesJsonPath** check by Json Path
* **matchesJsonSchema** check by Json Schema
* **includes** possible elements
* **hasExactly** exact elements 

### Templates

Templates are based on [mustache](https://mustache.github.io/) engine. There's support of default variable *request* based on request data.

* **request.id** - The unique ID of each request
* **request.url** - URL path and query
* **#request.queryFull.<key>** - values of a query parameter (zero indexed) e.g. *{{#request.queryFull.search}}{{.}}{{/request.queryFull.search}}*
* **request.query.<key>** - First value of a query parameter e.g. *request.query.search*
* **request.method** - request method e.g. *POST*
* **request.host** - hostname part of the URL e.g. *my.example.com*
* **request.port** - port number e.g. *8080*
* **request.scheme** - protocol part of the URL e.g. *https*
* **request.baseUrl** - URL up to the start of the path e.g. *https://my.example.com:8080*
* **#request.headersFull.<key>** - values of a header (zero indexed) e.g. *{{#request.headers.ManyThings}}{{.}}{{/request.headers.ManyThings}}*
* **request.headers.<key>** - first value of a request header e.g. *request.headers.X-Request-Id*
* **request.cookies.<key>** - First value of a request cookie e.g. *request.cookies.JSESSIONID*
* **request.body** - Request body text (avoid for non-text bodies)
* **request.bodyAsBase64** - The Base64 representation of the request body.

## To Be Implemented

### Comparation

* **matchesJsonPath** JSON matcher
* **matchesJsonSchema** JSON schema matcher

### Templates
* **xmlPath** and **jsonPath** helpers
* **request.parts** template for multipart files