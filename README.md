<p style="text-align: center"><img src="logo.png" alt="lhotse logo"/></p>
<h1 style="text-align: center">Lightweight HTTP server with controllable performance</h1>

<p style="text-align: center">
    <a href="https://github.com/oleiade/lhotse/releases"><img src="https://img.shields.io/github/release/oleiade/lhotse.svg" alt="release"></a>
    <a href="https://github.com/oleiade/lhotse/actions/workflows/build.yml"><img src="https://github.com/oleiade/lhotse/actions/workflows/build.yml/badge.svg" alt="Build status"></a>
</p>

Lhotse is a lightweight HTTP server designed for controlled performance testing.
It facilitates dictating the server's behavior, such as introducing specific latencies or generating responses of varying sizes.
Developed to complement [k6](https://github.com/grafana/k6), Lhotse is an excellent tool for measuring [k6](https://github.com/grafana/k6)'s  accuracy and assisting in its improvement and debugging.

## Installation

Install Lhotse using Go:

```bash
go install github.com/oleiade/lhotse@latest
```

## Usage & Examples

Lhotse provides functionalities such as latency simulation, data response control, and custom response generation.

Here are some examples:

### Latency Control

#### Fixed Duration

Command:
```bash
time curl -X GET -i http://localhost:3434/latency/500ms
```

Output:
```bash
HTTP/1.1 200 OK
Allow: GET, OPTIONS
Date: Tue, 01 Nov 2022 18:15:37 GMT
Content-Length: 12
Content-Type: text/plain; charset=utf-8

waited 500ms
```

#### Latency Range

Command:
```bash
curl -X GET -i http://localhost:3434/latency/500ms-1s
```

Output:
```bash
HTTP/1.1 200 OK
Allow: GET, OPTIONS
Date: Tue, 01 Nov 2022 18:16:43 GMT
Content-Length: 19
Content-Type: text/plain; charset=utf-8

waited 610.839158ms
```

### Data Response Control

#### Fixed Size

Command:
```bash
curl -X GET -i http://localhost:3434/data/8b
```

Output:
```bash
HTTP/1.1 200 OK
Allow: GET, OPTIONS
Date: Tue, 01 Nov 2022 18:14:12 GMT
Content-Length: 8
Content-Type: text/plain; charset=utf-8

TMtTCoaN
```

#### Size Range

Command:
```bash
curl -X GET -i http://localhost:3434/data/64b-128b
```

Output:
```bash
HTTP/1.1 200 OK
Allow: GET, OPTIONS
Date: Tue, 01 Nov 2022 18:13:02 GMT
Content-Length: 116
Content-Type: text/plain; charset=utf-8

whTHctcuAxhxKQFDaFpLSjFbcXoEFfRsWxPLDnJObCsNVlgTeMaPEZQleQYhYzRyJjPjzpfRFEgmotaFetHsbZRjxAwnwekrBEmfdzdcEkXBAkjQZLCt
```

### Custom Response Control

#### Custom Status Code

Command:
```bash
curl -X GET -i http://localhost:3434/response?status=204
```

Output:
```bash
HTTP/1.1 204 No Content
Content-Type: text/plain
Date: Sun, 14 Jan 2024 10:17:35 GMT
```

#### Custom content-type

Command:
```bash
curl -i -X GET --header 'Content-Type: application/json' http://localhost:3434/response
```

Output:
```bash
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 14 Jan 2024 10:18:36 GMT
Content-Length: 49

{"content_type":"application/json","status":200}
```

## API Reference

#### Latency Control

Endpoint `/latency/{duration}` simulates a response delay.

```http
  GET /latency/${duration}
```

| Parameter  | Type     | Description                                                                                   |
|:-----------|:---------|:----------------------------------------------------------------------------------------------|
| `duration` | `string` | Specifies the delay before responding. Use {value}{unit} or {lowerBound}-{upperBound} format. |

#### Data Volume Control

Endpoint `/data/{size}` controls the response size.

```http
  GET /data/${size}
```

| Parameter | Type     | Description                                                                                                                                                                                                                                                                                                                                                                                                                       |
|:----------|:---------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `size`    | `string` | Size of the payload produced by Lhotse. It can either be specified as a single size of the form {value}{unit}, or as a range {lowerBound}-{upperBound}. The value should always be an unsigned integer value. Valid units are `b`, `kb`, `mb`, and `gb`. When specifying a range, `lowerBound` needs to be less than `upperBound`, and as a result, the produced payload will be of a random size somewhere between those bounds. |

#### Custom Response Control

Endpoint `/response` allows customization of the response.

```http
  GET /response?status=${status}
```

##### Query Parameters

| Parameter | Type     | Description                                      |
|:----------|:---------|:-------------------------------------------------|
| `status`  | `string` | Specifies the HTTP status code for the response. |

##### Headers

| Header         | Type     | Description                                                                                                 |
|:---------------|:---------|:------------------------------------------------------------------------------------------------------------|
| `Content-Type` | `string` | Determines the `Content-Type` of the response. Currently `text/plain` and `application/json` are supported. |

## Contributing
Contributions to Lhotse are welcome! Whether it's bug reports, feature requests, or code contributions, please feel free to contribute. For more details, see CONTRIBUTING.md.

## License
Lhotse is open-sourced under the MIT License.