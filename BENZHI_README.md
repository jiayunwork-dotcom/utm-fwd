# utm-fwd — Go WGS84 UTM 正算地理坐标转换 HTTP 服务

本 WGS84 UTM 正算 HTTP 服务：给定经纬度算出带号与投影坐标，或按点求带；非法经纬或超出投影范围须报错。

## Build / Run / Test

```text
go build ./...
go run . serve -addr :8080
go run . example -file example/beijing.json
go test ./...
```

## Evaluation Image

Evaluation-specific files (do not overwrite project Dockerfile/README):

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md` (this file)

Build and verify in container:

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
