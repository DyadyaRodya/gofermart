# gofermart

Индивидуальный дипломный проект курса «Go-разработчик»

## Локальный запуск автотестов

Установить `gophermarttest` [отсюда](https://github.com/Yandex-Practicum/go-autotests?tab=readme-ov-file#%D1%82%D1%80%D0%B5%D0%BA-%D1%81%D0%B5%D1%80%D0%B2%D0%B8%D1%81-%D1%81%D0%BE%D0%BA%D1%80%D0%B0%D1%89%D0%B5%D0%BD%D0%B8%D1%8F-url)
Установить `statictest` [отсюда](https://github.com/Yandex-Practicum/go-autotests?tab=readme-ov-file#%D1%82%D1%80%D0%B5%D0%BA-%D1%81%D0%B5%D1%80%D0%B2%D0%B8%D1%81-%D1%81%D0%BE%D0%BA%D1%80%D0%B0%D1%89%D0%B5%D0%BD%D0%B8%D1%8F-url)

Для статического теста
```shell
make lint
```

Для запуска gophermarttest тестов
```shell
make test-proj
```

### Дополнительные команды

Для генерации моков
```shell
make mock
```

Для запуска unit тестов выполняем
```shell
make tests
```

Для запуска accrual-системы выполняем
```shell
make accrual-start
```

Для остановки accrual-системы выполняем
```shell
make accrual-stop
```

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.
