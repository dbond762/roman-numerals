# roman-numerals
[![Build Status](https://travis-ci.org/dbond762/roman-numerals.png?branch=master)](https://travis-ci.org/dbond762/roman-numerals)

## Описание
Конвертер римских и арабских чисел

## Как запустить
Для запуска сервера загрузить и запустить соответствующий бинарник https://github.com/dbond762/roman-numerals/releases

Для запуска фронтенда
```bash
git clone https://github.com/dbond762/roman-numerals.git
cd roman-numerals/frontend
npm install
npm run dev
```

Для того, что бы собрать сервер из исходников: ([бинарники уже скомпилированны](https://github.com/dbond762/roman-numerals/releases))
  - установить go 1.10
  - склонировать проект в $GOPATH/src
```bash
cd roman-numerals/backend
go get github.com/go-chi/chi
go get github.com/go-chi/cors
go build
./roman
```
