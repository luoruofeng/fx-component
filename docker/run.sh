docker run -d \
  --name redis \
  -p 6379:6379 \
  -v "${pwd}"/data:/data \
  -v "${pwd}"/log:/logs \
  -e REDIS_PASSWORD=luoruofeng \
  -e REDIS_BIND=0.0.0.0 \
  redis:latest
