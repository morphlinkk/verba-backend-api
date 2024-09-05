## Описание

Простое REST API на Go для управления задачами (ToDoList).

## Установка

### Требования

- Работоспособность проверялось на Go версии `1.23.0`, PGSQL `14.8.1`

### Порядок действий

1. Клонируйте репозиторий:

    ```bash
    git clone https://github.com/morphlinkk/verba-backend-api.git
    cd verba-backend-api
    ```

2. Установите зависимости:

    ```bash
    go mod tidy
    ```
3. Создайте файл .env в корне проекта
   
   К примеру:
    ```env
    PORT=8080
    DB_USER=POSTGRES
    DB_HOST=localhost
    DB_DATABASE=postgres
    DB_PASSWORD=password
    DB_PORT=5432
    ```
4. Создайте таблицу tasks в PGSQL

    ```sql
    CREATE TABLE tasks (
      id SERIAL PRIMARY KEY,
      title VARCHAR(255) NOT NULL,
      description TEXT,
      due_date TIMESTAMP WITH TIME ZONE,
      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP);
    ```
    
5. Запустите проект:

    ```bash
    go run main.go
    ```


## Использование

Для POST и PUT тело запроса необходимо отправлять в виде JSON, к примеру: 
```json
{
    "title": "Проверить проект",
    "description": "Проверить преокт на наличие ошибок",
    "due_date": "2024-09-05T17:30:00+03:00"
}
```
