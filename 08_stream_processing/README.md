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

![Sequence diagram](https://www.plantuml.com/plantuml/svg/XPJFRjf04CRl-nIZd5e9r3kgK4fIzLALI6qVm0rkQuK0jSqXjn0KvH0aL2ddLdq3OfNW1WolCFj6_NQyZjaKDv68pBVpvxVV3DwFOY-CXwTTvHt7_P1UTKw3Nl5i_YirG41m70oyNaojeDiDUb_84TjXnmyTvLq_3ZwVATNnOk-J5_dqCVB3wB1L2da45FLjw0zzElBEebROFq4X4Mo0bPIWKBtv0D2We_Rx3_JY8FIzs4Tv1DcpUjDsQn3Vy9flUSO9B_MKvxmYta1QwWalU0rvfYyWuv0u8VxEDpNI8-9CZtc9WqHV8qTABF79z0XPAulpibxVuccV_0Vpq1_qeo0-yAAlkp1_yZhduammtPQLp4zpo6WA2Fb2SeK5_qPjQRJ0mpq-2wNo0AILqm3S6kqs-lonYqlBSITuMtNqMBoyHSeQoboVbvmvaAMqpD3TXZ_uZjCICEM5WRWobNWEFNNQFvG1T-de2aU5bqWU8HRn1WW6v7L416Z6CtV2ysh3yRVLprbI0xx3AUa5Jqov7CAXDxqJvtfw0aBvad0AnZ4XEJ4dPlLaApDuyaVbpU7UOunVPNP5hbmGuaOwN9ih8kS2prFnOsUQdjPexzzjm8LmThcdB53f2mOXswOl9JwpqIYGygQYnJDh_HY6f93Dw2O4gV0i1UQaHdAO6bYKV5xyU9MIzrUzcGCqwCrhxFtR88h9GCimuGr0C9Whim_4zmP0Zzqzb-fvCMnjKKsFSiTaZIgcik_tEdYT_WK0)

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



