package main

import (
    "os"
    /*"time"
    "io/ioutil"*/
    "strings"
    /*"strconv"*/
    //"fmt"
    "log"
    /*"encoding/csv"
    "encoding/json"
    "database/sql"
    "net/http"
    "net/url"*/

    _ "github.com/lib/pq"
    "github.com/NicoNex/echotron/v3"
)

/* Объявление необходимых переменных */
var bot_token string
var admin_password string
/* Функция init() выполняется перед main() и
 * служит для определения переменных
 */
func init() {
    bot_token = os.Getenv("BOT_TOKEN")
    admin_password = os.Getenv("ADMIN_PASSWORD")
}

func main() {
/*    // connecting to database
    var conn_params string = fmt.Sprintf(
        "user=%s dbname=%s sslmode=disable host=DatabaseService password=%s",
        os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"))

    db, err := sql.Open("postgres", conn_params)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
*/
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
                    /* если юзер есть в бд - пишем "вы уже авторизированы как n. Вы можете выйти с помощью команды /exit" */
                    api.SendMessage("вошь", id, nil)
                case words[0] == "/exit":
                    /*удаление пользователя из бд (а если его и так нет в бд - пишем: "вы не авторизованы"*/
                    api.SendMessage("вышел", id, nil)
                case words[0] == admin_password && len(words) == 1:
                    /*проверка наличия пользователя в бд и если его там нет, то ставим is_admin==1*/
                    api.SendMessage("поздравляю, ты долбоёб", id, nil)
                /*case "":

                case "":

                case "":

                case "":

                case "":

                case "":

                case "":

                case "":*/

                default:
                    api.SendMessage("Я не знаю такой команды", id, nil)
            }
        }
    }
}
