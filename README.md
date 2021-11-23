# Uniswap V3 REST API via The Graph

**Objective:** Build a server-side REST API in Golang that uses The Graph's GraphQL API to provide Uniswap v3 information upon user request.

**Details:**
- using free, legacy [subgraph](https://thegraph.com/hosted-service/subgraph/ianlapham/uniswap-v3-alt)
- all monetary amounts returns will be in USD
- all time input and ouput will be Unix timestamps

**Routes:**
- `api/v1/pools/` - returns pools that include **given asset**
- `api/v1/volume/` - returns total volume of **given asset** swapped in **given time range**

**Bonus Routes:**
- `api/v1/swaps/{block_num}` - returns list of swaps that occurred during specific block
- `api/v1/assets/{block_num}` - returns list of all assets swapped during specific block

Route organization subject to change, but will continue to be answer the same questions.


