# Cryptos - Golang Hashing Service API

## Overview
This API provides a set of hash functions implemented in Golang. It allows clients to perform hashing operations on text inputs, returning the hashed result.

## Usage
To use the API, you need to make a `Get` request to the `/hashType` endpoint with a URL parameter named `text`. The text to be hashed should be URL encoded. Here's an example using curl:

```bash
curl "http://localhost:2323/sha1sum?text=example%20text"
```
The API will respond with a JSON string containing the original text and the hashed result.
