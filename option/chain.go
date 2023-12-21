package main

import "plumelog/log"

/**
1.创建一个处理器
2.创建下一个节点处理器
*/

type Handler interface {
	Handle(*Request)
	Next(*Request)
	Finish(*Request) bool
}

type Request struct {
	PayAmount float64 `json:"pay_amount"`
}

type GoodsHandler struct {
	Processor *OrderHandler
}

func (g *GoodsHandler) Handle(r *Request) {
	if g.Finish(r) {
		return
	} else {
		log.Info("商品流程处理结束交由下个流程处理")
		g.Next(r)
	}
}

func (g *GoodsHandler) Next(r *Request) {
	g.Processor = &OrderHandler{}
	g.Processor.Handle(r)
}

func (g *GoodsHandler) Finish(r *Request) bool {
	if r.PayAmount < 500 {
		return true
	}
	return false
}

type OrderHandler struct {
	Processor *ExpressHandler
}

func (o *OrderHandler) Handle(r *Request) {
	r.PayAmount = r.PayAmount*2 + 20
	if o.Finish(r) {
		return
	} else {
		log.Info("订单流程处理结束交由下个流程处理")
		o.Next(r)
	}
}

func (o *OrderHandler) Next(r *Request) {
	o.Processor = &ExpressHandler{}
	o.Processor.Handle(r)
}

func (o *OrderHandler) Finish(r *Request) bool {
	if r.PayAmount < 1000 {
		return true
	}
	return false
}

type ExpressHandler struct {
}

func (e *ExpressHandler) Handle(r *Request) {
	r.PayAmount = r.PayAmount - 1000
	if e.Finish(r) {
		log.Info("订单已完成，记得5星好评哦！")
		return
	} else {
		e.Next(r)
	}
}

func (e *ExpressHandler) Next(r *Request) {
	log.Info("流程异常")
}

func (e *ExpressHandler) Finish(r *Request) bool {
	if r.PayAmount <= 20 {
		return true
	}
	return false
}

//func main() {
//	g := &GoodsHandler{}
//	g.Handle(&Request{
//		PayAmount: 500,
//	})
//}
