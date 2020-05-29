FROM alpine:3.11.6

WORKDIR /app

ENV pfx_path="" pfx_password="" server_name="" edge="" catalog="" user_dn="" json_config_path="" catalog_file_name="06.catalog.json"

# Meant to be built from top directory
COPY dist/gm-cat_linux_amd64/gmcat /app/gmcat
COPY docker/docker-config.yaml /app/config.yaml
COPY certs /app/certs

ENTRYPOINT [ "./gmcat" ]
CMD [ "--help" ]