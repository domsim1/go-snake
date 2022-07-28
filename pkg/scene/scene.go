package scene

import "fmt"

type Scene interface {
	Setup()
	Update()
	Draw()
}

type SceneManager interface {
	Add(id string, scene Scene) error
	Get(id string) (Scene, error)
	Delete(id string) error
	Active() Scene
	Activate(id string) (Scene, error)
}

type sceneManager struct {
	scenes map[string]Scene
	active Scene
}

func NewSceneManager() SceneManager {
	return &sceneManager{
		scenes: make(map[string]Scene),
	}
}

func (sm *sceneManager) Add(id string, scene Scene) error {
	if sm.scenes[id] != nil {
		return fmt.Errorf("scene with id %s already exists", id)
	}
	sm.scenes[id] = scene
	return nil
}

func (sm *sceneManager) Get(id string) (Scene, error) {
	scene := sm.scenes[id]
	if scene == nil {
		return nil, fmt.Errorf("scene with id %s does not exist", id)
	}
	return scene, nil
}

func (sm *sceneManager) Delete(id string) error {
	if sm.scenes[id] == nil {
		return fmt.Errorf("scene with id %s does not exist", id)
	}
	delete(sm.scenes, id)
	return nil
}

func (sm *sceneManager) Active() Scene {
	return sm.active
}

func (sm *sceneManager) Activate(id string) (Scene, error) {
	scene, err := sm.Get(id)
	if err != nil {
		return nil, err
	}
	sm.active = scene
	scene.Setup()
	return scene, nil
}
