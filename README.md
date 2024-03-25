# Cryptos - Golang Hashing Service API

## Overview
This API provides a set of hash functions implemented in Golang. It allows clients to perform hashing operations on text inputs, returning the hashed result.

## Example Usage
To use the api, simply make a `Get` request to the endpoint, with the value you wish to encrypt/decrypt.

Here's an example:
```bash
curl "https://cryptos.up.railway.app/sha1/hello"
```

The API will respond with a JSON string containing the resuled hash.
```json
"d9e989f651cdd269d7f9bb8a215d024d8d283688"
```
