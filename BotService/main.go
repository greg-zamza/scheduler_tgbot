package main

import (
    "os"
    /*"time"
    "io/ioutil"*/
    "strings"
    /*"strconv"*/
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

/exit - выйти из аккаунта администратора`
}

func main() {
    // connecting to database
    var conn_params string = fmt.Sprintf(
        "user=%s dbname=%s sslmode=disable host=DatabaseService password=%s",
        os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"))

    db, err := sql.Open("postgres", conn_params)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // connecting to telegram API
    api := echotron.NewAPI(botToken)
    log.Println("started")
    for update := range echotron.PollingUpdates(botToken) {
        // помогает боту не сломаться от невалидного апдейта
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
                var unit string
                if err := rows.Scan(&unitName, &isAdmin); err != nil {
                    log.Println(err)
                }
            }

            switch {
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
                        //units
                    }

                case words[0] == adminPassword && len(words) == 1:
                    if isAdmin {
                        api.SendMessage("Вы уже вошли как администратор", id, nil)
                    } else if unitName != "" {
                         if _, err := db.Exec("UPDATE TABLE users SET is_admin = 1 WHERE id = $1", id); err != nil {
                            log.Println(err)
                         }
                        api.SendMessage(adminText, id, nil)
                    } else {
                        if _, err := db.Exec("INSERT INTO users(id, is_admin) VALUES($1, $2)", id, 1); err != nil {
                            log.Println(err)
                        }
                    }

                /* команды, доступные администратору */
                case words[0] == "/units":
                    rows, err := db.Query("SELECT name FROM groups")
                    if err != nil { log.Println(err) }
                    defer rows.Close()

                    var groupNames []string
                    for rows.Next() {
                        var groupName string
                        if err := rows.Scan(&groupName); err != nil {
                            log.Println(err)
                        }
                        groupNames = add(groupNames, groupName)
                    }

                    var sb strings.Builder
                    sb.WriteString("Список групп:\n")

                    for _, group := range groupNames {
                        sb.WriteString(group)
                        sb.WriteString("\n")
                    }
                    finalString := sb.String()
                    api.SendMessage(finalString, id, nil)

                case words[0] == "/delunit":
                    if (!isAdmin) {
                        api.SendMessage("Я не знаю такой команды", id, nil)
                    } else if len(words) != 2 {
                        api.SendMessage("Использование: /delunit <группа>")
                    } else {
                        // проверка наличия группы в бд
                        var exists bool
                        err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM units WHERE name = $1)", words[1]).Scan(&exists)
                        if err != nil { log.Println("not selected from db" }

                        if (!exists) {
                            msg := fmt.Sprintf("Группы %s не существует.\nВы можете посмотреть список групп с помощью /units", words[1])
                            api.SendMessage(msg, id, nil)
                        } else {
                            _, err := db.Exec("DELETE FROM TABLE units WHERE name = $1", words[1])
                            if err != nil { log.Println("not deleted") }
                            api.SendMessage(fmt.Sprintf("Группа %s и все занятия, связанные с ней, удалены.", words[1]), id, nil)
                        }
                    }

                case words[0] == "/addunit":
                    fmt.Println("q")

                case words[0] == "s":
                    //расписание
                    fmt.Println("q")

                case words[0] == "/add":
                    //добавить
                    fmt.Println("q")

                case words[0] == "/del":
                    //удалить
                    fmt.Println("q")

                /* команды, доступные студентам */
                case words[0] == "/unit":
                    //объявить свой класс
                    fmt.Println("Q")

                /* Команды "/s" и "/units" также доступны студентам, но
                 * обрабатываются в тех же кейсах, что и у администраторов
                 */

                default:
                    api.SendMessage("Я не знаю такой команды", id, nil)
            }
        }
    }
}
