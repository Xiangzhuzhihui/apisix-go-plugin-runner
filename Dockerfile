FROM apache/apisix:2.15.0-centos

# 添加软件包
ADD ./go-runner /plugins/xzzhAuth/go-runner

# 配置启动命令
CMD ["sh", "-c", "/usr/bin/apisix init && /usr/bin/apisix init_etcd && /usr/local/openresty/bin/openresty -p /usr/local/apisix -g 'daemon off;'"]

STOPSIGNAL SIGQUIT