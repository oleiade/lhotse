# Lhotse

Lhotse is a tiny HTTP server with controllable performance.

Lhotse lets you specify, directly in your request, the expected specific (or range of) performance of the server as it handles it. For instance, Lhotse supports setting the latency it should introduce as it produces a response to a request or lets you control the size of the response it's going to produce.

Lhotse was developed as a counterpart to [k6](https://github.com/grafana/k6), as a way of measuring its accuracy, as well as validating its behavior, to eventually facilitate improving and debugging it.

## Installation

Install Lhotse the standard Go way 

```bash
  go install github.com/oleiade/lhotse@latest
```

## Usage/Examples

### Controlling server's latency

#### Fixed duration

```bash
> time curl -X GET -i http://localhost:3434/latency/500ms
HTTP/1.1 200 OK
Allow: GET, OPTIONS
Date: Tue, 01 Nov 2022 18:15:37 GMT
Content-Length: 12
Content-Type: text/plain; charset=utf-8

waited 500ms
________________________________________________________
Executed in  508,63 millis    fish           external 
   usr time    7,72 millis  315,00 micros    7,41 millis 
   sys time    0,09 millis   86,00 micros    0,00 millis
```

#### Within a latency range

```bash
> time curl -X GET -i http://localhost:3434/latency/500ms-1s
HTTP/1.1 200 OK
Allow: GET, OPTIONS
Date: Tue, 01 Nov 2022 18:16:43 GMT
Content-Length: 19
Content-Type: text/plain; charset=utf-8

waited 610.839158ms
________________________________________________________
Executed in  620,86 millis    fish           external 
   usr time    9,02 millis  403,00 micros    8,62 millis 
   sys time    0,11 millis  110,00 micros    0,00 millis
```

### Producing a response payload


#### Of a fixed size

```bash
> curl -X GET -i http://localhost:3434/data/8b
HTTP/1.1 200 OK
Allow: GET, OPTIONS
Date: Tue, 01 Nov 2022 18:14:12 GMT
Content-Length: 8
Content-Type: text/plain; charset=utf-8

TMtTCoaN
```


#### With a size range

```bash
> curl -X GET -i http://localhost:3434/data/64b-128b
HTTP/1.1 200 OK
Allow: GET, OPTIONS
Date: Tue, 01 Nov 2022 18:13:02 GMT
Content-Length: 116
Content-Type: text/plain; charset=utf-8

whTHctcuAxhxKQFDaFpLSjFbcXoEFfRsWxPLDnJObCsNVlgTeMaPEZQleQYhYzRyJjPjzpfRFEgmotaFetHsbZRjxAwnwekrBEmfdzdcEkXBAkjQZLCt
```

## API Reference

#### Control latency

The `/latency` endpoint instructs Lhotse to take either the specified duration or a random duration within the specified bounds to respond to a request.

```http
  GET /latency/${duration}
```

| Parameter  | Type     | Description                                                                                                                                                                                                                                                                                                                                                                                                                |
| :--------- | :------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `duration` | `string` | **Required**. Duration Lhotse will sleep before replying to the request. It can either be specified as a single duration of the form {value}{unit}, or as a range {lowerBound}-{upperBound}. The value should always be an unsigned integer value. Valid units are `ns`, `us`, `ms`, `s`, `m', and `h`. When specifying a range, `lowerBound` needs to be less than `upperBound`; as a result, the time it takes for the server to respond will be a random duration between those bounds. |

#### Control data volume

The `/data` endpoint offers the user control of the type and size of the response produced by the server. The produced response payload will always be random. The size of the produced response payload will match either the exact value provided or will have an arbitrary size within the provided bounds. 

```http
  GET /data/${size}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `size`      | `string` | **Required**. Size of the payload produced by Lhotse. It can either be specified as a single size of the form {value}{unit}, or as a range {lowerBound}-{upperBound}. The value should always be an unsigned integer value. Valid units are `b`, `kb`, `mb`, and `gb`. When specifying a range, `lowerBound` needs to be less than `upperBound`, and as a result, the produced payload will be of a random size somewhere between those bounds. |
