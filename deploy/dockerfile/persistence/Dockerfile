FROM alpine:3.13

# 设置固定的项目路径
ENV WORKDIR /quickim

# 将可执行文件复制到目标目录
ADD ./persistence $WORKDIR/persistence

# 创建用于挂载的几个目录，添加可执行权限
RUN mkdir $WORKDIR/logs && \
  chmod +x $WORKDIR/persistence && \
  mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2


WORKDIR $WORKDIR
ENTRYPOINT [ "./persistence" ]