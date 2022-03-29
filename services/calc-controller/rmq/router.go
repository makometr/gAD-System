package rmq

import (
	"fmt"
)

type Router struct {
	input        <-chan Message
	routingTable map[MsgID]chan<- ExpressionWithID
}

func InitFilter(inputFromRMQ <-chan Message) Router {
	filter := Router{input: inputFromRMQ, routingTable: make(map[MsgID]chan<- ExpressionWithID)}

	go func() {
		for msg := range inputFromRMQ {
			sendChan := filter.routingTable[msg.MessageID]

			expr, err := ProtoToMsg(msg.Body) // TODO
			if err != nil {
				fmt.Printf("proto to msg failed(")
				expr = "success"
			}

			sendChan <- ExpressionWithID{Expr: expr, Id: msg.MessageID}
			delete(filter.routingTable, msg.MessageID)
			fmt.Println("Msg routed with id:", msg.MessageID)
		}
	}()

	return filter
}

func (r *Router) AddRoute(ID MsgID, goal chan<- ExpressionWithID) {
	if val, ok := r.routingTable[ID]; ok {
		fmt.Println("key existed in router-table!: ", val)
	}
	r.routingTable[ID] = goal
	fmt.Println("Msg id registrated in router:", ID)
}
