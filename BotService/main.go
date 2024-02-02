package main

import (
    "os"
    /*"time"
    "io/ioutil"*/
    "strings"
    /*"strconv"*/
    "fmt"
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

/* Функция init() выполняется перед main() и
 * служит для определения переменных
 */
func init() {
    bot_token = os.Getenv("BOT_TOKEN")
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
    //api := echotron.NewAPI(bot_token)
    fmt.Println("ok")
    for update := range echotron.PollingUpdates(bot_token) {
        // помогает боту не сломаться от невалидного апдейта
        if update.Message == nil {
            log.Println("Unhandled update")
        } else {
            var words []string = strings.Fields(update.Message.Text)
            fmt.Printf("Первое слово (команда): %s\n", words[0])
            fmt.Println(words)
        }
    }
}
