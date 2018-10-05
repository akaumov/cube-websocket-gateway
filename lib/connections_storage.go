package lib

import (
	"sync"
)

//TODO: optimize

type ConnectionsStats struct {
	NumberOfUsers                int
	NumberOfDevices              int
	NumberOfNotLoggedConnections int
}

type ConnectionsStorage struct {
	mutex                        sync.RWMutex
	connectionsById              map[ConnectionId]*Connection
	numberOfNotLoggedConnections int
}

func NewConnectionsStorage() *ConnectionsStorage {
	return &ConnectionsStorage{
		mutex:                        sync.RWMutex{},
		connectionsById:              make(map[ConnectionId]*Connection),
		numberOfNotLoggedConnections: 0,
	}
}

func (s *ConnectionsStorage) AddNewConnection(connection *Connection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.numberOfNotLoggedConnections++
	s.connectionsById[connection.id] = connection
}

func (s *ConnectionsStorage) RemoveConnection(connection *Connection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.removeConnection(connection)
}

func (s *ConnectionsStorage) removeConnection(connection *Connection) {

	connectionId, userId, _ := connection.GetInfo()

	connectionBefore := s.connectionsById[connectionId]
	if connectionBefore == nil {
		return
	}

	delete(s.connectionsById, connectionId)

	if userId == "" {
		s.numberOfNotLoggedConnections--
		return
	}
}

func (s *ConnectionsStorage) GetUserConnections(userId UserId) []*Connection {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	connections := []*Connection{}

	for _, connection := range s.connectionsById {
		if connection.userId == userId {
			connections = append(connections, connection)
		}
	}

	return connections
}

func (s *ConnectionsStorage) GetDeviceConnections(userId UserId, deviceId DeviceId) []*Connection {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	connections := []*Connection{}

	for _, connection := range s.connectionsById {
		if connection.deviceId == deviceId && connection.userId == userId {
			connections = append(connections, connection)
		}
	}

	return connections
}

func (s *ConnectionsStorage) GetConnectionById(connectionId ConnectionId) *Connection {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.connectionsById[connectionId]
}

func (s *ConnectionsStorage) GetStats() ConnectionsStats {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	stats := ConnectionsStats{
		// NumberOfDevices:              len(s.connectionsByDeviceId),
		// NumberOfUsers:                len(s.connectionsByUserId),
		NumberOfNotLoggedConnections: s.numberOfNotLoggedConnections,
	}

	return stats
}

func (s *ConnectionsStorage) RemoveIf(condition func(con *Connection) bool, afterRemove func(connections []*Connection)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	connections := []*Connection{}

	for id, connection := range s.connectionsById {
		if condition(connection) {
			delete(s.connectionsById, id)
			connections = append(connections, connection)
		}
	}

	afterRemove(connections)
}

func (s *ConnectionsStorage) RemoveDeviceConnections(userId UserId, deviceId DeviceId, afterRemove func(connections []*Connection)) {
	s.RemoveIf(func(con *Connection) bool {
		return con.deviceId == deviceId && con.userId == userId
	}, afterRemove)
}

func (s *ConnectionsStorage) RemoveUserConnections(userId UserId, afterRemove func(connections []*Connection)) {
	s.RemoveIf(func(con *Connection) bool {
		return con.userId == userId
	}, afterRemove)
}
