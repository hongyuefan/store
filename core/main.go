package core

import (
	"fmt"
	"store/context"
	"sync"
	"time"
)

type Core struct {
	Cont       *context.Context
	CodeArray  []string
	MapCorData map[string]*CoreData
	MapClass   map[uint32]string
	ChanExit   chan bool
	lock       sync.RWMutex
	lock_class sync.RWMutex
}

func NewCore() *Core {
	c := new(Core)
	c.Cont = context.NewContext(1, 5, 20, "8000", "root:350999@tcp(localhost:3306)/store?charset=utf8")
	c.MapCorData = make(map[string]*CoreData, 0)
	c.MapClass = make(map[uint32]string, 0)
	c.CodeArray = make([]string, 0)
	c.ChanExit = make(chan bool, 1)

	return c
}

type CoreData struct {
	ClassId   uint32
	Code      string
	Name      string
	Score     float64
	Open      string
	Close     string
	Now       string
	High      string
	Low       string
	DealNum   string
	DealMoney string
	Buy1      string
	BuyNum1   string
	Buy2      string
	BuyNum2   string
	Buy3      string
	BuyNum3   string
	Buy4      string
	BuyNum4   string
	Buy5      string
	BuyNum5   string
	Sell1     string
	SellNum1  string
	Sell2     string
	SellNum2  string
	Sell3     string
	SellNum3  string
	Sell4     string
	SellNum4  string
	Sell5     string
	SellNum5  string
	TotleBuy  int
	TotleSell int
	Date      string
	Time      string
}

func (c *Core) NewCoreData(classId int, code string, name string, score float64) *CoreData {
	return &CoreData{
		ClassId: uint32(classId),
		Code:    code,
		Name:    name,
		Score:   score,
	}
}

/*
板块与个股区分
委比、瞬时增长率、涨跌幅、板块增速、成交活跃度、各自权重值
*/

func (c *Core) SliceCodeArray(num int) (sliceArray [][]string) {

	if num <= 1 {
		sliceArray = append(sliceArray, c.CodeArray)
	}
	var n int

	sliceArray = make([][]string, num)

	nlen := len(c.CodeArray)

	result := nlen / num

	ys := nlen % num

	for i := 0; i < nlen-ys; i++ {
		sliceArray[i/result] = append(sliceArray[i/result], c.CodeArray[i])
	}

	for i := nlen - ys; i < nlen; i++ {
		sliceArray[n] = append(sliceArray[n], c.CodeArray[i])
		n++
	}

	return

}

func (c *Core) CallbackInitMap(classId int, code string, score float64, name, class string) error {

	c.InsertClass(uint32(classId), class)

	c.InsertCoreData(c.NewCoreData(classId, code, name, score))

	c.CodeArray = append(c.CodeArray, code)

	return nil
}

func (c *Core) GetClassbyId(id uint32) (class string) {
	c.lock_class.RLock()
	defer c.lock_class.RUnlock()

	class = c.MapClass[id]

	return
}

func (c *Core) InsertClass(id uint32, class string) {
	c.lock_class.Lock()
	defer c.lock_class.Unlock()

	c.MapClass[id] = class

	return
}
func (c *Core) GetCoreDatabyCode(code string) (coreData *CoreData) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	coreData = c.MapCorData[code]

	return
}

func (c *Core) UpdateCoreData(code string, coreData *CoreData) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.MapCorData[code] = coreData

	return
}

func (c *Core) InsertCoreData(coreData *CoreData) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.MapCorData[coreData.Code] = coreData

	return
}

type callbackRangeCoreData func(*CoreData)

func (c *Core) RangeCoreData(f callbackRangeCoreData) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, coreData := range c.MapCorData {
		f(coreData)
	}
}

func (c *Core) DelCoreData(code string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.MapCorData, code)

	return
}

type callbackupdateCoreData func([]*CoreData) error

func (c *Core) UpdateCoreDataCallfunc(f callbackupdateCoreData) {

	var (
		wg sync.WaitGroup
	)

	sliceArray := c.SliceCodeArray(c.Cont.Threads)

	for _, codeArray := range sliceArray {

		wg.Add(1)

		go func(codeArray []string) {

			defer wg.Done()

			var (
				count   int
				current int
			)

			chanTimer := time.Tick(time.Second * time.Duration(c.Cont.IntervalTime))

			for {

				//				if c.IsCloseTime() {
				//					continue
				//				}

				select {
				case <-chanTimer:

					var coreDataArray []*CoreData

					for count = 0; count < c.Cont.OnceUpdateCount; count++ {

						if current >= len(codeArray) {
							current = 0
							break
						}

						code := codeArray[current]

						coreData := c.GetCoreDatabyCode(code)

						coreDataArray = append(coreDataArray, coreData)

						current++

					}

					c.lock.Lock()
					if err := f(coreDataArray); err != nil {
						fmt.Println("UpdateCoreDataCall error:", err)
					}
					c.lock.Unlock()

				case <-c.ChanExit:
					return
				}
			}

		}(codeArray)

	}
	wg.Wait()

}
func (c *Core) Shutdown() {
	select {
	case <-c.ChanExit:
	default:
		close(c.ChanExit)
	}
}
func (c *Core) IsCloseTime() bool {

	weekDay := time.Now().Weekday()

	if "Sunday" == weekDay.String() || "Saturday" == weekDay.String() {
		return true
	}

	h, m, _ := time.Now().Clock()

	if h >= 15 || h < 9 || (h == 9 && m < 30) {
		return true
	}

	return false
}
