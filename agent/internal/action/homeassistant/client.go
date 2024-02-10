package homeassistant

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"sync"
	"time"

	"github.com/Doridian/go-haws"
	"github.com/Doridian/streamdeckpi/agent/internal/controller"
)

type haStateReceiver interface {
	OnState(entityID string, state haws.State) error
}

type haInstance struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`

	client *haws.Client

	stateReceiverMap map[string][]haStateReceiver
	stateLock        *sync.Mutex
	states           map[string]haws.State
}

var haInstances = map[string]*haInstance{}
var haInstanceLock sync.Mutex

func GetHomeAssistant(ctrl controller.Controller, name string) (*haInstance, error) {
	if name == "" {
		name = "default"
	}

	haInstanceLock.Lock()
	defer haInstanceLock.Unlock()

	instance, ok := haInstances[name]
	if !ok {
		instance = &haInstance{
			stateReceiverMap: make(map[string][]haStateReceiver),
			stateLock:        &sync.Mutex{},
			states:           make(map[string]haws.State),
		}

		path := path.Join("/global/homeassistant", fmt.Sprintf("%s.yml", name))
		path, err := ctrl.CleanPath(path)
		if err != nil {
			return nil, err
		}

		err = ctrl.LoadConfig(path, instance)
		if err != nil {
			return nil, err
		}

		instance.client = haws.NewClient(instance.Url, instance.Token, instance.onConnect, time.Duration(5)*time.Second)

		err = instance.client.Open()
		if err != nil {
			return nil, err
		}

		err = instance.client.AddEventHandler(haws.EventStateChanged, instance)
		if err != nil {
			instance.client.Close()
			return nil, err
		}

		haInstances[name] = instance
	}

	return instance, nil
}

func (i *haInstance) onConnect() {
	log.Printf("Websocket connection established")

	err := i.client.WaitAuth()
	if err != nil {
		i.client.Close()
		log.Printf("onConnect() error WaitAuth(): %v", err)
		return
	}

	log.Printf("Websocket connection authenticated")

	err = i.GetStates()
	if err != nil {
		i.client.Close()
		log.Printf("onConnect() error GetStates(): %v", err)
		return
	}

	log.Printf("Websocket connection handshake done")
}

func (i *haInstance) GetStates() error {
	states, err := i.client.GetStates()
	if err != nil {
		return err
	}

	i.stateLock.Lock()
	defer i.stateLock.Unlock()

	i.states = make(map[string]haws.State)
	for _, state := range states {
		i.states[state.EntityID] = state

		recvArr := i.stateReceiverMap[state.EntityID]
		if recvArr == nil {
			continue
		}

		stateCopyRef := state
		for _, recv := range recvArr {
			go i.OnStateHandleError(recv, &haws.StateChangeEvent{
				EntityID: stateCopyRef.EntityID,
				NewState: &stateCopyRef,
			})
		}
	}

	return nil
}

func (i *haInstance) OnStateHandleError(recv haStateReceiver, evt *haws.StateChangeEvent) {
	err := recv.OnState(evt.EntityID, *evt.NewState)
	if err != nil {
		log.Printf("Error handling state change event for %s: %v", evt.EntityID, err)
	}
}

func (i *haInstance) OnEvent(eventData *haws.EventData) {
	evt := &haws.StateChangeEvent{}
	err := json.Unmarshal(eventData.Data, evt)
	if err != nil {
		log.Printf("Invalid state change event JSON: %v: %s", err, eventData.Data)
		return
	}

	if evt.NewState == nil {
		log.Printf("Invalid state change event: missing NewState: %s", eventData.Data)
		return
	}

	i.stateLock.Lock()
	i.states[evt.EntityID] = *evt.NewState
	recvArr := i.stateReceiverMap[evt.EntityID]
	i.stateLock.Unlock()

	if recvArr == nil {
		return
	}
	for _, recv := range recvArr {
		go i.OnStateHandleError(recv, evt)
	}
}

func (i *haInstance) RegisterStateReceiver(recv haStateReceiver, entityID string) {
	i.stateLock.Lock()
	defer i.stateLock.Unlock()

	arr, ok := i.stateReceiverMap[entityID]
	if ok {
		arr = append(arr, recv)
	} else {
		arr = []haStateReceiver{recv}
	}
	i.stateReceiverMap[entityID] = arr

	state, ok := i.states[entityID]
	if ok {
		go i.OnStateHandleError(recv, &haws.StateChangeEvent{
			EntityID: entityID,
			NewState: &state,
		})
	}
}
