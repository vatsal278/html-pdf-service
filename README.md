# HtmlPdfService

[![Build](https://github.com/vatsal278/html-pdf-service/actions/workflows/build.yml/badge.svg)](https://github.com/vatsal278/html-pdf-service/actions/workflows/build.yml) [![Test Cases](https://github.com/vatsal278/html-pdf-service/actions/workflows/test.yml/badge.svg)](https://github.com/vatsal278/html-pdf-service/actions/workflows/test.yml) [![Codecov](https://codecov.io/gh/vatsal278/html-pdf-service/branch/main/graph/badge.svg)](https://codecov.io/gh/vatsal278/html-pdf-service)

* This service was created using Golang. This service functions conversion of a HTML template.
* This utilises "github.com/SebastiaanKlippert/go-wkhtmltopdf" service for converting html file to pdf file.
* This service has used clean code principle and appropriate go directory structure.
* This service is completely unit tested and all the errors have been handled

## Starting the html-pdf-service

```
docker compose up
```
## API Spec

You can test the api using post man, just import the [collection](./docs/html-to-pdf-svc.postman_collection.json) into your postman app.

<table>
    <th>Path</th>
    <th>HTTP Method</th>
    <th>Request</th>
    <th>Response</th>
    <th>Description</th>
<tr>
<td>

`/v1/register`
</td>
<td>

`POST`
</td>
<td>

**In Request Body:**<br>
HTML template file
</td>
<td>

```json

{
    "status":  201,
    "message": "SUCCESS",
    "data": {
        "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8" // UUID of registered file
    }
}

```
</td>
<td>
Generates a UUID for the file. Use this UUID to generate PDF files using this template.
</td>
</tr>
<tr>
<td>

`/v1/generate/{id}`
</td>
<td>

`POST`
</td>
<td>

**In URL Path{id}:**<br>
6ba7b810-9dad-11d1-80b4-00c04fd430c8

**In Request Body:**<br>

```json
{
    "values": { // key-value pairs for the placeholders used in the template
        "placeholder-1": "value",
        "placeholder-2": value,
    }
}
```
</td>
<td>
6ba7b810-9dad-11d1-80b4-00c04fd430c8.pdf
</td>
<td>
Generates a PDF file for the registered UUID of the HTML template file.

The request to this must be a map which has the keys as the placeholders and the values to be substituted in its place.
</td>
</tr>
<tr>
<td>

`/v1/register/{id}`
</td>
<td>

`PUT`
</td>
<td>

**In URL Path{id}:**<br>
6ba7b810-9dad-11d1-80b4-00c04fd430c8
 
**In Request Body:**<br>
HTML template file
</td>
<td>

```json
{
    "status":  200,
    "message": "SUCCESS",
    "data": {
        "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
    }
}
```
</td>
<td>
Updates the HTML template with the new one.
</td>
</tr>
<tr>
<td>

`/v1/health`
</td>
<td>

`GET`
</td>
<td></td>
<td>

```json
{
    "status":  200,
    "message": "OK",
    "data": null
}
```
</td>
<td>
Health check endpoint to see if the service is okay.
</td>
</tr>
</table>


### In order to use the SDK functions:
* `go get` the package 
```
go get github.com/vatsal278/html-pdf-service
```
* `import` the `sdk` package in the source code
```
import "github.com/vatsal278/html-pdf-service/pkg/sdk"
```
* Get an instance to the SDK Wrapper, Passing in the url to a running html-pdf-service.
```
s := sdk.NewHtmlToPdfSvc("html-pdf-service url")
```
* Read a html file and pass the Byte slice of file to Register a html template file. A Uuid of the publisher will be returned that needs to be used to push message to the `channel`.
```
fileBytes, _ := os.ReadFile("path to html file")
uuid, _ := s.Register(fileBytes)
```
* To Generate the pdf from html file. 
It takes in template data in map[string]interface format and `uuid` which was recieved at time of registration of template . 
```
_ = s.GeneratePdf((map[string]interface{}{"data": "anydata"}, `uuid`)
```
* To Replace the template file pass the byte slice of template file and Uuid to Replace function.
```
fileBytes, _ := os.ReadFile("path to new html file")
_ = s.Replace(`fileBytes`, `uuid`)
```
* Examples of the sdk usage can be found [here](./examples/test.go)

## Additional read
* [Docs](./docs/README 1.md)
* To check the code coverage 
```
cd docs
go tool cover -html=coverage
```
