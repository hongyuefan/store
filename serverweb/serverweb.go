package serverweb

import (
	"io/ioutil"
	"net/http"
	"store/core"
	"store/skeleton"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	BaseUrl   = "http://hq.sinajs.cn/list="
	DayUrl    = "http://image.sinajs.cn/newchart/daily/n/"
	MinUrl    = "http://image.sinajs.cn/newchart/min/n/"
	WeekUrl   = "http://image.sinajs.cn/newchart/weekly/n/"
	MonUrl    = "http://image.sinajs.cn/newchart/monthly/n/"
	SufGIF    = ".gif"
	Split_Sem = ";"
	Split_Com = ","
	Split_Quo = "\""
)
const (
	TYPE_GETBYCODE = 0X01
	TYPE_GETALL    = 0x02
)

type WebServer struct {
	core *core.Core
}

func NewWebServer() skeleton.Server {
	return &WebServer{}
}

func (w *WebServer) Run(core *core.Core) {

	w.core = core

	core.UpdateCoreDataCallfunc(w.handler)

	return
}

func (w *WebServer) handler(coreDataArray []*core.CoreData) error {

	var params string

	if len(coreDataArray) <= 0 {
		return nil
	}

	for _, coreData := range coreDataArray {
		params += coreData.Code + Split_Com
	}

	body, err := w.onRequest(params)

	if err != nil {
		return err
	}
	return w.parseParam(body, coreDataArray)
}

func (w *WebServer) onRequest(stores string) (body string, err error) {

	storeUrl := BaseUrl + stores

	rsp, err := http.Get(storeUrl)
	if err != nil {
		return
	}
	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}

	b, err := simplifiedchinese.GBK.NewDecoder().Bytes(buf)
	if err != nil {
		return
	}
	body = string(b)
	return
}

func (w *WebServer) parseParam(body string, coreDataArray []*core.CoreData) (err error) {

	bodyArr := strings.Split(body, Split_Sem)

	for i := 0; i < len(bodyArr)-1; i++ {

		nhead := strings.Index(bodyArr[i], Split_Quo)
		ntail := strings.LastIndex(bodyArr[i], Split_Quo)

		sub_one := bodyArr[i][nhead+1 : ntail]

		subArray := strings.Split(sub_one, Split_Com)

		if len(subArray) < 33 {
			continue
		}

		w.formCoreData(subArray, coreDataArray[i])

	}

	return

}

func (w *WebServer) formCoreData(stringArr []string, coreData *core.CoreData) {

	coreData.Open = stringArr[1]
	coreData.Close = stringArr[2]
	coreData.Now = stringArr[3]
	coreData.High = stringArr[4]
	coreData.Low = stringArr[5]

	coreData.DealNum = stringArr[8]
	coreData.DealMoney = stringArr[9]

	coreData.BuyNum1 = stringArr[10]
	coreData.Buy1 = stringArr[11]
	coreData.BuyNum2 = stringArr[12]
	coreData.Buy2 = stringArr[13]
	coreData.BuyNum3 = stringArr[14]
	coreData.Buy3 = stringArr[15]
	coreData.BuyNum4 = stringArr[16]
	coreData.Buy4 = stringArr[17]
	coreData.BuyNum5 = stringArr[18]
	coreData.Buy5 = stringArr[19]

	coreData.SellNum1 = stringArr[20]
	coreData.Sell1 = stringArr[21]
	coreData.SellNum2 = stringArr[22]
	coreData.Sell2 = stringArr[23]
	coreData.SellNum3 = stringArr[24]
	coreData.Sell3 = stringArr[25]
	coreData.SellNum4 = stringArr[26]
	coreData.Sell4 = stringArr[27]
	coreData.SellNum5 = stringArr[28]
	coreData.Sell5 = stringArr[29]

	b1, _ := strconv.Atoi(coreData.BuyNum1)
	b2, _ := strconv.Atoi(coreData.BuyNum2)
	b3, _ := strconv.Atoi(coreData.BuyNum3)
	b4, _ := strconv.Atoi(coreData.BuyNum4)
	b5, _ := strconv.Atoi(coreData.BuyNum5)

	coreData.TotleBuy = b1 + b2 + b3 + b4 + b5

	s1, _ := strconv.Atoi(coreData.SellNum1)
	s2, _ := strconv.Atoi(coreData.SellNum2)
	s3, _ := strconv.Atoi(coreData.SellNum3)
	s4, _ := strconv.Atoi(coreData.SellNum4)
	s5, _ := strconv.Atoi(coreData.SellNum5)

	coreData.TotleSell = s1 + s2 + s3 + s4 + s5

	coreData.Date = stringArr[30]
	coreData.Time = stringArr[31]

}

func (w *WebServer) HttpServerHanlde(c *gin.Context) {

	sType := c.Query("type")

	nType, _ := strconv.Atoi(sType)

	switch nType {
	case TYPE_GETBYCODE:
		code := c.Query("code")
		c.JSON(200, w.core.GetCoreDatabyCode(code))
		break
	case TYPE_GETALL:
		c.JSON(200, w.core.MapCorData)
		break
	default:
		c.String(400, "error request")
	}

	return
}

func (w *WebServer) Shutdown() {
	w.core.Shutdown()
}
