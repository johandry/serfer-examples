package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/johandry/serfer"
)

const timeoutTime = 60 * time.Minute

var (
	joinAddr string
	key      string
	debug    bool
)

var mu sync.Mutex

func init() {
	flag.StringVar(&joinAddr, "join", "", "Address of leader to join")
	flag.StringVar(&joinAddr, "j", "", "Address of leader to join")
	flag.StringVar(&key, "key", "", "Secret key to encrypt communication and securely join the cluster")
	flag.StringVar(&key, "k", "", "Secret key to encrypt communication and securely join the cluster")

	flag.BoolVar(&debug, "debug", false, "Debug mode")
}

func main() {
	flag.Parse()

	var keys []string

	// To generate the key with a known text (16 bytes) use this command:
	// echo -n "1234567890123456" | base64 # => "MTIzNDU2Nzg5MDEyMzQ1Ng=="
	// Or with 16 random character use this one:
	// head -c16 /dev/urandom | base64

	if len(key) != 0 {
		keys = make([]string, 1)
		keys[0] = key
	}

	serfer.UserEventHandleFunc("5secEvent", func(e serfer.Event) {
		log.Printf("Event '5secEvent' received")
		ue, ok := serfer.Event2UserEvent(e)
		if !ok {
			log.Printf("[ERROR] Event '5secEvent' is not an UserEvent type")
			return
		}
		log.Printf("event: { name: %q, payload: %q, ltime: %q, coalesce: %t, type: %q }", ue.Name, ue.Payload, ue.LTime, ue.Coalesce, ue.EventType().String())
	})

	serfer.QueryEventHandleFunc("7secQuery", func(e serfer.Event) {
		log.Printf("Query '7secQuery' received")
		q, ok := serfer.Event2Query(e)
		if !ok {
			log.Printf("[ERROR] Query '7secQuery' is not a Query type")
			return
		}
		resp := strings.ToUpper(string(q.Payload))
		if err := q.Respond([]byte(resp)); err != nil {
			log.Printf("[ERROR] Failed to respond query '7secQuery'. %s", err)
			return
		}
		log.Printf("query: { name: %q, payload: %q, ltime: %q, deadline: %q }", q.Name, q.Payload, q.LTime, q.Deadline())
		log.Printf("response: %q", resp)
	})

	serfer.MemberEventHandleFunc("join", func(e serfer.Event) {
		log.Printf("Member 'join' received")
		m, _, ok := serfer.Event2MemberEvent(e)
		if !ok {
			log.Printf("[ERROR] Member 'join' is not a Member-Join type")
			return
		}
		members := []string{}
		for _, memb := range m.Members {
			members = append(members, fmt.Sprintf("{ %s\t%s\t%s:%d\t[%v] }", memb.Name, memb.Status, memb.Addr, memb.Port, memb.Tags))
		}
		log.Printf("members: %s", strings.Join(members, ", "))
	})

	var role string
	if len(joinAddr) == 0 {
		role = "leader"
	} else {
		role = "worker"
	}
	tags := map[string]string{
		"role": role,
	}

	s, err := serfer.StartSecure(keys, tags, nil, joinAddr)
	if err != nil {
		panic(err)
	}
	defer s.Leave()

	if len(joinAddr) == 0 {
		log.Printf("Starting surfer on %s:%d as leader of this cluster", s.BindAddr, s.BindPort)
		log.Printf("Sending events to the cluster")
		sendEvents(s)
	} else {
		log.Printf("Starting surfer on %s:%d and joining to cluster lead by %s", s.BindAddr, s.BindPort, joinAddr)
		log.Printf("Receiving events")
		s.Wait()
	}

	if debug {
		log.Printf("Serfer:     %s", strings.Replace(fmt.Sprintf("%#v", s), ",", "\n"+strings.Repeat(" ", 46), -1))
		log.Printf("Serf:       %s", strings.Replace(fmt.Sprintf("%#v", s.Conf), ",", "\n"+strings.Repeat(" ", 44), -1))
		log.Printf("MemberList: %s", strings.Replace(fmt.Sprintf("%#v", s.Conf.MemberlistConfig), ",", "\n"+strings.Repeat(" ", 50), -1))
		log.Printf("Keyring:    %+v", s.Conf.MemberlistConfig.Keyring)
	}

}

func sendEvents(s *serfer.Serfer) {
	payload := []byte("event_payload")

	tick2s := time.NewTicker(2 * time.Second).C
	tick5s := time.NewTicker(5 * time.Second).C
	tick7s := time.NewTicker(6 * time.Second).C
	timeout := time.NewTimer(timeoutTime).C

	for {
		select {
		// Print members every 2 seconds
		case <-tick2s:
			members := s.Serf().Members()
			var memStr []string
			for _, m := range members {
				memStr = append(memStr, fmt.Sprintf("%s<%s> (%s:%d) [%s]", m.Name, m.Status, m.Addr, m.Port, m.Tags))
			}
			log.Printf("Total nodes: %d [%s]", len(members), strings.Join(memStr, ", "))

		// Send event every 5 seconds
		case <-tick5s:
			s.Event("5secEvent", payload)
			rotate(&payload)

		// Send a query every 7 seconds
		case <-tick7s:
			q, err := s.Query("7secQuery", payload, nil)
			if err != nil {
				log.Printf("[ERROR] sending query '7secQuery'. %s ", err)
				break
			}
			go printResponse(q)
			rotate(&payload)

		case <-timeout:
			s.Leave()
			return
		}
	}
}

func printResponse(q *serfer.QueryResponse) {
	if q == nil {
		return
	}
	r, err := q.GetResponse()
	if err != nil {
		log.Printf("[ERROR] getting response from query '7secQuery'. %s ", err)
		return
	}
	log.Printf("Response from query '7secQuery':")
	for from, resp := range r.Responses {
		log.Printf("  From %s: %s", from, resp)
	}
	log.Printf("Total Ack: %d, Total Response: %d", len(r.Acks), len(r.Responses))
}

func rotate(b *[]byte) {
	mu.Lock()
	defer mu.Unlock()
	i, r := (*b)[:1], (*b)[1:]
	*b = append(r, i...)
}
