##Task monster rpg

Концепция:

1. Юзер задает цель - изучить новый язык\доделать проект и т.д.

2. В цели прописываются ключевые особенности - сроки, источники информации, полезные ссылки, и т.д.

3. По полезным ссылкам парсятся какие-то "концепции" (пока просто рандомный небольшой текст, потом будет анализатор)

4. "Концепции" выдаются пользователю, он должен по ним либо что-то отписать, либо заучить, либо выполнить действие

5. Все это геймифицировано, таска - это рейд, действия несут экспу, кто будет монстром пока вопрос)

Отличия от habitRPG - наличие парсера для сбора информации/оформления ее в задания и мотивация пользователя на ее выполнение.
Смысл - всегда есть информация по проекту, которую внимание пользователя упустило. В некоторых цели порой сложно придумать, куда же двигаться дальше.

Контакт с пользователем через бота, пока телеграм.

Схема работы апи: - сет топика по токену (юзерИД или иное) и получение списка связанных тем

####todo
0. Сделать сервер с апи
1. Обдумать работу с chatGPT
2. Сделать на го взаимодействие с телегой
3. MVP взаимодействия с https://developers.google.com/custom-search/docs/tutorial/creatingcse
4. MVP геймификации выполнения задач

####hint
1. https://ru.wikipedia.org/wiki/%D0%9E%D0%B1%D0%BB%D0%B0%D0%BA%D0%BE_%D1%82%D0%B5%D0%B3%D0%BE%D0%B2 (облако тегов)
2. https://developers.google.com/search/docs/advanced/structured-data/search-gallery и https://cloud.google.com/natural-language и https://smacient.com/top-google-search-engine-apis/