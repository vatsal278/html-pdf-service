# HTML to PDF File Service
The service aims to provide the feature of generating a PDF file from a HTML template using [wkhtmltopdf](https://wkhtmltopdf.org/ "wkhtmltopdf").

This provides a RESTful service that allows for a HTML template file to be uploaded, which in turn can be used for subsequent calls to generate a PDF out of this.

The entire service is written in pure [Go](https://go.dev/ "Go") and makes use of Gos [template feature](https://www.practical-go-lessons.com/chap-32-templates "template feature").

Internally this service uses [Redis](https://redis.io/ "Redis") to store the HTML template during registration for retrieval when generating the PDF file.

### API Specification
- The response structure will always be as follows:
```json

{
    "status": <HTTP Status code>,
    "message": "<SOME MESSAGE>",
    "data": <object containing information for the request>
}

```
<table>
    <th>Path</th>
    <th>HTTP Method</th>
    <th>Request</th>
    <th>Response</th>
    <th>Description</th>
    <th>Comments</th>
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
<td>

- [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier "UUID")

- [Go UUID library](https://github.com/google/uuid "Go UUID library")

- [Go HTML Template](https://www.practical-go-lessons.com/chap-32-templates "Go HTML Template")

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
<td></td>
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
<td></td>
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
<td></td>
</tr>
</table>

### Considerations
1. Use basic [net/http](https://pkg.go.dev/net/http "net/http") package or at most [gorilla/mux](https://pkg.go.dev/github.com/gorilla/mux "gorilla/mux") of Go.
2. The service must be [dockerized](https://docs.docker.com/get-started/ "dockerized") with a [docker compose](https://docs.docker.com/compose/ "docker compose") file to manage the service(s).
3. Proper HTTP Headers must be sent in the response to identify the type of data being sent.
4. The default response from the service must be in JSON in the specified format, for all scenarios.
5. SDK functions may be used to prevent clubbing the functionality of various services(Redis, PDF generation). These functions must be fully unit test tested with minimum 85% code coverage.
6. Unit tested code with a minimum coverage of 80% and business scenarios to be demonstrated in the form of unit/system test cases.
7. Must follow proper CLEAN code principles with appropriate Go project structure.
8. There must be a timeout if the request is being processed for more than 5 seconds.
9. Verify that there is no access restriction to the file.

### References
- [HTTP request for file upload](https://stackoverflow.com/questions/14962592/whats-content-type-value-within-a-http-request-when-uploading-content)
- [Content-Disposition HEADER](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition)
- [go-wkhtmltopdf](https://github.com/SebastiaanKlippert/go-wkhtmltopdf)
- [Installing package on deb](https://wireframesketcher.com/support/install/installing-deb-package-on-ubuntu-debian.html)
- [Using wget to download files in deb](https://stackoverflow.com/questions/14306382/how-to-rename-the-downloaded-file-with-wget)
- [Creating Go HTML templates](https://www.practical-go-lessons.com/chap-32-templates "Creating Go HTML templates")