package impl

import (
	"errors"
	"log"
	"time"

	"github.com/Doridian/streamdeckpi/agent/internal/action"
	action_loader "github.com/Doridian/streamdeckpi/agent/internal/action/loader"

	_ "github.com/Doridian/streamdeckpi/agent/internal/action/homeassistant"
	_ "github.com/Doridian/streamdeckpi/agent/internal/action/misc"
	_ "github.com/Doridian/streamdeckpi/agent/internal/action/page"
)

type page struct {
	path    string
	timeout time.Duration
	actions []action.Action
	refCnt  int
}

func (c *controllerImpl) resolvePage(pageFile string) (*page, error) {
	pageFile, err := c.CleanPath(pageFile)
	if err != nil {
		return nil, err
	}

	c.pageCacheLock.Lock()
	defer c.pageCacheLock.Unlock()

	cachedPage, ok := c.pageCache[pageFile]
	if ok {
		cachedPage.refCnt++
		return cachedPage, nil
	}

	out := &pageSchema{}
	err = c.LoadConfig(pageFile, out)
	if err != nil {
		return nil, err
	}

	actionLen := c.dev.Columns * c.dev.Rows

	pageObj := &page{
		path:    pageFile,
		timeout: out.Timeout,
		actions: make([]action.Action, actionLen),
		refCnt:  0,
	}

	imageHelper := newPageImageHelper(c.imageHelper, pageObj)

	for _, actionSchema := range out.Actions {
		actionObj, err := action_loader.LoadAction(actionSchema.ActionName, &actionSchema.Parameters, imageHelper, c)
		if err != nil {
			return nil, err
		}

		mappedButton := actionSchema.Button[0] + (actionSchema.Button[1] * int(c.dev.Columns))
		pageObj.actions[mappedButton] = actionObj
	}

	pageObj.refCnt++
	c.pageCache[pageFile] = pageObj
	return pageObj, nil
}

func (c *controllerImpl) unrefPage(pageObj *page) {
	c.pageCacheLock.Lock()
	pageObj.refCnt--
	c.pageCacheLock.Unlock()
}

func (c *controllerImpl) doPageSwap(newPage *page) error {
	oldPage := c.pageTop
	c.pageTop = newPage

	// Reset press status of old page
	for _, action := range oldPage.actions {
		go c.runActionHandleError(action, false)
	}
	// Reset press status of new page
	for _, action := range newPage.actions {
		go c.runActionHandleError(action, false)
	}

	return nil
}

func (c *controllerImpl) SwapPage(pageFile string) error {
	pageObj, err := c.resolvePage(pageFile)
	if err != nil {
		return err
	}

	c.pageWait.Lock()
	defer c.pageWait.Unlock()

	c.unrefPage(c.pageTop)
	c.pageStack[len(c.pageStack)-1] = pageObj
	return c.doPageSwap(pageObj)
}

func (c *controllerImpl) PushPage(pageFile string) error {
	pageObj, err := c.resolvePage(pageFile)
	if err != nil {
		return err
	}

	c.pageWait.Lock()
	defer c.pageWait.Unlock()

	c.pageStack = append(c.pageStack, pageObj)
	return c.doPageSwap(pageObj)
}

func (c *controllerImpl) PopPage() error {
	c.pageWait.Lock()
	defer c.pageWait.Unlock()

	newLen := len(c.pageStack) - 1
	if newLen < 1 {
		return errors.New("page stack is empty")
	}

	c.unrefPage(c.pageTop)
	c.pageStack = c.pageStack[:newLen]
	return c.doPageSwap(c.pageStack[newLen-1])
}

func (c *controllerImpl) PreloadPage(pageFile string) error {
	go func() {
		log.Printf("Preloading page %s", pageFile)
		_, err := c.resolvePage(pageFile)
		if err != nil {
			log.Printf("Error preloading page %s: %v", pageFile, err)
		} else {
			log.Printf("Preloaded page %s", pageFile)
		}
	}()
	return nil
}
