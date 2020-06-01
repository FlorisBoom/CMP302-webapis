# CMP302 Web API's

* Get authorization, which returns an auth token <br>
`https://cpm302-webapis.herokuapp.com/authorize`

* Get all cars <br>
`GET https://cpm302-webapis.herokuapp.com/cars` `-H "Authorization:{token}"`

* Get car by id <br>
`GET https://cpm302-webapis.herokuapp.com/cars/{id}` `-H "Authorization:{token}"`

* Create car <br>
`POST https://cpm302-webapis.herokuapp.com/cars` `-H "Authorization:{token}"` `-d {raw json}`

* Update car <br>
`PUT https://cpm302-webapis.herokuapp.com/cars/{id}` `-H "Authorization:{token}"` `-d {raw json}`

* Delete car <br>
`DELETE https://cpm302-webapis.herokuapp.com/cars/{id}` `-H "Authorization:{token}"`

* Example <br>
`curl -X POST "https://cpm302-webapis.herokuapp.com/cars" -H "Authorization:{your access token}" -d "{"Brand":"Tesla","Model":"S","Year":2020,"Color":"Blue"}`
