
---
Made with ❤ by [AlRichy](https://github.com/AlRichy)
---

## Итоговое задание
### Go веб-сервис. Планировщик задач
1. Реализованы все задания (включая задачи со звездочкой)
2. Указывается один из адресов:
	2.1. http://127.0.0.1:7540/
	2.2. http://localhost:7540
3. Применен модульный подход при разработке. Настройки по-умолчанию:
```dockerfile
ENV  TODO_PORT=7540
ENV  TODO_PASSWORD=1234
ENV  TODO_DBFILE=data/scheduler.db
```
4. Запуск с помощью Docker:
	4.1. Docker-образ создается командой
	```bash
	docker build --progress plain --tag todofinal:v1 .
	```
	4.2. Запуск Docker-контейнера выполняется командой
	```bash
	docker run -d --name todogo -p 7540:7540 -e TODO_DBFILE=/data/scheduler.db -e TODO_PASSWORD=12345 -v data/scheduler.db todofinal:v1
	```
	4.3. Остановка контейнера сочетанием клавиш **Ctrl+C** в терминале.
5. Запуск через Go:
	5.1. Установка зависимостей
	```bash
	go mod tidy
	```
	5.2. Сборка
	```bash
	go build .
	```
	5.3. Запуск
		```bash
		./final-project-todo
		```
	5.4. Остановка контейнера сочетанием клавиш **Ctrl+C** в терминале.