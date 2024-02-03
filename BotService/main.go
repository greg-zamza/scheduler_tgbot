package main

import (
    "os"
    "time"
    "strings"
    "fmt"
    "log"
    "database/sql"

    _ "github.com/lib/pq"
    "github.com/NicoNex/echotron/v3"
)

/* Объявление необходимых переменных */
var botToken string
var adminPassword string
var adminText string
var userText string
var scheduleText string

/* Функция init() выполняется перед main() и
 * служит для определения переменных */
func init() {
    botToken = os.Getenv("BOT_TOKEN")
    adminPassword = os.Getenv("ADMIN_PASSWORD")
    adminText = `
Добро пожаловать в панель управления.

Доступные команды:
/s <дата> <группа> - посмотреть расписание
/add <группа> <дата> <номер пары> <название предмета> <аудитория> - добавить занятие
/del <группа> <дата> <номер пары> - удалить занятие

/units - список групп
/addunit <группа> - добавить группу
/delunit <группа> - удалить группу и все занятия, связанные с ней

/exit - выйти из аккаунта администратора
/start - войти как студент или как администратор

/help - справка`

    userText = `
Добро пожаловать в расписание.

Доступные команды:
/units - посмотреть список всех групп
/unit <группа> - изменить свою группу

/next - посмотреть ближайшие дни, когда есть занятия
/s <ДД.ММ.ГГГГ> - посмотреть расписание своей группы в определённый день

/exit - выйти
/start - зайти как студент или как администратор

/help - справка
    `

    scheduleText = `
[%s], группа %s.

1: %s, ауд. %s
2: %s, ауд. %s
3: %s, ауд. %s
4: %s, ауд. %s
5: %s, ауд. %s
6: %s, ауд. %s
7: %s, ауд. %s
8: %s, ауд. %s`
}

