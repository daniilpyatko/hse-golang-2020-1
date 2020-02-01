Это репозиторий курса "Разработка веб-сервисов на Golang" в ВШЭ.

Это факультатив. Запись на него - ваше сознательное решение, никто не заставлял.

Вы пришли сюда знаниями, а не за оценками. Поэтому все задания надо делать самостоятельно, без консультаций с другими студентами, без просмотра чужих решений, особенно если вы сами ничего не сделали епще сами, в крайнем случае - спрашивать препода.

При обнаружении списывайний вы отчисляетесь с курса.

Можно ходить на лекции и не делать домашки. Можно не делать домашку если она вам не нравится и перейти к другой.

Преподы на вашей стороне и всегда готовы помочь. Можно задавать вопросы в телеграме в личку и в общий чат.

Домашку необходимо выполнить в течении 3-х недель с момента лекции, на которой она была выдана.

-----

Прочее:

0. преподавательский состав оставляет за собой право дополнять правила
1. домашки сдаём тому преподавателю, который вёл занятие
2. хардкод ( код работающий под частное условие ) запрещён. за первый раз - предупреждение, за последующие -1 балл. можно спрашивать будет-ли что-то хардкодом до сдачи задания. весь код должен работать максимально универстально.
3. за ДЗ можно получить 10 баллов, после дедлайна - 5 баллов
4. домашки пишем там же гле лежит вводная (например, 1/99_hw/XXX), другие папки не создаём
5. тесты домашек править нельзя
6. вопросы задавать четко, конкретно: "я делаю Х, получаю Y, а хочу получить Z"
7. код и тем более решения домашек в паблик открывать нельзя, репозиторий должен быть приватным
8. студент должен иметь реальное имя-фамилию-фото в гитлабе и на портале. реальные фио и фото в телеграме так же желательны
9. домашку надо коммитить в свою репу, создавать merge request в основную репу не надо
10. домашки предназначены для выполнения индивидуально и самостоятельно. Это значить что нельзя делать их группой, нельзя обсуждать как делать, нельзя показывать свои решгения.
11. преподавательский состав оставляет за собой право не принимать мутые, некрасивые решения домашек. в этом случае необходимо поправить замечания без препирательств
12. консультации и проверки заданий даются в основном вечером
13. если вы пишите в 2 ночи - не надо писать в 9 утра вопросы "а вы посмотрели?" - в подобных случаях мы скорее всего посмотрим только вечером и надо напомнить про себя после 19 часов
14. халявы не будет, домашки сложные, придётся работать

* при добавлении вас в репозиторий необхдимо его форкнуть, сделать приватным ( должен быть таким по-умолчанию ), добавить туда преподавателей rvasily и DmitryDorofeev с уровне доступа maintainer. https://s.mail.ru/7XKz/fmJyoaZMA

-----

Последовательность выполнения ДЗ:
1. Забираем обнволения кода (см. ниже)
2. Читаем задание в `X`/99_hw/`X`.md
3. Доводим код до состояния прохождения тестов
4. Не забываем форматировать код (gofmt)
5. Стучимся в личку с просьбой поревьювить

Как забрать в свой репозиторий обновления из основного:
```
# будучи в своём репозитории
git pull https://gitlab.com/rvasily/hse-golang-2020-2.git master
git push origin master
```
