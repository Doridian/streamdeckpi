package impl

import (
	"errors"
	"time"

	"github.com/Doridian/streamdeckpi/agent/action"
	action_loader "github.com/Doridian/streamdeckpi/agent/action/loader"
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

	imageLoader := newPageImageLoader(c.imageLoader, pageObj)

	for _, actionSchema := range out.Actions {
		actionObj, err := action_loader.LoadAction(actionSchema.ActionName, &actionSchema.Parameters, imageLoader, c)
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

func (c *controllerImpl) SwapPage(pageFile string) error {
	pageObj, err := c.resolvePage(pageFile)
	if err != nil {
		return err
	}

	c.pageWait.Lock()
	defer c.pageWait.Unlock()

	c.unrefPage(c.pageTop)
	c.pageStack[len(c.pageStack)-1] = pageObj
	c.pageTop = pageObj

	return nil
}

func (c *controllerImpl) PushPage(pageFile string) error {
	pageObj, err := c.resolvePage(pageFile)
	if err != nil {
		return err
	}

	c.pageWait.Lock()
	defer c.pageWait.Unlock()

	c.pageStack = append(c.pageStack, pageObj)
	c.pageTop = pageObj

	return nil
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
	c.pageTop = c.pageStack[newLen-1]

	return nil
}
