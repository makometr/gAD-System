package rmq

import (
	"fmt"
	pr_result "gAD-System/internal/proto/result/event"
	"gAD-System/services/calc-controller/model"

	"google.golang.org/protobuf/proto"
)

type Router struct {
	input        <-chan MessageFromRMQ
	routingTable map[model.MsgID]chan<- pr_result.Event
}

func InitFilter(inputFromRMQ <-chan MessageFromRMQ) Router {
	filter := Router{input: inputFromRMQ, routingTable: make(map[model.MsgID]chan<- pr_result.Event)}

	go func() {
		for msg := range inputFromRMQ {
			sendChan := filter.routingTable[msg.MessageID]

			var result pr_result.Event
			if err := proto.Unmarshal(msg.Body, &result); err != nil {
				result.Result = &pr_result.Event_ErrorMsg{ErrorMsg: "error convert proto to struct"}
			}

			sendChan <- result
			delete(filter.routingTable, msg.MessageID)
		}
	}()

	return filter
}

func (r *Router) AddRoute(id model.MsgID, goal chan<- pr_result.Event) {
	if val, ok := r.routingTable[id]; ok {
		fmt.Println("key existed in router-table!: ", val)
	}
	r.routingTable[id] = goal
}
