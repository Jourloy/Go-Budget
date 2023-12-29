# Budget

## Описание

Сервис для ведения личных расходов и доходов

- [x] Авторизация через Логин/Пароль
- [ ] Авторизация через Google OAuth (отказался, так как закон запрещает)
- [x] Создание, настройка и удаление трекера
- [x] Создание, настройка и удаление дохода или траты
- [x] Создание, настройка и удаление кредита
- [x] Подсчет лимита на день
- [x] Подсчет оставшихся дней

## Начало

### Запуск

Запусти последнюю версию **бинарника** из релиза.

## Тесты

Практически все функции были проверены автотестами.

### Операции (spends)

Проверены вручную через Postman и frontend.


#### Cоздание

**url**
```
POST `/{bid}/`
```

**тело**
```
{
    cost: 100
    category: "food"
    isCredit: false
}
```

##### Коды
- Если пользователь не авторизован - 403
- Если передано некорректное тело - 400
- Если у пользователя нет бюджета - 404
- Если произошла ошибка во время добавление в БД - 500
- Если все хорошо - 200

#### Обновление

**url**
```
PATCH `/bid/sid/`
```

**тело**
```
{
    cost: 200
}
```

##### Коды
- Если пользователь не авторизован - 403
- Если передано некорректное тело - 400
- Если у пользователя нет бюджета - 404
- Если произошла ошибка во время добавление в БД - 500
- Если все хорошо - 200

#### Удаление

**url**
```
DELETE `/bid/sid/`
```

##### Коды
- Если пользователь не авторизован - 403
- Если у пользователя нет бюджета - 404
- Если произошла ошибка во время добавление в БД - 500
- Если все хорошо - 200