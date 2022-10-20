# Consumer rabitqueues and publish to rocket chat


for development use docker-compose.debug.yaml
```bash
docker-compose up -d -f docker-compose.debug.yaml
```

for normal use use docker-compose.yaml
```bash
docker-compose up -d
```

How to use this service in a bash script and jobs
```bash
#!/bin/bash
podman run --rm -it --env-file alertas_servicios.env -v /etc/localtime:/etc/localtime:ro domain:version >> alertas_servicios.log && \
```
