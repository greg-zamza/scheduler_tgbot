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
var bot_token string
var admin_password string

/* Функция init() выполняется перед main() и
 * служит для определения переменных */
func init() {
    bot_token = os.Getenv("BOT_TOKEN")
    admin_password = os.Getenv("ADMIN_PASSWORD")
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
    api := echotron.NewAPI(bot_token)
    log.Println("started")
    for update := range echotron.PollingUpdates(bot_token) {
        // помогает боту не сломаться от невалидного апдейта
        if update.Message == nil {
            log.Println("Unhandled update")
        } else {
            var id int64 = update.ChatID()
            var words []string = strings.Fields(update.Message.Text)
            api.SendMessage(words[0], id, nil)
            switch {
                case words[0] == "/start":
                    // проверка наличия пользователя в бд
                    rows, err := db.Query("SELECT 1 unit, is_admin FROM users WHERE id = $1", id)
                    if err != nil {
                        log.Println("not selected from db")
                    }

                    for rows.Next() {
                        var unit string
                        var isAdmin bool
                        if err := rows.Scan(&unit, &isAdmin); err != nil {
                            log.Fatal(err)
                        }

                        if isAdmin {
                            api.SendMessage("Вы уже вошли как администратор", id, nil)
                        } else if unit != "" {
                            api.SendMessage(fmt.Sprintf("Вы уже вошли как студент группы %s", unit), id, nil)
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
                    //проверка наличия пользователя в бд и удаление
                    
                    api.SendMessage("вышел", id, nil)

                case words[0] == "Войти":
                    api.SendMessage("студент", id, nil)


                case words[0] == admin_password && len(words) == 1:
                    /*проверка наличия пользователя в бд и если его там нет, то ставим is_admin==1*/
                    api.SendMessage("поздравляю, ты долбоёб", id, nil)

                /* команды, доступные администраторам */
                case words[0] == "/units":
                    fmt.Println("placeholder")
                    //select из бд и api.SendMessage

                case words[0] == "/delunit":
                    fmt.Println("q")
                    //remove from db where words[1]

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

                /* Команда "/s" также доступна студентам, но
                 * обрабатывается в том же кейсе, что и у администраторов
                 */

                default:
                    api.SendMessage("Я не знаю такой команды", id, nil)
            }
        }
    }
}
