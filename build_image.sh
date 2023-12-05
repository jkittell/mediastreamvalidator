docker buildx create \
--use \
--bootstrap \
#--config buildkit.toml

docker buildx build --push \
--platform linux/arm64/v8,linux/amd64 \
--tag jpkitt/streamvalidator:latest .