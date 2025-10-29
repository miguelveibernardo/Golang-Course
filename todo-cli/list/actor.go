package list

import (
	"errors"
	"sync"
)

var ErrActorStopped = errors.New("list actor has been stopped")

type commandType int

const (
	cmdAdd commandType = iota
	cmdUpdateDesc
	cmdUpdateStatus
	cmdDelete
	cmdGetAll
)

type command struct {
	cmdType commandType
	id      int
	value   string
	replyCh chan []Item
	errCh   chan error
}

// runs as a single actior go routine processing all commands
type ListActor struct {
	items  []Item
	cmdCh  chan command
	stopCh chan struct{}
	wg     sync.WaitGroup
}

func NewListActor(initial []Item) *ListActor {
	m := &ListActor{
		items:  initial,
		cmdCh:  make(chan command, 1000),
		stopCh: make(chan struct{}),
	}
	m.wg.Add(1)
	go m.run()
	return m
}

func (m *ListActor) run() {
	defer m.wg.Done()
	for {
		select {
		case cmd, ok := <-m.cmdCh:
			if !ok {
				return //close chanel, stop
			}
			switch cmd.cmdType {
			case cmdAdd:
				m.items = Add(m.items, cmd.value)
				cmd.replyCh <- m.items
				cmd.errCh <- nil

			case cmdUpdateDesc:
				updated, err := UpdateDescription(m.items, cmd.id, cmd.value)
				if err == nil {
					m.items = updated
				}
				cmd.replyCh <- m.items
				cmd.errCh <- err

			case cmdUpdateStatus:
				updated, err := UpdateStatus(m.items, cmd.id, cmd.value)
				if err == nil {
					m.items = updated
				}
				cmd.replyCh <- m.items
				cmd.errCh <- err

			case cmdDelete:
				m.items = Delete(m.items, cmd.id)
				cmd.replyCh <- m.items
				cmd.errCh <- nil

			case cmdGetAll:
				cmd.replyCh <- append([]Item{}, m.items...)
				cmd.errCh <- nil
			}
		case <-m.stopCh:
			return
		}
	}
}

func (m *ListActor) Stop() {
	//closeing channels will signal shutdown
	close(m.stopCh)
	close(m.cmdCh)
	m.wg.Wait()
}

func (m *ListActor) send(cmd command) ([]Item, error) {
	select {
	case <-m.stopCh:
		return nil, ErrActorStopped
	default:
	}
	select {
	case m.cmdCh <- cmd:
		return <-cmd.replyCh, <-cmd.errCh
	case <-m.stopCh:
		return nil, ErrActorStopped
	}
}

func (m *ListActor) Add(desc string) ([]Item, error) {
	reply := make(chan []Item, 1)
	errCh := make(chan error, 1)
	return m.send(command{cmdType: cmdAdd, value: desc, replyCh: reply, errCh: errCh})
}

func (m *ListActor) UpdateDescription(id int, desc string) ([]Item, error) {
	reply := make(chan []Item, 1)
	errCh := make(chan error, 1)
	return m.send(command{cmdType: cmdUpdateDesc, id: id, value: desc, replyCh: reply, errCh: errCh})
}

func (m *ListActor) UpdateStatus(id int, status string) ([]Item, error) {
	reply := make(chan []Item, 1)
	errCh := make(chan error, 1)
	return m.send(command{cmdType: cmdUpdateStatus, id: id, value: status, replyCh: reply, errCh: errCh})
}

func (m *ListActor) Delete(id int) ([]Item, error) {
	reply := make(chan []Item, 1)
	errCh := make(chan error, 1)
	return m.send(command{cmdType: cmdDelete, id: id, replyCh: reply, errCh: errCh})

}

func (m *ListActor) GetAll() ([]Item, error) {
	reply := make(chan []Item, 1)
	errCh := make(chan error, 1)
	return m.send(command{cmdType: cmdGetAll, replyCh: reply, errCh: errCh})
}
