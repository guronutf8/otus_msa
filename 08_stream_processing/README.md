# 08dz

## Сервис заказов на брокере

В наружу смотрят все 3 сервиса, между собой общаются через брокер NATS

### Сервис Billing (stateful):

Тут создаются пользователи, пополняется баланс, и списывается с баланса. Сервис ожидает поступление заказов в брокер,
при получении проверяет, если такой пользователь, хватает ли средств. Если
хвататет скидывает положительный результат брокер, если не хватает, то скидывает отрицательный.

Создать пользователя POST `/CreateUser`
`{
"user":"Pavel",
}`

Положить на баланс POST `/DepositCash`
`{
"user":"Pavel",
"sum":100
}`

### Сервис Order (stateless):

Тут создаются заказы, ничего не проверяется, т.к. сервис не знает о пользователях и их балансе. Просто создает заказ, и
скидывает в
очередь

Создать заказ POST `/Order`
`{
"user":"Pavel",
"sum":10
}`

### Сервис Notify (stateful):

Читает брокер, что бы условно отправить уведомление пользователю, и бонусом показывает что в БД(т.е. результат
положительный\отрицательный)

Посмотреть лог GET `/List`
`{
"user":"Pavel",
"sum":10
}`

## IDL

Запросы, сервисы, сообщения брокера описаны в `schema/*.proto`

## Sequence diagram

![Sequence diagram](https://www.plantuml.com/plantuml/svg/XPJFRjf04CRl-nIZd5e9r3kgK4fIzLALI6qVm0rkQuK0jSqXjn0KvH0aL2ddLdq3OfNW1Z2lCFj6_NPsOUCKcaY9x6_dps-_cNW_IOCuxP-sLN2STcDwrJgDEkjdyrUf7WAYuwWNT59gHErsrFb29R4PjyAufBphfj7dCwLQGHfy2fAGFYPXd1msl60ZeAXl_R0Vqb5mbZJYxWaAYU02dAG4XQhr3K2DZfhlFz2BWpWCqX3b4UJEgdNNhKPymmNVyeGpdkinJtb1l88qroEUSG5veYyWuv0u8_xEDnNI0-Ab7l8S1vc-HewKcFC4yW1P2uVpibfzmDCy-G_cOTZh9b5w42HVjs7-vC9oyL8moortCZ_D8QDPGyWBoHLc_9idjj4C3_VudIbb0v3cfm6u0krM-lonYtF7SITu3-heeNXvYvG2obIVUquBv2WTC_Homnzynxa7p7bc8AvC9Pv2pqljDsN0xHrTyDHxMUW1OX5lW620hICYG9SyAKzukjjmy-FxkINru5kSabxmn9Ip42LwqpihNK-FG2Wl6UTW71AICtDYLa-sSWbF_YZR7EuznlZbitVaoWKXhgJ3jRaYSYxmF1O_jgRfSONs_jq6N0XdbtlA0NbVC0ZPDdqfyKiNZG8vhtQqUEQi7yAG7B8PtGZ-R3pBWJcga1pcW3Kbdw__UPMIxrUrRW5gzEOrzlvjbAGaq1YcyGOW6CoLsGVYUmEWX-Kzb-fRBcJhOv6j8FSJ4sMpTWngF-oqy8Nz2m00)

## Тестирование

- клонировать репозиторий
- перейти в каталог 08_stream_processing
- создать namespace `kubectl create ns 08dz`
- выполнить `make i` это запустит helm (после тестов, удалить инсталляцию `make ui`)
- Импортировать коллекцию запросов `dz08.postman_collection.json` в постман
- выполнить запросы, `CreateUser`, `DepositCash`, `CreateOrder` (этот несколько раз, или изменить сумму заказа, что бы
  получилось что денег на оплату не хватает), `List notify`

### Результат:

![](img/img.png)



