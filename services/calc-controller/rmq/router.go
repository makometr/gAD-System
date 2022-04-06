package rmq

import (
	"encoding/json"
	"fmt"
	"gAD-System/services/calc-controller/model"
)

type Router struct {
	input        <-chan Message
	routingTable map[model.MsgID]chan<- model.ResultFromCalc
}

func InitFilter(inputFromRMQ <-chan Message) Router {
	filter := Router{input: inputFromRMQ, routingTable: make(map[model.MsgID]chan<- model.ResultFromCalc)}

	go func() {
		for msg := range inputFromRMQ {
			sendChan := filter.routingTable[msg.MessageID]
			// expr, err := ProtoToMsg(msg.Body) // TODO
			var result model.Result
			err := json.Unmarshal(msg.Body, &result)
			if err != nil {
				fmt.Println("error while unmarhsal result", err)
				result = model.Result{}
			}

			sendChan <- model.ResultFromCalc{Result: result, ID: msg.MessageID}
			delete(filter.routingTable, msg.MessageID)
		}
	}()

	return filter
}

func (r *Router) AddRoute(ID model.MsgID, goal chan<- model.ResultFromCalc) {
	if val, ok := r.routingTable[ID]; ok {
		fmt.Println("key existed in router-table!: ", val)
	}
	r.routingTable[ID] = goal
}
