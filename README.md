# Golang-REST-MongoDB
## Server Info
* **Language : Golang**
* **Database : MongoDB**
* **Note : Creation timestamp added to document**
## API Endpoints
### 1. Create an article
* Should be a POST request
* Use JSON request body
* URL should be ‘/articles’

### 2. Get an article using id
* Should be a GET request
* Id should be in the url parameter
* URL should be ‘/articles/<id here>’
### 3. List all articles
* Should be a GET request
* URL should be ‘/articles’
### 4. Search for an Article (search in title, subtitle, content) - TODO
* **Note: Stuck at text indexing**
* Should be a GET request
* Search term should be in the query parameter with key ‘q’
* URL should be ‘/articles/search?q=<search term here>’

## Features
### 1. Server thread safety
The server is thread safe. All multi-threaded operations will be managed by the inbuilt modules. 
### 2. Pagination - TODO
Feature almost complete
### 2. Unit Tests 
Unit Tests written for: 
* Create an article
* Get an article using id
* Get an article using id error
* List all articles
