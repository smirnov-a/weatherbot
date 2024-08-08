[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

# Планировщик задач по расписания из локального файла crontab (Russian)

## Содержание
 * [Общая информация](#general-info)
 * [Установка и настройка](#installation)
 * [TODO](#todo)

## <a name="general-info"></a>Общая информация
Данная программа предназначена для запуска задач по расписанию из файла crontab
Программа при запуске считывает файл crontab и дальше продолжает работать в фоне и выполнять задачи (т.е. явлеяется daemon'ом)
Также, программ отслеживает изменения в crontab, перечитывает новые задания и перезапускает внутренний планировщик
Формат crontab аналогичен стандартному файлу планировщика Cron:

| Параметр      | Допустимый интервал                           |
|---------------|-----------------------------------------------|
| минуты        | 0-59                                          |
| часы          | 0-23                                          |
| день месяца   | 1-31                                          |
| месяц         | 1-12                                          |
| день недели   | 0-7 (0-Вс,1-Пн,2-Вт,3-Ср,4-Чт,5-Пт,6-Сб,7-Вс) |

Пример находится в файле crontab.sample
Данную программу я писал для ежедневной отправки прогноза погоды в канал в Telegram
```cronexp
# min hour day month weekday command
30 8 * * * weather Moscow
```
Это означает ежедневно запускать локальную команду "weather" с параметром Moscow.
Город можно указать в виде Moscow[55.7558 37.6176], т.е. указать в скобка широту и долготу нужной точки.
В этом случае программа возьмет координаты прямо отсюда (из названия города, из скобок)
Это несколько ускорит работу, потому что не придется идти за координатами в geolocation
Также можно указывать несколько городов через пробел. Например:
```cronexp
# min hour day month weekday command
30 8 * * * weather Moscow Yekaterinburg
```
Прогноз для каждого города придет в телеграм отдельным сообщением

Список команд зашит в map в scheduler.go
```go
type cmdMapping map[string]interface{}

// cmdStorage contains commands for scheduled tasks
var cmdStorage = cmdMapping{
    "weather": providers.GetWeather,
    "test":    test,
}
```
Там же в scheduler.go происходит запуск уазанной функции из cmdStorage (с помощью reflect):
```go
// callFunc - call function with params via reflection
func callFunc(ctx *app.AppContext, funcName string, params ...interface{}) (result interface{}, err error) {
    f := reflect.ValueOf(cmdStorage[funcName])
    if (len(params) + 1) != f.Type().NumIn() {
        err = fmt.Errorf("the number of params is out of index. len:%d; num:%d", len(params), f.Type().NumIn())
        return
    }
    in := make([]reflect.Value, len(params)+1)
    // first param always app context
    in[0] = reflect.ValueOf(ctx)
    for k, param := range params {
        in[k+1] = reflect.ValueOf(param)
    }

    var res []reflect.Value
    res = f.Call(in)
    result = res[0].Interface()

    return
}
```
Внутри работа основана на гоуртинах: каждая задача в своей горутине.
```go
func RunTasks(app *app.AppContext, tasks []Task, cr *cron.Cron) {
	for _, task := range tasks {
		task := task // closure
		_, err := cr.AddFunc(task.Schedule, func() {
			go executeTask(app, task.Command)
		})
		if err != nil {
			app.Logger.Printf("Error adding cron task %s: %v", task.Schedule, err)
		}
	}
}
```
А обмен между горутинами осуществляется с помощью каналов:
```go
wg := &sync.WaitGroup{}
// chanData weather data
chanData := make(chan *weather.WeatherData)
// chanMessage channel for sending message to telegram
chanMessage := make(chan *weather.WeatherData)
once := &sync.Once{}
closeDataChan := func(ch chan *weather.WeatherData) {
	once.Do(func() {
		close(ch)
	})
}
defer func() {
	closeDataChan(chanData)
	close(chanMessage)
}()

sendMessageFunc := func(data *weather.WeatherData) {
	message.SendMessageToTelegram(app, data)
}
go worker(sendMessageFunc, chanMessage)

provider := getProvider(app.Cache)
for _, city := range cities {
	wg.Add(1)
	go provider.GetWeatherData(ctx, city, chanData, wg)
}

go func() {
	wg.Wait()
	// close data channel. when closed it will stop cycle below
	closeDataChan(chanData)
}()
```

## <a name="installation"></a>Установка и настройка
Скачайте репозиторий:
```shell
git clone https://github.com/smirnov-a/weatherbot.git
cd weatherbot
````
Установить пакеты и зависимости:
```shell
go mod tidy
```
Эта команда обновит файл go.mod и установит все зависимости из go.sum
Если вы хотите хранить пакеты в папке vendor, то воспользуйтесь командой:
```shell
go mod vendor 
```
После этого нужно собрать исполняемый файл:
```shell
make
```
или
```shell
go build -o weatherbot main.go
```
Настройка заключается в создании файла crontab (имя файла может быть любым, его можно указать параметром --crontab <filename>, по умолчанию "crontab")
Дальше нужно создать файл config/.env (.env.sample лежит в каталоге)
```text
WEATHER_PROVIDER="openweathermap"
#WEATHER_PROVIDER="weatherapi"
OPENWEATHERMAP_API_KEY="your-api-key"
#WEATHERAPI_API_KEY="your-api-key"

PROXY_URL="socks5://<username>:<password>@sock-server-address:port"

TELEGRAM_TOKEN="your-telegram-token"
# update/webhook
TELEGRAM_MODE="update"
#TELEGRAM_WEBHOOK="https://your-http-server-address/"
#WEBHOOK_CERT="config/cert/fullchain.pem"
#WEBHOOK_KEY="config/cert/privkey.pem"
#WEBHOOK_PORT=8443

TELEGRAM_CHAT_ID=-100<your-chat-id>

LANGUAGE="ru"
```
WEATHER_PROVIDER: задает через какого провайдера погоды работать (сейчас два значения "openweathermap" или "weatherapi")

OPENWEATHERMAP_API_KEY и WEATHERAPI_API_KEY: api-ключи для соответствующих сервисов

PROXY_URL: адрес прокси-сервера. Поддерживаются http и socks5 прокси. На момент создания программы сервис "openweathermap" из России недоступен напрямую, а только через прокси

TELEGRAM_TOKEN: токен вашего телеграм-бота

TELEGRAM_MODE: режим работы телеграм бота (update или webhook). Webhook создает меньше нагрузку, но требует поднятия https-сервера, доступного снаружи

TELEGRAM_CHAT_ID: ИД чат-группы в телеграм. Нужно в телеграм скопировать id и добавить префикс "-100" (это префикс у групп в телеграм)

LANGUAGE: язык локализации (используется в шаблоне, с помощью которого генерится выходная картинка с прогнозом погоды)

После успешной сборки и настройки, программа должна запуститься и начать выполнять задания. Вот пример работы для задачи:
"weather Moscow Yekaterinburg" (получение погоды для двух городов)

![](https://github.com/user-attachments/assets/f4b9081a-c17a-49f6-a595-e91fe50adffa "Погода для Екатеринбурга")
![](https://github.com/user-attachments/assets/b54aaec5-2a23-4212-ad99-01f2e9b93407 "Погода для Москвы")

## <a name="todo"></a>TODO
Сейчас провайдер, от которого забирать погоду, задается в конфиге (либо openweathermap, либо weatherapi). Хочу сделать перебор провайдеров, если первый не отвечает или отвечает с ошибкой.

Добавить возможность работы через очередь (внутреннюю или внешнюю, чтобы не терять задачи). Сейчас получение погоды происходит так:
```go
func DoRequestWithRetry(req *http.Request, maxRetires int, initialWait time.Duration) (*http.Response, error) {
}
```
т.е. пробует обратиться к сервису N-раз и в случае ошибки со стороны сервиса задача просто потеряется (до следующего дня)

Также планирую добавить рассылку event'ов для определенного города (предстоящие интересные события, которые предлстоят в городе). Нашел провайдера, который отдает по api для Екатеринбурга

## License

GPL-3.0-or-later
