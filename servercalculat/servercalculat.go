package servercalculat

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"store/core"
	"store/initstore"
	"store/skeleton"
	"strconv"
	"time"

	"store/util/log"

	"github.com/gin-gonic/gin"
)

const (
	Weight_Class = 0.3
	Weight_WeiBi = 0.3
	//Weight_RateSellNum = 0.125
	//Weight_RatBuyNum   = 0.125
	Weight_RateUping = 0.3
	Weight_RateUp    = 0.1
)

type CalculatServer struct {
	core             *core.Core
	mapCoreData      map[string]*CalculateData
	sortArray        SortStructArray
	mapCalculatClass map[uint32]*CalculatClass
}

type Calculat_Store struct {
	Weibi       float64
	RateSellNum float64
	RateBuyNum  float64
	RateUp      float64
	RateUping   float64
}

type SortStructArray []*SortStruct

type SortStruct struct {
	Code       string
	Score      float64
	Day        string
	Min        string
	ClassScore float64
	ClassName  string
}

type CalculatClass struct {
	ClassName string
	ClassId   uint32
	Score     float64
}

func (sort SortStructArray) Len() int {
	return len(sort)
}
func (sort SortStructArray) Less(i, j int) bool {

	if sort[i].Score > sort[j].Score {
		return true
	}
	return false
}
func (sort SortStructArray) Swap(i, j int) {
	var tmp *SortStruct = sort[i]
	sort[i] = sort[j]
	sort[j] = tmp
	return
}
func (c *CalculatServer) GetClassScore(code uint32) float64 {
	if class, ok := c.mapCalculatClass[code]; ok {
		return class.Score
	}
	return 0
}
func (c *CalculatServer) UpdateSort(code string, score, classcore float64) {

	for _, s := range c.sortArray {
		if s.Code == code {
			s.Score += score
			s.Score = Simulink(s.Score)
			s.ClassScore = classcore
			return
		}
	}

	s := new(SortStruct)
	s.Code = code
	s.Score = Simulink(c.core.MapCorData[code].Score + score)
	s.Day = "http://image.sinajs.cn/newchart/daily/n/" + code + ".gif"
	s.Min = "http://image.sinajs.cn/newchart/min/n/" + code + ".gif"
	s.ClassName = c.core.GetClassByCode(code)

	c.sortArray = append(c.sortArray, s)

	return

}

type CalculateData struct {
	preCore  *core.CoreData
	calCulat *Calculat_Store
}

func NewCalculateData() *CalculateData {
	return &CalculateData{
		preCore:  new(core.CoreData),
		calCulat: new(Calculat_Store),
	}
}

func NewCalculatServer() skeleton.Server {

	return &CalculatServer{
		mapCoreData:      make(map[string]*CalculateData, 0),
		sortArray:        make([]*SortStruct, 0),
		mapCalculatClass: make(map[uint32]*CalculatClass, 0),
	}
}

func (c *CalculatServer) UpdateCalculatClass(classId uint32, className string, score float64) {

	calClass, ok := c.mapCalculatClass[classId]
	if ok {
		calClass.Score += score
		return
	}
	calClass = new(CalculatClass)
	calClass.ClassId = classId
	calClass.ClassName = className
	calClass.Score = score
	c.mapCalculatClass[classId] = calClass

	return
}

func (c *CalculatServer) Run(core *core.Core) {

	c.core = core

	chanTimer := time.Tick(time.Second * time.Duration(10))

	for {
		//		if core.IsCloseTime() {
		//			continue
		//		}
		select {
		case <-chanTimer:
			c.CalculatOnce(core)
		case <-core.ChanExit:
			return
		}
	}
}

func (c *CalculatServer) CalculatOnce(core *core.Core) {
	c.mapCalculatClass = make(map[uint32]*CalculatClass, 0)
	core.RangeCoreData(c.calculateOneCode)
	sort.Sort(c.sortArray)
}

func (c *CalculatServer) Shutdown() {
	c.SaveScore()
	c.core.Shutdown()
}

const (
	TYPE_STORE = 0
	TYPE_CLASS = 1
)

func (c *CalculatServer) HttpServerHanlde(cont *gin.Context) {

	sType := cont.Query("type")

	nType, _ := strconv.Atoi(sType)

	switch nType {
	case TYPE_STORE:
		cont.HTML(http.StatusOK, "calculate.tmpl", gin.H{
			"SortArray": c.sortArray,
		})
		break
	case TYPE_CLASS:
		for _, m := range c.mapCalculatClass {
			cont.JSON(200, m)
		}
		break
	}

	return
}

