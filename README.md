# Cryptos - Golang Hashing Service API

## Overview
This API provides a set of hash functions implemented in Golang. It allows clients to perform hashing operations on text inputs, returning the hashed result.

## Example Usage
To use the api, simply make a `Get` request to the endpoint, providing the `text` parameter with the value you wish to encrypt/decrypt.

Here's an example:
```bash
curl "https://cryptos.up.railway.app/sha1?text=hello"
```
> [!NOTE]
> There is a common pattern, endpoints like `sha1` and `sha256` can take a url parameter and if it can then the parameter name is `text`

The API will respond with a JSON string containing the resuled hash.
```json
"d9e989f651cdd269d7f9bb8a215d024d8d283688"
```
