FROM rainbond/openresty:1.13.6.2

RUN apk add --no-cache bash net-tools
ADD . /run
ENV NGINX_CONFIG_TMPL=/run/nginxtmp
ENV NGINX_CUSTOM_CONFIG=/run/nginx/conf
ENV RELEASE_DESC=__RELEASE_DESC__
ENV OPENRESTY_HOME=/usr/local/openresty
ENV PATH="${PATH}:${OPENRESTY_HOME}/nginx/sbin"
EXPOSE 8080

ENTRYPOINT ["/run/entrypoint.sh"]