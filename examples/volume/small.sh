#!/bin/bash
curl -i -X GET -H Content-Type: application/json -d'{"id":"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48","start_timestamp":1637366400,"end_timestamp":1637712000}' http://127.0.0.1:8080/asset/volume
