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

###todo
0. Сделать сервер с апи
1. Обдумать работу с chatGPT
2. Сделать на го взаимодействие с телегой
3. MVP взаимодействия с https://developers.google.com/custom-search/docs/tutorial/creatingcse
4. MVP геймификации выполнения задач

###hint
1. https://ru.wikipedia.org/wiki/%D0%9E%D0%B1%D0%BB%D0%B0%D0%BA%D0%BE_%D1%82%D0%B5%D0%B3%D0%BE%D0%B2 (облако тегов)
2. https://developers.google.com/search/docs/advanced/structured-data/search-gallery и https://cloud.google.com/natural-language и https://smacient.com/top-google-search-engine-apis/

Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik1UaEVOVUpHTkVNMVFURTRNMEZCTWpkQ05UZzVNRFUxUlRVd1FVSkRNRU13UmtGRVFrRXpSZyJ9.eyJodHRwczovL2FwaS5vcGVuYWkuY29tL3Byb2ZpbGUiOnsiZW1haWwiOiJyb2FydXNAeWFuZGV4LnJ1IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImdlb2lwX2NvdW50cnkiOiJBUiJ9LCJodHRwczovL2FwaS5vcGVuYWkuY29tL2F1dGgiOnsidXNlcl9pZCI6InVzZXItNm5SWWhRazdWWVZ4d1hRTGVjUDM5NzlrIn0sImlzcyI6Imh0dHBzOi8vYXV0aDAub3BlbmFpLmNvbS8iLCJzdWIiOiJhdXRoMHw2M2MzMzEwMTQyOWQxM2ZhOTg5MGRkNGIiLCJhdWQiOlsiaHR0cHM6Ly9hcGkub3BlbmFpLmNvbS92MSIsImh0dHBzOi8vb3BlbmFpLmF1dGgwLmNvbS91c2VyaW5mbyJdLCJpYXQiOjE2NzQ3NjgyNDksImV4cCI6MTY3NTM3MzA0OSwiYXpwIjoiVGRKSWNiZTE2V29USHROOTVueXl3aDVFNHlPbzZJdEciLCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIG1vZGVsLnJlYWQgbW9kZWwucmVxdWVzdCBvcmdhbml6YXRpb24ucmVhZCBvZmZsaW5lX2FjY2VzcyJ9.LUcn9JjP7DWa9oPKHcKb0jXKWcQrcm3V5kMGEch4na8Y8GiScri3uJZuVGPOf0APHqPGXMt3-dKVWylNj8C7TcJjyjPkACp-9nv1UACbQ2j0ORN2cCXhfNmzmCOCWxxjZ2ACPagtblMRZrybxv8k3X7BU9eckGVVeWFpKhenihaNPrN4slusGMaqgX2b7z1NGUZC4MOHKTQqvsjAIXSERsDlvsJXO8BbS3G0PuDxqyookgd4ca30QaWf4xoEVIBoUpWyGEFfDtVwW18bByMICjPZLvHoxTIqCz92UeGnzsH2lZn7x86h7O06WHw85aRu9etqlAj8FNtRfbk5C5rj0w