FROM rabbitmq:3.9-management
ADD ./docker/rmq/conf/rabbitmq.conf /etc/rabbitmq/
ADD ./docker/rmq/conf/definitions.json /etc/rabbitmq/
ENV RABBITMQ_USER user
ENV RABBITMQ_PASSWORD password
EXPOSE 15672
RUN chown rabbitmq:rabbitmq /etc/rabbitmq/rabbitmq.conf /etc/rabbitmq/definitions.json
