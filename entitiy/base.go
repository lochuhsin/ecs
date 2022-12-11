package entitiy

import "github.com/google/uuid"

type EntityManager struct {
	ent map[string]map[string]any
}

func (e *EntityManager) Init() {
	e.ent = make(map[string]map[string]any)
}
func (e *EntityManager) AddNewEntity(compKey string, compValue any) {
	uuidStr := uuid.New().String()
	e.ent[compKey][uuidStr] = compValue
}

func (e *EntityManager) RegisterComponentsById(uuidStr string, compKey string, compValue any) {
	if _, ok := e.ent[compKey][uuidStr]; ok {
		return
	}
	e.ent[compKey][uuidStr] = compValue
}

func (e *EntityManager) DeleteEntity(uuidStr string) {
	for _, ents := range e.ent {
		if _, ok := ents[uuidStr]; ok {
			delete(ents, uuidStr)
		}
	}
}

func (e *EntityManager) GetEntities(compKey string) map[string]any {
	return e.ent[compKey]
}
