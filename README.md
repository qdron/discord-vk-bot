# discord-vk-bot

[![Go Report Card](https://goreportcard.com/badge/github.com/qdron/discord-vk-bot)](https://goreportcard.com/report/github.com/qdron/discord-vk-bot)

Бот для отображения личных сообщений сообощества ВК на канале в Discord.

Для запуска используются следущие параметры:

- **vk_token** API ключ для управления сообществом. Можно получить в разделе "Управление сообществом"
- **vk_groupid** Идентификатор группы. Для правильного отображения ответов от группы ВК в чате Discord.
- **discord_token** Токен приложения discord. Откройте или создайте новое приложение на [этой странице](https://discordapp.com/developers/applications/me) после чего в этом приложении в разделе Bot увидите его token.
- **discord_channelid** Идентификатор канала в который бот будет передавать сообщения. Вы можете увидеть его в режиме разработчика.
- **create** Если указан то сгенеррутеся файл конфигурации с указанными в аргументах параметрами. Если указан также параметр config то файл будет сохранен по указанному пути.
- **config** Путь к файлу с параметрами запуска (По умолчанию "./conifg.json")
- **log** Путь к файлу лога. На базе этого файла создается лог с ротацией внутри заданой папки (По умолчанию "./logs/bot.log")

Если в текущей папке уже был создан файл конфигурации то бот будет работать с параметрами из него, если параметром config не указан другой путь к файлу.

In english.

Bot for displaying personal messages of the VK community/group on channel in Discord.

The following parameters are used to start:

- **vk_token** API key for community management. You can get in the section "Managing the community"
- **vk_groupid** The identifier of the group. To correctly display the answers from the VK group in the Discord chat.
- **discord_token** The application discord token. Open or create a new application on [this page] (https://discordapp.com/developers/applications/me) and then in this application in the Bot section you will see its token.
- **discord_channelid** The identifier of the channel to which the bot will send messages. You can see them in developer mode.
- **create** If you specify the configuration file with the current parameters. If the config parameter is also specified, the file will be saved according to the specified path.
- **config** Path to the file with startup parameters (Default is "./config.json")
- **log** Путь к файлу лога. На базе этого файла создается лог с ротацией внутри заданой папки (По умолчанию "./logs/bot.log")

If a configuration file has already been created in the current folder, then the bot will work with the parameters from it if the **config** parameter does not specify a different path to the file.