func (c *CalculatServer) copyData(calStore *CalculateData, coreData *core.CoreData) {

	calStore.preCore.Buy1 = coreData.Buy1
	calStore.preCore.Buy2 = coreData.Buy2
	calStore.preCore.Buy3 = coreData.Buy3
	calStore.preCore.Buy4 = coreData.Buy4
	calStore.preCore.Buy5 = coreData.Buy5
	calStore.preCore.BuyNum1 = coreData.BuyNum1
	calStore.preCore.BuyNum2 = coreData.BuyNum2
	calStore.preCore.BuyNum3 = coreData.BuyNum3
	calStore.preCore.BuyNum4 = coreData.BuyNum4
	calStore.preCore.BuyNum5 = coreData.BuyNum5
	calStore.preCore.ClassId = coreData.ClassId
	calStore.preCore.Close = coreData.Close
	calStore.preCore.Code = coreData.Code
	calStore.preCore.ClassId = coreData.ClassId
	calStore.preCore.DealNum = coreData.DealNum
	calStore.preCore.DealMoney = coreData.DealMoney
	calStore.preCore.Date = coreData.Date
	calStore.preCore.High = coreData.High
	calStore.preCore.Low = coreData.Low
	calStore.preCore.Name = coreData.Name
	calStore.preCore.Now = coreData.Now
	calStore.preCore.Open = coreData.Open
	calStore.preCore.Sell1 = coreData.Sell1
	calStore.preCore.Sell2 = coreData.Sell2
	calStore.preCore.Sell3 = coreData.Sell3
	calStore.preCore.Sell4 = coreData.Sell4
	calStore.preCore.Sell5 = coreData.Sell5
	calStore.preCore.SellNum1 = coreData.SellNum1
	calStore.preCore.SellNum2 = coreData.SellNum2
	calStore.preCore.SellNum3 = coreData.SellNum3
	calStore.preCore.SellNum4 = coreData.SellNum4
	calStore.preCore.SellNum5 = coreData.SellNum5
	calStore.preCore.TotleBuy = coreData.TotleBuy
	calStore.preCore.TotleSell = coreData.TotleSell
	calStore.preCore.Time = coreData.Time
	return

}

func (c *CalculatServer) calculateOneCode(o *core.CoreData) {

	calculatStore, ok := c.mapCoreData[o.Code]

	if !ok {
		coreData := NewCalculateData()
		c.copyData(coreData, o)
		c.mapCoreData[o.Code] = coreData
		return
	}

	if calculatStore.preCore.TotleBuy == 0 || calculatStore.preCore.TotleSell == 0 {
		c.copyData(calculatStore, o)
		return
	}

	calculatStore.calCulat.Weibi = (float64(calculatStore.preCore.TotleBuy) - float64(calculatStore.preCore.TotleSell)) / (float64(calculatStore.preCore.TotleBuy) + float64(calculatStore.preCore.TotleSell))

	now, _ := strconv.ParseFloat(o.Now, 32)
	close, _ := strconv.ParseFloat(o.Close, 32)
	pnow, _ := strconv.ParseFloat(calculatStore.preCore.Now, 32)

	calculatStore.calCulat.RateUp = (now - close) / close
	calculatStore.calCulat.RateUping = (now - pnow) / pnow

	calculatStore.calCulat.RateBuyNum = (float64(o.TotleBuy) - float64(calculatStore.preCore.TotleBuy)) / float64(calculatStore.preCore.TotleBuy)
	calculatStore.calCulat.RateSellNum = (float64(o.TotleSell) - float64(calculatStore.preCore.TotleSell)) / float64(calculatStore.preCore.TotleSell)

	c.UpdateCalculatClass(o.ClassId, c.core.GetClassbyId(o.ClassId), calculatStore.calCulat.RateUping)

	score := c.GetClassScore(o.ClassId)*Weight_Class + calculatStore.calCulat.Weibi*Weight_WeiBi + calculatStore.calCulat.RateUping*Weight_RateUping + calculatStore.calCulat.RateUp*Weight_RateUp

	c.UpdateSort(o.Code, Simulink(score), c.GetClassScore(o.ClassId))

	log.GetLog().LogDebug("code:", o.Code, "weibi:", calculatStore.calCulat.Weibi, "RateUp:", calculatStore.calCulat.RateUp, "RateUping:", calculatStore.calCulat.RateUping, "Score:", score, "Score:", Simulink(score))

	return

}

func (c *CalculatServer) SaveScore() {
	for _, store := range c.sortArray {
		if err := initstore.SaveScore(store.Code, store.Score); err != nil {
			fmt.Println(err)
		}
	}
}
func Simulink(x float64) float64 {
	return 1/(1+math.Exp(-x)) - 0.5
}
