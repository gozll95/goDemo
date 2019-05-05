package model

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/zhu/qvm/server/enums"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Ticket_FindAndSave(t *testing.T) {
	ticket := NewTicketModel(
		rand.Uint32(),
		uuid.NewV4().String(),
		enums.ResourceTypeIp,
	)

	err := ticket.Save()
	assert.Nil(t, err)

	// duplicate save
	err = ticket.Save()
	assert.Nil(t, err)

	// find
	newTicket, err := Ticket.Find(ticket.Id.Hex())
	assert.Nil(t, err)
	assert.Equal(t, ticket.Uid, newTicket.Uid)
	assert.Equal(t, ticket.ResourceType, newTicket.ResourceType)
	assert.Equal(t, ticket.RegionId, newTicket.RegionId)
	assert.Equal(t, enums.TicketStatusTypeIdle, newTicket.Status)
}

func Test_Ticket_FindbyRegionIdAndResourcetype(t *testing.T) {
	ticket := NewTicketModel(
		rand.Uint32(),
		uuid.NewV4().String(),
		enums.ResourceTypeInstance,
	)

	err := ticket.Save()
	assert.Nil(t, err)

	// find
	newTicket, err := Ticket.FindbyRegionIdAndResourcetype(ticket.Uid, ticket.RegionId, ticket.ResourceType)
	assert.Nil(t, err)
	assert.Equal(t, ticket.Id, newTicket.Id)
	assert.False(t, newTicket.isNewRecord)

}

func Test_Ticket_Remove(t *testing.T) {
	ticket := NewTicketModel(
		rand.Uint32(),
		uuid.NewV4().String(),
		enums.ResourceTypeIp,
	)

	ticket.Save()

	err := ticket.Remove()
	assert.Nil(t, err)

	// find
	newTicket, err := Ticket.Find(ticket.Id.Hex())
	assert.Nil(t, newTicket)
	assert.Equal(t, ErrNotFound, err)
}

func Test_Ticket_FindAndLock(t *testing.T) {
	ticket1 := NewTicketModel(
		rand.Uint32(),
		uuid.NewV4().String(),
		enums.ResourceTypeIp,
	)
	ticket2 := NewTicketModel(
		rand.Uint32(),
		uuid.NewV4().String(),
		enums.ResourceTypeIp,
	)
	ticket3 := NewTicketModel(
		rand.Uint32(),
		uuid.NewV4().String(),
		enums.ResourceTypeIp,
	)
	ticket1.Save()
	ticket2.Save()
	ticket3.Save()

	// find and case more than 3 to test not found cases and sleep 3 seconds
	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(j int) {
			ticket, err := Ticket.FindAndLock()
			assert.Equal(t, err, ErrNotFound)
			assert.Nil(t, ticket)
			wg.Done()
		}(i)
	}

	wg.Wait()

	err := ticket1.Unlock()
	assert.Nil(t, err)
	err = ticket2.Unlock()
	assert.Nil(t, err)
	err = ticket3.Unlock()
	assert.Nil(t, err)
	// test time lock in three seconds
	ticket, err := Ticket.FindAndLock()
	assert.Nil(t, ticket)
	assert.Equal(t, ErrNotFound, err)

	// sleep 3 seconds
	//time.Sleep(time.Second * 3)
	//ticket, err = Ticket.FindAndLock()
	//assert.Nil(t, err)
	//assert.NotNil(t, ticket)

}