func main() {
    // подключение к БД
    var conn_params string = fmt.Sprintf(
        "user=%s dbname=%s sslmode=disable host=DatabaseService password=%s",
        os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"))

    db, err := sql.Open("postgres", conn_params)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // подключение к telegram API
    api := echotron.NewAPI(botToken)
    log.Println("started")
    for update := range echotron.PollingUpdates(botToken) {
        // помогает боту не сломаться от невалидного апдейта (добавление в чат/реакция)
        if update.Message == nil {
            log.Println("Unhandled update")
        } else {
            var id int64 = update.ChatID()
            var words []string = strings.Fields(update.Message.Text)
            var unitName string
            var isAdmin bool

            rows, err := db.Query("SELECT 1 unit, is_admin FROM users WHERE id = $1", id)
            if err != nil {
                log.Println(err)
            }

            for rows.Next() {
                if err := rows.Scan(&unitName, &isAdmin); err != nil {
                    log.Println(err)
                }
            }

            switch {
            case words[0] == "/help":
                if isAdmin { api.SendMessage(adminText, id, nil)
                } else { api.SendMessage(userText, id, nil) }

                case words[0] == "/start":
                        if isAdmin {
                            api.SendMessage("Вы уже вошли как администратор", id, nil)
                        } else if unitName != "" {
                            api.SendMessage(fmt.Sprintf("Вы уже вошли как студент группы %s", unitName), id, nil)
                        } else {
                            var markup1 echotron.ReplyMarkup = echotron.ReplyKeyboardMarkup{
                                Keyboard: [][]echotron.KeyboardButton{
                                    {{Text:"Войти как студент"}},
                                    {{Text:"Войти как администратор"}},
                                },
                            }
                            opts1 := echotron.MessageOptions{ReplyMarkup: markup1}
                            api.SendMessage("Добро пожаловать в расписание", id, &opts1)
                        }

                case words[0] == "/exit":
                    if (!isAdmin) && unitName == "" {
                        api.SendMessage("Вы не авторизованы", id, nil)
                    } else {
                        if _, err := db.Exec("DELETE FROM users WHERE id = $1", id); err != nil {
                            log.Println(err)
                        }
                        api.SendMessage("Вы больше не авторизованы", id, nil)
                    }

                case words[0] == "Войти" && len(words) == 3:
                    if words[2] == "администратор" {
                        api.SendMessage("Отправьте пароль", id, nil)
                    } else if words[2] == "студент" {
                        msg := "Посмотри список доступных групп с помощью /units и укажи свою группу с помощью /unit <группа>"
                        api.SendMessage(msg, id, nil)
                    }

                case words[0] == "/admin":
                    api.SendMessage("Отправьте пароль", id, nil)

                case words[0] == adminPassword && len(words) == 1:
                    if isAdmin {
                        api.SendMessage("Вы уже вошли как администратор", id, nil)
                    } else if unitName != "" {
                         if _, err := db.Exec("UPDATE users SET is_admin = 1 WHERE id = $1", id); err != nil {
                            log.Println(err)
                         }
                    } else {
                        if _, err := db.Exec("INSERT INTO users(id, is_admin) VALUES($1, $2)", id, 1); err != nil {
                            log.Println(err)
                        }
                    }
                    api.SendMessage(adminText, id, nil)

                case words[0] == "/units":
                    rows, err := db.Query("SELECT name FROM units")
                    if err != nil { log.Println(err) }
                    defer rows.Close()

                    var groupNames []string
                    for rows.Next() {
                        var groupName string
                        if err := rows.Scan(&groupName); err != nil {
                            log.Println(err)
                        }
                        groupNames = append(groupNames, groupName)
                    }

                    var sb strings.Builder
                    sb.WriteString("Список групп:\n")

                    for _, group := range groupNames {
                        sb.WriteString(group)
                        sb.WriteString("\n")
                    }
                    finalString := sb.String()
                    api.SendMessage(finalString, id, nil)

                /* команды, доступные только администратору */
                case words[0] == "/delunit":
                    if (!isAdmin) {
                        api.SendMessage("Я не знаю такой команды", id, nil)
                    } else if len(words) != 2 {
                        api.SendMessage("Использование: /delunit <группа>", id, nil)
                    } else {
                        // проверка наличия группы в бд
                        var exists bool
                        err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM units WHERE name = $1)", words[1]).Scan(&exists)
                        if err != nil { log.Println("not selected from db") }

                        if (!exists) {
                            msg := fmt.Sprintf("Группы %s не существует.\nВы можете посмотреть список групп с помощью /units",
                                            words[1])
                            api.SendMessage(msg, id, nil)
                        } else {
                            _, err := db.Exec("DELETE FROM TABLE units WHERE name = $1", words[1])
                            if err != nil { log.Println("not deleted") }
                            api.SendMessage(fmt.Sprintf("Группа %s и все занятия, связанные с ней, удалены.", words[1]), id, nil)
                        }
                    }

                case words[0] == "/addunit":
                    if (!isAdmin) {
                        api.SendMessage("Я не знаю такой команды", id, nil)
                    } else if len(words) != 2 {
                        api.SendMessage("Использование: /addunit <группа>", id, nil)
                    } else {
                        // проверка наличия группы в бд
                        var exists bool
                        err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM units WHERE name = $1)", words[1]).Scan(&exists)
                        if err != nil { log.Println("not selected from db") }

                        if (exists) {
                            msg := fmt.Sprintf("Группа %s уже существует.\nВы можете посмотреть список групп с помощью /units",
                                                words[1])
                            api.SendMessage(msg, id, nil)
                        } else {
                            _, err := db.Exec("INSERT INTO units(name) VALUES($1)", words[1])
                            if err != nil { log.Println("not inserted") }
                            api.SendMessage(fmt.Sprintf("Группа %s добавлена.", words[1]), id, nil)
                        }
                    }

                case words[0] == "/s":
                    resultMap := map[int][]string{
                        1: {"-", "-"},
                        2: {"-", "-"},
                        3: {"-", "-"},
                        4: {"-", "-"},
                        5: {"-", "-"},
                        6: {"-", "-"},
                        7: {"-", "-"},
                        8: {"-", "-"},
                        }
                    var groupName string

                    if isAdmin {
                        if len(words) != 3 {
                            api.SendMessage("Использование: /s <номер группы> <ДД.ММ.ГГГГ>", id, nil)
                            break
                        } else {
                            groupName = words[1]
                        }
                    } else if unitName != "" {
                        if len(words) !=2 {
                            api.SendMessage("Использование: /s <ДД.ММ.ГГГГ>", id, nil)
                            break
                        } else {
                            groupName = unitName
                        }
                    } else {
                        msg := "Укажи свою группу с помощью /unit <группа> или войди в аккаунт администратора c помощью /admin.\nСписок групп можно посмотреть с помощью /units"
                        api.SendMessage(msg, id, nil)
                        break
                    }
                    timeDay, erro := time.Parse("02.01.2006", words[2])
                    if erro != nil {
                        api.SendMessage("Некорректный формат даты", id, nil)
                        break
                    }
                    rows, err := db.Query("SELECT num, name, room FROM classes WHERE day = $1 AND unit = $2",
                                            timeDay, groupName)
                    if err != nil { log.Println(err) }

                    for rows.Next() {
                        var key int
                        var value1, value2 string
                        if err := rows.Scan(&key, &value1, value2); err != nil {
                            log.Println(err)
                        }
                        resultMap[key] = []string{value1, value2}
                    }

                    msg := fmt.Sprintf(scheduleText, words[2], words[1],
                                        resultMap[1][0], resultMap[1][1],
                                        resultMap[2][0], resultMap[2][1],
                                        resultMap[3][0], resultMap[3][1],
                                        resultMap[4][0], resultMap[4][1],
                                        resultMap[5][0], resultMap[5][1],
                                        resultMap[6][0], resultMap[6][1],
                                        resultMap[7][0], resultMap[7][1],
                                        resultMap[8][0], resultMap[8][1])
                    api.SendMessage(msg, id, nil)

                case words[0] == "/add":
                    if (!isAdmin) {
                        api.SendMessage("Я не знаю такой команды", id, nil)
                    } else if len(words) != 6 {
                        api.SendMessage("Использование: /add <группа> <ДД.ММ.ГГГГ> <номер пары> <название> <аудитория>", id, nil)
                    } else {
                        date, err := time.Parse("02.01.2006", words[2])
                        if err != nil {
                            api.SendMessage("Некорректный формат даты.", id, nil)
                        } else {
                            var exists bool
                            err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM classes WHERE unit = $1 AND day = $2 AND num = $3)",
                                                words[1], date, words[3]).Scan(&exists)
                            if err != nil { log.Println("not selected from db") }

                            if (exists) {
                                msg := fmt.Sprintf("Занятие у группы %s %s на %s паре уже есть.", words[1], words[2], words[3])
                                api.SendMessage(msg, id, nil)
                            } else {
                                _, err := db.Exec(
                                    "INSERT INTO classes(unit, day, num, name, room) VALUES($1, $2, $3, $4, $5)",
                                    words[1], words[2], words[3], words[4], words[5])

                                if err != nil { log.Println("not deleted") }
                                api.SendMessage(fmt.Sprintf(
                                                "Занятие добавлено\nГруппа: %s\nДата: %s, %s пара\nНазвание: %s\nАудитория: %s",
                                                words[1], words[2], words[3], words[4], words[5]), id, nil)
                            }
                        }
                    }

                case words[0] == "/del":
                    if (!isAdmin) {
                        api.SendMessage("Я не знаю такой команды", id, nil)
                    } else if len(words) != 4 {
                        api.SendMessage("Использование: /del <группа> <ДД.ММ.ГГГГ> <номер пары>", id, nil)
                    } else {
                        date, err := time.Parse("02.01.2006", words[2])
                        if err != nil {
                            api.SendMessage("Некорректный формат даты.", id, nil)
                        } else {
                            var exists bool
                            err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM classes WHERE unit = $1 AND day = $2 AND num = $3)",
                                                words[1], date, words[3]).Scan(&exists)
                            if err != nil { log.Println("not selected from db") }

                            if (!exists) {
                                msg := fmt.Sprintf("У группы %s нет занятия %s на %s паре.", words[1], words[2], words[3])
                                api.SendMessage(msg, id, nil)
                            } else {
                                _, err := db.Exec("DELETE FROM TABLE classes WHERE unit = $1 AND day = $2 AND num = $3",
                                                    words[1], words[2], words[3])
                                if err != nil { log.Println("not deleted") }
                                api.SendMessage(fmt.Sprintf("Занятие у группы %s %s на %s паре удалено.",
                                                words[1], words[2], words[3]), id, nil)
                            }
                        }
                    }

                /* команды, доступные только студентам */
                case words[0] == "/unit":
                    if len(words) != 2 {
                        api.SendMessage("Использование: /unit <группа>\nУбедись, что такая группа существует (/units)", id, nil)
                    } else {
                        var exists bool
                        err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM units WHERE name = $1)", words[1]).Scan(&exists)
                        if err != nil { log.Println("not selected from db") }

                        if (!exists) {
                            api.SendMessage(fmt.Sprintf("Группы %s не существует.", words[1]), id, nil)
                        } else if (!isAdmin) && unitName == "" {
                            _, err := db.Exec("INSERT INTO users(id, unit) VALUES($1, $2)", id, words[1])
                            if err != nil { log.Println("not inserted") }
                            api.SendMessage(fmt.Sprintf("Твоя группа - %s.", words[1]), id, nil)
                        } else {
                            if _, err := db.Exec("UPDATE users SET unit = $1 WHERE id = $2", words[1], id); err != nil {
                            log.Println(err)
                            }
                        }
                    }

            case words[0] == "/next":
                currentDate := time.Now()

                rows, err := db.Query("SELECT day FROM classes WHERE day >= $1", currentDate)
                if err != nil { log.Println(err) }
                defer rows.Close()

                var classDates []string
                for rows.Next() {
                    var classDate time.Time
                    if err := rows.Scan(&classDate); err != nil {
                        log.Println(err)
                    }
                    classDates = append(classDates, classDate.Format("02.01.2006"))
                }

                var sb strings.Builder
                sb.WriteString("Ближайшие дни, когда есть пары:\n")

                for _, day := range classDates {
                    sb.WriteString(day)
                    sb.WriteString("\n")
                }
                finalString := sb.String()
                api.SendMessage(finalString, id, nil)

            default:
                api.SendMessage("Я не знаю такой команды", id, nil)
            }
        }
    }
}
