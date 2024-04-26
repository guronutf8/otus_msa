# 08dz

## Сервис заказов на брокере с идемпотентностью

Основная информация описана в предыдущей работе, тут описано только введение идемпотентности.

Теперь что бы отправить заказ на оформление, его нужно отправлять с id корзины, который нужно получить заранее запрос `GetBasket`.
Сервис `billing` теперь пишет к пользователю: историю оплаченных корзин, так он понимает, обработали эту корзину или нет. А т.к. бд монга, операции с документов атомарные. Если ack упадет, то когда сообщение придет еще раз, сервис поймет, обрабатывал ли этот заказ, и если да, то просто ответит ack.         

### Сервис Order:

Сервис меняет статус корзины по паттерну compare-and-set

Получаем id корзины GET `/GetBasket` пока корзина не оплачена, всегда будет возвращаться один и тот же id
`{
"user":"Pavel"
}`

ответ
`{
"status": true,
"message": "662bd9a77c3d4e091e514e0f"
}
`


Создать заказ, теперь с basketId POST `/Order`
`{
"user":"Pavel",
"sum":10,
"basket":"662bd9a77c3d4e091e514e0f"
}`

### Сервис Billing:

Создать пользователя POST `/CreateUser`
`{
"user":"Pavel",
}`

Положить на баланс POST `/DepositCash` (тут изменений не делал, т.к. аналогично с getBasket)
`{
"user":"Pavel",
"sum":100
}`

### Сервис Notify:

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
- перейти в каталог 09_idempotent
- создать namespace `kubectl create ns 09dz`
- выполнить `make i` это запустит helm (после тестов, удалить инсталляцию `make ui`)
- Импортировать коллекцию запросов `dz09.postman_collection.json` в постман
- выполнить запросы, `CreateUser`, `DepositCash`, `CreateOrder` (этот несколько раз, или изменить сумму заказа, что бы
  получилось что денег на оплату не хватает), `List notify`

### Результат:

![](img/img.png)



