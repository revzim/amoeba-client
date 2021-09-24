# Go Amoeba Client

a [amoeba][1] client for golang updated for use with [amoeba][2]

[1]: https://github.com/NetEase/pomelo
[2]: https://github.com/revzim/amoeba
[3]: ./example/basic/js
### [How to use with GopherJS][3]:
`Run -> RunJS`

```
go run .\example\basic\main.go --id="test" --token='eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzIyNDU1MDIsImlhdCI6MTYzMjI0NDAwMiwiaWQiOiJ0ZXN0IiwibmFtZSI6InRlc3QgcGVyc29uIiwibmJmIjoxNjMyMjQzOTkyfQ.v07XyWCYX1ykMyoU2lbxlcpEzKyXw0sl40gyVqcD4Qc'
// Output --
2021/09/21 13:10:15 attempting to connect to: ws://127.0.0.1:8080/ws?id=test&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzIyNDU1MDIsImlhdCI6MTYzMjI0NDAwMiwiaWQiOiJ0ZXN0IiwibmFtZSI6InRlc3QgcGVyc29uIiwibmJmIjoxNjMyMjQzOTkyfQ.v07XyWCYX1ykMyoU2lbxlcpEzKyXw0sl40gyVqcD4Qc...
2021/09/21 13:10:16 200
2021/09/21 13:10:16 connected to server at: ws://127.0.0.1:8080/ws?id=test&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzIyNDU1MDIsImlhdCI6MTYzMjI0NDAwMiwiaWQiOiJ0ZXN0IiwibmFtZSI6InRlc3QgcGVyc29uIiwibmJmIjoxNjMyMjQzOTkyfQ.v07XyWCYX1ykMyoU2lbxlcpEzKyXw0sl40gyVqcD4Qc
2021/09/21 13:10:17 onMembers {"members":["5b4949004d"]}
2021/09/21 13:10:18 room join: {"code":0,"result":"success","username":"87367791cd"}
exit status 1
```

## Install

```shell
go get -u github.com/revzim/amoeba-client
```

## Example

see example/basic