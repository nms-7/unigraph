# Uniswap V3 REST API via The Graph

**Objective:** Build a server-side REST API in Golang that uses The Graph's GraphQL API to provide Uniswap v3 information upon user request.

**Details:**
- using free, legacy [subgraph](https://thegraph.com/hosted-service/subgraph/ianlapham/uniswap-v3-alt)
- all monetary amounts returns will be in USD
- all time input and ouput will be Unix timestamps

**Routes:**
- `api/v1/pools` - returns pools that include **given asset:**
	- `id`: assetId (string)
- `api/v1/volume` - returns total volume of **given asset** swapped in **given time range:**
	- `id`: assetId (string)
	- `start_timestamp`: start unix timestamp (int)
	- `end_timestamp`: end unix timestamp (int)

**Bonus Routes:**
- `api/v1/swaps` - returns list of swaps that occurred during specific block **given block number:**
	- `block_number` (int)
- `api/v1/assets` - returns list of all assets swapped during specific block **given block number:**
	- `block_number` (int)

Route organization subject to change, but will continue to be answer the same questions.

### To RUN

**with included build file:** `./unigraph`

try the example queries in the `examples` directory by running the individual shell files, or build your own queries with curl, postman, or any other http requests utility


