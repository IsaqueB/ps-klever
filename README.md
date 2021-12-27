# About this API
This API is a simple video-upvote API developed using Golang and comunicating using gRPC and HTTP.

This documentation describes which routes it accepts and their purpose.

# Vote
A VOTE is a document which stores:
* An unique ID 
* The video that was reacted
* The user that reacted
* The "upvote" value. 
  * True means that it was an upvote
  * False means that it was a downvote

# Routes
## HTTP
## Create an upvote
Creates a new `vote`
### Path
```http
POST /v1
```
### Body
```javascript
{
  "vote": {
    "video": string,
    "user": string,
    "upvote": boolean
  }
}
```
| Parameter| Description |
| :--- | :--- |
| `video` |  is the video that was reacted |
| `user` |  is the user that reacted |
| `upvote` |  is the value |

### Response
If success, the answer will be:
```javascript
{
  "id": string
}
```
| Parameter| Description |
| :--- | :--- |
| `id` |  is the objectId of the `vote` created |

If error, the answer will be:
```javascript
{
  "code": int,
  "message": string,
  "details": []
}
```
| Parameter| Description |
| :--- | :--- |
| `code` |  is the grpc error code |
| `message` |  is a description of the error |
| `details` |  are details to the error that occurred, if any |

## Get an upvote
Searches for an `vote` in the database with the id passed in the URL
### Path
```http
POST /v1/{id}
```
| Parameter| Description |
| :--- | :--- |
| `id` |  is the id of the `vote` |

### Response
If success, the answer will be:
```javascript
{
  "vote": {
    "id": string
    "video": string,
    "user": string,
    "upvote": boolean
  }
}
```
| Parameter| Description |
| :--- | :--- |
| `id` |  is the objectId of the `vote` found |
| `video` |  is the id of the video |
| `user` |  is the id of the user |
| `upvote` |  is the value of the `vote` found |

If error, the answer will be:
```javascript
{
  "code": int,
  "message": string,
  "details": []
}
```
| Parameter| Description |
| :--- | :--- |
| `code` |  is the grpc error code |
| `message` |  is a description of the error |
| `details` |  are details to the error that occurred, if any |

## Update an upvote
Updates the value of *upvote* of a `vote` to *new_value*
### Path
```http
PUT /v1
```
### Body
```javascript
{
  "vote": {
    "id": string,
    "new_value": boolean
  }
}
```
| Parameter| Description |
| :--- | :--- |
| `id` |  is the objectId of the `vote` you want to update |
| `new_value` |  is the new value for this upvote |

### Response
If success, the answer will be:
```javascript
{
  "matched": int,
  "modified": int
}
```
| Parameter| Description |
| :--- | :--- |
| `matched` |  the amount of documents that matched the id sent |
| `modified` |  the amount of document modified by the query |

If error, the answer will be:
```javascript
{
  "code": int,
  "message": string,
  "details": []
}
```
| Parameter| Description |
| :--- | :--- |
| `code` |  is the grpc error code |
| `message` |  is a description of the error |
| `details` |  are details to the error that occurred, if any |

## Delete an upvote
Deletes a `vote`
### Path
```http
DELETE /v1/{id}
```
| Parameter| Description |
| :--- | :--- |
| `id` |  is the id of the `vote` |

### Response
If success, the answer will be:
```javascript
{
  "deleted": int
}
```
| Parameter| Description |
| :--- | :--- |
| `deleted` |  is the amount of document deleted |

If error, the answer will be:
```javascript
{
  "code": int,
  "message": string,
  "details": []
}
```
| Parameter| Description |
| :--- | :--- |
| `code` |  is the grpc error code |
| `message` |  is a description of the error |
| `details` |  are details to the error that occurred, if any |

## List votes of a video
Finds all votes related to a video
### Path
```http
GET /v1/video/{id}
```
| Parameter| Description |
| :--- | :--- |
| `id` |  is the id of the video that was reacted |
### Response
If success, the answer will be:
```javascript
{
	"vote": []
}
```
| Parameter| Description |
| :--- | :--- |
| `vote` |  array with all the votes related to this video |

If error, the answer will be:
```javascript
{
  "code": int,
  "message": string,
  "details": []
}
```
| Parameter| Description |
| :--- | :--- |
| `code` |  is the grpc error code |
| `message` |  is a description of the error |
| `details` |  are details to the error that occurred, if any |

## List votes of an user
Finds all votes done by an user
### Path
```http
GET /v1/user/{id}
```
| Parameter| Description |
| :--- | :--- |
| `id` |  is the id of the user |
### Response
If success, the answer will be:
```javascript
{
	"vote": []
}
```
| Parameter| Description |
| :--- | :--- |
| `vote` |  array with all the votes related to this user |

If error, the answer will be:
```javascript
{
  "code": int,
  "message": string,
  "details": []
}
```
| Parameter| Description |
| :--- | :--- |
| `code` |  is the grpc error code |
| `message` |  is a description of the error |
| `details` |  are details to the error that occurred, if any |

## gRPC
