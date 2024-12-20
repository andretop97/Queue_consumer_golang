# Consumidor de fila com RabbitMQ

## Entendo RabbitMQ e sua importancia

Em sistemas distribuídos, o processamento assíncrono é frequentemente necessário para desacoplar tarefas e garantir que o fluxo principal não seja interrompido. Uma das abordagens mais eficazes para isso é o uso de sistemas de mensageria. Nesse contexto, é possível enviar mensagens para filas, onde elas ficam armazenadas até serem processadas posteriormente.

Este projeto utiliza o RabbitMQ, um dos serviços de mensageria mais populares. Além de implementar o conceito de filas, o RabbitMQ introduz o conceito de exchanges. Uma exchange atua como intermediária, recebendo mensagens enviadas pelos produtores e encaminhando-as para as filas apropriadas. Essa abordagem permite uma distribuição de mensagens mais dinâmica e direcionada.

As filas armazenam as mensagens até que os consumidores as processem. Os consumidores, por sua vez, recuperam essas mensagens e realizam as operações necessárias. Essa arquitetura garante maior flexibilidade, escalabilidade e eficiência em sistemas distribuídos.


## Arquitetura de dead letter queue

![Arquitetura da nossa dlq]("dlq.png")

